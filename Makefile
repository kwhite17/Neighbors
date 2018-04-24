install: test after-test
	GOBIN=${GOPATH}/bin go install ./cmd/neighbors.go

test: compile before-test
	go test ./pkg/*

compile:
	dep ensure && go build ./pkg/*

after-test:
	rm -r ./pkg/**/templates

before-test:
	mkdir ./pkg/users/templates ./pkg/items/templates && cp -r ./templates/items ./pkg/items/templates/ && cp -r ./templates/users ./pkg/users/templates/
