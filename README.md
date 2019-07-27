<p align="center"><img src="docs/assets/strict_logo.png" width="360"></p>

# Strict Development Kit

Strict is a statically typed multi-paradigm programming language that is
compiled to SIR.

### Building from source

When building from source you will first need some prerequisites:
 
  - Git (https://git-scm.com/) 
  - Golang (https://golang.org/doc/install)
  - Make (https://www.gnu.org/software/make/)

Once you have all prerequisites, execute following commands in your shell:

Linux/Unix: [install_nix.sh](/docs/install_linux.md)
```
git clone git@gitlab.com:strict-lang/sdk.git $GOPATH/gitlab.com/strict-lang/sdk
cd $GOPATH/gitlab.com/strict-lang/sdk
make deps
sudo make
```

Congrats! You just cloned and built the Strict SDK.

### Compiling and Running your first Strict proram

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