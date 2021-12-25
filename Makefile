.PHONY: build
build: .deps .build .copy

.PHONY: .deps
.deps:
	ls go.mod || go mod init

.PHONY: .build
.build:
	OOS=darwin GOARCH=amd64 go build -o bin/hey_grpc cmd/hey_grpc/main.go

.PHONY: .copy
.copy:
	cp bin/hey_grpc /usr/local/bin/hey_grpc