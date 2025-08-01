# 🧪 Testing Instructions for Cache Remover Utility

## 📋 **Prerequisites**
1. **Go installed** (version 1.19+)
   ```bash
   go version  # Should show go1.19 or higher
   ```
2. **Git repository cloned**
   ```bash
   git clone https://github.com/abudhahir/projects-cache-clean.git
   cd projects-cache-clean
   ```

## 🚀 **Quick Test Setup**

### 1. **Build the Application**
```bash
go build -o cache-remover
```

### 2. **Create Test Projects**
```bash
# Create a test directory with sample projects
mkdir -p test-area/{node-project,python-project,java-project}

# Node.js project
cd test-area/node-project
echo '{"name":"test-app","version":"1.0.0"}' > package.json
mkdir -p node_modules/.bin dist build
echo "dummy" > node_modules/test.js
echo "compiled" > dist/app.js

# Python project
cd ../python-project
echo "flask==2.0.1" > requirements.txt
mkdir -p __pycache__ .pytest_cache build
echo "dummy" > __pycache__/test.pyc
echo "cache" > .pytest_cache/test

# Java project
cd ../java-project
echo '<project></project>' > pom.xml
mkdir -p target/classes
echo "compiled" > target/classes/App.class

cd ../..
```

## 🧪 **Testing Scenarios**

### **Test 1: Dry Run Mode (Safe Preview)**
```bash
./cache-remover --dry-run test-area/
```
**Expected Output:**
- Shows found projects (node-project, python-project, java-project)
- Lists cache items that WOULD be removed
- Shows total space that would be reclaimed
- ❌ No files actually removed

### **Test 2: Verbose Mode**
```bash
./cache-remover --dry-run --verbose test-area/
```
**Expected Output:**
- More detailed output
- Shows "Skipping cache directory" messages
- Shows warnings for permission errors (if any)

### **Test 3: Interactive TUI Mode**
```bash
./cache-remover --ui test-area/
```
**Expected Behavior:**
1. **Loading Screen**: Shows spinner while scanning
2. **Project List**: Shows all 3 test projects with cache sizes
3. **Navigation**: 
   - Use ↑/↓ or j/k to move
   - Press Space to select/deselect projects
   - Press 'a' to select all
   - Press 'd' to deselect all
4. **Details View**: Press 'v' on a project to see cache breakdown
5. **Cleaning**: 
   - Press 'c' to clean selected projects
   - Confirm with 'y'
   - See "Cleaning [project-name]..." messages
   - View final statistics
6. **Exit**: Press 'q' or Esc

### **Test 4: Command Line Cleanup**
```bash
# First verify what will be removed
./cache-remover --dry-run test-area/

# Then actually clean (will ask for confirmation)
./cache-remover test-area/
```
**Expected:**
- Shows summary of what will be removed
- Asks "Continue? [y/N]:"
- After 'y', removes cache and shows statistics

### **Test 5: Configuration Testing**
```bash
# List all supported project types
./cache-remover --list-types

# Save configuration file
./cache-remover --save-config

# Check the generated file
cat cache-remover-config.json
```

### **Test 6: Performance Testing**
```bash
# Test with different worker counts
./cache-remover --workers 8 --verbose test-area/

# Test depth limiting
./cache-remover --max-depth 2 --verbose test-area/
```

### **Test 7: Error Handling**
```bash
# Test with non-existent directory
./cache-remover /non/existent/path

# Test with permission issues (if on Unix/Linux)
mkdir -p test-area/restricted
chmod 000 test-area/restricted
./cache-remover --verbose test-area/
chmod 755 test-area/restricted  # Restore permissions
```

## 🔍 **Verification Steps**

### **After Cleaning, Verify:**
```bash
# Check that cache directories are gone
ls -la test-area/node-project/  # node_modules, dist, build should be gone
ls -la test-area/python-project/  # __pycache__, .pytest_cache, build gone
ls -la test-area/java-project/  # target directory gone

# But project files remain
cat test-area/node-project/package.json  # Should still exist
cat test-area/python-project/requirements.txt  # Should still exist
cat test-area/java-project/pom.xml  # Should still exist
```

## 🧪 **Unit Tests**
```bash
# Run all tests
go test -v

# Run with race detection
go test -race -v

# Run specific test
go test -v -run TestFullWorkflowIntegration
```

## 📊 **Expected Test Results**

### **Good Signs ✅**
- All projects detected correctly
- Cache sizes calculated accurately
- Only cache directories removed, project files intact
- TUI responds smoothly to keyboard input
- Confirmation prompts work correctly
- Statistics show correct counts and sizes

### **Common Issues to Check ❌**
- Permission errors are handled gracefully
- Deep directory structures don't cause hangs
- Large cache directories are processed efficiently
- Cancellation (Ctrl+C) exits cleanly

## 🛠️ **Advanced Testing**

### **Create Larger Test Set**
```bash
# Create 10 projects with varying cache sizes
for i in {1..10}; do
    mkdir -p test-area/project-$i
    echo '{"name":"project-'$i'"}' > test-area/project-$i/package.json
    mkdir -p test-area/project-$i/node_modules
    # Create files of different sizes
    dd if=/dev/zero of=test-area/project-$i/node_modules/large-$i.dat bs=1M count=$i 2>/dev/null
done

# Test performance
time ./cache-remover --dry-run test-area/
```

### **Test Custom Configuration**
```bash
# Edit cache-remover-config.json to add custom project type
# Then test if it's detected
./cache-remover --dry-run test-area/
```

## 🧹 **Cleanup After Testing**
```bash
# Remove all test data
rm -rf test-area/
rm -f cache-remover-config.json
```

## 💡 **Tips**
- Always use `--dry-run` first to preview changes
- Test in a separate directory, not your actual projects
- The TUI mode provides the best user experience
- Check the exit code: `echo $?` (should be 0 for success)

Let me know if you need help with any specific test scenario!