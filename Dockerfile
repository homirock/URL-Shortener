FROM golang:1.19-alpine
WORKDIR /app
COPY . .
RUN go build -o url-shortener
EXPOSE 8084
CMD ["./url-shortener"]
