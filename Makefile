compile:
	go-bindata -pkg retriever -o templates.go templates/... && mv templates.go pkg/retriever && go build cmd/neighbors.go

compile-win:
	go-bindata -pkg assets -o assets.go assets/... && MOVE assets.go pkg/assets && go build cmd/neighbors.go

clean-win:
	DEL neighbors.exe

clean:
	rm neighbors.exe