.PHONY: build clean run

BINARY_NAME=dops

build:
	go build -o $(BINARY_NAME) main.go

run:
	./$(BINARY_NAME)

clean:
	rm -f $(BINARY_NAME)
