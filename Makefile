
run: fmt
	go run main.go
	mv site docs

build: fmt
	go build -ldflags '-w -extldflags "-static"' -tags netgo -o .

fmt: rereadme
	gofmt -w -s *.go
	fixjsstyle ./js/*.js


rereadme:
	#pandoc README.md -o ./README.rst