FROM golang:1.24-bullseye AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Install make for build process
RUN apt-get update && apt-get install -y make && rm -rf /var/lib/apt/lists/*

# Install Playwright and browsers properly
RUN go install github.com/playwright-community/playwright-go/cmd/playwright@latest
RUN playwright install chromium
RUN playwright install-deps

RUN make build-linux

FROM debian:bullseye-slim AS runtime

RUN apt-get update && apt-get install -y \
    ca-certificates \
    fonts-liberation \
    libasound2 \
    libatk-bridge2.0-0 \
    libdrm2 \
    libgtk-3-0 \
    libnspr4 \
    libnss3 \
    libxss1 \
    libgbm1 \
    libxcomposite1 \
    libxdamage1 \
    libxrandr2 \
    libxfixes3 \
    libxcursor1 \
    libxi6 \
    libxtst6 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy binary
COPY --from=builder /app/bin ./bin

# Copy Playwright browsers and driver
COPY --from=builder /root/.cache/ms-playwright /root/.cache/ms-playwright
COPY --from=builder /root/.cache/ms-playwright-go /root/.cache/ms-playwright-go

ENV TZ=UTC
ENV PLAYWRIGHT_BROWSERS_PATH=/root/.cache/ms-playwright

EXPOSE 8080

CMD ["./bin/unfurl"]