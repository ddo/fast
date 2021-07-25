# fast [![Github All Releases](https://img.shields.io/github/downloads/ddo/fast/total.svg?style=flat-square)]()
> Minimal zero-dependency utility for testing your internet download speed from terminal

*Powered by Fast.com - Netflix*

<p align="center"><a href="https://asciinema.org/a/80106"><img src="https://asciinema.org/a/80106.png" width="50%"></a></p>

## Installation

#### Bin

> replace the download link with your os one

> https://github.com/ddo/fast/releases

> below is ubuntu 64 bit example

```sh
curl -L https://github.com/ddo/fast/releases/download/v0.0.4/fast_linux_amd64 -o fast

# or wget
wget https://github.com/ddo/fast/releases/download/v0.0.4/fast_linux_amd64 -O fast

# then chmod
chmod +x fast

# run
./fast
```

#### Docker

> ~10 MB

```sh
docker run --rm ddooo/fast
```

#### Snap

```sh
snap install fast
```

#### Arch Linux (AUR)

```sh
yay -S fast || paru -S fast
```

#### Brew

> *soon*

#### For golang user

> golang user can install from the source code

```sh
go get -u github.com/ddo/fast
```

## Build

#### Docker

```sh
# build alpine binary file from root folder
docker run --rm -v "$PWD":/go/src/fast -w /go/src/fast golang:alpine go build -v

mv fast build/docker/
cd build/docker/
docker build -t ddooo/fast .
```

#### Snap

```sh
snapcraft
snapcraft push fast_*.snap
snapcraft release fast <revision> <channel>
```

## Bug

for bug report just open new issue
