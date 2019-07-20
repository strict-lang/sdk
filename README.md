<p align="center"><img src="docs/assets/strict_logo.png" width="360"></p>

# Strict Development Kit

###  Strict Is
- A computer language to be written by a computer programs (and humans) intelligently (not really calling this AI)
- Also readable by humans and convertible from other languages (starting with C#, must be TDD and functional, converting
  not functional .NET code is possible, but might lead to side effects and bad performance)
- Strongly typed, statically typed, while still looking dynamic and uses very strong references (null is never allowed,
  assignment is not allowed, reassigning to the same name is allowed, but otherwise a very functional like language)
- able to infer types automatically, you never have to declare types in methods as variables (members) will be evaluated 
- without access modifiers (public, private, protected, sealed, etc.) and has no static types, methods or members
- without classes or OOP. There are only types (like structs in Go) and methods (like functions) plus members in those 
  (which are all statements after all and are assigned statement values). Still supports polymorphism through type nesting! 
- data is stored in components, which are used together in processors that do stuff with them (derived from DeltaEngine)
- thread safe and uses threads automatically to execute any work (no need to worry about race conditions or locking)
- Uses fluent syntax (written as plain text with spaces, no . or camel case required) in Collapsed Mode

The purpose of Strict is not to be human like or Artificially Intelligent, but actually be useful and intelligent in its
own language context. The name Strict is used for the "intelligence" used by it to write code and evaluate it, which is 
what makes Strict understand its own code in its context. The whole point of this language is to give itself a way to 
change its own code automatically without human intervention, computers are taking over, but don't worry it won't do 
anything useful as it is safe like functional programs and without IO it cannot change any state.

"Strict language" or "Strict syntax" refers to the syntax of the language, usually we just talk about its Statements 
(there are no expressions in the language, everything is a statement, but in collapsed syntax we can make it look like usual expressions).

When a type like "size" is defined somewhere and we want to add new features, we do not have to derive from it. 
We can simply create a new size type matching the data and it can be used instead. A really bad example is C#, 
where you need many different size, point, etc. classes for different frameworks (Windows Forms, WPF, DirectX, XNA, 
DeltaEngine, OpenTK, etc.), but they are all identical except for maybe some additional methods here and there. Of course
 if we have access to the original "size" definition, we can just use that, which is the normal way to handle this. We 
 also cannot have two conflicting size defined at the same time (or any method or member, or even anything conflicting 
 with a sub context name, everything must be unique).

### The Name Strict
Why the name Strict? Well, the language is very strict about its input, much more than all other languages, 
even python is more flexible with empty lines, continuation lines, comments, etc. which are all absent in Strict. 
Even in collapsed mode the language looks always the same way, formatting is always the same and there is usually just
one way to do things. Collapsed mode is only there for easier reading by humans and programmers and it reflects the 
expanded mode always in the same manner. If you look at the resulting expanded mode it is ridiculously easy with 
supporting only a few build in types and a handful of statements (around 10). Normally languages have 20-50 expressions
supported build in and depending on the language 20-100 statements and build in keywords supported. Collapsed strict is
similar, but smaller, expanded strict is much simpler (10x less). This is only possible by limiting on what the language
can do and hiding away many facets (like threading, exception handling, no generics, no overloading, etc.). Strict is 
not limiting itself to hurt a programmer (which is not supposed to write much code in the language anyway), but instead 
to enable its Intelligence to write code inside a very strict and rigid system where it cannot make many mistakes. By 
following TDD problems are also always encountered in a small context and refactorings help to only work on one issue at a time. 

### Building from source

When building from source you will first need some prerequisites:
 
  - Git (https://git-scm.com/) 
  - Golang (https://golang.org/doc/install)
  - Make (https://www.gnu.org/software/make/)

Once you have all prerequisites, execute following commands in your shell:
```
git clone git@gitlab.com:strict-lang/sdk.git $GOPATH/gitlab.com/strict-lang/sdk
make deps
sudo make install
```

Congrats, you just cloned and built the Strict SDK.
