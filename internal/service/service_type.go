package service

// ServiceType 定义服务类型
type ServiceType string

const (
	// ServiceTypeLogin 登录服务
	ServiceTypeLogin ServiceType = "login"

	// ServiceTypeCenter 中心服务
	ServiceTypeCenter ServiceType = "center"

	// ServiceTypeGateway 网关服务
	ServiceTypeGateway ServiceType = "gateway"

	// ServiceTypeGame 游戏服务
	ServiceTypeGame ServiceType = "game"
)

// String 返回服务类型字符串
func (st ServiceType) String() string {
	return string(st)
}

// IsValid 检查服务类型是否有效
func (st ServiceType) IsValid() bool {
	switch st {
	case ServiceTypeLogin, ServiceTypeCenter, ServiceTypeGateway, ServiceTypeGame:
		return true
	default:
		return false
	}
}
