FROM golang:1.17-buster as builder

LABEL maintainer="team3"

WORKDIR /task
COPY . ./
RUN go mod download
RUN mkdir -p application && mkdir -p application/config && cp .env application/config/
RUN cd ./cmd/app/ && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /task/application/task

FROM alpine:3.15.4 as task
COPY --from=builder /task/application /task/application
COPY --from=builder /task/migrations /task/application/migrations
WORKDIR /task/application/
EXPOSE 3000
CMD ["/task/application/task"]
