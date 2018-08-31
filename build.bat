SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
SET GOPATH=%~dp0;%GOPATH%
go build -o bin/ProxyeeDownAdmin src/main.go