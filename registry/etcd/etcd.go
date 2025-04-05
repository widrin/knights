package etcd

import (
	"context"
	"fmt"
	"time"

	"github.com/widrin/knights/logger"
	"github.com/widrin/knights/registry"
	"go.etcd.io/etcd/clientv3"
)

type EtcdRegistry struct {
	client  *clientv3.Client
	leaseID clientv3.LeaseID
	config  *registry.Config
}

func New(config *registry.Config) (*EtcdRegistry, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   config.Endpoints,
		DialTimeout: time.Duration(config.Timeout) * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("etcd连接失败: %v", err)
	}
	return &EtcdRegistry{client: cli, config: config}, nil
}

func (r *EtcdRegistry) Register(serviceID string, meta map[string]string) error {
	ctx := context.Background()

	// 创建租约
	lease := clientv3.NewLease(r.client)
	grantResp, err := lease.Create(ctx, 10)
	if err != nil {
		return err
	}
	leaseResp := grantResp
	if err != nil {
		return err
	}
	r.leaseID = clientv3.LeaseID(leaseResp.ID)

	// 构建注册键值对
	key := fmt.Sprintf("/knights/services/%s", serviceID)
	val := fmt.Sprintf("%s|%s", meta["host"], meta["port"])

	// 注册服务
	_, err = r.client.Put(ctx, key, val, clientv3.WithLease(r.leaseID))
	if err != nil {
		return fmt.Errorf("服务注册失败: %v", err)
	}

	// 保持心跳
	ch, err := r.client.KeepAlive(ctx, r.leaseID)
	if err != nil {
		return fmt.Errorf("心跳维持失败: %v", err)
	}

	go func() {
		for {
			select {
			case <-ch:
				logger.Debug("etcd心跳续约成功")
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (r *EtcdRegistry) Discover(serviceName string) ([]string, error) {
	resp, err := r.client.Get(context.Background(),
		fmt.Sprintf("/knights/services/%s", serviceName),
		clientv3.WithPrefix(),
	)
	if err != nil {
		return nil, err
	}

	addrs := make([]string, 0)
	for _, kv := range resp.Kvs {
		addrs = append(addrs, string(kv.Value))
	}
	return addrs, nil
}

func (r *EtcdRegistry) Deregister(serviceID string) error {
	_, err := r.client.Revoke(context.Background(), r.leaseID)
	return err
}
