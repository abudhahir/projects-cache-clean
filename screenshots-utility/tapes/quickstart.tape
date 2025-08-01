# Quickstart Demo - Combined essential workflow
Output ../screenshots/quickstart.gif
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
Type "# 🧹 Cache Remover - 2-Minute Quickstart"
Enter
Sleep 1s

# Step 1: Build
Type "# Step 1: Build the tool"
Type "go build -o cache-remover"
Enter
Sleep 1s

# Step 2: Safe preview
Type "# Step 2: Safe preview (always start here!)"
Enter
Type "./cache-remover -dry-run demo-projects"
Enter
Sleep 4s

# Step 3: Actual cleanup
Type "# Step 3: Clean with confirmation"
Enter
Type "./cache-remover demo-projects"
Enter
Sleep 2s
Type "y"
Enter
Sleep 3s

# Step 4: Interactive mode
Type "# Step 4: Try interactive TUI mode"
Enter
Type "./cache-remover -ui demo-projects"
Enter
Sleep 1s

# Show completion
Type "# 🎉 That's it! You're now ready to reclaim disk space!"
Enter
Type "# 💡 Tip: Always use -dry-run first to preview changes"
Enter
Sleep 3s