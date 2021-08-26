
## 环境准备

### MySQL

```shell
$ docker pull mysql:5.6
$ docker run -p 3306:3306 --name mysql -e MYSQL_ROOT_PASSWORD=123456 -d mysql:5.6
```

### etcd 

```shell
$ docker pull bitnami/etcd:latest
$ docker run -d --name etcd-server \
    --publish 2379:2379 \
    --publish 2380:2380 \
    --env ALLOW_NONE_AUTHENTICATION=yes \
    --env ETCD_ADVERTISE_CLIENT_URLS=http://127.0.0.1:2379 \
    bitnami/etcd:latest
```

### 链路追踪

```shell
$ docker pull jaegertracing/all-in-one:1.16
$ docker run -d --name jaeger \
-e COLLECTOR_ZIPKIN_HTTP_PORT=9411 \
-p 5775:5775/udp \
-p 6831:6831/udp \
-p 6832:6832/udp \
-p 5778:5778 \
-p 16686:16686 \
-p 14268:14268 \
-p 9411:9411 \
jaegertracing/all-in-one:1.16
```

[Jaeger UI: 链路追踪](http://127.0.0.1:16686/search)


