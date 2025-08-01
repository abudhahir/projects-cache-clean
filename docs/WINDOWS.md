# 🪟 Windows Installation & Usage Guide

Complete guide for installing and using Cache Remover Utility on Windows.

## 🚀 Quick Installation (Recommended)

### Method 1: Automated Installation (Easiest)
```cmd
# Download or clone the repository
git clone https://github.com/abudhahir/projects-cache-clean.git
cd projects-cache-clean

# Run the automated installer
install.bat
```

**What this does:**
- ✅ Builds the application
- ✅ Installs to appropriate directory (Program Files or Local)
- ✅ Adds to Windows PATH
- ✅ Creates Start Menu shortcut
- ✅ Optional desktop shortcut
- ✅ Tests the installation

### Method 2: Manual Build + Copy
```cmd
# Build the application
build.bat

# Copy to a directory in your PATH
copy cache-remover-utility.exe C:\Windows\System32\
```

## 📋 Prerequisites

- **Windows 10/11** (Windows 7+ should work)
- **Go 1.19+** - Download from [golang.org](https://golang.org/dl/)
- **PowerShell 5.0+** (for advanced scripts)
- **Admin rights** (for system-wide installation)

### Installing Go on Windows
```cmd
# Download from https://golang.org/dl/
# Or use Chocolatey:
choco install golang

# Or use Scoop:
scoop install go

# Verify installation:
go version
```

## 🛠️ Build Options

### Using Windows Batch Scripts
```cmd
# Simple build
build.bat

# Build and install
install.bat

# Uninstall
uninstall.bat
```

### Using Windows Makefile
```cmd
# Build application
make -f Makefile.windows build

# Install system-wide (requires admin)
make -f Makefile.windows install

# Install for current user only
make -f Makefile.windows install-user

# Build portable version
make -f Makefile.windows portable

# Clean build artifacts
make -f Makefile.windows clean

# Show all options
make -f Makefile.windows help
```

### Manual Go Commands
```cmd
# Basic build
go build -o cache-remover-utility.exe

# Optimized build (smaller size)
go build -ldflags="-s -w" -o cache-remover-utility.exe

# Portable build (no external dependencies)
set CGO_ENABLED=0
go build -ldflags="-s -w" -a -installsuffix cgo -o cache-remover-utility.exe
```

## 📦 Installation Methods

### System-Wide Installation (Requires Admin)
- **Location**: `C:\Program Files\CacheRemover\`
- **Available to**: All users
- **PATH**: Added to system PATH
- **Command**: `make -f Makefile.windows install` or `install.bat`

### User Installation (No Admin Required)
- **Location**: `%LOCALAPPDATA%\CacheRemover\`
- **Available to**: Current user only
- **PATH**: Added to user PATH
- **Command**: `make -f Makefile.windows install-user`

### Portable Installation
- **Location**: Any directory
- **PATH**: Manual or relative paths
- **Usage**: `.\cache-remover-utility.exe`

## 🎯 Usage Examples

### Command Line Usage
```cmd
# Show help
cache-remover-utility.exe --help

# Scan current directory (dry-run)
cache-remover-utility.exe --dry-run .

# Scan specific directory
cache-remover-utility.exe --dry-run C:\Projects

# Interactive TUI mode
cache-remover-utility.exe --ui

# Scan with TUI
cache-remover-utility.exe --ui C:\Projects

# Clean with confirmation
cache-remover-utility.exe C:\Projects

# List supported project types
cache-remover-utility.exe --list-types

# Save default configuration
cache-remover-utility.exe --save-config
```

### PowerShell Usage
```powershell
# Scan all projects
cache-remover-utility.exe --dry-run $env:USERPROFILE\Projects

# Clean specific project types
cache-remover-utility.exe --verbose C:\Code\JavaScript

# Use with variables
$projectPath = "C:\Development\MyProjects"
cache-remover-utility.exe --ui $projectPath
```

### Batch Script Usage
```cmd
@echo off
echo Cleaning all development projects...
cache-remover-utility.exe --dry-run C:\Dev
pause

echo.
echo Continue with cleanup? Press Ctrl+C to cancel.
pause

cache-remover-utility.exe C:\Dev
echo Cleanup complete!
pause
```

## 🧪 Testing & Demo

### Create Demo Projects
```cmd
# Using PowerShell
powershell -ExecutionPolicy Bypass -File scripts\create-demo-windows.ps1

# Using Makefile
make -f Makefile.windows demo
```

### Test Installation
```cmd
# Verify binary works
cache-remover-utility.exe --version

# Test with demo projects
cache-remover-utility.exe --dry-run demo-projects

# Test TUI mode
cache-remover-utility.exe --ui demo-projects
```

## 🔧 Configuration

### Default Configuration Location
- **System**: `C:\Program Files\CacheRemover\`
- **User**: `%LOCALAPPDATA%\CacheRemover\`
- **Current Directory**: `.\config.json` or `.\cache-remover-config.json`

### Create Custom Configuration
```cmd
# Generate default config file
cache-remover-utility.exe --save-config

# Edit the generated cache-remover-config.json
notepad cache-remover-config.json
```

### Windows-Specific Paths
```json
{
  "project_types": [
    {
      "name": "Visual Studio",
      "indicators": ["*.sln", "*.csproj"],
      "cache_config": {
        "directories": ["bin", "obj", "packages", ".vs"]
      }
    },
    {
      "name": "Unity",
      "indicators": ["Assets", "ProjectSettings"],
      "cache_config": {
        "directories": ["Library", "Temp", "obj"]  
      }
    }
  ]
}
```

## 🛡️ Windows Security & Permissions

### Windows Defender
If Windows Defender flags the executable:
```cmd
# Add exclusion for the installation directory
# Windows Security → Virus & threat protection → Exclusions
# Add: C:\Program Files\CacheRemover\
```

### UAC and Admin Rights
- **System installation**: Requires admin rights (one-time)
- **User installation**: No admin rights needed
- **Portable usage**: No special permissions needed
- **PATH modification**: Admin rights for system PATH, user rights for user PATH

### Execution Policy (PowerShell)
```powershell
# If PowerShell scripts are blocked:
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser

# Or run with bypass:
powershell -ExecutionPolicy Bypass -File script.ps1
```

## 🔍 Troubleshooting

### Common Issues

#### "Go not found"
```cmd
# Install Go from https://golang.org/dl/
# Or add Go to PATH manually:
set PATH=%PATH%;C:\Go\bin
```

#### "Access denied" during installation
```cmd
# Run Command Prompt as Administrator:
# Right-click Command Prompt → "Run as administrator"
# Then run: install.bat
```

#### "Command not found" after installation
```cmd
# Restart Command Prompt/PowerShell
# Or refresh PATH:
set PATH=%PATH%;C:\Program Files\CacheRemover

# Check installation:
where cache-remover-utility
```

#### PowerShell execution policy
```powershell
# Enable script execution:
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser

# Or run with bypass:
powershell -ExecutionPolicy Bypass -File install.bat
```

### Build Issues

#### "Module not found"
```cmd
# Download dependencies:
go mod download
go mod tidy

# Verify modules:
go mod verify
```

#### Antivirus blocking build
```cmd
# Add project directory to antivirus exclusions
# Then rebuild:
go clean -cache
build.bat
```

## 🗑️ Uninstallation

### Automated Uninstall
```cmd
# Run the uninstaller
uninstall.bat

# Or using Makefile
make -f Makefile.windows uninstall
```

### Manual Uninstall
```cmd
# Remove executable
del "C:\Program Files\CacheRemover\cache-remover-utility.exe"
rmdir "C:\Program Files\CacheRemover"

# Remove from PATH (manual)
# System Properties → Advanced → Environment Variables
# Edit PATH and remove: C:\Program Files\CacheRemover

# Remove shortcuts
del "%APPDATA%\Microsoft\Windows\Start Menu\Programs\Cache Remover Utility.lnk"
del "%USERPROFILE%\Desktop\Cache Remover.lnk"
```

## 📊 Windows-Specific Features

### File System Support
- ✅ **NTFS**: Full support
- ✅ **FAT32**: Basic support (file size limitations)
- ✅ **exFAT**: Full support
- ✅ **Network drives**: Supported with proper permissions

### Path Handling
- ✅ **Long paths**: Supported (>260 characters)
- ✅ **Unicode**: Full Unicode path support
- ✅ **Drive letters**: C:\, D:\, etc.
- ✅ **UNC paths**: \\server\share\path
- ✅ **Spaces in paths**: Properly handled

### Performance Considerations
- **SSD drives**: Optimal performance
- **HDD drives**: Good performance
- **Network drives**: Slower, depends on network speed
- **Memory usage**: ~10-50MB during scanning

## 💡 Tips & Best Practices

### Performance Tips
```cmd
# Use SSD for best performance
# Exclude from antivirus real-time scanning during build
# Use --workers flag to adjust parallelism:
cache-remover-utility.exe --workers 8 --dry-run C:\Projects
```

### Integration with IDEs
```cmd
# Add as External Tool in Visual Studio:
# Tools → External Tools → Add
# Command: cache-remover-utility.exe
# Arguments: --ui $(SolutionDir)
```

### Batch Processing
```cmd
# Clean multiple project directories:
for /d %%i in (C:\Projects\*) do (
    echo Cleaning %%i...
    cache-remover-utility.exe --dry-run "%%i"
)
```

### Scheduled Cleanup
```cmd
# Create scheduled task for weekly cleanup:
schtasks /create /tn "Weekly Cache Cleanup" /tr "cache-remover-utility.exe --ui C:\Projects" /sc weekly
```

## 🆘 Getting Help

### Documentation
- **This file**: Complete Windows guide
- **README.md**: General usage and features
- **TESTING.md**: Testing procedures
- **Built-in help**: `cache-remover-utility.exe --help`

### Support Channels
- **GitHub Issues**: [Report bugs/feature requests](https://github.com/abudhahir/projects-cache-clean/issues)
- **Built-in commands**: `--list-types`, `--version`, `--help`

### Diagnostic Information
```cmd
# Get system info for bug reports:
cache-remover-utility.exe --version
go version
echo %PATH%
echo %GOPATH%
echo %GOOS% %GOARCH%
```

## 🎉 Success! 

After installation, you should be able to:

1. ✅ Run `cache-remover-utility` from any Command Prompt
2. ✅ Use the Start Menu shortcut
3. ✅ Access the TUI with `cache-remover-utility --ui`
4. ✅ Clean cache from any Windows project directory

**Happy cache cleaning! 🗑️✨**