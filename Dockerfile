FROM alpine:3.14

ENV LANG=C.UTF-8

RUN apk upgrade --no-cache && apk add --no-cache restic dumb-init

COPY entrypoint.sh /app/
RUN chmod +x /app/entrypoint.sh

ENTRYPOINT ["/usr/bin/dumb-init", "--", "/app/entrypoint.sh"]
CMD ["run", "help"]
