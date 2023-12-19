FROM golang:1.20-bullseye

RUN apt update

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

# Change .env with docker
COPY .env.docker ./.env

# Build the Go app
RUN go build -o main .

# Expose port 3002 to the outside world
EXPOSE 3002

# Run the executable
CMD ["./main"]