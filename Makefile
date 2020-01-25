CMD = ./cmd/xin.go
RUN = go run -race ${CMD}

all: run

run:
	${RUN}


# build for specific OS target
build-%:
	GOOS=$* GOARCH=amd64 go build ${LDFLAGS} -o ink-$* ${CMD}


# install on host system
install:
	cp util/xin.vim ~/.vim/syntax/xin.vim


# build for all OS targets, useful for releases
build: build-linux build-darwin build-windows
