# path: Dockerfile
FROM golang:1.22 AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
# Proto code should already be generated on host (make proto)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/rt-price-engine ./cmd/aggregator

FROM gcr.io/distroless/static:nonroot
WORKDIR /app
COPY --from=build /out/rt-price-engine /app/rt-price-engine
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/app/rt-price-engine"]
