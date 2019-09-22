heroku:
	go build cmd/neighbors.go

compile-win: clean-win
	go-bindata -pkg assets -o assets.go assets/... && IF NOT EXIST pkg\assets MKDIR pkg\assets && MOVE assets.go pkg\assets\assets.go && go build cmd/neighbors.go

compile: clean
	go-bindata -pkg assets -o assets.go assets/... && mkdir -p pkg/assets && cp assets.go pkg/assets/assets.go && rm assets.go && go build cmd/neighbors.go

clean-win:
	IF EXIST neighbors.exe cmd \/C DEL neighbors.exe && IF EXIST pkg\assets\assets.go cmd \/C RMDIR /S /Q pkg\assets

clean:
	rm -rf *.exe neighbors */**/assets.go ./pkg/assets
