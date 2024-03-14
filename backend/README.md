基于GO重写的Server端。

项目目录结构参考自：

- https://github.com/go-eagle/eagle.git

### generate swagger documents

```shell
swag init -g ./cmd/main.go
```

### generate dependency injection for fusion

```shell
cd ./internal/app && wire && cd ../../
```
