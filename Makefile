# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

# Run tests
test: fmt vet
	go test ./...  -coverprofile=coverage.out
	go run main.go generate -q

