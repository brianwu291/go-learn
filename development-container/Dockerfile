FROM golang:1.23.4-alpine

# Add git and necessary build tools
RUN apk add --no-cache git curl

# Install Air for hot-reloading
RUN curl -sSf https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY ../go.mod ../go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY .. .

# Expose the port your application runs on
EXPOSE 8080

# Use Air for hot-reloading in development
CMD ["air"]
