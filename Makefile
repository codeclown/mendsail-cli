build:
	go build -o bin src/mendsail.go src/send.go

test:
	go test ./...
