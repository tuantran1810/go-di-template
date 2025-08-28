FROM golang:1.24-bookworm AS builder

WORKDIR /app
COPY . .
RUN make build
EXPOSE 8080 9090
ENTRYPOINT [ "./go-di-template", "start-server" ]
