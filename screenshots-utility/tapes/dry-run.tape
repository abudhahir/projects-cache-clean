# Dry Run Demo - Safe preview mode
Output ../screenshots/dry-run.gif
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
Type "# Cache Remover - Dry Run Demo (Safe Preview)"
Enter
Sleep 1s

# Explain dry run
Type "# Shows what WOULD be removed without actually removing"
Enter
Sleep 2s

# Run dry run
Type "./cache-remover -dry-run demo-projects"
Enter
Sleep 5s

# Highlight safety
Type "# 🛡️ No files were actually removed - completely safe!"
Enter
Sleep 2s