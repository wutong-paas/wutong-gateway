#!/bin/bash

docker build . -t wutongpaas/wutong-gateway:v0.1.0-amd64
docker push wutongpaas/wutong-gateway:v0.1.0-amd64