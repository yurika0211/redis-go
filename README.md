# 📦 go-redis-clone

`go-redis-clone` 是一个用 **Go语言** 实现的轻量级 Redis 克隆项目，旨在学习 Redis 的底层原理与高性能网络服务设计。  
项目严格遵循 Redis 的 **RESP 协议**，实现了常见的五种核心数据结构与命令分发机制。

---

## 🧩 项目结构

```
go-redis-clone/
├── cmd/
│ └── server/
│ └── main.go # 程序入口：启动 Redis 服务
├── internal/
│ ├── core/
│ │ ├── server.go # TCP 服务端逻辑（监听/解析请求）
│ │ ├── client.go # 客户端状态维护（事务/DB/连接）
│ │ ├── command.go # 命令注册与分发
│ │ └── db.go # 内存数据库实现（map[string]Value）
│ ├── data/
│ │ ├── string.go # String 类型实现
│ │ ├── list.go # List 类型实现
│ │ ├── hash.go # Hash 类型实现
│ │ ├── set.go # Set 类型实现
│ │ └── zset.go # SortedSet 类型实现
│ ├── protocol/
│ │ ├── resp_parser.go # RESP 协议解析器
│ │ └── resp_writer.go # RESP 协议响应生成
│ └── utils/
│ ├── logger.go # 日志工具
│ └── errors.go # 错误定义
├── go.mod
└── README.md
```


---

## ⚙️ 功能特性
- 🧠 支持五种核心数据结构（String、List、Hash、Set、SortedSet）  
- 🔌 使用 `net` 库实现 TCP 服务，支持多客户端并发  
- 🗣️ 完全遵循 Redis RESP 协议，可用 Redis-cli 连接  
- 📖 命令注册与动态分发机制  
- 🧰 模块化设计，便于扩展更多命令或持久化机制  

---

## 🚀 启动方式
```bash
cd cmd/server
go run main.go
```

然后使用 redis-cli 连接：

```
redis-cli -p 6379
```

📊 逻辑结构图（ASCII）

```
          ┌────────────────────┐
          │   redis-cli 客户端 │
          └─────────┬──────────┘
                    │  RESP 协议
                    ▼
          ┌────────────────────┐
          │  TCP Server (Go)   │
          ├────────────────────┤
          │ 请求解析  │  响应封装 │
          └──────┬─────────────┘
                 ▼
        ┌───────────────────────┐
        │      Command 执行器   │
        └──────┬───────────────┘
               ▼
        ┌───────────────────────┐
        │    内存数据库 (map)   │
        │ String/List/Set/ZSet  │
        └───────────────────────┘

```

💡 未来计划

 实现 RDB/AOF 持久化

 支持事务（MULTI/EXEC）

 实现订阅发布（Pub/Sub）

 支持异步 IO 模型
