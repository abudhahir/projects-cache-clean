@echo off
REM Cache Remover Utility - Windows Build Script
REM Simple build script for Windows users

echo ========================================
echo  Cache Remover Utility - Build Script
echo ========================================
echo.

REM Check if Go is installed
where go >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: Go is not installed or not in PATH!
    echo Please install Go from https://golang.org/dl/
    echo.
    pause
    exit /b 1
)

REM Display Go version
echo Detected Go version:
go version
echo.

REM Clean previous builds
if exist cache-remover-utility.exe (
    echo Cleaning previous build...
    del /F /Q cache-remover-utility.exe
)

REM Build the application
echo Building Cache Remover Utility...
go build -ldflags="-s -w" -o cache-remover-utility.exe

if %ERRORLEVEL% NEQ 0 (
    echo.
    echo ERROR: Build failed!
    pause
    exit /b 1
)

REM Verify build
if exist cache-remover-utility.exe (
    echo.
    echo Build successful!
    echo.
    echo Binary created: cache-remover-utility.exe
    for %%I in (cache-remover-utility.exe) do echo Size: %%~zI bytes
    echo.
    echo You can now run:
    echo   cache-remover-utility.exe --help
    echo   cache-remover-utility.exe --ui
    echo   cache-remover-utility.exe --dry-run C:\Projects
) else (
    echo.
    echo ERROR: Build completed but executable not found!
    pause
    exit /b 1
)

echo.
echo Press any key to exit...
pause >nul