package consul

import (
	"fmt"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/widrin/knights/logger"
	"github.com/widrin/knights/registry"
)

type ConsulRegistry struct {
	client *api.Client
	config *registry.Config
}

func New(config *registry.Config) (*ConsulRegistry, error) {
	cfg := api.DefaultConfig()
	cfg.Address = config.Endpoints[0]
	client, err := api.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("Consul连接失败: %v", err)
	}
	return &ConsulRegistry{client: client, config: config}, nil
}

func (r *ConsulRegistry) Register(serviceID string, meta map[string]string) error {
	registration := &api.AgentServiceRegistration{
		ID:   serviceID,
		Name: "knights-service",
		Port: 8080,
		Check: &api.AgentServiceCheck{
			TTL:                            "10s",
			DeregisterCriticalServiceAfter: "30s",
		},
	}

	if err := r.client.Agent().ServiceRegister(registration); err != nil {
		return fmt.Errorf("服务注册失败: %v", err)
	}

	go r.maintainHeartbeat(serviceID)
	return nil
}

func (r *ConsulRegistry) maintainHeartbeat(serviceID string) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := r.client.Agent().UpdateTTL(serviceID, "", "pass"); err != nil {
				logger.Warn("Consul心跳更新失败: %v", err)
			}
		}
	}
}

func (r *ConsulRegistry) Discover(serviceName string) ([]string, error) {
	services, _, err := r.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, err
	}

	addrs := make([]string, 0)
	for _, service := range services {
		addrs = append(addrs, fmt.Sprintf("%s:%d",
			service.Service.Address,
			service.Service.Port))
	}
	return addrs, nil
}

func (r *ConsulRegistry) Deregister(serviceID string) error {
	return r.client.Agent().ServiceDeregister(serviceID)
}
