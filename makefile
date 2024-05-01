RM = rm -f


deps:
	go mod tidy

run:
	go run text-clipper.go

build:
	go build text-clipper.go

clean:
	$(RM) text-clipper
