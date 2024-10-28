#!/bin/bash

curl -X POST localhost:8080/message \
     -H "Content-Type: application/json" \
     -d '{"title":"test", "message":"Simple message", "enqueue":false, "time":200}'

curl -X POST localhost:8080/message \
     -H "Content-Type: application/json" \
     -d '{"title":"test", "message":">Quoted message", "enqueue":false}'

curl -X POST localhost:8080/message \
     -H "Content-Type: application/json" \
     -d '{"title":"test", "message":"Message with custom button", "enqueue":false, "button": {"text": "Click to open", "url": "https://google.com"}}'
