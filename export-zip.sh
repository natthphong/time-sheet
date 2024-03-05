GOARCH=arm64 GOOS=linux go build -o bootstrap main.go
zip awslambda.zip bootstrap
rm -f bootstrap