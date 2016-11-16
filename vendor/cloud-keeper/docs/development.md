#devel

## auto-gen code

```
docker run -it -v/home/seanchann/pgitlab/cloud/src/cloud-keeper:/go/src/cloud-keeper  --rm gcr.io/google_containers/kube-cross:v1.7.1-2 /bin/bash
```

### deepcopy-gen:

- bounding-dirs: input-dirs中的数据结构所依赖的包列表
- input-dirs:要生成的包的列表
- go-header-file:要生成的文件的文件头定义


example for api with gofreezer:
```
./deep_copy --bounding-dirs="gofreezer/pkg/api" --input-dirs="cloud-keeper/pkg/api" --go-header-file="/go/src/cloud-keeper/contrib/gengo/boilerplate/boilerplate.go.txt" -v=6
```

example for v1:

```
./deep_copy --input-dirs="cloud-keeper/pkg/api/v1" --go-header-file="/go/src/cloud-keeper/contrib/gengo/boilerplate/boilerplate.go.txt" -v=6
```

*现阶段，必须在kube-cross的容器中进行命令操作，由于golang的版本问题*


### conversion-gen:

- extra-peer-dirs: input-dirs中的数据结构所依赖的包列表
- input-dirs:要生成的包的列表
- go-header-file:要生成的文件的文件头定义

example:

```
./conversion_gen  --extra-peer-dirs="cloud-keeper/pkg/api/v1" --input-dirs="cloud-keeper/pkg/api/v1" --go-header-file="/go/src/cloud-keeper/contrib/gengo/boilerplate/boilerplate.go.txt" -v=6
```

*现阶段，必须在kube-cross的容器中进行命令操作，由于golang的版本问题*




### default-gen:

- extra-peer-dirs: input-dirs中的数据结构所依赖的包列表
- input-dirs:要生成的包的列表
- go-header-file:要生成的文件的文件头定义

example:

```
./default_gen --bounding-dirs="gofreezer/pkg/api" --input-dirs="cloud-keeper/pkg/api" --go-header-file="/go/src/cloud-keeper/contrib/gengo/boilerplate/boilerplate.go.txt" -v=6
```

*现阶段，必须在kube-cross的容器中进行命令操作，由于golang的版本问题*
