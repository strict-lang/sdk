# How to install the Strict SDK on Linux

## Installing Glide (Dependency-Manager)
### Ubuntu
```
sudo apt install golang-glide
```
### Other Distributions
```
curl https://glide.sh/get | sh
```
## Installing Mage (Build-Tool)

```
go get -u -d github.com/magefile/mage
cd $GOPATH/src/github.com/magefile/mage
go run bootstrap.go
```
## Installing the SDK
```
git clone git@gitlab.com:strict-lang/sdk.git $GOPATH/gitlab.com/strict-lang/sdk
cd $GOPATH/gitlab.com/strict-lang/sdk
mage install
```