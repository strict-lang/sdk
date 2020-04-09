<p align="center"><img src="docs/assets/banner.png" width="500"></p>

# Strict Development Kit
![Build Status](https://api.travis-ci.org/strict-lang/sdk.svg?branch=master)

Strict is a statically typed multi-paradigm programming language that is
compiled to SIR.

### Building from Source

In order to build, you will need the [latest version of go](https://golang.org/).
Clone the repository and run:

```shell script
go get -d -v all
go install ./cmd/strict
```

### Compiling and Running your first Strict program

Prerequisites:
  - Working C++ Compiler
  - Strict SDK

This small tutorial shows you how to write, compile and run a small strict
program. Once you have all prerequisites, open up your text editor of choice
and write some simple strict code:

```strict
log("Hello, World!")
```

Save the code into a file called `hello_world.strict` (or any other filename
ending with `.strict`) and invoke the strict compiler:

```
strict build hello_world.strict
```

The strict build command now compiles the strict source code into C++ code. The
generated code is written into the `hello_world.cc` file and has to be compiled
using a C++ compiler in order for you to execute it.

Compile the program using your C++ compiler of choice and run the
generated binary.

```
c++ hello_world.cc -o hello_world.exe
./hello_world.exe
```

Congrats! You have just written, compiled and run your first strict program.
