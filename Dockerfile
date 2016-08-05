FROM busybox

COPY docker-util /
COPY tests/ /tests

RUN chmod +x /tests/entrypoint.sh

ENTRYPOINT ["/tests/entrypoint.sh"]
