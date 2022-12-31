FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY mmap-server .

CMD ["./mmap-server"]