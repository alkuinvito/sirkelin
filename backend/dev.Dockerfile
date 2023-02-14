FROM golang:1.19-alpine
WORKDIR /go/src/github.com/alkuinvito/backend
RUN go install github.com/githubnemo/CompileDaemon@latest
COPY go.mod ./
COPY go.sum ./
COPY .env ./
COPY main.go ./
RUN go mod download
EXPOSE 8080