# Use an official Go runtime as a parent image
FROM golang:1.19-alpine

# Set the working directory to /app
WORKDIR /opt/code

# Copy the current directory contents into the container at /app
COPY . /opt/code

RUN apk update && apk upgrade && \
    apk add --no-cache git

RUN go mod download

# Build the Go app
RUN go build -o bin/bta cmd/apiserver/main.go

# Expose port 8080 for the application
EXPOSE 8085

# Define the command to run the executable
ENTRYPOINT ["/opt/code/bin/bta"]
