FROM scratch

ARG version

LABEL maintainer="hbm@kassisol.com"
LABEL version=$version
LABEL description="An application acting as a CA (Certification Authority)."

COPY build/tsad /tsad

EXPOSE 80
EXPOSE 443

ENTRYPOINT ["/tsad"]
