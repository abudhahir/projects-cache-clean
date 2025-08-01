# Code Review Checklist

## Security Review
- [ ] No global mutable state variables (see ADR-001)
- [ ] No race condition possibilities 
- [ ] Proper mutex usage if shared state required
- [ ] No memory leaks or resource management issues
- [ ] Input validation for all user data
- [ ] No hardcoded secrets or credentials

## Architecture Review  
- [ ] Follows Bubble Tea patterns correctly (model-update-view)
- [ ] State encapsulated in appropriate structs
- [ ] Message passing used for component communication
- [ ] Single responsibility principle followed
- [ ] No anti-patterns (global state, mixed paradigms)

## Code Quality
- [ ] Functions are focused and small
- [ ] Clear variable and function names
- [ ] Proper error handling throughout
- [ ] Resource cleanup with defer patterns
- [ ] No unused imports or variables

## Testing
- [ ] Unit tests cover new functionality
- [ ] Integration tests for user workflows
- [ ] Race condition tests for concurrent code
- [ ] Resource leak detection in tests
- [ ] Edge cases and error conditions tested

## Performance
- [ ] No unnecessary goroutine creation
- [ ] Timer cleanup mechanisms present
- [ ] Memory usage patterns validated
- [ ] No blocking operations in UI thread

## Documentation
- [ ] Public APIs documented
- [ ] Complex logic explained with comments
- [ ] README updated if needed
- [ ] Breaking changes noted

## Regression Prevention
- [ ] Changes don't reintroduce known issues
- [ ] Performance doesn't degrade significantly
- [ ] All existing tests still pass
- [ ] New functionality doesn't break existing features

---
**Note:** This checklist must be completed for all PRs. Any "No" items must be addressed before merge.