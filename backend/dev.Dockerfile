FROM golang:1.19-alpine
WORKDIR /go/src/github.com/alkuinvito/backend
ENV APP_MODE=debug
RUN go install github.com/githubnemo/CompileDaemon@latest
COPY go.mod go.sum .env main.go ./
RUN go mod download
EXPOSE 8080