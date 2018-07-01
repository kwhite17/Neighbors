install: test
	GOBIN=${GOPATH}/bin go install ./cmd/neighbors.go

test: compile
	go test ./pkg/items && go test ./pkg/users && go test ./pkg/login

compile: pre-compile
	go build ./pkg/*

pre-compile:
	dep ensure && go-bindata -pkg utils templates/... && mv bindata.go ./pkg/utils/
