FROM golang:1.24.6 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN mkdir /build && CGO_ENABLED=0 GOOS=linux go build -o /build/app ./cmd/app

FROM debian:bookworm-slim
WORKDIR /build  

COPY --from=builder /build/app /build/app

EXPOSE 8080

CMD ["/build/app"]