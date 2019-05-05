compile-win: clean-win
	go-bindata -pkg assets -o assets.go assets/... && MOVE assets.go pkg/assets && go build cmd/neighbors.go

compile: clean
	go-bindata -pkg assets -o assets.go assets/templates/... && mv assets.go pkg/assets && go build cmd/neighbors.go

clean-win:
	cmd \/C DEL neighbors.exe pkg\assets\assets.go

clean:
	rm -f neighbors.exe neighbors */**/assets.go
