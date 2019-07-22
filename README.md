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
```
git clone git@gitlab.com:strict-lang/sdk.git $GOPATH/gitlab.com/strict-lang/sdk
make deps
sudo make install
```

Congrats, you just cloned and built the Strict SDK.
