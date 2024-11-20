FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY main.go .
COPY frenchdiceware.txt .
COPY diceware-fr-alt.txt .
RUN touch go.mod && \
    echo 'module diceware' > go.mod && \
    go mod tidy && \
    go build -o diceware

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/diceware .
COPY --from=builder /app/frenchdiceware.txt .
COPY --from=builder /app/diceware-fr-alt.txt .
ENTRYPOINT ["./diceware"]
