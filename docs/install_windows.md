# How to install the Strict SDK on Windows

The windows installation is slightly more complicated than
the one for linux, since *Glide* does not provide an installation
script for windows.

## Installing Glide (Dependency-Manager)
Pick the binary that matches your platform:
https://github.com/Masterminds/glide/releases?ts=2
- most likely this one: glide-v0.13.3-windows-amd64.zip
- extract the zip file and put the glide.exe file into your %GOBIN% path (e.g. c:\go\bin)

## Installing Mage (Build-Tool)

```
go get -u -d github.com/magefile/mage
cd %GOPATH%/src/github.com/magefile/mage
go run bootstrap.go
```
## Installing the SDK
```
<<<<<<< docs/install_windows.md
git clone https://gitlab.com/strict-lang/sdk.git %GOPATH%/src/gitlab.com/strict-lang/sdk
cd %GOPATH%/src/gitlab.com/strict-lang/sdk
mage install
```

## Test strict
Now you should have strict.exe in c:\go\bin
- Type "strict -v" to see your strict version