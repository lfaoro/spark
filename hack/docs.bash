#!/usr/bin/env bash

DOCS=$(ls proto/api/vault/*.swagger.json | xargs -n 1 basename)
for doc in $DOCS; do
  jq -s '.[0] * .[1]' proto/api/vault/$doc web/doc/info.swagger.json >web/doc/$doc
done

swagger mixin web/doc/*.swagger.* > web/doc/spark_swagger.json
