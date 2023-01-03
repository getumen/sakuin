compile:
	protoc fieldindex/*.proto --go_out=. --go_opt=paths=source_relative --proto_path=.

test:
	go test -timeout 30s -race ./...
