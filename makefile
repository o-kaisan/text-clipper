.PHONY: run

RM = rm -f
export TEXT_CLIPPER_PATH := $(CURDIR)
export TEXT_CLIPPER_DEBUG := true

deps:
	go mod tidy

run:
	@go run text-clipper.go

build:
	go build text-clipper.go

clean:
	$(RM) text-clipper
