<p align="center"><img src="docs/assets/banner.png" width="500"></p>

# Strict Development Kit

Strict is a statically typed multi-paradigm programming language that is
compiled to SIR.

### Building from source

You need Bazel in order to build the SDK. 
https://docs.bazel.build/versions/master/install.html 

After bazel has been installed, run following commands in the root directory.
*If you are building on windows, you will have to run them from a bash command line.*

`bazel build ...`

To test the built source code, run the following:

`bazel test ...`

You can also directly run the strict program with bazel:
`bazel run //cmd/strict:strict`


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
