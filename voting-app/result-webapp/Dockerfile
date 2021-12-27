FROM golang:1.17-alpine as builder
WORKDIR /opt/build
COPY go.* ./
RUN go mod download
COPY . ./
RUN go build -v -o ./result-app

FROM alpine:3.14
WORKDIR /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /opt/build/result-app /
CMD ["/result-app"]