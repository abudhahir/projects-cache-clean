package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	itemStyle = lipgloss.NewStyle().PaddingLeft(4)

	selectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(lipgloss.Color("170"))

	paginationStyle = list.DefaultStyles().PaginationStyle.PaddingLeft(4)

	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

	statsStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(1, 2).
			Margin(1, 0)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("46"))

	warningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("208"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196"))
)

type keyMap struct {
	Up         key.Binding
	Down       key.Binding
	Left       key.Binding
	Right      key.Binding
	Select     key.Binding
	Quit       key.Binding
	Help       key.Binding
	Clean      key.Binding
	Refresh    key.Binding
	Details    key.Binding
	SelectAll  key.Binding
	DeselectAll key.Binding
	ToggleView key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Select, k.Clean, k.Details, k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right, k.Select},
		{k.SelectAll, k.DeselectAll, k.ToggleView},
		{k.Clean, k.Details, k.Refresh, k.Help, k.Quit},
	}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "collapse/back"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "expand/forward"),
	),
	Select: key.NewBinding(
		key.WithKeys(" ", "enter"),
		key.WithHelp("space/enter", "select/expand"),
	),
	SelectAll: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "select all"),
	),
	DeselectAll: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "deselect all"),
	),
	ToggleView: key.NewBinding(
		key.WithKeys("t"),
		key.WithHelp("t", "toggle tree/list"),
	),
	Clean: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "clean selected"),
	),
	Details: key.NewBinding(
		key.WithKeys("v"),
		key.WithHelp("v", "view details"),
	),
	Refresh: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "refresh"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

type ProjectItem struct {
	Project     *Project
	Selected    bool
	CacheItems  []CacheItem
	TotalSize   int64
	ItemCount   int
}

func (p ProjectItem) FilterValue() string {
	return p.Project.Name
}

func (p ProjectItem) Title() string {
	icon := "○"
	if p.Selected {
		icon = "●"
	}
	
	sizeStr := formatBytes(p.TotalSize)
	countStr := fmt.Sprintf("%d items", p.ItemCount)
	
	status := successStyle.Render("✓ Clean")
	if p.ItemCount > 0 {
		status = warningStyle.Render(fmt.Sprintf("🗑 %s (%s)", countStr, sizeStr))
	}
	
	return fmt.Sprintf("%s %s [%s] - %s", icon, p.Project.Name, p.Project.Type, status)
}

func (p ProjectItem) Description() string {
	return fmt.Sprintf("Path: %s", p.Project.Path)
}

type Project struct {
	Name string
	Path string
	Type string
}

// TreeNode represents a node in the project directory tree
type TreeNode struct {
	Name           string         // Directory or project name
	Path           string         // Full path
	IsProject      bool           // true = project, false = directory
	Project        *ProjectItem   // Only set if IsProject == true
	Children       []*TreeNode    // Child nodes (subdirectories/projects)
	Parent         *TreeNode      // Parent directory node
	Expanded       bool           // Directory expansion state
	Selected       bool           // Selection state
	Level          int            // Depth level (for indentation)
	ChildProjects  int            // Number of projects in subtree
	ChildCacheSize int64          // Total cache size in subtree
}

// TreeModel manages the tree structure and display
type TreeModel struct {
	Root          *TreeNode      // Root directory node
	FlatView      []*TreeNode    // Flattened view for display
	CurrentIndex  int            // Currently selected item index
	VisibleStart  int            // First visible item index for scrolling
}

type AppState int

const (
	StateLoading AppState = iota
	StateProjectList
	StateDetails
	StateCleaning
	StateResults
	StateConfirm
)

type model struct {
	state       AppState
	keys        keyMap
	list        list.Model
	spinner     spinner.Model
	progress    progress.Model
	projects    []ProjectItem
	tree        *TreeModel      // Tree structure for projects
	useTreeView bool            // Toggle between tree and list view
	loading     bool
	err         error
	
	// Cleaning state
	cleaningIndex      int
	cleaningProgress   float64
	cleaningResults    CleanupStats
	currentProject     string
	totalProjects      int
	projectsCompleted  int
	
	// Details view
	detailsProject *ProjectItem
	
	// Confirmation
	confirmMessage string
	confirmAction  func() tea.Cmd
	
	width  int
	height int
}

type loadProjectsMsg struct {
	projects []ProjectItem
	err      error
}

type cleanProgressMsg struct {
	index          int
	progress       float64
	results        CleanupStats
	currentProject string
	totalProjects  int
	completed      bool
}

type cleanCompleteMsg struct {
	results CleanupStats
}


func initialModel(rootDir string) model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	items := []list.Item{}
	l := list.New(items, itemDelegate{}, 0, 0)
	l.Title = "🧹 Cache Remover - Project Scanner"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.KeyMap.Quit.SetKeys()
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{keys.Select, keys.Clean, keys.Help}
	}

	prog := progress.New(progress.WithDefaultGradient())

	m := model{
		state:       StateLoading,
		keys:        keys,
		list:        l,
		spinner:     s,
		progress:    prog,
		useTreeView: true, // Enable tree view by default
		loading:     true,
	}

	return m
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		loadProjects("."),
	)
}

func loadProjects(rootDir string) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		projects := []ProjectItem{}
		
		err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				// Log error but continue walking other paths
				return nil
			}
			
			if !info.IsDir() {
				return nil
			}
			
			// Skip deep nesting and hidden directories
			depth := strings.Count(strings.TrimPrefix(path, rootDir), string(os.PathSeparator))
			if depth > 10 || strings.Contains(path, "/.") {
				return filepath.SkipDir
			}
			
			// Skip descending into cache directories - they're meant to be removed as units
			if isCacheDirectory(info.Name()) {
				return filepath.SkipDir
			}
			
			if projectType := detectProjectType(path); projectType != nil {
				cacheItems := findCacheItems(path, projectType.CacheConfig)
				totalSize := int64(0)
				for _, item := range cacheItems {
					totalSize += item.Size
				}
				
				project := ProjectItem{
					Project: &Project{
						Name: filepath.Base(path),
						Path: path,
						Type: projectType.Name,
					},
					Selected:   false,
					CacheItems: cacheItems,
					TotalSize:  totalSize,
					ItemCount:  len(cacheItems),
				}
				projects = append(projects, project)
			}
			
			return nil
		})
		
		// Sort projects by cache size (largest first)
		sort.Slice(projects, func(i, j int) bool {
			return projects[i].TotalSize > projects[j].TotalSize
		})
		
		return loadProjectsMsg{projects: projects, err: err}
	})
}

// buildProjectTree creates a tree structure from a flat list of projects
func buildProjectTree(projects []ProjectItem, rootPath string) *TreeModel {
	// Create root node
	root := &TreeNode{
		Name:           filepath.Base(rootPath),
		Path:           rootPath,
		IsProject:      false,
		Children:       make([]*TreeNode, 0),
		Expanded:       true, // Root is always expanded
		Level:          0,
		ChildProjects:  0,
		ChildCacheSize: 0,
	}
	
	// Build tree by inserting each project
	for i := range projects {
		insertProjectIntoTree(root, &projects[i], rootPath)
	}
	
	// Calculate aggregate statistics
	calculateTreeStats(root)
	
	// Create tree model
	treeModel := &TreeModel{
		Root:         root,
		FlatView:     make([]*TreeNode, 0),
		CurrentIndex: 0,
		VisibleStart: 0,
	}
	
	// Generate initial flat view
	treeModel.rebuildFlatView()
	
	return treeModel
}

// insertProjectIntoTree inserts a project into the appropriate position in the tree
func insertProjectIntoTree(root *TreeNode, project *ProjectItem, rootPath string) {
	// Get relative path from root
	relativePath, err := filepath.Rel(rootPath, project.Project.Path)
	if err != nil {
		relativePath = project.Project.Path
	}
	
	// Split path into components
	pathComponents := strings.Split(relativePath, string(filepath.Separator))
	if pathComponents[0] == "." {
		pathComponents = pathComponents[1:]
	}
	
	// Navigate/create directory structure
	currentNode := root
	currentPath := rootPath
	
	// Create intermediate directory nodes
	for i, component := range pathComponents {
		if i == len(pathComponents)-1 {
			// This is the project itself
			projectNode := &TreeNode{
				Name:           component,
				Path:           project.Project.Path,
				IsProject:      true,
				Project:        project,
				Children:       make([]*TreeNode, 0),
				Parent:         currentNode,
				Selected:       project.Selected,
				Level:          currentNode.Level + 1,
				ChildProjects:  1,
				ChildCacheSize: project.TotalSize,
			}
			currentNode.Children = append(currentNode.Children, projectNode)
		} else {
			// This is an intermediate directory
			currentPath = filepath.Join(currentPath, component)
			
			// Check if directory node already exists
			var dirNode *TreeNode
			for _, child := range currentNode.Children {
				if !child.IsProject && child.Name == component {
					dirNode = child
					break
				}
			}
			
			// Create directory node if it doesn't exist
			if dirNode == nil {
				dirNode = &TreeNode{
					Name:           component,
					Path:           currentPath,
					IsProject:      false,
					Children:       make([]*TreeNode, 0),
					Parent:         currentNode,
					Expanded:       false, // Directories start collapsed
					Level:          currentNode.Level + 1,
					ChildProjects:  0,
					ChildCacheSize: 0,
				}
				currentNode.Children = append(currentNode.Children, dirNode)
			}
			
			currentNode = dirNode
		}
	}
}

// calculateTreeStats calculates aggregate statistics for each directory node
func calculateTreeStats(node *TreeNode) {
	if node.IsProject {
		// Project nodes already have their stats set
		return
	}
	
	// Reset stats
	node.ChildProjects = 0
	node.ChildCacheSize = 0
	
	// Calculate stats from children
	for _, child := range node.Children {
		calculateTreeStats(child) // Recursive call
		
		if child.IsProject {
			node.ChildProjects += 1
			node.ChildCacheSize += child.Project.TotalSize
		} else {
			node.ChildProjects += child.ChildProjects
			node.ChildCacheSize += child.ChildCacheSize
		}
	}
}

// rebuildFlatView creates a flattened view of the tree for display purposes
func (tm *TreeModel) rebuildFlatView() {
	tm.FlatView = tm.FlatView[:0] // Clear existing view
	tm.flattenNode(tm.Root)
}

// flattenNode recursively flattens tree nodes that should be visible
func (tm *TreeModel) flattenNode(node *TreeNode) {
	if node == tm.Root {
		// Don't show root in the display, but process its children
		for _, child := range node.Children {
			tm.flattenNode(child)
		}
		return
	}
	
	// Add this node to flat view
	tm.FlatView = append(tm.FlatView, node)
	
	// If it's a directory and expanded, add its children
	if !node.IsProject && node.Expanded {
		for _, child := range node.Children {
			tm.flattenNode(child)
		}
	}
}

// Tree navigation functions

// getCurrentNode returns the currently selected node
func (tm *TreeModel) getCurrentNode() *TreeNode {
	if tm.CurrentIndex < 0 || tm.CurrentIndex >= len(tm.FlatView) {
		return nil
	}
	return tm.FlatView[tm.CurrentIndex]
}

// moveUp moves the selection up by one item
func (tm *TreeModel) moveUp() {
	if tm.CurrentIndex > 0 {
		tm.CurrentIndex--
	}
}

// moveDown moves the selection down by one item
func (tm *TreeModel) moveDown() {
	if tm.CurrentIndex < len(tm.FlatView)-1 {
		tm.CurrentIndex++
	}
}

// toggleExpansion expands or collapses the current directory node
func (tm *TreeModel) toggleExpansion() {
	node := tm.getCurrentNode()
	if node != nil && !node.IsProject && len(node.Children) > 0 {
		node.Expanded = !node.Expanded
		tm.rebuildFlatView()
		
		// Ensure current index is still valid after rebuild
		if tm.CurrentIndex >= len(tm.FlatView) {
			tm.CurrentIndex = len(tm.FlatView) - 1
		}
	}
}

// expandNode expands the current directory node
func (tm *TreeModel) expandNode() {
	node := tm.getCurrentNode()
	if node != nil && !node.IsProject && len(node.Children) > 0 && !node.Expanded {
		node.Expanded = true
		tm.rebuildFlatView()
	}
}

// collapseNode collapses the current directory node
func (tm *TreeModel) collapseNode() {
	node := tm.getCurrentNode()
	if node != nil && !node.IsProject && node.Expanded {
		node.Expanded = false
		tm.rebuildFlatView()
		
		// Ensure current index is still valid after rebuild
		if tm.CurrentIndex >= len(tm.FlatView) {
			tm.CurrentIndex = len(tm.FlatView) - 1
		}
	}
}

// toggleSelection toggles the selection state of the current node
func (tm *TreeModel) toggleSelection() {
	node := tm.getCurrentNode()
	if node != nil {
		if node.IsProject {
			// Toggle project selection
			node.Selected = !node.Selected
			node.Project.Selected = node.Selected
		} else {
			// Toggle directory selection (affects all child projects)
			newSelectionState := !tm.isDirectorySelected(node)
			tm.setDirectorySelection(node, newSelectionState)
		}
	}
}

// isDirectorySelected checks if a directory should be considered selected
// (true if all child projects are selected)
func (tm *TreeModel) isDirectorySelected(node *TreeNode) bool {
	if node.IsProject {
		return node.Selected
	}
	
	hasProjects := false
	for _, child := range node.Children {
		if child.IsProject {
			hasProjects = true
			if !child.Selected {
				return false
			}
		} else if tm.hasProjectsInSubtree(child) {
			hasProjects = true
			if !tm.isDirectorySelected(child) {
				return false
			}
		}
	}
	
	return hasProjects // Return true only if has projects and all are selected
}

// hasProjectsInSubtree checks if a directory has any projects in its subtree
func (tm *TreeModel) hasProjectsInSubtree(node *TreeNode) bool {
	if node.IsProject {
		return true
	}
	
	for _, child := range node.Children {
		if tm.hasProjectsInSubtree(child) {
			return true
		}
	}
	
	return false
}

// setDirectorySelection sets the selection state for all projects in a directory
func (tm *TreeModel) setDirectorySelection(node *TreeNode, selected bool) {
	if node.IsProject {
		node.Selected = selected
		node.Project.Selected = selected
		return
	}
	
	for _, child := range node.Children {
		tm.setDirectorySelection(child, selected)
	}
}

// getSelectedProjects returns all currently selected projects
func (tm *TreeModel) getSelectedProjects() []ProjectItem {
	var selected []ProjectItem
	tm.collectSelectedProjects(tm.Root, &selected)
	return selected
}

// collectSelectedProjects recursively collects selected projects
func (tm *TreeModel) collectSelectedProjects(node *TreeNode, selected *[]ProjectItem) {
	if node.IsProject && node.Selected {
		*selected = append(*selected, *node.Project)
		return
	}
	
	for _, child := range node.Children {
		tm.collectSelectedProjects(child, selected)
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height - 4)
		m.progress.Width = msg.Width - 4

	case loadProjectsMsg:
		m.loading = false
		m.err = msg.err
		m.projects = msg.projects
		
		// Build tree structure from projects
		if len(m.projects) > 0 {
			// Use the common ancestor path as root, or current directory
			rootPath := "."
			if len(m.projects) > 0 {
				rootPath = filepath.Dir(m.projects[0].Project.Path)
				// Find common root for all projects
				for _, project := range m.projects[1:] {
					for !strings.HasPrefix(project.Project.Path, rootPath) {
						rootPath = filepath.Dir(rootPath)
						if rootPath == "/" || rootPath == "." {
							break
						}
					}
				}
			}
			m.tree = buildProjectTree(m.projects, rootPath)
		}
		
		// Set up list view (for backward compatibility)
		items := make([]list.Item, len(m.projects))
		for i, project := range m.projects {
			items[i] = project
		}
		
		m.list.SetItems(items)
		m.state = StateProjectList
		
		totalSize := int64(0)
		totalItems := 0
		for _, p := range m.projects {
			totalSize += p.TotalSize
			totalItems += p.ItemCount
		}
		
		m.list.Title = fmt.Sprintf("🧹 Cache Remover - %d Projects (%s potential cleanup)",
			len(m.projects), formatBytes(totalSize))

	case cleanProgressMsg:
		m.cleaningIndex = msg.index
		m.cleaningProgress = msg.progress
		m.cleaningResults = msg.results
		m.currentProject = msg.currentProject
		m.totalProjects = msg.totalProjects
		if msg.completed {
			m.projectsCompleted = msg.index
		}

	case cleanCompleteMsg:
		m.cleaningResults = msg.results
		m.state = StateResults

	case tea.KeyMsg:
		switch m.state {
		case StateProjectList:
			switch {
			case key.Matches(msg, m.keys.Quit):
				return m, tea.Quit
				
			case key.Matches(msg, m.keys.ToggleView):
				// Toggle between tree and list view
				m.useTreeView = !m.useTreeView
				
			case key.Matches(msg, m.keys.Up):
				if m.useTreeView && m.tree != nil {
					m.tree.moveUp()
				}
				
			case key.Matches(msg, m.keys.Down):
				if m.useTreeView && m.tree != nil {
					m.tree.moveDown()
				}
				
			case key.Matches(msg, m.keys.Left):
				if m.useTreeView && m.tree != nil {
					m.tree.collapseNode()
				}
				
			case key.Matches(msg, m.keys.Right):
				if m.useTreeView && m.tree != nil {
					m.tree.expandNode()
				}
				
			case key.Matches(msg, m.keys.Select):
				if m.useTreeView && m.tree != nil {
					// Tree view: toggle selection or expand/collapse
					node := m.tree.getCurrentNode()
					if node != nil {
						if node.IsProject {
							m.tree.toggleSelection()
						} else {
							m.tree.toggleExpansion()
						}
					}
				} else if selectedItem, ok := m.list.SelectedItem().(ProjectItem); ok {
					// List view: toggle selection
					for i, project := range m.projects {
						if project.Project.Path == selectedItem.Project.Path {
							m.projects[i].Selected = !m.projects[i].Selected
							// Update the list item
							items := m.list.Items()
							items[m.list.Index()] = m.projects[i]
							m.list.SetItems(items)
							break
						}
					}
				}
				
			case key.Matches(msg, m.keys.SelectAll):
				if m.useTreeView && m.tree != nil {
					// Tree view: select all projects
					m.tree.setDirectorySelection(m.tree.Root, true)
				} else {
					// List view: select all projects
					for i := range m.projects {
						if m.projects[i].ItemCount > 0 {
							m.projects[i].Selected = true
						}
					}
					items := make([]list.Item, len(m.projects))
					for i, project := range m.projects {
						items[i] = project
					}
					m.list.SetItems(items)
				}
				
			case key.Matches(msg, m.keys.DeselectAll):
				if m.useTreeView && m.tree != nil {
					// Tree view: deselect all projects
					m.tree.setDirectorySelection(m.tree.Root, false)
				} else {
					// List view: deselect all projects
					for i := range m.projects {
						m.projects[i].Selected = false
					}
					items := make([]list.Item, len(m.projects))
					for i, project := range m.projects {
						items[i] = project
					}
					m.list.SetItems(items)
				}
				
			case key.Matches(msg, m.keys.Details):
				var detailProject *ProjectItem
				
				if m.useTreeView && m.tree != nil {
					// Tree view: get current project
					node := m.tree.getCurrentNode()
					if node != nil && node.IsProject {
						detailProject = node.Project
					}
				} else if selectedItem, ok := m.list.SelectedItem().(ProjectItem); ok {
					// List view: get selected project
					for i, project := range m.projects {
						if project.Project.Path == selectedItem.Project.Path {
							detailProject = &m.projects[i]
							break
						}
					}
				}
				
				if detailProject != nil {
					m.detailsProject = detailProject
					m.state = StateDetails
				}
				
			case key.Matches(msg, m.keys.Clean):
				var selectedProjects []ProjectItem
				
				if m.useTreeView && m.tree != nil {
					// Tree view: get selected projects from tree
					selectedProjects = m.tree.getSelectedProjects()
				} else {
					// List view: get selected projects from list
					for _, project := range m.projects {
						if project.Selected {
							selectedProjects = append(selectedProjects, project)
						}
					}
				}
				
				if len(selectedProjects) > 0 {
					totalSize := int64(0)
					totalItems := 0
					for _, project := range selectedProjects {
						totalSize += project.TotalSize
						totalItems += project.ItemCount
					}
					
					m.confirmMessage = fmt.Sprintf(
						"Clean %d projects?\nThis will remove %d cache items (%s)\n\nPress 'y' to confirm, 'n' to cancel",
						len(selectedProjects), totalItems, formatBytes(totalSize))
					m.confirmAction = func() tea.Cmd {
						// Initialize progress tracking
						m.totalProjects = len(selectedProjects)
						m.projectsCompleted = 0
						m.currentProject = "Starting..."
						m.cleaningProgress = 0.0
						return cleanSelectedProjects(selectedProjects)
					}
					m.state = StateConfirm
				}
				
			case key.Matches(msg, m.keys.Refresh):
				m.loading = true
				m.state = StateLoading
				return m, tea.Batch(m.spinner.Tick, loadProjects("."))
			}
			
		case StateDetails:
			switch {
			case key.Matches(msg, m.keys.Quit):
				return m, tea.Quit
			case msg.String() == "esc":
				m.state = StateProjectList
			}
			
		case StateConfirm:
			switch msg.String() {
			case "y", "Y":
				m.state = StateCleaning
				m.cleaningIndex = 0
				m.cleaningProgress = 0
				return m, m.confirmAction()
			case "n", "N", "esc":
				m.state = StateProjectList
			}
			
		case StateCleaning:
			switch {
			case key.Matches(msg, m.keys.Quit):
				return m, tea.Quit
			}
			
		case StateResults:
			switch {
			case key.Matches(msg, m.keys.Quit):
				return m, tea.Quit
			case msg.String() == "esc", msg.String() == "enter":
				m.state = StateProjectList
				// Refresh the project list
				m.loading = true
				m.state = StateLoading
				return m, tea.Batch(m.spinner.Tick, loadProjects("."))
			}
		}
	}

	// Update components
	switch m.state {
	case StateLoading:
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	case StateProjectList:
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)
	case StateCleaning:
		// Update both spinner and progress bar during cleaning
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
		
		var progressModel tea.Model
		progressModel, cmd = m.progress.Update(msg)
		if pm, ok := progressModel.(progress.Model); ok {
			m.progress = pm
		}
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func cleanSelectedProjects(projects []ProjectItem) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		var wg sync.WaitGroup
		results := CleanupStats{}
		
		for _, project := range projects {
			// Send progress update
			tea.Printf("Cleaning %s...\n", project.Project.Name)
			
			removedItems, removedSize := removeCacheItems(project.CacheItems, false)
			results.TotalCacheItems += removedItems
			results.TotalSizeRemoved += removedSize
			results.TotalProjects++
		}
		
		wg.Wait()
		return cleanCompleteMsg{results: results}
	})
}

func (m model) View() string {
	switch m.state {
	case StateLoading:
		return fmt.Sprintf("\n\n   %s Loading projects...\n\n", m.spinner.View())
		
	case StateProjectList:
		if m.useTreeView && m.tree != nil {
			return m.renderTreeView()
		} else {
			// Original list view
			selectedCount := 0
			selectedSize := int64(0)
			for _, project := range m.projects {
				if project.Selected {
					selectedCount++
					selectedSize += project.TotalSize
				}
			}
			
			statusBar := ""
			if selectedCount > 0 {
				statusBar = warningStyle.Render(fmt.Sprintf(
					" Selected: %d projects (%s) - Press 'c' to clean ",
					selectedCount, formatBytes(selectedSize)))
			} else {
				statusBar = helpStyle.Render(" Use ↑/↓ to navigate, Space to select, 'c' to clean, 't' to toggle tree view ")
			}
			
			return m.list.View() + "\n" + statusBar
		}
		
	case StateDetails:
		if m.detailsProject == nil {
			return "No project selected"
		}
		
		details := fmt.Sprintf("📁 %s (%s)\n", m.detailsProject.Project.Name, m.detailsProject.Project.Type)
		details += fmt.Sprintf("Path: %s\n\n", m.detailsProject.Project.Path)
		details += fmt.Sprintf("Cache Items (%d):\n", len(m.detailsProject.CacheItems))
		
		for _, item := range m.detailsProject.CacheItems {
			itemType := "📄"
			if item.Type == "directory" {
				itemType = "📁"
			}
			details += fmt.Sprintf("  %s %s (%s)\n", itemType, filepath.Base(item.Path), formatBytes(item.Size))
		}
		
		details += fmt.Sprintf("\nTotal Size: %s\n", formatBytes(m.detailsProject.TotalSize))
		details += helpStyle.Render("\nPress ESC to go back")
		
		return details
		
	case StateConfirm:
		return fmt.Sprintf("\n%s\n", warningStyle.Render(m.confirmMessage))
		
	case StateCleaning:
		// Create animated spinner for active cleanup
		spinnerView := m.spinner.View()
		
		// Status message
		status := fmt.Sprintf("%s Cleaning cache files...", spinnerView)
		if m.currentProject != "" && m.currentProject != "Starting..." {
			status = fmt.Sprintf("%s Cleaning: %s", spinnerView, m.currentProject)
		}
		
		// Progress information
		progressInfo := "Initializing cleanup..."
		if m.totalProjects > 0 {
			progressInfo = fmt.Sprintf("Project %d of %d", m.projectsCompleted+1, m.totalProjects)
		}
		
		// Progress bar - always show some progress during cleanup
		displayProgress := m.cleaningProgress
		if displayProgress == 0.0 && m.totalProjects > 0 {
			// Show some progress even when starting
			displayProgress = 0.1
		}
		progressBar := m.progress.ViewAs(displayProgress)
		
		// Current statistics
		stats := fmt.Sprintf(
			"📊 Status: Processing cache directories\n"+
			"🗑️  Items removed: %d\n"+
			"💾 Space reclaimed: %s",
			m.cleaningResults.TotalCacheItems,
			formatBytes(m.cleaningResults.TotalSizeRemoved))
		
		return fmt.Sprintf("\n\n   %s\n   %s\n\n   %s\n\n   %.1f%% Complete\n\n%s\n\n   💡 Tip: This may take a moment for large cache directories\n   Press 'q' to cancel\n\n",
			status,
			progressInfo,
			progressBar,
			displayProgress*100,
			stats)
		
	case StateResults:
		results := statsStyle.Render(fmt.Sprintf(
			"✅ Cleanup Complete!\n\n"+
				"Projects cleaned: %d\n"+
				"Cache items removed: %d\n"+
				"Space reclaimed: %s\n\n"+
				"Press ENTER to continue",
			m.cleaningResults.TotalProjects,
			m.cleaningResults.TotalCacheItems,
			formatBytes(m.cleaningResults.TotalSizeRemoved)))
		
		return fmt.Sprintf("\n%s\n", results)
	}
	
	return ""
}

// renderTreeView renders the tree structure with proper styling and navigation
func (m model) renderTreeView() string {
	if m.tree == nil {
		return "No tree data available"
	}
	
	// Calculate selected projects and total size
	selectedProjects := m.tree.getSelectedProjects()
	selectedSize := int64(0)
	for _, project := range selectedProjects {
		selectedSize += project.TotalSize
	}
	
	// Build the tree display
	var output strings.Builder
	
	// Title with tree icon
	title := titleStyle.Render("🌳 Cache Remover - Tree View")
	output.WriteString(title + "\n\n")
	
	// Render visible tree nodes
	visibleHeight := m.height - 6 // Leave space for title and status bar
	startIdx := 0
	endIdx := len(m.tree.FlatView)
	
	// Simple scrolling logic
	if len(m.tree.FlatView) > visibleHeight {
		if m.tree.CurrentIndex >= visibleHeight/2 {
			startIdx = m.tree.CurrentIndex - visibleHeight/2
		}
		endIdx = startIdx + visibleHeight
		if endIdx > len(m.tree.FlatView) {
			endIdx = len(m.tree.FlatView)
			startIdx = endIdx - visibleHeight
			if startIdx < 0 {
				startIdx = 0
			}
		}
	}
	
	// Render each visible node
	for i := startIdx; i < endIdx; i++ {
		if i >= len(m.tree.FlatView) {
			break
		}
		
		node := m.tree.FlatView[i]
		line := m.renderTreeNode(node, i == m.tree.CurrentIndex)
		output.WriteString(line + "\n")
	}
	
	// Status bar
	var statusBar string
	if len(selectedProjects) > 0 {
		statusBar = warningStyle.Render(fmt.Sprintf(
			" Selected: %d projects (%s) - Press 'c' to clean ",
			len(selectedProjects), formatBytes(selectedSize)))
	} else {
		statusBar = helpStyle.Render(" ↑/↓:navigate ←/→:expand/collapse Space:select 'c':clean 't':list view ")
	}
	
	output.WriteString("\n" + statusBar)
	
	return output.String()
}

// renderTreeNode renders a single tree node with proper indentation and styling
func (m model) renderTreeNode(node *TreeNode, isSelected bool) string {
	var output strings.Builder
	
	// Add indentation based on level
	indent := strings.Repeat("  ", node.Level)
	output.WriteString(indent)
	
	// Add tree structure symbols
	if node.Level > 0 {
		output.WriteString("├─ ")
	}
	
	if node.IsProject {
		// Render project node
		selectionIcon := "○"
		if node.Selected {
			selectionIcon = "●"
		}
		
		projectInfo := fmt.Sprintf("%s %s [%s]", 
			selectionIcon, 
			node.Name, 
			node.Project.Project.Type)
		
		// Add cache information
		if node.Project.ItemCount > 0 {
			cacheInfo := fmt.Sprintf(" - 🗑 %d items (%s)", 
				node.Project.ItemCount, 
				formatBytes(node.Project.TotalSize))
			projectInfo += warningStyle.Render(cacheInfo)
		} else {
			projectInfo += successStyle.Render(" - ✓ Clean")
		}
		
		// Highlight current selection
		if isSelected {
			projectInfo = selectedItemStyle.Render("> " + projectInfo)
		} else {
			projectInfo = itemStyle.Render(projectInfo)
		}
		
		output.WriteString(projectInfo)
		
	} else {
		// Render directory node
		var dirIcon string
		if len(node.Children) == 0 {
			dirIcon = "📁"
		} else if node.Expanded {
			dirIcon = "📂"
		} else {
			dirIcon = "📁"
		}
		
		// Check if directory should be considered selected
		isDirectorySelected := m.tree.isDirectorySelected(node)
		selectionIcon := "○"
		if isDirectorySelected {
			selectionIcon = "●"
		}
		
		directoryInfo := fmt.Sprintf("%s %s %s", 
			selectionIcon,
			dirIcon, 
			node.Name)
		
		// Add aggregate statistics
		if node.ChildProjects > 0 {
			stats := fmt.Sprintf(" [%d projects, %s]", 
				node.ChildProjects, 
				formatBytes(node.ChildCacheSize))
			directoryInfo += helpStyle.Render(stats)
		}
		
		// Highlight current selection
		if isSelected {
			directoryInfo = selectedItemStyle.Render("> " + directoryInfo)
		} else {
			directoryInfo = itemStyle.Render(directoryInfo)
		}
		
		output.WriteString(directoryInfo)
	}
	
	return output.String()
}

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 2 }
func (d itemDelegate) Spacing() int                            { return 1 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(ProjectItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i.Title())

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
	fmt.Fprint(w, "\n"+itemStyle.Render(i.Description()))
}

func runInteractiveUI(rootDir string) error {
	p := tea.NewProgram(initialModel(rootDir), tea.WithAltScreen())
	_, err := p.Run()
	return err
}