# go-redis-server Tester

这是一个轻量的集成测试工具，用来对本地运行的 go-redis-server 做端到端基本功能验证（PING、SET/GET、HMSET/HGET、LPUSH/LGET、SADD/SMEMBERS）。

位置
- 可执行入口: `cmd/tester/main.go`

前提
- 机器须安装 Go（建议 1.20+）并能运行项目。
- 目标服务器（go-redis-server）必须在运行并监听一个 TCP 地址（默认 `127.0.0.1:6379`）。

快速开始
1. 在项目根（含 `go.mod` 的模块）下启动服务器：

```bash
# 在项目模块目录下
go run ./cmd/server/main.go &
```

2. 运行测试器：

```bash
go run ./cmd/tester -addr 127.0.0.1:6379
```

3. 期望输出（示例）:

```
PING ok
SET ok
GET ok
HMSET ok
HGET ok
LPUSH reply ok
LGET ok
SADD ok
SMEMBERS ok

Summary: 9 tests, 0 failures
```

参数
- `-addr` : 指定服务器地址（默认 `127.0.0.1:6379`），例如 `-addr 0.0.0.0:16379`。

故障排查
- `connect failed: dial tcp ...: connection refused` ：确认服务器已经启动并监听指定地址（使用 `ss -ltnp | grep :6379` 或 `netstat -ltnp`）。
- 如果测试显示命令失败，请查看服务器日志（如果你使用 `nohup` 重定向，查看 `/tmp/redis-server.log`），并把服务器和测试器的交互输出贴出来以便排查。
- 若出现编译错误，先运行 `go build ./...` 或 `go test ./...` 查看具体报错，再修复代码。

扩展
- 你可以修改 `cmd/tester/main.go` 中的测试用例以覆盖更多命令或更复杂的边界情况。

授权
- 这是一个开发辅助工具，仅用于本地或受控环境的测试，请勿在生产环境直接运行未经审查的测试脚本。
s