FROM golang:1.24-bullseye AS builder 

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN make build-linux


FROM mcr.microsoft.com/playwright:v1.52.0-jammy AS runtime

WORKDIR /app

COPY --from=builder /app/bin ./bin

ENV TZ=UTC
ENV PLAYWRIGHT_BROWSERS_PATH=/ms-playwright

EXPOSE 8080

CMD ["./bin/unfurl"]