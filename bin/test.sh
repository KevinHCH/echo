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

curl -X POST localhost:8080/message \
     -H "Content-Type: application/json" \
     -d '{"title":"test", "message":"*this is working*", "enqueue":false}'

curl -X POST localhost:8080/message \
     -H "Content-Type: application/json" \
     -d '{"title":"test", "message":"*Title:* Experienced Freelancer Needed for n8n Installation and Configuration on AWS or Google Cloud\n\n*Posted At:* Posted 1 hour ago\n\n*Price:* no\n\n*Job Type:* Hourly 2000  5000\n\n*Duration:* Est time 1 to 3 months Hours to be determined\n\n*Experience Level:* Intermediate\n\n", "enqueue":false}'
