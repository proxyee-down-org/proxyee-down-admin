SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
SET GOPATH=%~dp0;%GOPATH%
go build -o bin/ProxyeeDownExtensionServer src/main.go