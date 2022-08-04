# Wutong Gateway
欢迎使用 Wutong Gateway！
Wutong Gateway 是一款基于 Openresty 的 Ingress Controller。
实现通过 http，tcp 方式将流量转发到位于 Kubernetes 的 Service 中。
## 安装
```bash
kubectl apply -f https://raw.githubusercontent.com/wutong-paas/wutong-gateway/master/deploy/manifests.yaml
```
## 使用示例
- 如果集群中的 Ingress APIVersion 是 `networking.k8s.io/v1`，将 Ingress 的模型字段 `ingress.spec.ingressClassName` 值设置为 `wutong-gateway` 生效;
- 如果集群中的 Ingress APIVersion 是 `extensions/v1beta1`，将 Ingress Annotation `kubernetes.io/ingress.class` 值设置为 `wutong-gateway` 生效。

**创建 Deployment & Service**
```bash
kubectl create deployment nginx --image=nginx
kubectl expose deployment nginx --type=ClusterIP --port=80 --target-port=80
```
**创建 Ingress** 

*nginx-tcp.yaml*
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: nginx-tcp
  annotations:
    nginx.ingress.kubernetes.io/l4-enable: "true"
    nginx.ingress.kubernetes.io/l4-host: 0.0.0.0
    nginx.ingress.kubernetes.io/l4-port: "31080"
spec:
  ingressClassName: wutong-gateway
  defaultBackend:
    service:
      name: nginx
      port:
        number: 80
```
*nginx-http.yaml*
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: nginx-http
spec:
  ingressClassName: wutong-gateway
  rules:
    - host: nginx.wutong-gateway-sample.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: nginx
                port:
                  number: 80
```
```bash
kubectl apply -f nginx-tcp.yaml
kubectl apply -f nginx-http.yaml
```
> wutong-gateway 作为 DaemonSet 部署在集群中的 Master 节点中，此时可以进入节点终端进行验证：
```bash
# tcp
curl localhost:31080

# http
curl localhost:80 -H 'Host: nginx.wutong-gateway-sample.com'
```