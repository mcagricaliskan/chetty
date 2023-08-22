FROM golang:1.21-alpine3.18 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/the-game-backend .

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /bin/the-game-backend /bin/the-game-backend

EXPOSE 33333

# Run
CMD ["/bin/the-game-backend"]