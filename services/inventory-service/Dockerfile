# Use the official Golang image as the base image
FROM golang:1.23

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod .
COPY go.sum .

# Download dependencies
RUN go mod tidy 

# Copy the source code
COPY . .

# Build the application
RUN go build -o inventory-service .

# Run the application
CMD ["./inventory-service"]
