FROM golang:alpine

COPY go.mod .
COPY go.sum .
ENV GOPATH=/
RUN go mod download

#build appliction
COPY . .
RUN go build -o voting-service ./cmd/main/app.go

CMD ["./voting-service"]