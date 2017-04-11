FROM alpine

ARG version

LABEL maintainer="hbm@kassisol.com"
LABEL version=$version
LABEL description="An application acting as a CA (Certificate Authority)."

COPY build/tsa /usr/local/sbin/tsa
COPY entrypoint.sh /entrypoint.sh

RUN apk --update --no-cache add bash

EXPOSE 443

ENTRYPOINT ["/entrypoint.sh"]
