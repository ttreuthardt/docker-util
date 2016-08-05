#!/bin/sh

/docker-util -config /tests/config.json
rc=$?

ls -Rla /tests

exit $rc