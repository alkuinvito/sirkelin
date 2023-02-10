FROM golang:1.19-alpine
WORKDIR /go/src/github.com/alkuinvito/backend/
COPY go.mod ./
COPY go.sum ./
COPY .env ./
COPY main.go ./
RUN go install github.com/githubnemo/CompileDaemon@latest
RUN go mod download
CMD ["CompileDaemon", "-command='./app'"]