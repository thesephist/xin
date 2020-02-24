CMD = ./xin.go
XIN = go run -race ${CMD}

all: init gen test

# initialize development workspace
init:
	go get github.com/rakyll/statik


test: gen
	go build -race -o ./xin
	./xin ./samples/first.xin
	./xin ./samples/hello.xin
	./xin ./samples/fact.xin
	./xin ./samples/fib.xin
	./xin ./samples/prime.xin
	./xin ./samples/collatz.xin
	./xin ./samples/map.xin
	./xin ./samples/list.xin
	./xin ./samples/async.xin
	./xin ./samples/stream.xin
	./xin ./samples/file.xin
	./xin ./samples/nest-import.xin
	# we echo in some input for prompt.xin testing stdin
	echo "Linus" | ./xin ./samples/prompt.xin
	./xin ./samples/test.xin
	rm ./xin


run-test: gen
	${XIN} ./samples/test.xin


# start interactive repl
repl: gen
	${XIN}


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
