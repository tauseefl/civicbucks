# Start from golang:1.16-alpine base image
FROM golang:1.16-alpine

# The latest alpine images don't have some tools like (`git` and `bash`).
# Adding git, bash to the image
RUN apk add --no-cache bash git

# Set the Current Working Directory inside the container
WORKDIR /

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependancies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o main .

# Expose port 8082 to the outside world
EXPOSE 8082

# Run the executable
ENTRYPOINT ["./main"]