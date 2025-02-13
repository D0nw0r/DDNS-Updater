FROM golang:1.20-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o myapp

# Step 2: Create a smaller image for running the Go app
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Go app from the builder image
COPY --from=builder /app/myapp .

# Expose the port the app will run on
EXPOSE 8080

# Run the Go app
CMD ["./myapp"]
