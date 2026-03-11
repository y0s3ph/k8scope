FROM alpine:3.20

RUN apk add --no-cache ca-certificates
COPY k8scope /usr/local/bin/k8scope

ENTRYPOINT ["k8scope"]
