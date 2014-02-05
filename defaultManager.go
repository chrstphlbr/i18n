package i18n

import (
	"encoding/json"
	"fmt"
	"github.com/chrstphlbr/resource"
	"log"
	"sync"
	"time"
)

type DefaultManager struct {
	accessLock sync.RWMutex
	// paths to directories where files are located
	repository         resource.Repository
	mapping            keys
	mappingConstructed time.Time
	language           string
	defaultLanguage    string
}

func (m DefaultManager) Get(key string) (value string, err error) {
	m.accessLock.RLock()
	defer m.accessLock.RUnlock()
	if m.language == "" {
		err = fmt.Errorf("language is not set")
		return
	}
	value, err = m.getWithoutLock(key, m.language)
	return
}

func (m DefaultManager) GetByLanguage(key string, language string) (value string, err error) {
	m.accessLock.RLock()
	defer m.accessLock.RUnlock()
	value, err = m.getWithoutLock(key, language)
	return
}

func (m DefaultManager) getWithoutLock(key, language string) (value string, err error) {
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

func (m DefaultManager) GetAll(key string) (values map[string]string, err error) {
	m.accessLock.RLock()
	defer m.accessLock.RUnlock()
	values, err = m.getAllWithoutLock(key)
	return
}

func (m DefaultManager) getAllWithoutLock(key string) (values map[string]string, err error) {
	values, ok := m.mapping[key]
	if !ok {
		err = fmt.Errorf("could not find mapping for key (%s)", key)
	}
	return
}

func (m *DefaultManager) SetLanguage(language string) {
	m.accessLock.Lock()
	defer m.accessLock.Unlock()
	m.language = language
}

func (m DefaultManager) GetLanguage() string {
	m.accessLock.RLock()
	defer m.accessLock.RUnlock()
	return m.language
}

func (m *DefaultManager) SetDefaultLanguage(language string) {
	m.accessLock.Lock()
	defer m.accessLock.Unlock()
	m.defaultLanguage = language
}

func (m DefaultManager) GetDefaultLanguage() string {
	m.accessLock.RLock()
	defer m.accessLock.RUnlock()
	return m.defaultLanguage
}

func (m *DefaultManager) constructMapping() {
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

func NewDefaultManager(repository resource.Repository) *DefaultManager {
	defaultManager := &DefaultManager{repository: repository}
	defaultManager.constructMapping()
	return defaultManager
}
