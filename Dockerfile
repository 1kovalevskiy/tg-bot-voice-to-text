# Step 1: Modules caching
FROM golang:1.22-alpine as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Step 2: Builder
FROM golang:1.22-alpine as builder
COPY --from=modules /go/pkg /go/pkg
WORKDIR /app
COPY ./cmd /app/cmd
COPY ./config /app/config
COPY ./pkg /app/pkg
COPY go.mod go.sum /app/
COPY ./internal /app/internal
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /bin/app ./cmd/app

# Step 3: Final
FROM scratch
COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs/
COPY --from=builder /bin/app /app
WORKDIR /
