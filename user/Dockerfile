FROM golang:1.25-alpine
WORKDIR /app
COPY . /app/
RUN go mod download
EXPOSE 3002
ENTRYPOINT [ "go", "run", "src/main.go"]
