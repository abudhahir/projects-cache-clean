package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type WebProject struct {
	ID        int           `json:"id"`
	Name      string        `json:"name"`
	Path      string        `json:"path"`
	Type      string        `json:"type"`
	CacheSize int64         `json:"cacheSize"`
	ItemCount int           `json:"itemCount"`
	Items     []WebCacheItem `json:"items"`
	Selected  bool          `json:"selected"`
}

type WebCacheItem struct {
	Path string `json:"path"`
	Size int64  `json:"size"`
	Type string `json:"type"`
}

type WebServer struct {
	projects []WebProject
	mutex    sync.RWMutex
}

func (ws *WebServer) loadProjects(rootDir string) {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()
	
	ws.projects = []WebProject{}
	projectID := 0
	
	filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || !info.IsDir() {
			return nil
		}
		
		depth := strings.Count(strings.TrimPrefix(path, rootDir), string(filepath.Separator))
		if depth > 10 || strings.Contains(path, "/.") {
			return filepath.SkipDir
		}
		
		if projectType := detectProjectType(path); projectType != nil {
			cacheItems := findCacheItems(path, projectType.CacheConfig)
			totalSize := int64(0)
			webItems := []WebCacheItem{}
			
			for _, item := range cacheItems {
				totalSize += item.Size
				webItems = append(webItems, WebCacheItem{
					Path: item.Path,
					Size: item.Size,
					Type: item.Type,
				})
			}
			
			ws.projects = append(ws.projects, WebProject{
				ID:        projectID,
				Name:      filepath.Base(path),
				Path:      path,
				Type:      projectType.Name,
				CacheSize: totalSize,
				ItemCount: len(cacheItems),
				Items:     webItems,
				Selected:  false,
			})
			projectID++
		}
		
		return nil
	})
}

func (ws *WebServer) handleIndex(w http.ResponseWriter, r *http.Request) {
	tmpl := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>🧹 Cache Remover Utility</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background: #f5f5f5; }
        .container { max-width: 1200px; margin: 0 auto; padding: 20px; }
        .header { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; padding: 30px; border-radius: 10px; margin-bottom: 30px; text-align: center; }
        .stats { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 20px; margin-bottom: 30px; }
        .stat-card { background: white; padding: 20px; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); text-align: center; }
        .stat-value { font-size: 2em; font-weight: bold; color: #667eea; }
        .controls { background: white; padding: 20px; border-radius: 10px; margin-bottom: 20px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        .btn { padding: 10px 20px; border: none; border-radius: 5px; cursor: pointer; font-size: 14px; margin: 5px; transition: all 0.3s; }
        .btn-primary { background: #667eea; color: white; }
        .btn-success { background: #28a745; color: white; }
        .btn-danger { background: #dc3545; color: white; }
        .btn-secondary { background: #6c757d; color: white; }
        .btn:hover { transform: translateY(-2px); box-shadow: 0 4px 15px rgba(0,0,0,0.2); }
        .project-grid { display: grid; gap: 20px; }
        .project-card { background: white; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); overflow: hidden; transition: all 0.3s; }
        .project-card:hover { transform: translateY(-5px); box-shadow: 0 10px 25px rgba(0,0,0,0.15); }
        .project-header { padding: 20px; border-bottom: 1px solid #eee; display: flex; justify-content: space-between; align-items: center; }
        .project-info h3 { color: #333; margin-bottom: 5px; }
        .project-type { color: #667eea; font-weight: 500; }
        .project-path { color: #666; font-size: 0.9em; margin-top: 5px; }
        .project-stats { display: flex; gap: 20px; margin-top: 10px; }
        .project-stat { text-align: center; }
        .project-stat-value { font-weight: bold; color: #28a745; }
        .project-items { padding: 0 20px 20px; }
        .cache-item { display: flex; justify-content: between; align-items: center; padding: 8px 0; border-bottom: 1px solid #f0f0f0; }
        .cache-item:last-child { border-bottom: none; }
        .item-icon { margin-right: 8px; }
        .item-size { color: #666; font-size: 0.9em; }
        .checkbox { width: 20px; height: 20px; cursor: pointer; }
        .selected { background-color: #e3f2fd; }
        .progress-bar { width: 100%; height: 20px; background: #f0f0f0; border-radius: 10px; overflow: hidden; margin: 10px 0; }
        .progress-fill { height: 100%; background: linear-gradient(90deg, #667eea, #28a745); transition: width 0.3s; }
        .hidden { display: none; }
        #cleaningModal { position: fixed; top: 0; left: 0; width: 100%; height: 100%; background: rgba(0,0,0,0.8); display: flex; justify-content: center; align-items: center; z-index: 1000; }
        .modal-content { background: white; padding: 30px; border-radius: 10px; text-align: center; max-width: 500px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🧹 Cache Remover Utility</h1>
            <p>Efficiently clean up project cache files with visual controls</p>
        </div>

        <div class="stats">
            <div class="stat-card">
                <div class="stat-value" id="totalProjects">0</div>
                <div>Projects Found</div>
            </div>
            <div class="stat-card">
                <div class="stat-value" id="totalSize">0 B</div>
                <div>Total Cache Size</div>
            </div>
            <div class="stat-card">
                <div class="stat-value" id="selectedProjects">0</div>
                <div>Selected Projects</div>
            </div>
            <div class="stat-card">
                <div class="stat-value" id="selectedSize">0 B</div>
                <div>Selected Cache Size</div>
            </div>
        </div>

        <div class="controls">
            <button class="btn btn-primary" onclick="refreshProjects()">🔄 Refresh</button>
            <button class="btn btn-secondary" onclick="selectAll()">✅ Select All</button>
            <button class="btn btn-secondary" onclick="deselectAll()">❌ Deselect All</button>
            <button class="btn btn-success" onclick="showCleanDialog()" id="cleanBtn" disabled>🧹 Clean Selected</button>
        </div>

        <div class="project-grid" id="projectGrid">
            <!-- Projects will be loaded here -->
        </div>
    </div>

    <!-- Cleaning Modal -->
    <div id="cleaningModal" class="hidden">
        <div class="modal-content">
            <h3>🧹 Cleaning Projects...</h3>
            <div class="progress-bar">
                <div class="progress-fill" id="progressFill" style="width: 0%"></div>
            </div>
            <p id="cleaningStatus">Preparing to clean...</p>
            <p id="cleaningResults" class="hidden"></p>
            <button class="btn btn-primary hidden" id="doneBtn" onclick="closeCleaningModal()">Done</button>
        </div>
    </div>

    <script>
        let projects = [];

        async function loadProjects() {
            try {
                const response = await fetch('/api/projects');
                projects = await response.json();
                renderProjects();
                updateStats();
            } catch (error) {
                console.error('Error loading projects:', error);
            }
        }

        function renderProjects() {
            const grid = document.getElementById('projectGrid');
            grid.innerHTML = projects.map(project => {
                const items = project.items.slice(0, 5); // Show first 5 items
                const hasMore = project.items.length > 5;
                
                return `
                    <div class="project-card ${project.selected ? 'selected' : ''}">
                        <div class="project-header">
                            <div class="project-info">
                                <h3>
                                    <input type="checkbox" class="checkbox" ${project.selected ? 'checked' : ''} 
                                           onchange="toggleProject(${project.id})">
                                    ${project.name}
                                </h3>
                                <div class="project-type">${project.type}</div>
                                <div class="project-path">${project.path}</div>
                                <div class="project-stats">
                                    <div class="project-stat">
                                        <div class="project-stat-value">${project.itemCount}</div>
                                        <div>Items</div>
                                    </div>
                                    <div class="project-stat">
                                        <div class="project-stat-value">${formatBytes(project.cacheSize)}</div>
                                        <div>Size</div>
                                    </div>
                                </div>
                            </div>
                        </div>
                        ${project.items.length > 0 ? `
                            <div class="project-items">
                                <h4>Cache Items:</h4>
                                ${items.map(item => `
                                    <div class="cache-item">
                                        <span class="item-icon">${item.type === 'directory' ? '📁' : '📄'}</span>
                                        <span style="flex-grow: 1">${item.path.split('/').pop()}</span>
                                        <span class="item-size">${formatBytes(item.size)}</span>
                                    </div>
                                `).join('')}
                                ${hasMore ? `<div class="cache-item"><em>... and ${project.items.length - 5} more items</em></div>` : ''}
                            </div>
                        ` : ''}
                    </div>
                `;
            }).join('');
        }

        function toggleProject(id) {
            const project = projects.find(p => p.id === id);
            if (project) {
                project.selected = !project.selected;
                renderProjects();
                updateStats();
            }
        }

        function selectAll() {
            projects.forEach(p => p.selected = p.itemCount > 0);
            renderProjects();
            updateStats();
        }

        function deselectAll() {
            projects.forEach(p => p.selected = false);
            renderProjects();
            updateStats();
        }

        function updateStats() {
            const totalProjects = projects.length;
            const totalSize = projects.reduce((sum, p) => sum + p.cacheSize, 0);
            const selectedProjects = projects.filter(p => p.selected).length;
            const selectedSize = projects.filter(p => p.selected).reduce((sum, p) => sum + p.cacheSize, 0);

            document.getElementById('totalProjects').textContent = totalProjects;
            document.getElementById('totalSize').textContent = formatBytes(totalSize);
            document.getElementById('selectedProjects').textContent = selectedProjects;
            document.getElementById('selectedSize').textContent = formatBytes(selectedSize);
            
            document.getElementById('cleanBtn').disabled = selectedProjects === 0;
        }

        function formatBytes(bytes) {
            if (bytes === 0) return '0 B';
            const k = 1024;
            const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
            const i = Math.floor(Math.log(bytes) / Math.log(k));
            return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i];
        }

        async function showCleanDialog() {
            const selectedCount = projects.filter(p => p.selected).length;
            const selectedSize = projects.filter(p => p.selected).reduce((sum, p) => sum + p.cacheSize, 0);
            
            if (confirm(`Clean ${selectedCount} projects?\nThis will remove cache files totaling ${formatBytes(selectedSize)}.`)) {
                await cleanSelected();
            }
        }

        async function cleanSelected() {
            const modal = document.getElementById('cleaningModal');
            const progressFill = document.getElementById('progressFill');
            const status = document.getElementById('cleaningStatus');
            const results = document.getElementById('cleaningResults');
            const doneBtn = document.getElementById('doneBtn');
            
            modal.classList.remove('hidden');
            
            const selectedProjects = projects.filter(p => p.selected).map(p => p.id);
            
            try {
                const response = await fetch('/api/clean', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ projectIds: selectedProjects })
                });
                
                const result = await response.json();
                
                progressFill.style.width = '100%';
                status.textContent = 'Cleanup completed!';
                results.innerHTML = `
                    <strong>✅ Cleanup Results:</strong><br>
                    Projects cleaned: ${result.projectsCleaned}<br>
                    Items removed: ${result.itemsRemoved}<br>
                    Space reclaimed: ${formatBytes(result.sizeReclaimed)}
                `;
                results.classList.remove('hidden');
                doneBtn.classList.remove('hidden');
                
                // Reload projects
                await loadProjects();
                
            } catch (error) {
                status.textContent = 'Error during cleanup: ' + error.message;
                doneBtn.classList.remove('hidden');
            }
        }

        function closeCleaningModal() {
            document.getElementById('cleaningModal').classList.add('hidden');
            document.getElementById('progressFill').style.width = '0%';
            document.getElementById('cleaningStatus').textContent = 'Preparing to clean...';
            document.getElementById('cleaningResults').classList.add('hidden');
            document.getElementById('doneBtn').classList.add('hidden');
        }

        async function refreshProjects() {
            await loadProjects();
        }

        // Load projects on page load
        loadProjects();
    </script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html")
	tmpl = strings.ReplaceAll(tmpl, "\n    ", "\n")
	fmt.Fprint(w, tmpl)
}

func (ws *WebServer) handleAPIProjects(w http.ResponseWriter, r *http.Request) {
	ws.mutex.RLock()
	defer ws.mutex.RUnlock()
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ws.projects)
}

func (ws *WebServer) handleAPIClean(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	var request struct {
		ProjectIds []int `json:"projectIds"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	
	ws.mutex.Lock()
	defer ws.mutex.Unlock()
	
	results := struct {
		ProjectsCleaned int   `json:"projectsCleaned"`
		ItemsRemoved    int   `json:"itemsRemoved"`
		SizeReclaimed   int64 `json:"sizeReclaimed"`
	}{}
	
	for _, projectID := range request.ProjectIds {
		for i, project := range ws.projects {
			if project.ID == projectID && project.Selected {
				// Convert WebCacheItem to CacheItem for cleaning
				cacheItems := make([]CacheItem, len(project.Items))
				for j, item := range project.Items {
					cacheItems[j] = CacheItem{
						Path: item.Path,
						Size: item.Size,
						Type: item.Type,
					}
				}
				
				removedItems, removedSize := removeCacheItems(cacheItems, false)
				results.ProjectsCleaned++
				results.ItemsRemoved += removedItems
				results.SizeReclaimed += removedSize
				
				// Update project state
				ws.projects[i].Selected = false
				ws.projects[i].ItemCount = 0
				ws.projects[i].CacheSize = 0
				ws.projects[i].Items = []WebCacheItem{}
				break
			}
		}
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func runWebUI(rootDir string, port int) error {
	ws := &WebServer{}
	ws.loadProjects(rootDir)
	
	http.HandleFunc("/", ws.handleIndex)
	http.HandleFunc("/api/projects", ws.handleAPIProjects)
	http.HandleFunc("/api/clean", ws.handleAPIClean)
	
	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("🌐 Web UI available at: http://localhost%s\n", addr)
	fmt.Printf("📱 Mobile-friendly interface with visual controls\n")
	fmt.Println("Press Ctrl+C to stop the server")
	
	return http.ListenAndServe(addr, nil)
}