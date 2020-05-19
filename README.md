<p align="center"><img src="docs/assets/banner.png" width="500"></p>

# Strict Development Kit
[![Build Status](https://api.travis-ci.org/strict-lang/sdk.svg?branch=master)](https://travis-ci.org/github/strict-lang/sdk)

Strict is a statically typed multi-paradigm programming language. The strict CLI is used to generate .silk (Strict Intermediate Language Kit) packages and run the Strict Virtual Machine Bloom for fast execution or generate code for one of the supported backends (C++, C#, Java, Arduino, JavaScript, WebAssembly, etc.).

To write Strict code use the [Strict IDE IntelliJ plugin](https://github.com/strict-lang/strict-intellij), the command line tools are never required for normal users. 

### Building from Source

In order to build, you will need the [latest version of go](https://golang.org/).
Clone the repository and run:

```shell script
go get -d -v all
go install ./cmd/strict
```

To check if everything worked just call strict:

```shell script
strict version
```

### Compiling and Running your first Strict program

This small tutorial shows you how to write, compile and run a small strict
program. Once you have the strict command line tool from above and a working C++ Compiler (e.g. gcc or Visual Studio), open up your text editor of choice
and write some simple strict code ([see the official docs for more tips](https://strict.dev/docs/Overview)):

```strict
implement App
has log Log

method Run()
	log("Hello World")
```

Save the code into a file called `src/Hello.strict` (or any other filename ending with `.strict`) and invoke the strict compiler:

```
strict build Hello.strict -b c++ -r pretty-json
```

If you run into trouble here because the App framework is missing, just write some library code instead:

```strict
method Double(number Number) returns Number
	return number * 2
```

The strict build command now compiles the strict source code into C++ code. The
generated code is written into `Hello.h` and `Hello.cpp` files and has to be compiled
using a C++ compiler in order for you to execute it.

Compile the program using your C++ compiler of choice and run the
generated binary.

```
c++ Hello.cpp -o Hello.exe
./Hello.exe
```

Congrats! You have just written, compiled and run your first strict program.
