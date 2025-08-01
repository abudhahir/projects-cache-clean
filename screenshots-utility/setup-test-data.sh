#!/bin/bash

# 📁 Test Data Setup for Screenshot Generation
# Creates consistent, realistic project structure for demos

set -e

echo "🏗️  Setting up test data for screenshot generation..."

# Clean up existing demo-projects
rm -rf demo-projects

# Create demo projects directory
mkdir -p demo-projects

# =============================================================================
# 🟢 Node.js Project
# =============================================================================
echo "📦 Creating Node.js project..."
mkdir -p demo-projects/react-app

# package.json
cat > demo-projects/react-app/package.json << 'EOF'
{
  "name": "react-app",
  "version": "1.0.0",
  "description": "A sample React application",
  "main": "src/index.js",
  "scripts": {
    "start": "react-scripts start",
    "build": "react-scripts build",
    "test": "react-scripts test"
  },
  "dependencies": {
    "react": "^18.2.0",
    "react-dom": "^18.2.0",
    "react-scripts": "5.0.1"
  },
  "devDependencies": {
    "@testing-library/react": "^13.4.0",
    "eslint": "^8.45.0"
  }
}
EOF

# Create realistic node_modules structure (sample files to simulate size)
mkdir -p demo-projects/react-app/node_modules/{react,react-dom,react-scripts}
echo "// React library" > demo-projects/react-app/node_modules/react/index.js
echo "// React DOM" > demo-projects/react-app/node_modules/react-dom/index.js
echo "// React Scripts" > demo-projects/react-app/node_modules/react-scripts/index.js

# Add some realistic files to simulate size
for i in {1..50}; do
    echo "// Generated module file $i" > "demo-projects/react-app/node_modules/module-$i.js"
done

# Create build directory
mkdir -p demo-projects/react-app/build/static/{css,js}
echo "/* Compiled CSS */" > demo-projects/react-app/build/static/css/main.css
echo "// Compiled JS" > demo-projects/react-app/build/static/js/main.js
echo "<!DOCTYPE html><html><head><title>React App</title></head></html>" > demo-projects/react-app/build/index.html

# Create source code (kept)
mkdir -p demo-projects/react-app/src
echo "import React from 'react';" > demo-projects/react-app/src/index.js
echo "function App() { return <h1>Hello World</h1>; }" > demo-projects/react-app/src/App.js

# =============================================================================
# 🐍 Python Project  
# =============================================================================
echo "🐍 Creating Python project..."
mkdir -p demo-projects/python-api

# requirements.txt
cat > demo-projects/python-api/requirements.txt << 'EOF'
fastapi==0.68.0
uvicorn==0.15.0
pydantic==1.8.2
pytest==6.2.5
requests==2.25.1
EOF

# Create __pycache__ directory
mkdir -p demo-projects/python-api/__pycache__
for i in {1..10}; do
    echo "# Compiled Python bytecode for module$i" > "demo-projects/python-api/__pycache__/module$i.cpython-39.pyc"
done

# Create .pytest_cache
mkdir -p demo-projects/python-api/.pytest_cache/{v,d}
echo "pytest cache data" > demo-projects/python-api/.pytest_cache/README.md
echo "cache data" > demo-projects/python-api/.pytest_cache/v/cache_data
echo "session data" > demo-projects/python-api/.pytest_cache/d/session_data

# Create build directory
mkdir -p demo-projects/python-api/{build,dist}
echo "# Build artifacts" > demo-projects/python-api/build/lib.txt
echo "# Distribution package" > demo-projects/python-api/dist/package.whl

# Create source code (kept)
mkdir -p demo-projects/python-api/src
cat > demo-projects/python-api/src/main.py << 'EOF'
from fastapi import FastAPI

app = FastAPI()

@app.get("/")
def read_root():
    return {"Hello": "World"}
EOF

cat > demo-projects/python-api/src/models.py << 'EOF'
from pydantic import BaseModel

class Item(BaseModel):
    name: str
    description: str
    price: float
EOF

# =============================================================================
# ☕ Java/Maven Project
# =============================================================================
echo "☕ Creating Java/Maven project..."
mkdir -p demo-projects/java-service

# pom.xml
cat > demo-projects/java-service/pom.xml << 'EOF'
<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0"
         xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
         xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 
         http://maven.apache.org/xsd/maven-4.0.0.xsd">
    <modelVersion>4.0.0</modelVersion>
    
    <groupId>com.example</groupId>
    <artifactId>java-service</artifactId>
    <version>1.0.0</version>
    <packaging>jar</packaging>
    
    <properties>
        <maven.compiler.source>11</maven.compiler.source>
        <maven.compiler.target>11</maven.compiler.target>
    </properties>
    
    <dependencies>
        <dependency>
            <groupId>org.springframework.boot</groupId>
            <artifactId>spring-boot-starter-web</artifactId>
            <version>2.7.0</version>
        </dependency>
    </dependencies>
</project>
EOF

# Create target directory (Maven build output)
mkdir -p demo-projects/java-service/target/{classes,test-classes,surefire-reports}
echo "// Compiled Java class" > demo-projects/java-service/target/classes/Main.class
echo "// Test class" > demo-projects/java-service/target/test-classes/MainTest.class
echo "Test results" > demo-projects/java-service/target/surefire-reports/results.xml

# Add realistic JAR files
for i in {1..20}; do
    echo "JAR content $i" > "demo-projects/java-service/target/dependency-$i.jar"
done

# Create source code (kept)
mkdir -p demo-projects/java-service/src/{main,test}/java/com/example
cat > demo-projects/java-service/src/main/java/com/example/Main.java << 'EOF'
package com.example;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@SpringBootApplication
public class Main {
    public static void main(String[] args) {
        SpringApplication.run(Main.class, args);
    }
}
EOF

# =============================================================================
# 🦀 Rust Project (Bonus)
# =============================================================================
echo "🦀 Creating Rust project..."
mkdir -p demo-projects/rust-cli

# Cargo.toml
cat > demo-projects/rust-cli/Cargo.toml << 'EOF'
[package]
name = "rust-cli"
version = "0.1.0"
edition = "2021"

[dependencies]
clap = "4.0"
serde = { version = "1.0", features = ["derive"] }
tokio = { version = "1.0", features = ["full"] }
EOF

# Create target directory
mkdir -p demo-projects/rust-cli/target/{debug,release}
echo "// Debug binary" > demo-projects/rust-cli/target/debug/rust-cli
echo "// Release binary" > demo-projects/rust-cli/target/release/rust-cli

# Add build artifacts
for i in {1..15}; do
    echo "Build artifact $i" > "demo-projects/rust-cli/target/debug/deps/artifact_$i"
done

# Create source code (kept)
mkdir -p demo-projects/rust-cli/src
cat > demo-projects/rust-cli/src/main.rs << 'EOF'
use clap::Parser;

#[derive(Parser)]
#[command(name = "rust-cli")]
#[command(about = "A sample Rust CLI application")]
struct Cli {
    #[arg(short, long)]
    name: Option<String>,
}

fn main() {
    let cli = Cli::parse();
    
    match cli.name {
        Some(name) => println!("Hello, {}!", name),
        None => println!("Hello, World!"),
    }
}
EOF

# =============================================================================
# 📊 Generate Size Information
# =============================================================================
echo "📊 Calculating demo project sizes..."

# Function to get directory size
get_size() {
    du -sh "$1" 2>/dev/null | cut -f1 | tr -d '\t'
}

echo ""
echo "📁 Demo Project Structure Created:"
echo "├── react-app (Node.js)"
echo "│   ├── node_modules/ (~$(get_size demo-projects/react-app/node_modules 2>/dev/null || echo "N/A"))"
echo "│   └── build/ (~$(get_size demo-projects/react-app/build 2>/dev/null || echo "N/A"))"
echo "├── python-api (Python)"  
echo "│   ├── __pycache__/ (~$(get_size demo-projects/python-api/__pycache__ 2>/dev/null || echo "N/A"))"
echo "│   ├── .pytest_cache/ (~$(get_size demo-projects/python-api/.pytest_cache 2>/dev/null || echo "N/A"))"
echo "│   └── build/ (~$(get_size demo-projects/python-api/build 2>/dev/null || echo "N/A"))"
echo "├── java-service (Java/Maven)"
echo "│   └── target/ (~$(get_size demo-projects/java-service/target 2>/dev/null || echo "N/A"))"
echo "└── rust-cli (Rust)"
echo "    └── target/ (~$(get_size demo-projects/rust-cli/target 2>/dev/null || echo "N/A"))"
echo ""

echo "✅ Test data setup complete!"
echo "💡 You can now run: ./generate.sh to create all screenshots"