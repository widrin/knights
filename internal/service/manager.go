package service

import (
	"fmt"

	"github.com/widrin/knights/internal/actor"
	"github.com/widrin/knights/internal/service/center"
	"github.com/widrin/knights/internal/service/gateway"
	"github.com/widrin/knights/internal/service/login"
	"github.com/widrin/knights/internal/service/player"
)

// Manager 服务管理器
type Manager struct {
	system      *actor.ActorSystem
	services    map[ServiceType]*actor.PID
	serviceType ServiceType
}

// NewManager 创建服务管理器
func NewManager(system *actor.ActorSystem, serviceType ServiceType) *Manager {
	return &Manager{
		system:      system,
		services:    make(map[ServiceType]*actor.PID),
		serviceType: serviceType,
	}
}

// StartService 启动指定类型的服务
func (m *Manager) StartService(stype ServiceType) (*actor.PID, error) {
	if !stype.IsValid() {
		return nil, fmt.Errorf("invalid service type: %s", stype)
	}

	var props *actor.Props

	switch stype {
	case ServiceTypeLogin:
		props = actor.NewProps(func() actor.Actor {
			return login.NewLoginActor()
		})

	case ServiceTypeCenter:
		props = actor.NewProps(func() actor.Actor {
			return center.NewCenterActor()
		})

	case ServiceTypeGateway:
		props = actor.NewProps(func() actor.Actor {
			return gateway.NewGatewayActor()
		})

	case ServiceTypeGame:
		// 游戏服务启动玩家管理器
		props = actor.NewProps(func() actor.Actor {
			return player.NewPlayerManager()
		})

	default:
		return nil, fmt.Errorf("unsupported service type: %s", stype)
	}

	// 使用服务类型作为名称创建Actor
	pid, err := m.system.SpawnNamed(props, string(stype))
	if err != nil {
		return nil, fmt.Errorf("failed to spawn %s service: %w", stype, err)
	}

	m.services[stype] = pid
	return pid, nil
}

// GetService 获取指定类型的服务PID
func (m *Manager) GetService(stype ServiceType) (*actor.PID, bool) {
	pid, ok := m.services[stype]
	return pid, ok
}

// StopService 停止指定类型的服务
func (m *Manager) StopService(stype ServiceType) error {
	pid, ok := m.services[stype]
	if !ok {
		return fmt.Errorf("service not found: %s", stype)
	}

	m.system.Stop(pid)
	delete(m.services, stype)
	return nil
}

// StopAll 停止所有服务
func (m *Manager) StopAll() {
	for stype, pid := range m.services {
		m.system.Stop(pid)
		delete(m.services, stype)
	}
}

// GetServiceType 获取当前服务类型
func (m *Manager) GetServiceType() ServiceType {
	return m.serviceType
}
