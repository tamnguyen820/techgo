FROM golang:1.22.2-alpine AS build
WORKDIR /app
COPY . .
RUN go mod download
WORKDIR /app/cmd/techgo
RUN go build -o /bin/techgo .

FROM alpine:3
COPY --from=build /bin/techgo /app/techgo
COPY ./config.yml /app/config.yml
WORKDIR /app
CMD ["./techgo"]
