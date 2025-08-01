package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

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
	Up       key.Binding
	Down     key.Binding
	Select   key.Binding
	Quit     key.Binding
	Help     key.Binding
	Clean    key.Binding
	Refresh  key.Binding
	Details  key.Binding
	SelectAll key.Binding
	DeselectAll key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Select, k.Clean, k.Details, k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Select, k.SelectAll, k.DeselectAll},
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
	Select: key.NewBinding(
		key.WithKeys(" ", "enter"),
		key.WithHelp("space/enter", "select/deselect"),
	),
	SelectAll: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "select all"),
	),
	DeselectAll: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "deselect all"),
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

// Step B: Cleanup state tracker for progress updates
type cleanupState struct {
	projects       []ProjectItem
	currentIndex   int
	currentProject string
	results        CleanupStats
	isActive       bool
	mutex          sync.RWMutex
}

var globalCleanupState *cleanupState

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
		state:    StateLoading,
		keys:     keys,
		list:     l,
		spinner:  s,
		progress: prog,
		loading:  true,
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
				
			case key.Matches(msg, m.keys.Select):
				if selectedItem, ok := m.list.SelectedItem().(ProjectItem); ok {
					// Toggle selection
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
				
			case key.Matches(msg, m.keys.DeselectAll):
				for i := range m.projects {
					m.projects[i].Selected = false
				}
				items := make([]list.Item, len(m.projects))
				for i, project := range m.projects {
					items[i] = project
				}
				m.list.SetItems(items)
				
			case key.Matches(msg, m.keys.Details):
				if selectedItem, ok := m.list.SelectedItem().(ProjectItem); ok {
					for i, project := range m.projects {
						if project.Project.Path == selectedItem.Project.Path {
							m.detailsProject = &m.projects[i]
							m.state = StateDetails
							break
						}
					}
				}
				
			case key.Matches(msg, m.keys.Clean):
				selectedProjects := []ProjectItem{}
				totalSize := int64(0)
				totalItems := 0
				
				for _, project := range m.projects {
					if project.Selected {
						selectedProjects = append(selectedProjects, project)
						totalSize += project.TotalSize
						totalItems += project.ItemCount
					}
				}
				
				if len(selectedProjects) > 0 {
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

// Step B: Cleanup with shared state and periodic progress updates
func cleanSelectedProjects(projects []ProjectItem) tea.Cmd {
	// Initialize global cleanup state
	globalCleanupState = &cleanupState{
		projects:       projects,
		currentIndex:   0,
		currentProject: "Starting cleanup...",
		results:        CleanupStats{},
		isActive:       true,
	}
	
	return tea.Batch(
		// Send initial progress message immediately
		func() tea.Msg {
			return cleanProgressMsg{
				index:          0,
				progress:       0.0,
				results:        CleanupStats{},
				currentProject: "Starting cleanup...",
				totalProjects:  len(projects),
				completed:      false,
			}
		},
		// Start periodic progress updates
		tea.Every(200*time.Millisecond, func(t time.Time) tea.Msg {
			if globalCleanupState == nil {
				return nil
			}
			
			globalCleanupState.mutex.RLock()
			defer globalCleanupState.mutex.RUnlock()
			
			if !globalCleanupState.isActive {
				return nil
			}
			
			return cleanProgressMsg{
				index:          globalCleanupState.currentIndex,
				progress:       float64(globalCleanupState.currentIndex) / float64(len(globalCleanupState.projects)),
				results:        globalCleanupState.results,
				currentProject: globalCleanupState.currentProject,
				totalProjects:  len(globalCleanupState.projects),
				completed:      globalCleanupState.currentIndex >= len(globalCleanupState.projects),
			}
		}),
		// Start the actual cleanup process
		func() tea.Msg {
			results := CleanupStats{}
			
			for i, project := range projects {
				// Update shared state
				globalCleanupState.mutex.Lock()
				globalCleanupState.currentIndex = i
				globalCleanupState.currentProject = project.Project.Name
				globalCleanupState.mutex.Unlock()
				
				// Perform cleanup for each project
				removedItems, removedSize := removeCacheItems(project.CacheItems, false)
				results.TotalCacheItems += removedItems
				results.TotalSizeRemoved += removedSize
				results.TotalProjects++
				
				// Update results in shared state
				globalCleanupState.mutex.Lock()
				globalCleanupState.results = results
				globalCleanupState.mutex.Unlock()
			}
			
			// Mark cleanup as complete
			globalCleanupState.mutex.Lock()
			globalCleanupState.isActive = false
			globalCleanupState.mutex.Unlock()
			
			return cleanCompleteMsg{results: results}
		},
	)
}

func (m model) View() string {
	switch m.state {
	case StateLoading:
		return fmt.Sprintf("\n\n   %s Loading projects...\n\n", m.spinner.View())
		
	case StateProjectList:
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
			statusBar = helpStyle.Render(" Use ↑/↓ to navigate, Space to select, 'c' to clean, '?' for help ")
		}
		
		return m.list.View() + "\n" + statusBar
		
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