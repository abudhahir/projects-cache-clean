@echo off
REM Cache Remover Utility - Windows Uninstaller
REM Removes the utility and cleans up PATH

setlocal EnableDelayedExpansion

echo ==========================================
echo  Cache Remover Utility - Uninstaller
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

REM Define installation directories
set SYSTEM_DIR=C:\Program Files\CacheRemover
set USER_DIR=%LOCALAPPDATA%\CacheRemover

REM Check what's installed
set FOUND_SYSTEM=0
set FOUND_USER=0

if exist "%SYSTEM_DIR%\cache-remover-utility.exe" (
    set FOUND_SYSTEM=1
    echo Found system installation: %SYSTEM_DIR%
)

if exist "%USER_DIR%\cache-remover-utility.exe" (
    set FOUND_USER=1
    echo Found user installation: %USER_DIR%
)

if %FOUND_SYSTEM%==0 if %FOUND_USER%==0 (
    echo No Cache Remover Utility installations found.
    echo.
    pause
    exit /b 0
)

echo.
echo The following will be removed:
if %FOUND_SYSTEM%==1 echo - System installation: %SYSTEM_DIR%
if %FOUND_USER%==1 echo - User installation: %USER_DIR%
echo - Start Menu shortcuts
echo - Desktop shortcuts
echo - PATH entries
echo.

choice /C YN /M "Continue with uninstallation"
if %ERRORLEVEL%==2 (
    echo Uninstallation cancelled.
    pause
    exit /b 0
)

echo.
echo Starting uninstallation...

REM Remove system installation
if %FOUND_SYSTEM%==1 (
    echo.
    echo Removing system installation...
    if exist "%SYSTEM_DIR%\cache-remover-utility.exe" (
        del /F /Q "%SYSTEM_DIR%\cache-remover-utility.exe" >nul 2>&1
        if %ERRORLEVEL% EQU 0 (
            echo Removed executable from system directory
        ) else (
            echo Warning: Could not remove executable from system directory
        )
    )
    
    if exist "%SYSTEM_DIR%" (
        rmdir "%SYSTEM_DIR%" >nul 2>&1
        if %ERRORLEVEL% EQU 0 (
            echo Removed system directory
        )
    )
    
    REM Remove from system PATH (requires admin)
    if %ADMIN%==1 (
        echo Removing from system PATH...
        for /f "skip=2 tokens=3*" %%a in ('reg query "HKLM\SYSTEM\CurrentControlSet\Control\Session Manager\Environment" /v PATH 2^>nul') do (
            set SysPath=%%b
            set NewSysPath=!SysPath:%SYSTEM_DIR%;=!
            set NewSysPath=!NewSysPath:;%SYSTEM_DIR%=!
            set NewSysPath=!NewSysPath:%SYSTEM_DIR%=!
            
            if not "!SysPath!"=="!NewSysPath!" (
                setx /M PATH "!NewSysPath!" >nul 2>&1
                echo Removed from system PATH
            )
        )
    ) else (
        echo Warning: Cannot remove from system PATH (requires admin privileges)
    )
)

REM Remove user installation
if %FOUND_USER%==1 (
    echo.
    echo Removing user installation...
    if exist "%USER_DIR%\cache-remover-utility.exe" (
        del /F /Q "%USER_DIR%\cache-remover-utility.exe" >nul 2>&1
        echo Removed executable from user directory
    )
    
    if exist "%USER_DIR%" (
        rmdir "%USER_DIR%" >nul 2>&1
        if %ERRORLEVEL% EQU 0 (
            echo Removed user directory
        )
    )
    
    REM Remove from user PATH
    echo Removing from user PATH...
    for /f "skip=2 tokens=3*" %%a in ('reg query HKCU\Environment /v PATH 2^>nul') do (
        set UserPath=%%b
        set NewUserPath=!UserPath:%USER_DIR%;=!
        set NewUserPath=!NewUserPath:;%USER_DIR%=!
        set NewUserPath=!NewUserPath:%USER_DIR%=!
        
        if not "!UserPath!"=="!NewUserPath!" (
            setx PATH "!NewUserPath!" >nul 2>&1
            echo Removed from user PATH
        )
    )
)

REM Remove shortcuts
echo.
echo Removing shortcuts...

REM Start Menu shortcut
if exist "%APPDATA%\Microsoft\Windows\Start Menu\Programs\Cache Remover Utility.lnk" (
    del /F /Q "%APPDATA%\Microsoft\Windows\Start Menu\Programs\Cache Remover Utility.lnk" >nul 2>&1
    echo Removed Start Menu shortcut
)

REM Desktop shortcut
if exist "%USERPROFILE%\Desktop\Cache Remover.lnk" (
    del /F /Q "%USERPROFILE%\Desktop\Cache Remover.lnk" >nul 2>&1
    echo Removed desktop shortcut
)

REM Check for any remaining traces
echo.
echo Checking for remaining files...
set REMAINING=0

if exist "%SYSTEM_DIR%" (
    echo Warning: System directory still exists: %SYSTEM_DIR%
    set REMAINING=1
)

if exist "%USER_DIR%" (
    echo Warning: User directory still exists: %USER_DIR%
    set REMAINING=1
)

echo.
echo ==========================================
if %REMAINING%==0 (
    echo  Uninstallation Complete!
) else (
    echo  Uninstallation Complete with Warnings
)
echo ==========================================
echo.

if %REMAINING%==0 (
    echo Cache Remover Utility has been completely removed.
) else (
    echo Cache Remover Utility has been removed, but some files may remain.
    echo You may need to manually delete remaining directories.
)

echo.
echo IMPORTANT: Restart your Command Prompt or PowerShell
echo for PATH changes to take effect!
echo.
pause