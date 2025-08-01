# Interactive Mode Demo - Per-project confirmation
Output ../screenshots/interactive.gif
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
Type "# Cache Remover - Interactive Mode Demo"
Enter
Sleep 1s

# Explain interactive mode
Type "# Asks for confirmation before cleaning each project"
Enter
Sleep 2s

# Run interactive mode
Type "./cache-remover -interactive demo-projects"
Enter
Sleep 2s

# First project - yes
Type "y"
Enter
Sleep 2s

# Second project - no
Type "n"
Enter
Sleep 2s

# Third project - yes
Type "y"
Enter
Sleep 3s

# Show selective cleanup
Type "# ✅ Cleaned only selected projects - full control!"
Enter
Sleep 2s