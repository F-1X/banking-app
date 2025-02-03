FROM golang:1.20.14-alpine3.19 AS build
WORKDIR /src
COPY . .
COPY go.mod .
COPY go.sum .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/

FROM alpine:3.19
WORKDIR /app
COPY --from=build /src .
EXPOSE 8080
CMD ["./app"]