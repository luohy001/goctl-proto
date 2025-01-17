# goctl-proto

## 简介

通过api文件生成proto文件。

## 安装

```go
go install github.com/liferod/goctl-proto/cmd/goctl-proto@latest
```

## 示例

当前目录下有一个api文件service.api，内容如下：

```
syntax = "v1"

type Mock {
	Id         int64                 `json:"id,optional"` // ID
	Name       string                `json:"name,optional"` // 名称
	Type       int32                 `json:"type,optional"` // 类型
	Tags       []string              `json:"tags,optional"` // 标签列表
	ExtraInfos map[string]*ExtraInfo `json:"extra_infos,optional"` // 额外信息
}

type ExtraInfo {
	Content string `json:"content,optional"` // 内容
}

service mocker-api {
	@doc "创建"
	@handler CreateMock
	post /create (Mock) returns (Mock)
}
```

使用goctl-proto生成proto文件的命令：

```
goctl-proto proto --input ./service.api --output .
```

或者配合goctl插件命令使用：

```
goctl api plugin -plugin goctl-proto="proto" -api ./service.api -dir .
```

在当前目录下就得到了一个service.proto文件，内容如下：

```
syntax = "proto3";

package mocker.api;
option go_package = "/protoc-gen-go";

message ExtraInfo {
    // 内容
    string Content = 1 [json_name = "content"];
}

message Mock {
    // ID
    int64 Id = 1 [json_name = "id"];
    // 名称
    string Name = 2 [json_name = "name"];
    // 类型
    int32 Type = 3 [json_name = "type"];
    // 标签列表
    repeated string Tags = 4 [json_name = "tags"];
    // 额外信息
    map<string,ExtraInfo> ExtraInfos = 5 [json_name = "extra_infos"];
}

service MockerApi {
    // 创建
    rpc CreateMock (Mock) returns (Mock);
}
```

使用该proto文件即可调用goctl的rpc命令生成rpc代码：

```
goctl rpc protoc ./service.proto --go_out=./service --go-grpc_out=./service --zrpc_out=./service --client=true
```

更多示例[点击此处](https://github.com/liferod/goctl-proto/tree/main/example)。

## 存在的问题

1. 支持生成的protobuf只支持proto3。

2. 一些在protobuf中不存在/不支持的数据类型会转换为bytes，例如interface{}、any等等。

3. 不支持按照api中的import的文件分别生成proto文件，并且不会生成类似google.protobuf.Any和google.protobuf.Empty这种类型的数据类型，但不排除以后支持，原因见[这里](https://go-zero.dev/docs/tutorials/proto/faq#2-%E4%B8%BA%E4%BB%80%E4%B9%88%E4%BD%BF%E7%94%A8-goctl-%E7%94%9F%E6%88%90-grpc-%E4%BB%A3%E7%A0%81%E6%97%B6-proto-%E4%B8%8D%E6%94%AF%E6%8C%81%E4%BD%BF%E7%94%A8%E5%8C%85%E5%A4%96-proto-%E5%92%8C-service)。

4. 在api文件中声明的type类型的备注信息不会被复制到proto文件中，原因是目前使用的goctl v1.6.5并没有支持解析这类备注，但不排除以后支持。

5. 所有protobuf的map键不支持的类型会转化为string，不支持的值会转化为bytes，例如map[float32]\[]string会转为map\<string,bytes>

## TODO

- 支持对proto中的service分组 ✅
- 支持导出字段除json外的其余tag
- 导出rpc时更深度的清理未被使用的message ✅
- 支持导出描述信息中包含@goctl-proto的rpc ✅
