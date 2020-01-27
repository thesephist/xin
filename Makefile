CMD = ./xin.go
RUN = go run -race ${CMD}

all: run


# initialize development workspace
init:
	go get github.com/rakyll/statik


run:
	${RUN}


gen:
	statik -src=./lib


# build for specific OS target
build-%:
	GOOS=$* GOARCH=amd64 go build -o xin-$* ${CMD}


# install on host system
install:
	cp util/xin.vim ~/.vim/syntax/xin.vim
	go install ${CMD}


# build for all OS targets, useful for releases
build: build-linux build-darwin build-windows
