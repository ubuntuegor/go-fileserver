# syntax=docker/dockerfile:1

# build stage
FROM golang:1.21 as builder

WORKDIR /app

COPY main.go ./
RUN go build -o go-fileserver main.go

# run stage
FROM gcr.io/distroless/base-debian12

COPY --from=builder /app/go-fileserver /go-fileserver

WORKDIR /

EXPOSE 8080

ENV PORT=8080
ENV SERVE_FROM_FOLDER=content
ENTRYPOINT [ "/go-fileserver" ]
