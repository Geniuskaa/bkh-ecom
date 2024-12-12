FROM golang:1.23 AS build

WORKDIR /app
COPY ./go.mod ./go.sum ./

RUN go mod download && go mod tidy

COPY . ./
RUN CGO_ENABLED=0 go build -o ./build/main ./cmd

FROM alpine

WORKDIR /app
COPY --from=build /app/build/main ./main
RUN mkdir /configs

EXPOSE 9900
CMD ["./main"]

