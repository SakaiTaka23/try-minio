FROM golang:1.17-alpine3.13 as dev

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download
RUN go get -u github.com/cosmtrek/air
COPY . .

EXPOSE 5000

CMD ["air"]
