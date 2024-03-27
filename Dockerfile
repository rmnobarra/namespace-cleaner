FROM golang:1.20 AS build

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o namespace-cleaner main.go

FROM alpine:3.13

RUN apk --no-cache add ca-certificates

COPY --from=build /app/namespace-cleaner /namespace-cleaner

CMD ["/namespace-cleaner"]
