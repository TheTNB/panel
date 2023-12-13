package internal

type Setting interface {
	Get(key string, defaultValue ...string) string
	Set(key, value string) error
	Delete(key string) error
}
