package registry

// Registry provides service discovery
type Registry interface {
	Register(serviceName string, address string) error
	Deregister(serviceName string) error
	Discover(serviceName string) ([]string, error)
}
