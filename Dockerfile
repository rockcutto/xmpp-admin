FROM golang:1.27-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOFLAGS="-mod=readonly" go test ./...     && CGO_ENABLED=0 GOFLAGS="-mod=readonly" go build -buildvcs=false -trimpath -ldflags="-s -w" -o /out/xmpp-admin .

FROM alpine:3.24

RUN apk add --no-cache ca-certificates tzdata     && adduser -D -H -u 10001 xmpp-admin

WORKDIR /app
COPY --from=build /out/xmpp-admin /usr/local/bin/xmpp-admin

USER xmpp-admin
ENV ADMIN_LISTEN=0.0.0.0:8090
ENV EJABBERD_CONFIG=/etc/ejabberd/ejabberd.yml

EXPOSE 8090
ENTRYPOINT ["/usr/local/bin/xmpp-admin"]
