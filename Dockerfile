ARG GO_VERSION=1.26.5
FROM node:24-bookworm-slim AS assets-builder

WORKDIR /usr/src/app

RUN corepack enable && corepack prepare pnpm@11.3.0 --activate
COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
RUN pnpm install --frozen-lockfile

ADD https://github.com/tailwindlabs/tailwindcss/releases/download/v4.3.2/tailwindcss-linux-x64 ./bin/tailwindcli
RUN echo "5036c4fb4328e0bcdbb6065c70d8ac9452e0d4c947113a788a8f94fd390425c1  ./bin/tailwindcli" | sha256sum -c - \
    && chmod +x ./bin/tailwindcli
COPY css ./css
COPY resources ./resources
COPY views ./views
COPY vite.config.ts tsconfig.json components.json ./

RUN ./bin/tailwindcli -i ./css/base.css -o ./assets/css/style.css --minify \
    && mkdir -p ./assets/css/files \
    && cp ./node_modules/@fontsource-variable/noto-sans/files/*-wght-normal.woff2 ./assets/css/files/ \
    && cp ./node_modules/@fontsource-variable/roboto/files/*-wght-normal.woff2 ./assets/css/files/
RUN pnpm build

FROM golang:${GO_VERSION}-bookworm AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
COPY --from=assets-builder /usr/src/app/assets/css ./assets/css
COPY --from=assets-builder /usr/src/app/assets/dist ./assets/dist

RUN CGO_ENABLED=0 GOOS=linux go build -v -o /run-app ./cmd/app

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /run-app /usr/local/bin/run-app

WORKDIR /app

EXPOSE 8080

CMD ["run-app"]
