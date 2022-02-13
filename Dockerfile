FROM golang:1.17-alpine as builder
RUN mkdir /app
WORKDIR /app
COPY . .
RUN go build

FROM alpine
COPY --from=builder /app/workflow-job-handler /usr/local/bin/workflow-job-handler
ENTRYPOINT ["/usr/local/bin/workflow-job-handler"]
