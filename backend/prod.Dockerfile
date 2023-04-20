FROM golang:1.19-alpine AS builder
WORKDIR /go/src/github.com/alkuinvito/backend/
ENV APP_MODE=release
COPY go.mod go.sum .env main.go ./
RUN go mod download
COPY . .
RUN go build -o app .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/alkuinvito/backend/app .
COPY --from=builder /go/src/github.com/alkuinvito/backend/.env .
COPY --from=builder /go/src/github.com/alkuinvito/backend/adminsdk-key.json .
EXPOSE 8080