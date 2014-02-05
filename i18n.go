package i18n

import (
	"encoding/json"
	"fmt"
	"github.com/chrstphlbr/resource"
	"log"
	"sync"
	"time"
)

const defaultDirectory = "/files"

var manager struct {
	sync.Mutex
	m I18nManager
}

func Manager() I18nManager {
	manager.Lock()
	defer manager.Unlock()
	if manager.m == nil {
		repository := resource.NewFileRepository(defaultDirectory)
		manager.m = NewDefaultI18nManager(repository)
	}
	return manager.m
}

type I18nManager interface {
	Get(key string, language string) (value string, err error)
	GetAll(key string) (values map[string]string, err error)
	SetDefaultLanguage(language string)
}

// mapping types
// key map from keys to values
type keys map[string]values

// values map from laguage to actual values
type values map[string]string

type DefaultI18nManager struct {
	accessLock sync.RWMutex
	// paths to directories where files are located
	repository         resource.Repository
	mapping            keys
	mappingConstructed time.Time
	defaultLanguage    string
}

func (m DefaultI18nManager) Get(key string, language string) (value string, err error) {
	m.accessLock.RLock()
	defer m.accessLock.RUnlock()
	value, err = m.getWithoutLock(key, language)
	return
}

func (m DefaultI18nManager) getWithoutLock(key, language string) (value string, err error) {
	values, err := m.getAllWithoutLock(key)
	if err != nil {
		return
	}

	// check for language
	value, ok := values[language]
	if ok {
		// found correct value
		return
	}
	// did not find language
	// check if default language specified
	if m.defaultLanguage != "" {
		// default language set
		value, ok = values[m.defaultLanguage]
		if !ok {
			// did not find default language
			err = fmt.Errorf("did not find value for key (%s) in language (%s) and default language (%s)", key, language, m.defaultLanguage)
		}
	} else {
		// no default language set
		err = fmt.Errorf("did not find value for key (%s) in language (%s). No default language set.", key, language)
	}
	return
}

func (m DefaultI18nManager) GetAll(key string) (values map[string]string, err error) {
	m.accessLock.RLock()
	defer m.accessLock.RUnlock()
	values, err = m.getAllWithoutLock(key)
	return
}

func (m DefaultI18nManager) getAllWithoutLock(key string) (values map[string]string, err error) {
	values, ok := m.mapping[key]
	if !ok {
		err = fmt.Errorf("could not find mapping for key (%s)", key)
	}
	return
}

func (m *DefaultI18nManager) SetDefaultLanguage(language string) {
	m.accessLock.Lock()
	defer m.accessLock.Unlock()
	m.defaultLanguage = language
}

func (m *DefaultI18nManager) constructMapping() {
	m.accessLock.Lock()
	defer m.accessLock.Unlock()
	m.mapping = make(keys)
	var keys keys

	// function that adds json decoded keys (currently stored in keys) to mapping
	addKeysToMapping := func() {
		for key, value := range keys {
			m.mapping[key] = value
		}
	}

	// update the repository
	m.repository.Update()
	for _, res := range m.repository.Resources() {
		// get reader
		jsonReader, err := res.Get()
		// check for errors when getting reader to resource
		if err != nil {
			log.Printf("could not open resource: %v", err)
		}
		// create json decoder from reader
		jsonDecoder := json.NewDecoder(jsonReader)
		// decode json
		err = jsonDecoder.Decode(&keys)
		// check for decoding errors
		if err != nil {
			log.Printf("could not unmarshal json resource: %v\n", err)
			continue
		}
		// decoding was successful
		addKeysToMapping()
	}
	m.mappingConstructed = time.Now()
}

func NewDefaultI18nManager(repository resource.Repository) *DefaultI18nManager {
	defaultManager := &DefaultI18nManager{repository: repository}
	defaultManager.constructMapping()
	return defaultManager
}
