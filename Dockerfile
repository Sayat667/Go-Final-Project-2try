# Use the official Golang image as a base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the source code into the container
COPY . .

# Build the Go application and name the executable as "main"
RUN go build -o main .

# Expose the port your application listens on
EXPOSE 8080

