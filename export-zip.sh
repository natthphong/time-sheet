GOARCH=amd64 GOOS=linux go build main.go
zip main.zip main
rm -f main