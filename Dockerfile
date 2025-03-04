# Copyright 2020 Coinbase, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Compile golang
FROM ubuntu:22.04 AS golang-builder

RUN mkdir -p /app \
  && chown -R nobody:nogroup /app
WORKDIR /app

RUN apt-get update && apt-get install -y curl make g++ git
ENV GOLANG_VERSION 1.24.0
ENV GOLANG_DOWNLOAD_SHA256 dea9ca38a0b852a74e81c26134671af7c0fbe65d81b0dc1c5bfe22cf7d4c8858
ENV GOLANG_DOWNLOAD_URL https://golang.org/dl/go$GOLANG_VERSION.linux-amd64.tar.gz

RUN curl -fsSL "$GOLANG_DOWNLOAD_URL" -o golang.tar.gz \
  && echo "$GOLANG_DOWNLOAD_SHA256  golang.tar.gz" | sha256sum -c - \
  && tar -C /usr/local -xzf golang.tar.gz \
  && rm golang.tar.gz

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

# Compile geth
FROM golang-builder as go-quai-builder

# VERSION: go-ethereum v.1.10.16
RUN git clone https://github.com/djadih/go-quai \
  && cd go-quai \
  && git checkout secp-lib

RUN cd go-quai \
  && make go-quai

RUN mv go-quai/build/bin/go-quai /app/go-quai \
  && rm -rf go-quai

# Compile rosetta-ethereum
FROM golang-builder as rosetta-builder

# Use native remote build context to build in any directory
COPY . src
RUN cd src \
  && go build

RUN mv src/rosetta-ethereum /app/rosetta-ethereum \
  && mkdir /app/quai \
  && mv src/quai/call_tracer.js /app/quai/call_tracer.js \
  && mv src/quai/geth.toml /app/quai/geth.toml \
  && rm -rf src

## Build Final Image
FROM ubuntu:22.04

RUN apt-get update && apt-get install -y ca-certificates && update-ca-certificates

RUN mkdir -p /app \
  && chown -R nobody:nogroup /app \
  && mkdir -p /data \
  && chown -R nobody:nogroup /data

WORKDIR /app

# Copy binary from geth-builder
COPY --from=go-quai-builder /app/go-quai /app/go-quai

# Copy binary from rosetta-builder
COPY --from=rosetta-builder /app/quai /app/quai
COPY --from=rosetta-builder /app/rosetta-ethereum /app/rosetta-ethereum

# Set permissions for everything added to /app
RUN chmod -R 755 /app/*

CMD ["/app/rosetta-ethereum", "run"]
