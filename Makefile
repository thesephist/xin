CMD = ./cmd/xin.go
RUN = go run -race ${CMD}

all: run

run:
	${RUN}


# build for specific OS target
build-%:
	GOOS=$* GOARCH=amd64 go build ${LDFLAGS} -o ink-$* ${CMD}


# build for all OS targets, useful for releases
build: build-linux build-darwin build-windows
