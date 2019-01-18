digest: generate
	go build github.com/farnasirim/digest/cmd/digest

generate: 
	go generate ./...

test: generate
	go test ./...
