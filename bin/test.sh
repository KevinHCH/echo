#!/bin/bash

curl -X POST localhost:8080/message \
     -H "Content-Type: application/json" \
     -d '{"title":"test", "message":"<a href='www.google.com'>hello from golang service</a>", "enqueue":true, "time":200}'
