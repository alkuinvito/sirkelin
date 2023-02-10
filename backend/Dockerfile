FROM golang:1.19-alpine AS builder
WORKDIR /go/src/github.com/alkuinvito/malakh-api/
RUN go install github.com/cespare/reflex@latest
COPY go.mod ./
COPY go.sum ./
COPY .env ./
COPY main.go ./
RUN go mod download
COPY . .
RUN go build -o app .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/alkuinvito/malakh-api/app .
COPY --from=builder /go/src/github.com/alkuinvito/malakh-api/.env .
EXPOSE 8080
CMD ["./app"]