# Build-Stage
FROM golang:1.21 AS build-env

WORKDIR /app

# Kopieren Sie Mod-Dateien und laden Sie die Abhängigkeiten herunter
COPY go.mod go.sum ./
RUN go mod download

# Kopieren Sie den Quellcode und kompilieren Sie die Anwendung
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# Final Stage
FROM alpine:latest

WORKDIR /app

# Kopieren Sie die kompilierte Anwendung vom Build-Stage
COPY --from=build-env /app/main /app/

# Port und CMD
EXPOSE 3000
CMD ["./main"]
