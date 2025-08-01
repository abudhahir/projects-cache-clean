# TUI Demo - Interactive terminal interface (simulated)
Output ../screenshots/ui-demo.gif
Set FontSize 14
Set Width 1200
Set Height 800
Set Theme "Dracula"
Set PlaybackSpeed 1.0
Set TypingSpeed 50ms

# Clear screen and show prompt
Hide
Type "clear"
Enter
Show

# Title
Type "# Cache Remover - Interactive TUI Demo"
Enter
Sleep 1s

# Explain TUI mode
Type "# Beautiful terminal interface with project selection"
Enter
Sleep 2s

# Note: Since we can't actually run the TUI in VHS, we'll simulate the experience
Type "./cache-remover -ui demo-projects"
Enter
Sleep 1s

# Show simulated TUI interface
Type "echo '🚀 Launching Interactive TUI Cache Remover...'"
Enter
Sleep 1s

# Simulate TUI layout (using echo commands to mimic the interface)
Type "echo '┌─ 🧹 Cache Remover - Project Selection ─────────────────────────┐'"
Enter
Type "echo '│                                                                 │'"
Enter  
Type "echo '│  📁 react-app (Node.js)                    [✓] 245.2 MB       │'"
Enter
Type "echo '│     └─ 2 cache items: node_modules, build                      │'"
Enter
Type "echo '│                                                                 │'"
Enter
Type "echo '│  📁 python-api (Python)                    [ ] 12.8 MB        │'"
Enter
Type "echo '│     └─ 3 cache items: __pycache__, .pytest_cache, build       │'"
Enter
Type "echo '│                                                                 │'"
Enter
Type "echo '│  📁 java-service (Java/Maven)              [✓] 156.4 MB       │'"
Enter
Type "echo '│     └─ 1 cache items: target                                   │'"
Enter
Type "echo '│                                                                 │'"
Enter
Type "echo '│  📊 Statistics: 3 projects, 2 selected, 401.6 MB reclaimable  │'"
Enter
Type "echo '│                                                                 │'"
Enter
Type "echo '│  🔧 Controls: ↑↓=navigate, space=select, c=clean, q=quit       │'"
Enter
Type "echo '└─────────────────────────────────────────────────────────────────┘'"
Enter
Sleep 3s

# Show controls
Type "# 🎮 Use arrow keys to navigate, space to select, 'c' to clean"
Enter
Sleep 2s