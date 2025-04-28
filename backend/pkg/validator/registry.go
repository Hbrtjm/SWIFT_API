package validator

import (
	"fmt"
	"sync"
)

var (
	registry = make(map[string]GeneralValidator)
	mutex    = &sync.RWMutex{}
)

func Register(name string, validator GeneralValidator) {
	mutex.Lock()
	defer mutex.Unlock()

	if _, exists := registry[name]; exists {
		panic(fmt.Sprintf("GeneralValidator %s already registered", name))
	}
	registry[name] = validator
}

func GetValidator(name string) (GeneralValidator, error) {
	mutex.RLock()
	defer mutex.RUnlock()

	GeneralValidator, exists := registry[name]
	if !exists {
		return nil, fmt.Errorf("GeneralValidator %s not found", name)
	}
	return GeneralValidator, nil
}
