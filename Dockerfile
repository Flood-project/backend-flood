FROM golang:1.24.4-alpine3.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

WORKDIR /app/cmd

RUN CGO_ENABLE=0 go build -ldflags="-s -w" -o /app/server ./cmd/main.go

FROM scratch

COPY --from=builder /app/server /app/server

WORKDIR /app
CMD ["/app/server"]
