FROM golang:1.23.1-alpine3.20

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

COPY cmd/       /app/cmd/
COPY config/    /app/config/
COPY docs/      /app/docs/
COPY internal/  /app/internal/
COPY pkg/       /app/pkg/
COPY go.mod     /app/go.mod
COPY go.sum     /app/go.sum

WORKDIR /app

RUN go mod download
RUN go build -o ./build/go-rest-api ./cmd/go-rest-api

CMD [ "./build/go-rest-api" ]