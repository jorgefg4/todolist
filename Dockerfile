FROM golang:1.18.3 as development

# Spedify the working directory
WORKDIR /app

# Install dependencies
COPY go.mod  go.sum ./

RUN go mod download

# Copy files to the image
COPY . .

# Port to listen to
EXPOSE 8000

# Start the app
#CMD ["go run", "./cmd/todolist"]
CMD go run ./cmd/todolist