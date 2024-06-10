# Use an official Golang runtime as a parent image
FROM golang:1.21.10-alpine3.20
# Set environment variables
ENV user=myuser
ENV uid=1000
ENV home=/home/$user
ENV app=/go/src/app

# Create a non-root user
RUN addgroup -S $user && adduser -S -G $user -u $uid -h $home $user

# Create the application directory
RUN mkdir $app

# Set the working directory to /go/src/app
WORKDIR /go/src/app

# Copy the current directory contents into the container at /go/src/app
COPY . /go/src/app

# Change ownership of the application directory to the non-root user
RUN chown -R $user:$user $app

# Download and install any required dependencies
RUN go mod download

# Set the working directory to /go/src/app
WORKDIR /go/src/app/cmd/api

# Build the application
RUN go build -o main .

USER $user

# Expose port 3004 for the application to listen on
EXPOSE $APP_PORT

# Run the application
CMD ["./main"]
