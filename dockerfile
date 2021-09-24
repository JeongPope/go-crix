### Builder
FROM golang:1.17-alpine3.13 as builder
RUN apk update && apk add git

WORKDIR /tmp/crix
COPY . .

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-s -w' -o server .

### Excutable image
FROM alpine:3.13
COPY --from=builder /tmp/crix /

CMD [ "/server" ]