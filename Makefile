compile: clean
	go-bindata -pkg assets -o assets.go assets/... && MKDIR pkg\assets && MOVE assets.go pkg\assets\assets.go && go build cmd/neighbors.go

compile:
	go-bindata -pkg assets -o assets.go assets/templates/... && mv assets.go pkg/assets/assets.go && go build cmd/neighbors.go

clean:
	IF EXIST neighbors.exe cmd \/C DEL neighbors.exe && IF EXIST pkg\assets\assets.go cmd \/C RMDIR /S /Q pkg\assets

clean:
	rm -f *.exe neighbors */**/assets.go
