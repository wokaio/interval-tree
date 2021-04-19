go mod init github.com/miczone/interval-tree
go mod tidy
go mod vendor
go mod verify
go mod edit -replace github.com/miczone/interval-tree/pkg=../pkg
go mod vendor