@echo off
REM Cache Remover Utility - Windows Installation Script
REM Installs the utility with proper PATH configuration

setlocal EnableDelayedExpansion

echo ==========================================
echo  Cache Remover Utility - Windows Installer
echo ==========================================
echo.

REM Check for admin rights
net session >nul 2>&1
if %ERRORLEVEL% == 0 (
    set ADMIN=1
    echo Running with Administrator privileges
) else (
    set ADMIN=0
    echo Running without Administrator privileges
)
echo.

REM Build first
echo Step 1: Building application...
call build.bat
if %ERRORLEVEL% NEQ 0 (
    echo Build failed! Aborting installation.
    pause
    exit /b 1
)
echo.

REM Choose installation directory based on admin rights
if %ADMIN%==1 (
    set INSTALL_DIR=C:\Program Files\CacheRemover
    echo Step 2: Installing system-wide to %INSTALL_DIR%
) else (
    set INSTALL_DIR=%LOCALAPPDATA%\CacheRemover
    echo Step 2: Installing for current user to %INSTALL_DIR%
)

REM Create installation directory
if not exist "%INSTALL_DIR%" (
    echo Creating installation directory...
    mkdir "%INSTALL_DIR%"
)

REM Copy executable
echo Copying executable...
copy /Y cache-remover-utility.exe "%INSTALL_DIR%\" >nul
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: Failed to copy executable!
    pause
    exit /b 1
)

REM Add to PATH
echo.
echo Step 3: Adding to PATH...

if %ADMIN%==1 (
    REM System PATH
    echo Adding to system PATH...
    for /f "skip=2 tokens=3*" %%a in ('reg query "HKLM\SYSTEM\CurrentControlSet\Control\Session Manager\Environment" /v PATH') do set SysPath=%%b
    
    REM Check if already in PATH
    echo !SysPath! | find /i "%INSTALL_DIR%" >nul
    if %ERRORLEVEL% NEQ 0 (
        setx /M PATH "!SysPath!;%INSTALL_DIR%" >nul 2>&1
        echo Successfully added to system PATH
    ) else (
        echo Already in system PATH
    )
) else (
    REM User PATH
    echo Adding to user PATH...
    for /f "skip=2 tokens=3*" %%a in ('reg query HKCU\Environment /v PATH 2^>nul') do set UserPath=%%b
    
    REM If no user PATH exists, create it
    if "!UserPath!"=="" set UserPath=%PATH%
    
    REM Check if already in PATH
    echo !UserPath! | find /i "%INSTALL_DIR%" >nul
    if %ERRORLEVEL% NEQ 0 (
        setx PATH "!UserPath!;%INSTALL_DIR%" >nul 2>&1
        echo Successfully added to user PATH
    ) else (
        echo Already in user PATH
    )
)

REM Create Start Menu shortcut
echo.
echo Step 4: Creating Start Menu shortcut...
powershell -Command "$WshShell = New-Object -comObject WScript.Shell; $Shortcut = $WshShell.CreateShortcut([Environment]::GetFolderPath('Programs') + '\Cache Remover Utility.lnk'); $Shortcut.TargetPath = '%INSTALL_DIR%\cache-remover-utility.exe'; $Shortcut.WorkingDirectory = '%USERPROFILE%'; $Shortcut.Description = 'Clean cache directories from your projects'; $Shortcut.Save()" >nul 2>&1

if %ERRORLEVEL% EQU 0 (
    echo Start Menu shortcut created
) else (
    echo Warning: Could not create Start Menu shortcut
)

REM Create desktop shortcut (optional)
echo.
choice /C YN /M "Create desktop shortcut"
if %ERRORLEVEL%==1 (
    powershell -Command "$WshShell = New-Object -comObject WScript.Shell; $Shortcut = $WshShell.CreateShortcut([Environment]::GetFolderPath('Desktop') + '\Cache Remover.lnk'); $Shortcut.TargetPath = '%INSTALL_DIR%\cache-remover-utility.exe'; $Shortcut.Arguments = '--ui'; $Shortcut.WorkingDirectory = '%USERPROFILE%'; $Shortcut.Description = 'Clean cache directories from your projects'; $Shortcut.Save()" >nul 2>&1
    echo Desktop shortcut created
)

REM Test installation
echo.
echo Step 5: Testing installation...
"%INSTALL_DIR%\cache-remover-utility.exe" --version >nul 2>&1
if %ERRORLEVEL% EQU 0 (
    echo Installation test passed!
) else (
    echo Warning: Installation test failed
)

REM Display summary
echo.
echo ==========================================
echo  Installation Complete!
echo ==========================================
echo.
echo Installed to: %INSTALL_DIR%
echo.
echo You can now use Cache Remover Utility by:
echo.
echo 1. From Start Menu: Cache Remover Utility
echo 2. From Command Prompt (after restart):
echo    cache-remover-utility --ui
echo    cache-remover-utility --dry-run C:\Projects
echo    cache-remover-utility --help
echo.
echo IMPORTANT: You must restart your Command Prompt
echo or PowerShell for PATH changes to take effect!
echo.
pause