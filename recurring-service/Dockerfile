# syntax=docker/dockerfile:1

FROM rust:bookworm AS base
WORKDIR /app

RUN rustup toolchain install nightly && rustup default nightly

RUN apt-get update && apt-get install -y --no-install-recommends \
    protobuf-compiler pkg-config libssl-dev ca-certificates \
 && rm -rf /var/lib/apt/lists/*

RUN cargo +nightly install cargo-chef

FROM base AS planner
WORKDIR /app
COPY . .
RUN cargo +nightly chef prepare --recipe-path recipe.json

FROM base AS builder
WORKDIR /app

COPY --from=planner /app/recipe.json recipe.json
RUN cargo +nightly chef cook --release --recipe-path recipe.json

COPY . .
ARG BIN_NAME=recurring-service
RUN cargo +nightly build --release --locked --bin ${BIN_NAME}

FROM debian:bookworm-slim AS runtime
WORKDIR /app
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates tzdata libssl3 \
 && rm -rf /var/lib/apt/lists/*

ARG BIN_NAME=recurring-service
COPY --from=builder /app/target/release/${BIN_NAME} /usr/local/bin/${BIN_NAME}

ENV RUST_LOG=info \
    PORT=50053

EXPOSE 50053
CMD ["recurring-service"]