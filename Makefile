CMD = ./xin.go
XIN = go run -race ${CMD}

all: init gen test

# initialize development workspace
init:
	go get github.com/rakyll/statik


test: gen
	go build -race -o ./xin
	./xin run ./samples/first.xin
	./xin run ./samples/hello.xin
	./xin run ./samples/fact.xin
	./xin run ./samples/fib.xin
	./xin run ./samples/prime.xin
	./xin run ./samples/collatz.xin
	./xin run ./samples/map.xin
	./xin run ./samples/list.xin
	./xin run ./samples/async.xin
	./xin run ./samples/stream.xin
	./xin run ./samples/file.xin
	./xin run ./samples/nest-import.xin
	# we echo in some input for prompt.xin testing stdin
	echo "Linus" | ./xin run ./samples/prompt.xin
	./xin run ./samples/test.xin
	rm ./xin


# start interactive repl
repl: gen
	${XIN} repl


# re-generate static files
# that are bundled into the executable (standard library)
gen:
	statik -src=./lib


# build for specific OS target
build-%: gen
	GOOS=$* GOARCH=amd64 go build -o xin-$* ${CMD}


# install on host system
install: gen
	cp util/xin.vim ~/.vim/syntax/xin.vim
	go install ${CMD}
	ls -l `which xin`


# build for all OS targets, useful for releases
build: build-linux build-darwin build-windows

# clean any generated files
clean:
	rm -rvf xin xin-*
