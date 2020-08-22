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


test-release: test goreleaser
	goreleaser --rm-dist --skip-publish --snapshot

release: test goreleaser
	goreleaser --rm-dist

goreleaser:
ifeq (, $(shell which goreleaser))
 $(shell go get github.com/goreleaser/goreleaser)
endif