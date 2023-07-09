FROM golang:alpine AS builder

WORKDIR /app

ENV GO111MODULE on
ENV CGO_ENABLED=1
ENV GOOS=linux

COPY go.mod .
COPY go.sum .
COPY . .

RUN go mod tidy
RUN go mod download
RUN go mod verify
RUN go mod vendor
RUN go fmt ./...

RUN apk update
RUN apk add gcc libc-dev

RUN go build -mod=vendor -a -installsuffix cgo -tags musl -o main ./cmd/main.go

FROM alpine:latest AS release

COPY --from=builder /app/main /app/cmd/

RUN chmod +x /app/cmd/main

WORKDIR /app

EXPOSE 8080

CMD ["cmd/main"]