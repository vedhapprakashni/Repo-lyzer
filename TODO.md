# Menu Navigation Code Duplication Refactoring

## Current Status
- [x] Analyze menu.go for code duplication in cursor navigation
- [x] Create refactoring plan with helper functions
- [x] Get user approval for plan

## Implementation Steps
- [ ] Add helper functions: moveCursorUp(), moveCursorDown(), moveCursorHome(), moveCursorEnd()
- [ ] Refactor "up"/"k"/"w"/"W" key handling to use moveCursorUp()
- [ ] Refactor "down"/"j"/"S" key handling to use moveCursorDown()
- [ ] Refactor "home"/"g" key handling to use moveCursorHome()
- [ ] Refactor "end"/"G" key handling to use moveCursorEnd()
- [ ] Test navigation works correctly in both main menu and submenus
- [ ] Verify edge cases (empty menus, single item menus)
