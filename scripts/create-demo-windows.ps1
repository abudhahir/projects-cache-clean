# Cache Remover Utility - Windows Demo Project Creator
# Creates demo projects with various cache directories for testing

param(
    [string]$DemoDir = "demo-projects"
)

Write-Host "========================================" -ForegroundColor Cyan
Write-Host " Cache Remover - Demo Project Creator" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host

# Clean up previous demo if it exists
if (Test-Path $DemoDir) {
    Write-Host "Cleaning up previous demo projects..." -ForegroundColor Yellow
    Remove-Item -Recurse -Force $DemoDir
}

# Create demo directory
Write-Host "Creating demo projects in: $DemoDir" -ForegroundColor Green
New-Item -ItemType Directory -Force -Path $DemoDir | Out-Null

# Node.js Projects
Write-Host "Creating Node.js projects..." -ForegroundColor Blue

# Project 1: React App
$reactDir = "$DemoDir\react-todo-app"
New-Item -ItemType Directory -Force -Path $reactDir | Out-Null
@"
{
  "name": "react-todo-app",
  "version": "1.0.0",
  "dependencies": {
    "react": "^18.0.0",
    "react-dom": "^18.0.0"
  },
  "scripts": {
    "start": "react-scripts start",
    "build": "react-scripts build"
  }
}
"@ | Out-File -FilePath "$reactDir\package.json" -Encoding UTF8

# Create node_modules with dummy content
$nodeModulesDir = "$reactDir\node_modules"
New-Item -ItemType Directory -Force -Path $nodeModulesDir | Out-Null
New-Item -ItemType Directory -Force -Path "$nodeModulesDir\react" | Out-Null
New-Item -ItemType Directory -Force -Path "$nodeModulesDir\.cache" | Out-Null

# Create dummy files (1-50MB total)
$dummyContent = "// Dummy module file`n" * 1000
$dummyContent | Out-File -FilePath "$nodeModulesDir\react\index.js" -Encoding UTF8
$dummyContent | Out-File -FilePath "$nodeModulesDir\.cache\babel-loader.json" -Encoding UTF8

# Create build directory
$buildDir = "$reactDir\build"
New-Item -ItemType Directory -Force -Path $buildDir | Out-Null
"Built application files" | Out-File -FilePath "$buildDir\index.html" -Encoding UTF8

# Project 2: Express API
$expressDir = "$DemoDir\express-api"
New-Item -ItemType Directory -Force -Path $expressDir | Out-Null
@"
{
  "name": "express-api",
  "version": "1.0.0",
  "dependencies": {
    "express": "^4.18.0",
    "lodash": "^4.17.21"
  }
}
"@ | Out-File -FilePath "$expressDir\package.json" -Encoding UTF8

# Create large node_modules
$expressNodeModules = "$expressDir\node_modules"
New-Item -ItemType Directory -Force -Path $expressNodeModules | Out-Null
New-Item -ItemType Directory -Force -Path "$expressNodeModules\express" | Out-Null
New-Item -ItemType Directory -Force -Path "$expressNodeModules\lodash" | Out-Null

# Create larger dummy files
$largeContent = "module.exports = " + ("'dummy';" * 2000)
$largeContent | Out-File -FilePath "$expressNodeModules\express\index.js" -Encoding UTF8
$largeContent | Out-File -FilePath "$expressNodeModules\lodash\index.js" -Encoding UTF8

# Python Projects
Write-Host "Creating Python projects..." -ForegroundColor Blue

# Project 1: Data Science Project
$dataDir = "$DemoDir\data-science-project"
New-Item -ItemType Directory -Force -Path $dataDir | Out-Null
@"
pandas==1.5.0
numpy==1.24.0
scikit-learn==1.1.0
matplotlib==3.6.0
jupyter==1.0.0
"@ | Out-File -FilePath "$dataDir\requirements.txt" -Encoding UTF8

# Create virtual environment
$venvDir = "$dataDir\venv"
New-Item -ItemType Directory -Force -Path $venvDir | Out-Null
New-Item -ItemType Directory -Force -Path "$venvDir\Lib" | Out-Null
New-Item -ItemType Directory -Force -Path "$venvDir\Scripts" | Out-Null
New-Item -ItemType Directory -Force -Path "$venvDir\Lib\site-packages" | Out-Null

# Create dummy packages in venv
$packages = @("pandas", "numpy", "sklearn", "matplotlib")
foreach ($pkg in $packages) {
    $pkgDir = "$venvDir\Lib\site-packages\$pkg"
    New-Item -ItemType Directory -Force -Path $pkgDir | Out-Null
    $dummyPython = "# Dummy Python package`n" * 500
    $dummyPython | Out-File -FilePath "$pkgDir\__init__.py" -Encoding UTF8
}

# Create __pycache__ directories
$pycacheDir = "$dataDir\__pycache__"
New-Item -ItemType Directory -Force -Path $pycacheDir | Out-Null
"dummy bytecode" | Out-File -FilePath "$pycacheDir\main.cpython-39.pyc" -Encoding UTF8

# Create build directory
$buildDir = "$dataDir\build"
New-Item -ItemType Directory -Force -Path $buildDir | Out-Null
"Built distribution files" | Out-File -FilePath "$buildDir\dist.tar.gz" -Encoding UTF8

# Project 2: ML Project with Conda
$mlDir = "$DemoDir\ml-project"
New-Item -ItemType Directory -Force -Path $mlDir | Out-Null
@"
name: ml-env
dependencies:
  - python=3.9
  - tensorflow=2.9
  - pytorch=1.12
  - scikit-learn
"@ | Out-File -FilePath "$mlDir\environment.yml" -Encoding UTF8

"torch>=1.12.0" | Out-File -FilePath "$mlDir\requirements.txt" -Encoding UTF8

# Create conda environment
$condaDir = "$mlDir\conda"
New-Item -ItemType Directory -Force -Path $condaDir | Out-Null
New-Item -ItemType Directory -Force -Path "$condaDir\lib" | Out-Null
New-Item -ItemType Directory -Force -Path "$condaDir\lib\python3.9" | Out-Null

# Create large ML packages
$tfDir = "$condaDir\lib\python3.9\tensorflow"
New-Item -ItemType Directory -Force -Path $tfDir | Out-Null
$tfContent = "# TensorFlow dummy content`n" * 3000
$tfContent | Out-File -FilePath "$tfDir\__init__.py" -Encoding UTF8

# Java Projects
Write-Host "Creating Java projects..." -ForegroundColor Blue

# Maven Project
$mavenDir = "$DemoDir\spring-boot-api"
New-Item -ItemType Directory -Force -Path $mavenDir | Out-Null
@"
<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0">
    <modelVersion>4.0.0</modelVersion>
    <groupId>com.example</groupId>
    <artifactId>spring-boot-api</artifactId>
    <version>1.0.0</version>
    <packaging>jar</packaging>
    
    <dependencies>
        <dependency>
            <groupId>org.springframework.boot</groupId>
            <artifactId>spring-boot-starter-web</artifactId>
            <version>2.7.0</version>
        </dependency>
    </dependencies>
</project>
"@ | Out-File -FilePath "$mavenDir\pom.xml" -Encoding UTF8

# Create target directory with compiled classes
$targetDir = "$mavenDir\target"
New-Item -ItemType Directory -Force -Path $targetDir | Out-Null
New-Item -ItemType Directory -Force -Path "$targetDir\classes" | Out-Null
"Compiled Java bytecode" | Out-File -FilePath "$targetDir\classes\Application.class" -Encoding UTF8
"Built JAR file content" | Out-File -FilePath "$targetDir\spring-boot-api-1.0.0.jar" -Encoding UTF8

# Gradle Project
$gradleDir = "$DemoDir\android-app"
New-Item -ItemType Directory -Force -Path $gradleDir | Out-Null
@"
plugins {
    id 'com.android.application'
}

android {
    compileSdk 32
}

dependencies {
    implementation 'androidx.appcompat:appcompat:1.5.0'
    implementation 'com.google.android.material:material:1.6.1'
}
"@ | Out-File -FilePath "$gradleDir\build.gradle" -Encoding UTF8

# Create build and .gradle directories
$buildDir = "$gradleDir\build"
$gradleCacheDir = "$gradleDir\.gradle"
New-Item -ItemType Directory -Force -Path $buildDir | Out-Null
New-Item -ItemType Directory -Force -Path $gradleCacheDir | Out-Null

"Compiled Android app" | Out-File -FilePath "$buildDir\app-release.apk" -Encoding UTF8
"Gradle cache data" | Out-File -FilePath "$gradleCacheDir\buildOutputCleanup.lock" -Encoding UTF8

# Go Project
Write-Host "Creating Go project..." -ForegroundColor Blue

$goDir = "$DemoDir\go-microservice"
New-Item -ItemType Directory -Force -Path $goDir | Out-Null
@"
module go-microservice

go 1.19

require (
    github.com/gin-gonic/gin v1.8.1
    github.com/lib/pq v1.10.6
)
"@ | Out-File -FilePath "$goDir\go.mod" -Encoding UTF8

# Create vendor directory
$vendorDir = "$goDir\vendor"
New-Item -ItemType Directory -Force -Path $vendorDir | Out-Null
New-Item -ItemType Directory -Force -Path "$vendorDir\github.com" | Out-Null
"Vendored dependencies" | Out-File -FilePath "$vendorDir\github.com\gin-gonic.go" -Encoding UTF8

# Flutter Project  
Write-Host "Creating Flutter project..." -ForegroundColor Blue

$flutterDir = "$DemoDir\flutter-mobile-app"
New-Item -ItemType Directory -Force -Path $flutterDir | Out-Null
@"
name: flutter_mobile_app
description: A demo Flutter application

dependencies:
  flutter:
    sdk: flutter
  http: ^0.13.5
  provider: ^6.0.3

dev_dependencies:
  flutter_test:
    sdk: flutter
"@ | Out-File -FilePath "$flutterDir\pubspec.yaml" -Encoding UTF8

# Create build and .dart_tool directories
$buildDir = "$flutterDir\build"
$dartToolDir = "$flutterDir\.dart_tool"
New-Item -ItemType Directory -Force -Path $buildDir | Out-Null
New-Item -ItemType Directory -Force -Path $dartToolDir | Out-Null

"Compiled Flutter app for Android" | Out-File -FilePath "$buildDir\app.apk" -Encoding UTF8
"Dart analysis cache" | Out-File -FilePath "$dartToolDir\package_config.json" -Encoding UTF8

# Calculate and display demo statistics
Write-Host
Write-Host "Demo project creation complete!" -ForegroundColor Green
Write-Host
Write-Host "Created projects:" -ForegroundColor Yellow
Write-Host "  📦 Node.js projects: 2 (react-todo-app, express-api)" -ForegroundColor White
Write-Host "  🐍 Python projects: 2 (data-science-project, ml-project)" -ForegroundColor White
Write-Host "  ☕ Java projects: 2 (spring-boot-api, android-app)" -ForegroundColor White
Write-Host "  🚀 Go project: 1 (go-microservice)" -ForegroundColor White
Write-Host "  📱 Flutter project: 1 (flutter-mobile-app)" -ForegroundColor White
Write-Host

# Calculate total size
try {
    $totalSize = (Get-ChildItem -Recurse $DemoDir | Measure-Object -Property Length -Sum).Sum
    $totalSizeMB = [math]::Round($totalSize / 1MB, 2)
    Write-Host "Total demo size: $totalSizeMB MB" -ForegroundColor Cyan
} catch {
    Write-Host "Could not calculate total size" -ForegroundColor Yellow
}

Write-Host
Write-Host "You can now test the cache remover with:" -ForegroundColor Green
Write-Host "  cache-remover-utility.exe --dry-run $DemoDir" -ForegroundColor White
Write-Host "  cache-remover-utility.exe --ui $DemoDir" -ForegroundColor White
Write-Host