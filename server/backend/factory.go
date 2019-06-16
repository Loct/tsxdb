package backend

import (
	"fmt"
)

func Factory(backendType string, opts map[string]interface{}) IAbstractBackend {
	switch backendType {
	case MemoryType.String():
		return NewMemoryBackend()
	case RedisType.String():
		return NewRedisBackend(ExtractRedisOpts(opts))
	default:
		panic(fmt.Sprintf("backend %s not supported", backendType))
	}
}

func ExtractRedisOpts(opts map[string]interface{}) *RedisOpts {
	// @todo connection details etc
	return &RedisOpts{}
}
