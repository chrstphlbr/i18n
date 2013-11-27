package i18n

import (
	"fmt"
	"strings"
	"sync"
)

const defaultDirectory = "/files"
const keyNotFoundValue = "no value for key "

var lock sync.RWMutex
var locales map[string]*Locale
var defaultLocale *Locale
var checkInDefaultLocale bool

func init() {
	lock.Lock()
	defer lock.Unlock()
	locales = make(map[string]locales, 3)
	checkInDefaultLocale = true
}

// AddLocale adds a new locale to the i18n library by specifying a name and a directory
func AddLocale(name, directories []string) {
	lock.Lock()
	defer lock.Unlock()
	if cap(locales) == length(locales) {
		// capacity reached
		// extend capacity
	}

	// check if directories are specified, else
	if directories == nil || length(directories) == 0 {
		directories = []string{defaultDirectory}
	}

	lowerName := strings.ToLower(name)

	// do insert here
	// extract key value pairs from files and store them in the mapping attribute
}

// GetLocale returns the appropiate Locale for the provided input
func GetLocale(name string) *Locale {
	lock.RLock()
	defer lock.RUnlock()
	lowerName = strings.ToLower(name)
	l, ok := locales[lowerName]
	// TODO
}

func SetDefaultLocale(defaultLocale *Locale) {

}

// CheckInDefaultLocale changes if a Locale should check a key an not found in the default locale
func CheckInDefaultLocale(check bool) {
	lock.Lock()
	defer lock.Unlock()
	checkInDefaultLocale = check
}
