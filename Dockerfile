FROM golang:1.23

RUN go version
ENV GOPATH=/

COPY ./ ./
RUN go mod download
RUN go build -o todo-list-service ./cmd/main.go

cmd ["./todo-list-service"]