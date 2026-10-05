# 新蜂资产管理平台 — 服务开发模板

为平台新增服务提供 go-zero API、RPC、Ent 及权限集成的代码骨架。模板包含未替换的占位符，必须完成服务命名、协议和配置调整后才能编译，不能直接部署。

仓库：[coder-lulu/newbee-service-template](https://github.com/coder-lulu/newbee-service-template) · [平台工作区](https://github.com/coder-lulu/newbee)

## 获取代码

推荐通过完整工作区开发，保留兄弟模块目录及本地 `replace` 依赖。以下命令使用 Bash；Go 工作区要求 Go 1.25.1 或更高版本。

```bash
git clone --recurse-submodules https://github.com/coder-lulu/newbee.git
cd newbee/templates/service-template
```

已有工作区执行 `git submodule update --init --recursive`。单独克隆模块时，需要自行补齐 `go.mod` 中的本地依赖路径。

## 目录导航

| 路径 | 用途 |
| --- | --- |
| `api/` | API 入口、接口定义和中间件骨架 |
| `rpc/` | RPC 入口、协议和 Ent 骨架 |
| `api/etc/io.yaml.example` | API 配置模板 |
| `rpc/etc/io.yaml.example` | RPC 配置模板 |
| `api/Makefile`、`rpc/Makefile` | 构建与代码生成目标 |

## 创建新服务

1. 将模板复制到工作区中的新服务目录，保留 API/RPC 分层和版权声明。
2. 替换 `{{SERVICE_NAME}}`、`{{SERVICE_NAME_UPPER}}`、`{{API_PORT}}` 等全部占位符，检查模块路径、Go import、协议包名及生成目标。
3. 将模板中的 `io.go`、`io.proto`、`types/io`、`ioclient` 等名称调整为新服务名称；同步 Makefile、入口及生成脚本，不只修改 go.mod。
4. 根据新服务接口和 Ent Schema 重新生成代码，校验公共租户与权限中间件的接入。
5. 将配置模板整理为新服务的 `etc/<服务名>.yaml.example`，为本地运行创建不含未替换占位符的配置；设置数据库、Redis、RPC 地址和独立监听端口。
6. 在平台根目录使用 `go work use ./<新服务>/api ./<新服务>/rpc` 加入工作区，并核对本地 `replace` 路径。

模板顶层没有 `go.mod`。实例化后，在新服务的 API 和 RPC 目录分别执行：

```bash
go build .
go test ./...
go vet ./...
```

运行入口沿用 `-f` 参数，例如在生成后的模块中执行 `go run . -f etc/<服务名>.yaml`。尖括号表示要替换的名称，不能原样执行。生产配置、密钥和生成的二进制不要加入仓库。

## 参考实现

- [核心服务](https://github.com/coder-lulu/newbee-core)
- [统一 I/O API](https://github.com/coder-lulu/newbee-io-api)
- [统一 I/O RPC](https://github.com/coder-lulu/newbee-io-rpc)
- [公共库](https://github.com/coder-lulu/newbee-common)

## 许可证与来源

本仓库采用 [Apache-2.0](LICENSE)。沿用现有服务模板的上游许可和文件版权声明。第三方依赖遵循各自许可证，保留原有版权与许可声明。
