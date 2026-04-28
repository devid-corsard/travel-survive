build: main.go
	go build -o bin/game main.go

run: build
	./bin/game