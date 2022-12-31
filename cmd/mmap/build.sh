# build binary
GOOS=linux GOARCH=amd64 go build -o mmap-server

# docker build
docker build -t mmap-server:v0.1 .