RM = rm -f


deps:
	go mod tidy

run:
	go run ./cmd/text-clipper.go

build:
	go build ./cmd/text-clipper.go

clean:
	$(RM) text-clipper
