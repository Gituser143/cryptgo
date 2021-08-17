FROM golang:1.16 as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux \
    go build --ldflags "-s -w" -a -o ./output/cryptgo ./cryptgo.go

FROM alpine:latest
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser
COPY --from=builder /app/output/cryptgo /app/cryptgo

ENTRYPOINT ["/app/cryptgo"]
