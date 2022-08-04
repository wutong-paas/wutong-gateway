#!/bin/bash

docker buildx build --push --platform=linux/amd64,linux/arm64 . -t wutongpaas/wutong-gateway:v0.1.0