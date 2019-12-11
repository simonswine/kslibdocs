FROM golang:1.13.5 as golang

WORKDIR /go/src/github.com/ksonnet/kslibdocs

ADD go.mod go.sum ./
RUN go mod download

ADD cmd ./cmd
ADD pkg ./pkg

RUN go build ./cmd/kslibdocgen

FROM node:8.16.2 as node

COPY --from=golang /go/src/github.com/ksonnet/kslibdocs/kslibdocgen .

RUN curl -LO https://github.com/ksonnet/ksonnet-lib/raw/master/ksonnet.beta.4/k8s.libsonnet

RUN ./kslibdocgen --k8sLib=k8s.libsonnet --outDir /_output && \
    find /_output -type d -exec chmod 0755 {} \; && \
    find /_output -type f -exec chmod 0644 {} \;

FROM nginx:1.17.6
COPY --from=node --chown=33:33 /_output /usr/share/nginx/html
