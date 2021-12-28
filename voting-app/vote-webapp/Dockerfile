FROM golang:1.17-alpine as builder
WORKDIR /opt/build
COPY go.* ./
RUN go mod download
COPY . ./
RUN go build -v -o ./vote-app

FROM node:16.5-alpine3.14 as node-builder
WORKDIR /opt/build
COPY vote-app-frontend ./vote-app-frontend
WORKDIR /opt/build/vote-app-frontend
RUN yarn install
RUN yarn build

FROM alpine:3.14
WORKDIR /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /opt/build/vote-app /
COPY --from=node-builder /opt/build/vote-app-frontend/dist /static
CMD ["/vote-app"]