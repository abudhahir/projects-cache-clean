# Basic Usage Demo - Standard cleanup with confirmation
Output ../screenshots/basic-usage.gif
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
Type "# Cache Remover - Basic Usage Demo"
Enter
Sleep 1s

# Show directory structure first
Type "ls -la demo-projects/"
Enter
Sleep 2s

# Run basic cleanup
Type "./cache-remover demo-projects"
Enter
Sleep 3s

# Confirm cleanup
Type "y"
Enter
Sleep 4s

# Show completion
Type "# ✅ Cleanup completed successfully!"
Enter
Sleep 2s