build: bin
	go build -o bin/iam

bin:
	mkdir -p bin

.PHONY: build