FROM rust:1.67 AS builder

WORKDIR /app

COPY . .

RUN cargo build --release --bin mortenvistisen_blog

FROM debian:buster-slim

WORKDIR /app

RUN apt-get update -y \
    && apt-get install -y --no-install-recommends openssl \
    && apt-get install ca-certificates \
    # Clean up
    && apt-get autoremove -y \
    && apt-get clean -y \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/target/release/mortenvistisen_blog mortenvistisen_blog
COPY --from=builder /app/templates templates
COPY --from=builder /app/static static
COPY --from=builder /app/posts posts

ENTRYPOINT ["./mortenvistisen_blog"]
EXPOSE 8080
