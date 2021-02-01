FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . /app

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

EXPOSE 4560

#Command to run the executable
ENTRYPOINT ["./app"]
