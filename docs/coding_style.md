<p align="center"><img src="docs/assets/strict_logo.png" width="360"></p>

# Coding Style

Go coding style is based on https://golang.org/doc/effective_go.html<br/>
- Execute fmt for proper formatting before committing, Ctrl+Shift+Alt+P to format whole project (and use Ctrl+Shift+Alt+F when saving files)
- See .editorconfig for rules (tabs, 2 spaces used)

Use go test -cover (or GoLand 'Run with Coverage') to check all important packages have at least 90% coverage: https://blog.golang.org/cover
- Excluded are tools, visual programs or things based on already tested low level parts (either own or external)

Another guideline for general style (mostly for c#) http://deltaengine.net/learn/codingstyle

# Example from Modern C++ Challenge book
1. Sum of naturals divisible by 3 and 5 Write a program that calculates and prints the sum of all the natural numbers divisible by either 3 or 5, up to a given limit entered by the user.

TODO: code in Strict, see old examples: https://strict.fogbugz.com/f/page?W9
```cpp
int main()
{
  unsigned int limit = 0;
  std::cout << "Upper limit:";
  std::cin >> limit;
  unsigned long long sum = 0;
  for (unsigned int i = 3; i < limit; ++i)
  { 
    if (i % 3 == 0 || i % 5 == 0)
      sum += i; 
  }
  std::cout << "sum=" << sum << std::endl;
}
```

-> what comes out as ast<br />
-> output
