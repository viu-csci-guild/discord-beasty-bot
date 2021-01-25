FROM golang:1.15
# TODO: remove root container access
WORKDIR /app
COPY . .

RUN go get -d ./...
RUN go install ./...
RUN go build