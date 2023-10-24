# Official Golang image
FROM golang:1.20-alpine

# Run Go ver
RUN go version

# Set working directory
WORKDIR /app

# Copy the entire code
COPY . .

# Download and install dependencies
RUN go get -d -v ./...

# Build the Go app
RUN go build -o app ./cmd/app

# Expose the port
EXPOSE 8000

# Run the executable
CMD ["./app"]


