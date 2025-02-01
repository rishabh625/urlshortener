# Use official Golang image as a base image
FROM golang:1.23-alpine

# Set the working directory
WORKDIR /app

# Copy go modules file and download dependencies
COPY go.mod go.sum ./
RUN go mod tidy
RUN go mod vendor

# Copy the entire project into the container
COPY . .

# Expose the port that the server will run on
EXPOSE 8080

#Set GOFlags
RUN export GOFLAGS=""

# Build the application
RUN go build -o urlshortener-service cmd/main.go

# Command to run the application
CMD ["./urlshortener-service"]