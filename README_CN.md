# Knights - 基于Actor模型的游戏服务器框架

Knights 是一个用 Go 语言编写的高性能游戏服务器框架，基于 Actor 模型构建。它为构建多人在线游戏提供了一个可扩展、容错的基础架构。

## 特性

- **Actor 模型**: 轻量级、并发的 Actor 系统，基于邮箱的消息传递机制
- **高性能**: 异步 I/O、高效的消息调度和最小的开销
- **容错性**: 监督策略提供自动错误恢复
- **可扩展**: 通过集群支持水平扩展（规划中）
- **灵活的网络**: 支持 TCP、WebSocket 和多种编解码格式
- **游戏就绪**: 内置玩家管理、房间系统和匹配系统

## 架构

```
knights/
├── cmd/                     # 应用程序入口
│   ├── server/             # 游戏服务器
│   └── tools/              # 工具和实用程序
├── internal/               # 私有应用代码
│   ├── actor/              # Actor 框架核心
│   ├── game/               # 游戏逻辑（玩家、房间、匹配）
│   ├── network/            # 网络层（TCP、WebSocket）
│   ├── cluster/            # 分布式集群支持
│   ├── config/             # 配置管理
│   └── ...
├── pkg/                    # 公共库
│   ├── proto/              # 协议定义
│   ├── errors/             # 错误码
│   └── constants/          # 常量
└── api/                    # 外部 API（HTTP、gRPC）
```

## 快速开始

### 前置要求

- Go 1.21 或更高版本
- Make（可选，用于使用 Makefile 命令）

### 安装

```bash
# 克隆仓库
git clone https://github.com/widrin/knights.git
cd knights

# 安装依赖
go mod download

# 构建服务器
make build
# 或者
go build -o bin/server cmd/server/main.go
```

### 运行服务器

```bash
# 直接运行
make run

# 或运行二进制文件
./bin/server
```

### 配置

编辑 `configs/server.yaml` 来自定义服务器设置：

```yaml
server:
  name: "knights-server"
  address: "0.0.0.0"
  port: 8080

game:
  max_players: 10000
  tick_rate: 20
  room_max_players: 4
```

## Actor 系统

Knights 的核心是其 Actor 系统，它提供：

### 创建 Actor

```go
import "github.com/widrin/knights/internal/actor"

type MyActor struct {}

func (a *MyActor) Receive(ctx actor.Context) {
    switch msg := ctx.Message().(type) {
    case *actor.Started:
        // Actor 启动
    case string:
        // 处理字符串消息
        ctx.Respond("你好: " + msg)
    }
}

// 创建 actor 系统
system := actor.NewActorSystem("game")

// 生成 actor
props := actor.NewProps(func() actor.Actor {
    return &MyActor{}
})
pid := system.Spawn(props)

// 发送消息
system.Send(pid, "Hello")
```

### 核心概念

- **Actor**: 处理消息的轻量级并发实体
- **PID**: 唯一标识 actor 的进程 ID
- **Mailbox**: 每个 actor 的消息队列
- **Supervisor**: 管理子 actor 并处理故障
- **Context**: 提供消息处理和 actor 生命周期方法

## 游戏组件

### 玩家管理

```go
// Player actor 处理单个玩家逻辑
playerProps := actor.NewProps(func() actor.Actor {
    return player.NewPlayerActor(playerID)
})
playerPID := system.Spawn(playerProps)
```

### 房间系统

```go
// Room actor 管理游戏房间/战斗
roomProps := actor.NewProps(func() actor.Actor {
    return room.NewRoomActor(roomID, maxPlayers)
})
roomPID := system.Spawn(roomProps)
```

### 匹配系统

```go
// Matchmaker actor 处理玩家匹配
matchmakerProps := actor.NewProps(func() actor.Actor {
    return match.NewMatchmakerActor()
})
matchmakerPID := system.Spawn(matchmakerProps)
```

## 网络层

Knights 支持多种网络协议和编解码器：

```go
import (
    "github.com/widrin/knights/internal/network"
    "github.com/widrin/knights/internal/network/codec"
)

// 使用 JSON 编解码器创建服务器
server := network.NewServer(&network.ServerConfig{
    Address:     ":8080",
    Codec:       codec.NewJSONCodec(),
    Handler:     handler,
    ActorSystem: system,
})

server.Start()
```

## 开发

### 运行测试

```bash
make test
```

### 代码格式化

```bash
make fmt
```

### 生成 Protobuf 代码

```bash
make proto
```

## 项目状态

🚧 **正在积极开发中**

已实现功能：
- ✅ 核心 Actor 系统
- ✅ 基础游戏组件（玩家、房间、匹配）
- ✅ 支持 TCP 的网络层
- ✅ 配置管理
- ✅ 日志和指标

计划功能：
- ⏳ 分布式部署的集群支持
- ⏳ WebSocket 支持
- ⏳ 持久化层
- ⏳ 完整的匹配算法
- ⏳ 管理仪表板
- ⏳ 性能基准测试

## 详细文档

### 核心文档
- [架构设计](docs/architecture.md) - 系统架构详解
- [快速入门](docs/getting_started.md) - 详细的入门指南

### 代码示例

#### 创建自定义 Actor

```go
package game

import "github.com/widrin/knights/internal/actor"

// 定义你的 Actor
type GameActor struct {
    score int
}

func NewGameActor() actor.Actor {
    return &GameActor{score: 0}
}

// 实现消息处理
func (g *GameActor) Receive(ctx actor.Context) {
    switch msg := ctx.Message().(type) {
    case *actor.Started:
        // Actor 启动时的初始化
        g.score = 0

    case *ScoreMessage:
        // 处理得分消息
        g.score += msg.Points
        ctx.Respond(&ScoreResponse{
            TotalScore: g.score,
        })

    case *actor.Stopping:
        // Actor 停止前的清理
        g.cleanup()
    }
}
```

#### 使用监督策略

```go
// 创建带有监督策略的 Actor
props := actor.NewProps(func() actor.Actor {
    return NewGameActor()
}).WithSupervisor(
    actor.NewOneForOneStrategy(3, time.Minute), // 1分钟内最多重启3次
)

pid := system.Spawn(props)
```

#### 使用路由器进行负载均衡

```go
// 创建工作池
workers := make([]*actor.PID, 10)
for i := 0; i < 10; i++ {
    workers[i] = system.Spawn(workerProps)
}

// 创建轮询路由器
router := actor.NewRoundRobinRouter(workers)

// 消息会被均匀分配到所有工作者
for _, msg := range messages {
    workerPID := router.Route(msg)
    system.Send(workerPID, msg)
}
```

#### 构建完整的游戏服务器

```go
package main

import (
    "log"
    "os"
    "os/signal"
    "syscall"

    "github.com/widrin/knights/internal/actor"
    "github.com/widrin/knights/internal/config"
    "github.com/widrin/knights/internal/game/player"
    "github.com/widrin/knights/internal/network"
    "github.com/widrin/knights/internal/network/codec"
)

func main() {
    // 加载配置
    cfg, err := config.LoadConfig("configs/server.yaml")
    if err != nil {
        log.Fatal("加载配置失败:", err)
    }

    // 创建 Actor 系统
    system := actor.NewActorSystem(cfg.Server.Name)
    defer system.Shutdown()

    // 创建玩家管理器
    playerMgrProps := actor.NewProps(func() actor.Actor {
        return player.NewPlayerManager()
    })
    playerMgr := system.Spawn(playerMgrProps)

    // 创建消息处理器
    handler := NewGameHandler(system, playerMgr)

    // 创建网络服务器
    server := network.NewServer(&network.ServerConfig{
        Address:     fmt.Sprintf(":%d", cfg.Server.Port),
        Codec:       codec.NewJSONCodec(),
        Handler:     handler,
        ActorSystem: system,
    })

    // 启动服务器
    if err := server.Start(); err != nil {
        log.Fatal("启动服务器失败:", err)
    }

    log.Printf("服务器运行在 %s:%d\n", cfg.Server.Address, cfg.Server.Port)

    // 等待关闭信号
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    log.Println("服务器正在关闭...")
    server.Stop()
}
```

## Makefile 命令

```bash
make build      # 构建服务器
make run        # 运行服务器
make test       # 运行测试
make clean      # 清理构建产物
make proto      # 生成 protobuf 代码
make deps       # 安装依赖
make fmt        # 格式化代码
make lint       # 运行代码检查
make help       # 显示所有可用命令
```

## 项目结构说明

### cmd/ - 应用程序入口
包含所有可执行程序的入口点：
- `server/` - 主游戏服务器
- `tools/` - 工具和辅助程序

### internal/ - 内部包
框架的核心实现，不对外暴露：
- `actor/` - Actor 模型核心实现
- `game/` - 游戏业务逻辑
- `network/` - 网络通信层
- `config/` - 配置管理
- `logger/` - 日志系统
- `metrics/` - 监控指标

### pkg/ - 公共包
可以被外部项目引用的包：
- `proto/` - Protobuf 协议定义
- `errors/` - 错误码定义
- `constants/` - 常量定义

### api/ - 外部 API
对外提供的 API 接口：
- `http/` - HTTP REST API
- `grpc/` - gRPC API

### configs/ - 配置文件
服务器配置文件存放目录

### scripts/ - 脚本
构建、部署等脚本

### docs/ - 文档
项目文档

### test/ - 测试
集成测试和性能测试

## 性能特性

- **并发处理**: 基于 Actor 模型，天然支持高并发
- **异步消息**: 所有消息传递都是异步的，不阻塞发送者
- **无锁设计**: 每个 Actor 独立处理消息，避免锁竞争
- **资源池化**: 使用 goroutine 池和对象池减少开销
- **高效序列化**: 支持 Protobuf 等高效二进制协议

## 安全特性

- **消息验证**: 在网络边界验证所有消息
- **会话管理**: 安全的客户端会话管理
- **速率限制**: 防止恶意客户端攻击（规划中）
- **加密通信**: TLS 支持（规划中）

## 贡献

欢迎贡献！请随时提交 Pull Request。

### 贡献指南

1. Fork 本仓库
2. 创建你的特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交你的更改 (`git commit -m '添加一些很棒的特性'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启一个 Pull Request

## 常见问题

### 如何调试 Actor 消息？

可以添加日志中间件：

```go
func loggingMiddleware(next actor.ReceiveFunc) actor.ReceiveFunc {
    return func(ctx actor.Context) {
        log.Printf("Actor %s 收到消息: %T", ctx.Self(), ctx.Message())
        next(ctx)
    }
}

props := actor.NewProps(producer).WithMiddleware(loggingMiddleware)
```

### 如何处理 Actor 崩溃？

使用监督策略自动重启：

```go
props := actor.NewProps(producer).
    WithSupervisor(actor.NewOneForOneStrategy(5, time.Minute))
```

### 如何进行性能优化？

1. 使用对象池减少内存分配
2. 批量处理消息
3. 使用 Protobuf 替代 JSON
4. 调整 dispatcher 工作线程数
5. 使用路由器分散负载

## 致谢

灵感来源于：
- [Proto.Actor](https://proto.actor/) - Go 的 Actor 模型实现
- [Akka](https://akka.io/) - JVM 的 Actor 框架
- [Orleans](https://dotnet.github.io/orleans/) - 虚拟 Actor 模型

## 许可证

MIT License - 详见 LICENSE 文件

## 联系方式

- 问题反馈: https://github.com/widrin/knights/issues
- 项目主页: https://github.com/widrin/knights

---

**注意**: 本项目目前处于积极开发阶段，API 可能会发生变化。建议在生产环境使用前进行充分测试。
