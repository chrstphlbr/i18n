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
	Get(key string, language string) (value string, err error)
	GetAll(key string) (values map[string]string, err error)
	SetDefaultLanguage(language string)
}

// mapping types
// key map from keys to values
type keys map[string]values

// values map from laguage to actual values
type values map[string]string
