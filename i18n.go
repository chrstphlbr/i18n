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

type I18nManager interface {
	Get(key string, language string) (value string, err error)
	SetDefaultLanguage(language string)
}

// mapping types
// key map from keys to values
type keys map[string]values

// values map from laguage to actual values
type values map[string]string

type DefaultI18nManager struct {
	accessLock sync.RWMutex
	// pathes to directories where files are located
	resources          []resource.Repository
	mapping            keys
	mappingConstructed time.Time
	defaultLanguage    string
}

func (m DefaultI18nManager) Get(key string, language string) (value string, err error) {
	m.accessLock.RLock()
	defer m.accessLock.RUnlock()
	// check if there
	values, ok := m.mapping[key]
	if !ok {
		err = fmt.Errorf("could not find mapping for key (%s)", key)
		return
	}

	// check for language
	value, ok = values[language]
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

	for _, repo := range m.resources {
		// update the repository
		repo.Update()
		for _, res := range repo.Resources() {
			// get reader
			jsonReader, err := res.Get()
			// check for errors when getting reader to resource
			if err != nil {
				log.Printf("could not open ressorce: %v", err)
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
	}
	m.mappingConstructed = time.Now()
}

func NewDefaultI18nManager(repositories []resource.Repository) *DefaultI18nManager {
	defaultManager := &DefaultI18nManager{resources: repositories}
	defaultManager.constructMapping()
	return defaultManager
}
