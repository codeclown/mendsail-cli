OUT ?= mendsail
LIBS = src/mendsail.go src/send.go src/post.go

build:
	mkdir -p bin && go build -o bin/$(OUT) $(LIBS)

test:
	go test ./...
