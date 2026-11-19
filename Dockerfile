FROM golang:1.27-alpine3.24 AS build

ENV GOFLAGS=-mod=mod
ARG MAXIO_LDFLAGS="-s -w"
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="${MAXIO_LDFLAGS}" -o /maxio .

FROM alpine:3.24
RUN apk upgrade --no-cache && apk add --no-cache ca-certificates && \
	addgroup -S -g 1000 maxio && adduser -S -u 1000 -G maxio -h /data maxio
COPY --from=build /maxio /usr/local/bin/maxio
USER maxio
EXPOSE 9000 9001
VOLUME ["/data"]
ENTRYPOINT ["/usr/local/bin/maxio"]
CMD ["server", "/data", "--console-address", ":9001"]
