

FROM golang:1.24

WORKDIR /app

# Copy only the compiled binary from builder
COPY . .

RUN go build -o app

EXPOSE 8000

CMD ["./app"]