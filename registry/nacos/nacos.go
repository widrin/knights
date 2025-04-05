package nacos

import (
	"fmt"
	"time"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/widrin/knights/logger"
	"github.com/widrin/knights/registry"
)

type NacosRegistry struct {
	client naming_client.INamingClient
	config *registry.Config
}

func New(config *registry.Config) (*NacosRegistry, error) {
	clientConfig := vo.NacosClientParam{}
	clientConfig.ClientConfig = &constant.ClientConfig{
		TimeoutMs:      uint64(config.Timeout * 1000),
		ListenInterval: 10000,
	}
	clientConfig.ServerConfigs = []constant.ServerConfig{
		{
			IpAddr: config.Endpoints[0],
			Port:   8848,
		},
	}

	client, err := clients.NewNamingClient(clientConfig)
	if err != nil {
		return nil, fmt.Errorf("Nacos连接失败: %v", err)
	}
	// 问题在于需要传递指针类型，因为 DeregisterInstance 方法有指针接收器
	return &NacosRegistry{client: client, config: config}, nil
}

func (r *NacosRegistry) Register(serviceID string, meta map[string]string) error {
	_, err := r.client.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          meta["host"],
		Port:        8080,
		ServiceName: "knights-service",
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    meta,
	})
	if err != nil {
		return fmt.Errorf("服务注册失败: %v", err)
	}

	go r.maintainHeartbeat(serviceID, meta)
	return nil
}

func (r *NacosRegistry) maintainHeartbeat(serviceID string, meta map[string]string) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			_, err := r.client.RegisterInstance(vo.RegisterInstanceParam{
				Ip:          meta["host"],
				Port:        8080,
				ServiceName: "knights-service",
				Weight:      10,
				Enable:      true,
				Healthy:     true,
				Ephemeral:   true,
				Metadata:    meta,
			})
			if err != nil {
				logger.Warn("Nacos心跳更新失败: %v", err)
			}
		}
	}
}

func (r *NacosRegistry) Discover(serviceName string) ([]string, error) {
	instances, err := r.client.SelectInstances(vo.SelectInstancesParam{
		ServiceName: serviceName,
		Clusters:    []string{"DEFAULT"},
		HealthyOnly: true,
	})
	if err != nil {
		return nil, err
	}

	addrs := make([]string, 0)
	for _, instance := range instances {
		addrs = append(addrs, fmt.Sprintf("%s:%d", instance.Ip, instance.Port))
	}
	return addrs, nil
}

func (r *NacosRegistry) Deregister(serviceID string) error {
	_, err := r.client.DeregisterInstance(vo.DeregisterInstanceParam{
		ServiceName: "knights-service",
		Ip:          "",
		Port:        8080,
		Ephemeral:   true,
	})
	return err
}
