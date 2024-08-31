## How to run

### build docker image for unbound in which cache is enabled

```shell
docker build -t unbound:with-cache -f Dockerfile.unbound-with-cache .
```

### build docker image for unbound in which cache is disabled

```shell
docker build -t unbound:without-cache -f Dockerfile.unbound-without-cache .
```

### run docker image for unbound in which cache is enabled
```shell
docker run --rm -it -d -p 53:53/tcp -p 53:53/udp unbound:with-cache
```

### run docker image for unbound in which cache is disabled
```shell
docker run --rm -it -d -p 53:53/tcp -p 53:53/udp unbound:without-cache
```

## Reference

- <https://github.com/obi12341/docker-unbound/blob/master/README.md>
