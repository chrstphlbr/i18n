package i18n

import (
	"fmt"
)

// possible localeError values
const (
	errNoEntry localeError = localeError(iota)
	errConsMap
	errMapNotCons
)

// localeError indicates an error with handling an object of type Locale
// possible errors are:
//		1	no entry for given key found
//		2	error while contructing mapping
//		3	mapping not constructed
type localeError int

// Error returns a string representation of localeError
// this method is needed in order for localeError to satisfy the error interface
func (le localeError) Error() string {
	return fmt.Sprintf("LocaleError: #%d", int(le))
}

// Locale represents a locale (a specific language) with a Name and Directories where the locale files can be found.
type Locale struct {
	Name               string
	Directories        []string
	mapping            map[string]string
	mappingConstructed bool
}

// Get takes a key and returns the corresponding value in this Locale (language).
// If everything went good value is returned and err is nil.
// If there was an internal error err has a value. (see localeError for more information)
// If there was no match for the given key an err is also returned.
func (l *Locale) Get(key string) (value string, err localeError) {
	if !mappingConstructed {
		err = errMapNotCons
		return
	}
	// find key in mapping for this locale
	// if key is not available, check in default locale (e.g. en)
}

// constructMapping constructs the internally used mapping object with the currently specified Directories.
// if some error occurs while constructing the mapping err is returned (check localeError for information in error codes).
func (l *Locale) constructMapping() {

}
