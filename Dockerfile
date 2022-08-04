FROM golang:1.17-alpine AS builder
WORKDIR /app
COPY . .
ENV GOPROXY=https://goproxy.cn
RUN CGO_ENABLED=1 go build -o /wutong-gateway cmd/gateway.go

FROM wutongpaas/openresty:1.19.3.2-alpine
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && apk add --no-cache bash net-tools curl tzdata && \
        cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
        echo "Asia/Shanghai" >  /etc/timezone && \
        date && apk del --no-cache tzdata
COPY ./hack/ /run/
COPY --from=builder /wutong-gateway /wutong-gateway

ENV NGINX_CONFIG_TMPL=/run/nginxtmp
ENV NGINX_CUSTOM_CONFIG=/run/nginx/conf
ENV OPENRESTY_HOME=/usr/local/openresty
ENV PATH="${PATH}:${OPENRESTY_HOME}/nginx/sbin"
EXPOSE 8080

ENTRYPOINT ["/wutong-gateway"]
