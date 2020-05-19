<img src="assets/banner.png" width="1000" alt="logo">

# Coding Style

Go coding style is based on https://golang.org/doc/effective_go.html
- Execute fmt for proper formatting before committing, Ctrl+Shift+Alt+P to format whole project (and use Ctrl+Shift+Alt+F when saving files)
- See .editorconfig for rules (tabs, 2 spaces used)

Use go test -cover (or GoLand 'Run with Coverage') to check all important packages have at least 90% coverage: https://blog.golang.org/cover
- Excluded are tools, visual programs or things based on already tested low level parts (either own or external)

Another guideline for general style (mostly for c#) http://deltaengine.net/learn/codingstyle