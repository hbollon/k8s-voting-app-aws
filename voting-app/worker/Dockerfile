FROM golang:1.17-alpine as builder
WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./
RUN go build -v -o /voting-app-worker

FROM alpine:3.14
COPY --from=builder /voting-app-worker /voting-app-worker
EXPOSE 8080
CMD [ "/voting-app-worker" ]