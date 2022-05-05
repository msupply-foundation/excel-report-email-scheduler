ARG GRAFANA_VERSION="8.4.4"
FROM grafana/grafana:${GRAFANA_VERSION}

MAINTAINER Sworup Shakya <sworup@susol.net>

USER root

RUN apk --no-cache update && \
    apk add --no-cache curl gettext-base

COPY entrypoint.sh .
RUN chmod +x entrypoint.sh

# Todo: Just for debugging purposes, can be removed
# Make sure that the init flag is not there
RUN rm -rf /var/lib/grafana/.init

ENTRYPOINT ["./entrypoint.sh"]
