#!/bin/bash

docker build . -t swr.cn-southwest-2.myhuaweicloud.com/wutong/wutong-gateway:v0.1.0-amd64
docker push swr.cn-southwest-2.myhuaweicloud.com/wutong/wutong-gateway:v0.1.0-amd64

docker tag swr.cn-southwest-2.myhuaweicloud.com/wutong/wutong-gateway:v0.1.0-amd64 wutongpaas/wutong-gateway:v0.1.0-amd64
docker push wutongpaas/wutong-gateway:v0.1.0-amd64