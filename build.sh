go build fastqsplit.go
cp fastqsplit fastqsplit_osx
GOOS=linux go build fastqsplit.go
cp fastqsplit fastqsplit_linux
