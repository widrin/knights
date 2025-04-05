package registry

type ServiceRegistry interface {
	Register(serviceID string, meta map[string]string) error
	Discover(serviceName string) ([]string, error)
	Deregister(serviceID string) error
}

// 基础配置结构应与现有db模块配置保持一致
type Config struct {
	Endpoints []string `yaml:"endpoints"`
	Timeout   int      `yaml:"timeout"`
}
