# 服务架构说明

## 服务类型

Knights 框架支持多种服务类型，每种服务都作为独立的 Actor 运行：

### 1. 中心服务 (Center Service)

**作用**: 服务器管理中心，负责服务注册、发现和负载均衡

**Actor**: `CenterActor`

**主要功能**:
- 服务器注册与注销
- 服务器列表管理
- 负载均衡（选择最佳服务器）
- 服务器状态更新
- 心跳检测

**配置文件**: `configs/center.yaml`

**启动命令**:
```bash
# 使用 Makefile
make run-center

# 或直接运行
./bin/server -config=configs/center.yaml
```

**端口**: 9001

### 2. 登录服务 (Login Service)

**作用**: 处理用户登录认证

**Actor**: `LoginActor`

**主要功能**:
- 用户登录验证
- Token 生成与验证
- 用户注销
- 在线人数统计

**配置文件**: `configs/login.yaml`

**启动命令**:
```bash
# 使用 Makefile
make run-login

# 或直接运行
./bin/server -config=configs/login.yaml
```

**端口**: 9002

### 3. 网关服务 (Gateway Service)

**作用**: 客户端连接网关，负责消息转发和路由

**Actor**: `GatewayActor`

**主要功能**:
- 客户端会话管理
- 消息转发（客户端 <-> 服务器）
- 消息路由
- 广播消息
- 在线用户管理

**配置文件**: `configs/server.yaml` (默认)

**启动命令**:
```bash
# 使用 Makefile
make run-gateway

# 或直接运行
./bin/server -config=configs/server.yaml
```

**端口**: 8080

### 4. 游戏服务 (Game Service)

**作用**: 游戏逻辑处理服务器

**Actor**: `PlayerManager`, `RoomManager`, `MatchmakerActor`

**主要功能**:
- 玩家管理
- 房间管理
- 匹配系统
- 游戏逻辑处理

**配置文件**: `configs/game.yaml`

**启动命令**:
```bash
# 使用 Makefile
make run-game

# 或直接运行
./bin/server -config=configs/game.yaml
```

**端口**: 9003

## 服务架构图

```
┌─────────────────────────────────────────────────────────┐
│                      客户端层                            │
│                    Game Clients                         │
└─────────────────────────────────────────────────────────┘
                           │
                           ↓
┌─────────────────────────────────────────────────────────┐
│                    网关服务 (8080)                       │
│                   Gateway Service                       │
│  - 会话管理                                              │
│  - 消息路由                                              │
│  - 消息转发                                              │
└─────────────────────────────────────────────────────────┘
                           │
            ┌──────────────┼──────────────┐
            ↓              ↓              ↓
┌─────────────────┐ ┌─────────────┐ ┌──────────────┐
│  登录服务(9002)  │ │ 游戏服务     │ │  其他服务     │
│  Login Service  │ │ Game Service │ │  ...         │
│  - 用户认证      │ │  - 玩家管理  │ │              │
│  - Token管理    │ │  - 房间管理  │ │              │
└─────────────────┘ │  - 匹配系统  │ └──────────────┘
                    └─────────────┘
            ↓              ↓              ↓
┌─────────────────────────────────────────────────────────┐
│                  中心服务 (9001)                         │
│                  Center Service                         │
│  - 服务注册                                              │
│  - 服务发现                                              │
│  - 负载均衡                                              │
│  - 心跳检测                                              │
└─────────────────────────────────────────────────────────┘
```

## 通信流程

### 玩家登录流程

```
1. 客户端 -> 网关服务: 登录请求
2. 网关服务 -> 登录服务: 转发登录请求
3. 登录服务: 验证用户名密码
4. 登录服务 -> 网关服务: 返回 Token
5. 网关服务 -> 客户端: 返回登录成功
6. 网关服务: 绑定会话
```

### 服务器启动流程

```
1. 启动中心服务
   - 监听服务注册请求
   - 初始化服务器列表

2. 启动登录服务
   - 向中心服务注册
   - 开始监听客户端连接

3. 启动网关服务
   - 向中心服务注册
   - 开始监听客户端连接
   - 注册路由到其他服务

4. 启动游戏服务
   - 向中心服务注册
   - 初始化游戏管理器
   - 准备接收玩家
```

## 配置说明

### 服务类型配置

在配置文件中通过 `server.service_type` 指定服务类型：

```yaml
server:
  service_type: "gateway"  # login, center, gateway, game
```

### 各服务配置

#### 中心服务

```yaml
center:
  address: "localhost:9001"
  heartbeat_interval: 30  # 心跳间隔(秒)
  timeout: 90             # 超时时间(秒)
```

#### 登录服务

```yaml
login:
  address: "localhost:9002"
  token_expire: 3600      # token过期时间(秒)
```

#### 网关服务

```yaml
gateway:
  address: "0.0.0.0"
  port: 8080
  max_connections: 10000  # 最大连接数
```

#### 游戏服务

```yaml
game:
  server_id: "game-1"
  max_players: 10000
  tick_rate: 20
  room_max_players: 4
```

## 启动所有服务

### 方式一：使用 Makefile

```bash
# 构建
make build

# 启动所有服务
make run-all
```

### 方式二：使用启动脚本

```bash
# 构建
make build

# 启动所有服务
bash scripts/start_all.sh
```

### 方式三：分别启动

```bash
# 1. 启动中心服务
./bin/server -config=configs/center.yaml &

# 2. 启动登录服务
./bin/server -config=configs/login.yaml &

# 3. 启动网关服务
./bin/server -config=configs/server.yaml &

# 4. 启动游戏服务
./bin/server -config=configs/game.yaml &
```

## 服务扩展

### 添加新的服务类型

1. 在 `internal/service/service_type.go` 中定义新的服务类型常量

```go
const (
    ServiceTypeNewService ServiceType = "new_service"
)
```

2. 创建新服务的 Actor

```go
// internal/service/newservice/newservice_actor.go
package newservice

type NewServiceActor struct {}

func NewNewServiceActor() actor.Actor {
    return &NewServiceActor{}
}

func (n *NewServiceActor) Receive(ctx actor.Context) {
    // 处理消息
}
```

3. 在 `Manager.StartService()` 中添加启动逻辑

```go
case ServiceTypeNewService:
    props = actor.NewProps(func() actor.Actor {
        return newservice.NewNewServiceActor()
    })
```

4. 创建配置文件 `configs/newservice.yaml`

5. 更新 `cmd/server/main.go` 添加初始化逻辑

## Actor 消息示例

### 登录服务消息

```go
// 登录请求
loginReq := &login.LoginRequest{
    Username:  "player1",
    Password:  "password123",
    IPAddress: "192.168.1.100",
}

// 发送到登录服务
system.Send(loginPID, loginReq)
```

### 中心服务消息

```go
// 注册服务器
registerReq := &center.RegisterServerRequest{
    ServerID:   "game-1",
    ServerType: "game",
    Address:    "localhost",
    Port:       9003,
    MaxPlayers: 1000,
}

system.Send(centerPID, registerReq)
```

### 网关服务消息

```go
// 转发消息到客户端
forwardReq := &gateway.ForwardToClientRequest{
    UserID:  "player1",
    Message: gameData,
}

system.Send(gatewayPID, forwardReq)
```

## 监控与运维

### 查看服务状态

```bash
# 查看所有服务进程
ps aux | grep server

# 查看服务日志
tail -f logs/center.log
tail -f logs/login.log
tail -f logs/gateway.log
tail -f logs/game.log
```

### 停止服务

```bash
# 停止所有服务
killall server

# 或发送 SIGTERM 信号
kill -TERM <PID>
```

## 性能优化建议

1. **中心服务**: 使用缓存减少查询开销
2. **登录服务**: 实现 Token 缓存，使用 Redis
3. **网关服务**: 使用连接池，限制最大连接数
4. **游戏服务**: 根据负载启动多个实例

## 故障处理

### 服务崩溃

- 使用 Supervisor 策略自动重启 Actor
- 配置进程守护工具（如 systemd、supervisor）

### 网络分区

- 实现心跳检测机制
- 自动重连机制
- 状态同步

### 负载过高

- 启动多个游戏服务实例
- 使用中心服务进行负载均衡
- 实现玩家迁移机制
