FROM golang:1.16-alpine as build
RUN mkdir /app
WORKDIR /app
COPY ./ ./
RUN go mod download
RUN cd ./cmd/beasty && go build

FROM golang:1.16-alpine as runtime
RUN mkdir /app
WORKDIR /app
COPY --from=build /app/cmd/beasty/beasty ./
# TODO: pack the statics as binary
RUN mkdir /app/client
COPY --from=build /app/cmd/beasty/client/config.yaml ./client/
RUN mkdir /app/responses
COPY --from=build /app/cmd/beasty/responses/responses.yaml ./responses/

CMD ["./beasty"]