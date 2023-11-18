FROM golang:1.21 AS builder
ARG VERSION
WORKDIR /src
COPY myapp .
RUN CGO_ENABLED=0 go build -ldflags "-X main.Version=${VERSION}" -o /usr/bin/myapp main.go

FROM scratch AS runner
COPY --from=builder /usr/bin/myapp /usr/bin/myapp
CMD ["/usr/bin/myapp"]
