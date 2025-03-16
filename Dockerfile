# Setup Node & Pnpm
FROM node:20-slim AS vite-builder
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable
RUN pnpm install -g pnpm

# Build the Vite bundle
COPY ./web /app/web
WORKDIR /app/web
RUN rm -rf node_modules
RUN pnpm install --frozen-lockfile
#RUN pnpm run build

# Build the Go app #
FROM golang:1.23.1 AS go-builder
COPY ./internal /app/internal
COPY ./cmd /app/cmd
COPY ./common /app/common
COPY ./go.mod /app/go.mod
COPY ./go.sum /app/go.sum
COPY --from=vite-builder /app/web/dist /app/internal/server/dist
WORKDIR /app
RUN go mod download
RUN go build -o forkman ./cmd/forkman/main.go

# Setup the final image
FROM ubuntu:22.04 AS final
WORKDIR /app

# Install ca-certificates
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

COPY --from=go-builder /app/forkman /app/forkman
RUN chmod +x /app/forkman

# Run the app
EXPOSE 8080
CMD ["./forkman"]
