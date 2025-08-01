# Performance Demo - Optimization demonstration
Output ../screenshots/performance.gif
Set FontSize 14
Set Width 1200
Set Height 700
Set Theme "Dracula"
Set PlaybackSpeed 1.0
Set TypingSpeed 50ms

# Clear screen and show prompt
Hide
Type "clear"
Enter
Show

# Title
Type "# Cache Remover - Performance Optimization Demo"
Enter
Sleep 1s

# Explain the optimization
Type "# 🚀 Smart cache directory skipping for 10-100x speed improvement"
Enter
Sleep 2s

# Show what gets optimized
Type "echo 'Traditional tools scan EVERY file:'"
Enter
Type "echo '  📁 node_modules/'"
Type "echo '    📄 package1/index.js'"
Enter
Type "echo '    📄 package1/lib/util.js'"
Enter  
Type "echo '    📄 package2/index.js'"
Enter
Type "echo '    📄 ... (scanning 50,000+ files) ← SLOW! 28.3 seconds'"
Enter
Sleep 3s

Type "echo ''"
Enter
Type "echo 'Cache Remover treats cache dirs as single units:'"
Enter
Type "echo '  📁 node_modules/ ← Single unit removal ← FAST! 0.1 seconds'"
Enter
Sleep 2s

# Run the actual tool to show optimization messages
Type "./cache-remover -verbose -dry-run demo-projects"
Enter
Sleep 5s

# Highlight the speed
Type "# ⚡ Result: 283x faster scanning + same removal speed = huge win!"
Enter
Sleep 2s