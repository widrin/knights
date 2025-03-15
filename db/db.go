// db/db.go
package db

// 数据库通用接口
type Database interface {
	Connect(config interface{}) error
	Close() error
	Query(query interface{}, result interface{}) error
	Exec(command interface{}) (interface{}, error)
	BeginTransaction() (Transaction, error)
}

type Transaction interface {
	Commit() error
	Rollback() error
	Exec(command interface{}) (interface{}, error)
}

// 数据库配置接口
type DBConfig interface {
	Validate() error
	DriverName() string
}

// 驱动注册管理
var drivers = make(map[string]func() Database)

func RegisterDriver(name string, factory func() Database) {
	drivers[name] = factory
}

func NewDatabase(config DBConfig) (Database, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}
	factory, ok := drivers[config.DriverName()]
	if !ok {
		return nil, ErrDriverNotRegistered
	}
	return factory(), nil
}
