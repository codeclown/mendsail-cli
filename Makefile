build:
	go build -o bin src/mendsail.go src/send.go src/post.go

test:
	go test ./...
