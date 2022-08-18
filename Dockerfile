FROM golang:1.17-buster as builder

LABEL maintainer="team3"

WORKDIR /team3-task
COPY . ./
RUN go mod download
RUN mkdir -p application && mkdir -p application/config && cp .env application/config/
RUN cd ./cmd/app/ && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /team3-task/application/task

FROM alpine:3.15.4 as task
COPY --from=builder /team3-task/application /team3-task/application
COPY --from=builder /team3-task/migrations /team3-task/application/migrations
WORKDIR /team3-task/application/
EXPOSE 3000
CMD ["/team3-task/application/task"]
