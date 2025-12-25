# Knights - Actor-Based Game Server Framework

Knights is a high-performance game server framework built in Go, based on the Actor model. It provides a scalable, fault-tolerant foundation for building multiplayer games.

## Features

- **Actor Model**: Lightweight, concurrent actor system with mailbox-based message passing
- **High Performance**: Asynchronous I/O, efficient message dispatching, and minimal overhead
- **Fault Tolerance**: Supervisor strategies for automatic error recovery
- **Scalable**: Horizontal scaling support through clustering (planned)
- **Flexible Networking**: Support for TCP, WebSocket, and multiple codec formats
- **Game-Ready**: Built-in player management, room system, and matchmaking

## Architecture

```
knights/
├── cmd/                     # Application entry points
│   ├── server/             # Game server
│   └── tools/              # Utilities and tools
├── internal/               # Private application code
│   ├── actor/              # Actor framework core
│   ├── game/               # Game logic (player, room, match)
│   ├── network/            # Network layer (TCP, WebSocket)
│   ├── cluster/            # Distributed clustering support
│   ├── config/             # Configuration management
│   └── ...
├── pkg/                    # Public libraries
│   ├── proto/              # Protocol definitions
│   ├── errors/             # Error codes
│   └── constants/          # Constants
└── api/                    # External APIs (HTTP, gRPC)
```

## Quick Start

### Prerequisites

- Go 1.21 or higher
- Make (optional, for using Makefile commands)

### Installation

```bash
# Clone the repository
git clone https://github.com/widrin/knights.git
cd knights

# Install dependencies
go mod download

# Build the server
make build
# or
go build -o bin/server cmd/server/main.go
```

### Running the Server

```bash
# Run directly
make run

# Or run the binary
./bin/server
```

### Configuration

Edit `configs/server.yaml` to customize server settings:

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

## Actor System

The core of Knights is its Actor system, which provides:

### Creating an Actor

```go
import "github.com/widrin/knights/internal/actor"

type MyActor struct {}

func (a *MyActor) Receive(ctx actor.Context) {
    switch msg := ctx.Message().(type) {
    case *actor.Started:
        // Actor started
    case string:
        // Handle string message
        ctx.Respond("Hello: " + msg)
    }
}

// Create actor system
system := actor.NewActorSystem("game")

// Spawn actor
props := actor.NewProps(func() actor.Actor {
    return &MyActor{}
})
pid := system.Spawn(props)

// Send message
system.Send(pid, "Hello")
```

### Key Concepts

- **Actor**: Lightweight concurrent entity that processes messages
- **PID**: Process ID that uniquely identifies an actor
- **Mailbox**: Message queue for each actor
- **Supervisor**: Manages child actors and handles failures
- **Context**: Provides message handling and actor lifecycle methods

## Game Components

### Player Management

```go
// Player actors handle individual player logic
playerProps := actor.NewProps(func() actor.Actor {
    return player.NewPlayerActor(playerID)
})
playerPID := system.Spawn(playerProps)
```

### Room System

```go
// Room actors manage game rooms/battles
roomProps := actor.NewProps(func() actor.Actor {
    return room.NewRoomActor(roomID, maxPlayers)
})
roomPID := system.Spawn(roomProps)
```

### Matchmaking

```go
// Matchmaker actor handles player matching
matchmakerProps := actor.NewProps(func() actor.Actor {
    return match.NewMatchmakerActor()
})
matchmakerPID := system.Spawn(matchmakerProps)
```

## Network Layer

Knights supports multiple network protocols and codecs:

```go
import (
    "github.com/widrin/knights/internal/network"
    "github.com/widrin/knights/internal/network/codec"
)

// Create server with JSON codec
server := network.NewServer(&network.ServerConfig{
    Address:     ":8080",
    Codec:       codec.NewJSONCodec(),
    Handler:     handler,
    ActorSystem: system,
})

server.Start()
```

## Development

### Running Tests

```bash
make test
```

### Code Formatting

```bash
make fmt
```

### Generate Protobuf Code

```bash
make proto
```

## Project Status

🚧 **Under Active Development**

Currently implemented:
- ✅ Core Actor system
- ✅ Basic game components (Player, Room, Match)
- ✅ Network layer with TCP support
- ✅ Configuration management
- ✅ Logging and metrics

Planned features:
- ⏳ Cluster support for distributed deployment
- ⏳ WebSocket support
- ⏳ Persistence layer
- ⏳ Complete matchmaking algorithm
- ⏳ Admin dashboard
- ⏳ Performance benchmarks

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see LICENSE file for details

## Acknowledgments

Inspired by:
- [Proto.Actor](https://proto.actor/) - Actor model for Go
- [Akka](https://akka.io/) - Actor framework for JVM
- [Orleans](https://dotnet.github.io/orleans/) - Virtual Actor Model
