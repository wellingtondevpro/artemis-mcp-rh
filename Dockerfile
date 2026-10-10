FROM golang:alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o mcp-server main.go

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/mcp-server .
EXPOSE 8080
ENTRYPOINT ["./mcp-server", "-transport", "sse", "-port", "8080"]
