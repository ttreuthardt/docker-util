FROM busybox

COPY docker-util /
COPY tests/ /tests

RUN chown -R 10000 /tests

USER 10000

ENTRYPOINT ["/docker-util", "-config", "/tests/config.json"]
