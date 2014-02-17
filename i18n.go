package i18n

import (
	"github.com/chrstphlbr/resource"
	"sync"
)

const defaultDirectory = "/files"

var manager struct {
	sync.Mutex
	m Manager
}

func GetManager() Manager {
	manager.Lock()
	defer manager.Unlock()
	if manager.m == nil {
		repository := resource.NewFileRepository(defaultDirectory)
		manager.m = NewDefaultManager(repository)
	}
	return manager.m
}

func SetManager(m Manager) {
	manager.Lock()
	defer manager.Unlock()
	manager.m = m
}

type Manager interface {
	ManagerMinimal
	ManagerWithLanguage
	ManagerWithLogging
}

type ManagerMinimal interface {
	// Get searches in the imnternal repository for a value which matches the provided key and language and returns if found the value (value) or if not found
	// an error (err). The parameter value can be a general string or can have the form for Accept-Language HTTP-Header.
	Get(key string, language string) (value string, err error)
	GetAll(key string) (values map[string]string, err error)
	SetDefaultLanguage(language string)
	GetDefaultLanguage() string
}

type ManagerWithLanguage interface {
	GetByLanguage(key string) (value string, err error)
	SetLanguage(language string)
	GetLanguage() string
}

type ManagerWithLogging interface {
	GetOrLog(key string, language string) (value string)
	GetAllOrLog(key string) (values map[string]string)
}

// mapping types
// key map from keys to values
type keys map[string]values

// values map from laguage to actual values
type values map[string]string
