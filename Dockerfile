# Setup Node & Pnpm
FROM node:20-slim AS vite-builder
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN corepack enable
RUN pnpm install -g pnpm

# Build the Vite bundle
WORKDIR /app
COPY . .
WORKDIR /app/web
RUN pnpm install --frozen-lockfile
RUN pnpm run build

# Build the Go app
FROM golang:1.22 AS go-builder
COPY --from=vite-builder /app /app
WORKDIR /app
RUN go mod download
RUN CGO_ENABLED=0 go build -o sentinel ./cmd/sentinel/main.go

# Run the app
EXPOSE 8080
CMD ["./sentinel"]
