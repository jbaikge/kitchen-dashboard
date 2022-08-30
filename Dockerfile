FROM golang:1.19

WORKDIR /usr/local/src/kitchen-dashboard

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/kitchen-dashboard .

CMD ["kitchen-dashboard"]
