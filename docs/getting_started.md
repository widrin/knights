# Getting Started with Knights

## Introduction

Knights is an Actor-based game server framework written in Go. This guide will help you get started with building your first game server using Knights.

## Installation

### Prerequisites

- Go 1.21 or higher
- Git

### Clone and Setup

```bash
# Clone the repository
git clone https://github.com/widrin/knights.git
cd knights

# Download dependencies
go mod download

# Build the server
make build
```

## Running Your First Server

### 1. Start the Server

```bash
# Using make
make run

# Or directly
./bin/server
```

The server will start and listen on the default port (8080).

### 2. Configuration

Edit `configs/server.yaml` to customize settings:

```yaml
server:
  name: "my-game-server"
  address: "0.0.0.0"
  port: 8080

game:
  max_players: 1000
  tick_rate: 20
```

## Basic Concepts

### Actor

An Actor is a concurrent entity that processes messages:

```go
type MyGameActor struct {
    state int
}

func (a *MyGameActor) Receive(ctx actor.Context) {
    switch msg := ctx.Message().(type) {
    case *actor.Started:
        // Initialize actor
        a.state = 0

    case int:
        // Handle integer messages
        a.state += msg
        ctx.Respond(a.state)

    case string:
        // Handle string messages
        fmt.Println("Received:", msg)
    }
}
```

### Creating Actors

```go
// Create actor system
system := actor.NewActorSystem("my-game")

// Define actor properties
props := actor.NewProps(func() actor.Actor {
    return &MyGameActor{}
})

// Spawn actor
pid := system.Spawn(props)

// Send message
system.Send(pid, 42)
```

### Named Actors

```go
// Create a named actor for easy lookup
pid, err := system.SpawnNamed(props, "game-manager")
if err != nil {
    log.Fatal(err)
}
```

## Building a Simple Game

### Step 1: Create Player Actor

```go
package player

import "github.com/widrin/knights/internal/actor"

type PlayerActor struct {
    playerID string
    hp       int
    position Position
}

func NewPlayerActor(id string) actor.Actor {
    return &PlayerActor{
        playerID: id,
        hp:       100,
        position: Position{X: 0, Y: 0},
    }
}

func (p *PlayerActor) Receive(ctx actor.Context) {
    switch msg := ctx.Message().(type) {
    case *MoveMessage:
        p.position = msg.NewPosition
        ctx.Respond(&MoveResponse{Success: true})

    case *AttackMessage:
        // Handle attack logic
        p.hp -= msg.Damage
        if p.hp <= 0 {
            ctx.StopSelf()
        }
    }
}
```

### Step 2: Setup Network Server

```go
package main

import (
    "github.com/widrin/knights/internal/actor"
    "github.com/widrin/knights/internal/network"
    "github.com/widrin/knights/internal/network/codec"
)

func main() {
    // Create actor system
    system := actor.NewActorSystem("game")

    // Create message handler
    handler := network.NewDefaultHandler(system)

    // Create network server
    server := network.NewServer(&network.ServerConfig{
        Address:     ":8080",
        Codec:       codec.NewJSONCodec(),
        Handler:     handler,
        ActorSystem: system,
    })

    // Start server
    if err := server.Start(); err != nil {
        log.Fatal(err)
    }

    // Wait for shutdown signal
    select {}
}
```

### Step 3: Handle Client Messages

```go
type GameHandler struct {
    actorSystem *actor.ActorSystem
    playerMgr   *actor.PID
}

func (h *GameHandler) HandleMessage(session *Session, msg interface{}) error {
    switch m := msg.(type) {
    case *LoginMessage:
        // Create player actor
        props := actor.NewProps(func() actor.Actor {
            return player.NewPlayerActor(m.PlayerID)
        })
        playerPID := h.actorSystem.Spawn(props)

        // Store PID in session
        session.SetUserData("player_pid", playerPID)

        // Respond to client
        session.Send(&LoginResponse{Success: true})

    case *MoveMessage:
        // Get player PID from session
        if pid, ok := session.GetUserData("player_pid"); ok {
            h.actorSystem.Send(pid.(*actor.PID), m)
        }
    }
    return nil
}
```

## Project Structure

```
your-game/
├── cmd/
│   └── server/
│       └── main.go          # Entry point
├── internal/
│   ├── game/
│   │   ├── player/
│   │   │   └── player.go    # Player logic
│   │   └── battle/
│   │       └── battle.go    # Battle logic
│   └── handler/
│       └── handler.go       # Message handler
└── configs/
    └── server.yaml          # Configuration
```

## Advanced Features

### Supervisor Strategies

Automatically restart actors on failure:

```go
props := actor.NewProps(producer).
    WithSupervisor(actor.NewOneForOneStrategy(3, time.Minute))
```

### Actor Routers

Distribute load across multiple actors:

```go
// Create worker pool
workers := make([]*actor.PID, 10)
for i := 0; i < 10; i++ {
    workers[i] = system.Spawn(workerProps)
}

// Create router
router := actor.NewRoundRobinRouter(workers)

// Router will distribute messages
for _, msg := range messages {
    workerPID := router.Route(msg)
    system.Send(workerPID, msg)
}
```

### Message Middleware

Add logging or metrics to message processing:

```go
func loggingMiddleware(next actor.ReceiveFunc) actor.ReceiveFunc {
    return func(ctx actor.Context) {
        log.Printf("Actor %s received: %T", ctx.Self(), ctx.Message())
        next(ctx)
    }
}

props := actor.NewProps(producer).
    WithMiddleware(loggingMiddleware)
```

## Testing

### Unit Testing Actors

```go
func TestPlayerActor(t *testing.T) {
    system := actor.NewActorSystem("test")
    defer system.Shutdown()

    props := actor.NewProps(func() actor.Actor {
        return player.NewPlayerActor("test-player")
    })
    pid := system.Spawn(props)

    // Send test message
    system.Send(pid, &player.MoveMessage{
        NewPosition: Position{X: 10, Y: 20},
    })

    // Assert results...
}
```

## Next Steps

- Read the [Architecture Guide](architecture.md) for system design details
- Explore example games in the `examples/` directory
- Join our community on Discord

## Troubleshooting

### Port Already in Use

Change the port in `configs/server.yaml`:
```yaml
server:
  port: 8081
```

### Actor Not Responding

- Check that actor is spawned successfully
- Verify message types match what actor expects
- Enable debug logging to trace messages

## Getting Help

- GitHub Issues: https://github.com/widrin/knights/issues
- Documentation: https://knights-framework.dev
- Discord: https://discord.gg/knights
