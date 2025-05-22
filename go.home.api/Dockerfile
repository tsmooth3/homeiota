FROM golang:1.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o go.home.api main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/go.home.api .

CMD ["./go.home.api"] 

EXPOSE 8080