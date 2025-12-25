package persistence

// Repository provides data persistence abstraction
type Repository interface {
	Save(key string, data interface{}) error
	Load(key string) (interface{}, error)
	Delete(key string) error
}
