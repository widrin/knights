# Knights Architecture

## Overview

Knights is a game server framework built on the Actor model, providing high concurrency, fault tolerance, and scalability for multiplayer games.

## Core Components

### 1. Actor System ([internal/actor](../internal/actor))

The foundation of the framework, implementing a lightweight Actor model:

- **Actor**: Base interface for all actors
- **ActorSystem**: Manages actor lifecycle and message routing
- **PID**: Process identifier for actor addressing
- **Mailbox**: Asynchronous message queue for each actor
- **Context**: Actor execution context with messaging APIs
- **Supervisor**: Fault tolerance through supervision strategies
- **Dispatcher**: Goroutine pool for message processing
- **Router**: Load balancing and message routing patterns

### 2. Game Logic ([internal/game](../internal/game))

Game-specific business logic implemented as actors:

#### Player System
- **PlayerActor**: Represents an individual player
- **PlayerManager**: Manages all player actors
- Handles login, logout, movement, and player state

#### Room System
- **RoomActor**: Game room/battle instance
- **RoomManager**: Manages all room actors
- Supports join, leave, and broadcasting

#### Matchmaking
- **MatchmakerActor**: Player matching service
- **MatchQueue**: Queue-based matchmaking
- Configurable matching algorithms

#### World
- **WorldActor**: Global game state management
- Server time, online player count, etc.

### 3. Network Layer ([internal/network](../internal/network))

Handles client connections and message routing:

- **Server**: TCP server with configurable codec
- **Session**: Client connection wrapper
- **SessionManager**: Manages all active sessions
- **MessageHandler**: Routes network messages to actors
- **Codec**: Pluggable encoding/decoding (JSON, Protobuf)

### 4. Supporting Systems

#### Configuration ([internal/config](../internal/config))
- YAML-based configuration
- Hot-reload support (planned)

#### Logging ([internal/logger](../internal/logger))
- Structured logging interface
- Multiple backend support

#### Metrics ([internal/metrics](../internal/metrics))
- Real-time performance metrics
- Prometheus integration (planned)

#### Clustering ([internal/cluster](../internal/cluster))
- Distributed actor system (planned)
- Remote actor communication
- Node discovery via gossip protocol

## Message Flow

```
Client -> Network Layer -> Session -> MessageHandler -> Actor System -> Game Actors
                                                            ↓
                                                         Mailbox
                                                            ↓
                                                         Actor.Receive()
```

## Actor Hierarchy

```
ActorSystem
├── PlayerManager
│   ├── PlayerActor (player-1)
│   ├── PlayerActor (player-2)
│   └── ...
├── RoomManager
│   ├── RoomActor (room-1)
│   ├── RoomActor (room-2)
│   └── ...
├── MatchmakerActor
└── WorldActor
```

## Concurrency Model

- Each actor processes messages sequentially from its mailbox
- No shared state between actors
- All communication through asynchronous messages
- Goroutine pool for efficient CPU utilization

## Fault Tolerance

### Supervision Strategies

1. **OneForOne**: Restart only the failed actor
2. **AllForOne**: Restart all sibling actors

### Error Handling
- Automatic actor restart on failure
- Configurable retry policies
- Error escalation to parent supervisor

## Scalability

### Vertical Scaling
- Configurable dispatcher thread pool
- Efficient message batching
- Lock-free data structures

### Horizontal Scaling (Planned)
- Cluster-aware actor system
- Location transparency
- Consistent hashing for actor placement

## Performance Considerations

- Lock-free mailbox implementation
- Zero-copy message passing where possible
- Object pooling for frequent allocations
- Efficient serialization formats

## Security

- Message validation at network boundary
- Rate limiting per session
- Authentication/authorization hooks
- Encrypted communication (TLS support planned)

## Extensibility

### Custom Actors
Implement the `Actor` interface:
```go
type MyActor struct {}

func (a *MyActor) Receive(ctx actor.Context) {
    // Handle messages
}
```

### Custom Codecs
Implement the `Codec` interface:
```go
type MyCodec struct {}

func (c *MyCodec) Encode(conn net.Conn, msg interface{}) error
func (c *MyCodec) Decode(conn net.Conn) (interface{}, error)
```

### Middleware
Add message processing middleware:
```go
props := actor.NewProps(producer).
    WithMiddleware(loggingMiddleware).
    WithMiddleware(metricsMiddleware)
```

## Deployment

### Standalone Mode
Single-node deployment for development and small games

### Cluster Mode (Planned)
Multi-node deployment with:
- Service discovery
- Load balancing
- Hot code deployment

## Monitoring

- Health check endpoints
- Metrics export
- Distributed tracing (planned)
- Admin dashboard (planned)
