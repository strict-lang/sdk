# How to install the Strict SDK on Windows

The windows installation is slightly more complicated than
the one for linux, since *Glide* does not provide an installation
script for windows.

## Installing Glide (Dependency-Manager)
Pick the binary that matches your platform:
https://github.com/Masterminds/glide/releases?ts=2

## Installing Mage (Build-Tool)

```
go get -u -d github.com/magefile/mage
cd $GOPATH/src/github.com/magefile/mage
go run bootstrap.go
```
## Installing the SDK
```
git clone https://gitlab.com/strict-lang/sdk.git $GOPATH/gitlab.com/strict-lang/sdk
cd $GOPATH/gitlab.com/strict-lang/sdk
mage install
```