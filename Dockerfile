# https://www.bytesizego.com/blog/production-go-docker-image
FROM golang:1.23.2-bookworm as builder

COPY . /workdir
WORKDIR /workdir

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

ENV CGO_CPPFLAGS="-D_FORTIFY_SOURCE=2 -fstack-protector-all"
ENV GOFLAGS="-buildmode=pie"

RUN go build -ldflags "-s -w" -trimpath -o /workdir/echo ./api/*.go

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends redis-server libatomic1 && \
  apt-get clean && rm -rf /var/lib/apt/lists/*

COPY --from=builder /workdir/echo /bin/echo

USER 65534

# service ports
EXPOSE 8080

CMD ["sh", "-c", "/usr/bin/redis-server --daemonize yes && /bin/echo"]
