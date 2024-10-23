FROM golang:1.23.2-bookworm as builder

COPY . /workdir
WORKDIR /workdir

ENV CGO_CPPFLAGS="-D_FORTIFY_SOURCE=2 -fstack-protector-all"
ENV GOFLAGS="-buildmode=pie"

RUN go build -ldflags "-s -w" -trimpath .

FROM gcr.io/distroless/base-debian12:nonroot
COPY --from=builder /workdir/echo /bin/echo

USER 65534

ENTRYPOINT ["/bin/echo"]