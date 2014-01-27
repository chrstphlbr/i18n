package i18n

import (
	"github.com/chrstphlbr/ressource"
	"github.com/chrstphlbr/testHelpers"
	"os"
	"testing"
	"time"
)

const (
	filesDirectory = "./temp"
	greetingJson   = `{
		"hello": {
			"en": "hello",
			"de": "hallo",
			"ru": "Привет"
		}
	}`
	greetingFileName = "./temp/greeting.json"
)

func setUp(t *testing.T) (repo ressource.Repository) {
	os.Mkdir(filesDirectory, 0700)
	testHelpers.CreateFile(t, greetingFileName, greetingJson)

	repo = ressource.NewFileRepository(filesDirectory)
	return
}

func tearDown(t *testing.T) {
	testHelpers.RemoveFile(t, filesDirectory)
}

func TestDefaultLanguage(t *testing.T) {
	const defaultLanaguage = "en"

	manager := &DefaultI18nManager{ressources: []ressource.Repository{}}
	manager.SetDefaultLanguage(defaultLanaguage)

	if manager.defaultLanguage != "en" {
		t.Fatalf("did not set default language (\"%s\") to \"%s\"", manager.defaultLanguage, defaultLanaguage)
	}
}

func TestConstructMapping(t *testing.T) {
	repo := setUp(t)
	defer tearDown(t)

	manager := NewDefaultI18nManager([]ressource.Repository{repo})

	checkElement := func(t *testing.T, mapping keys, element string) *values {
		el, ok := mapping[element]
		if !ok {
			t.Fatal("element (hello) from file not included in mapping")
		}
		return &el
	}

	checkValue := func(t *testing.T, values values, element string, expectedResult string) {
		el, ok := values[element]
		if !ok {
			t.Fatalf("value (%s) from file not included in mapping", element)
		} else if el != expectedResult {
			t.Fatalf("%s (%s) has not the correct value (%s)", element, el, expectedResult)
		}
	}

	before := time.Now()
	// run constructMapping
	manager.constructMapping()

	// check if time (mappingConstructed) is set and after before
	if manager.mappingConstructed.IsZero() {
		t.Fatal("mappingConstructed is not set")
	} else if !manager.mappingConstructed.After(before) {
		t.Fatalf("mappingConstructed is set (%v), but not after before (%v)", manager.mappingConstructed, before)
	}

	// check mapping
	if manager.mapping == nil {
		t.Fatal("mapping not initialized")
	}

	// check if there is just one element
	mappingLength := len(manager.mapping)
	const expectedLength = 1
	if mappingLength != expectedLength {
		t.Fatalf("mapping length wrong (%d), expected %d element(s)", mappingLength, expectedLength)
	}

	el := checkElement(t, manager.mapping, "hello")

	checkValue(t, *el, "en", "hello")
	checkValue(t, *el, "de", "hallo")
	checkValue(t, *el, "ru", "Привет")

	t.Logf("mapping: %+v", manager.mapping)
}

func TestGetFound(t *testing.T) {
	repo := setUp(t)
	defer tearDown(t)

	manager := NewDefaultI18nManager([]ressource.Repository{repo})
	manager.SetDefaultLanguage("en")

	const (
		vEn = "hello"
		vDe = "hallo"
		vRu = "Привет"
	)

	value, err := manager.Get("hello", "en")
	if err != nil {
		t.Fatalf("should have found value (err: %v)", err)
	} else if value != vEn {
		t.Fatalf("found value (%s), but expected \"%s\"", value, vEn)
	}
	value, err = manager.Get("hello", "de")
	if err != nil {
		t.Fatalf("should have found value (err: %v)", err)
	} else if value != vDe {
		t.Fatalf("found value (%s), but expected \"%s\"", value, vDe)
	}
	value, err = manager.Get("hello", "ru")
	if err != nil {
		t.Fatalf("should have found value (err: %v)", err)
	} else if value != vRu {
		t.Fatalf("found value (%s), but expected \"%s\"", value, vRu)
	}
}

func TestGetFoundWithWrongDefaultLang(t *testing.T) {
	repo := setUp(t)
	defer tearDown(t)

	manager := NewDefaultI18nManager([]ressource.Repository{repo})
	manager.SetDefaultLanguage("es")

	const vEn = "hello"

	value, err := manager.Get("hello", "en")
	if err != nil {
		t.Fatalf("should have found value (err: %v)", err)
	} else if value != vEn {
		t.Fatalf("found value (%s), but expected \"%s\"", value, vEn)
	}
}

func TestGetWrongKey(t *testing.T) {
	repo := setUp(t)
	defer tearDown(t)

	manager := NewDefaultI18nManager([]ressource.Repository{repo})

	value, err := manager.Get("huhu", "en")
	if err == nil {
		t.Fatalf("result found (%s). should have returned error", value)
	}
}

func TestGetWrongLanguage(t *testing.T) {
	repo := setUp(t)
	defer tearDown(t)

	manager := NewDefaultI18nManager([]ressource.Repository{repo})

	value, err := manager.Get("hello", "es")
	if err == nil {
		t.Fatalf("result found (%s). should have returned error", value)
	}
}

func TestGetWrongLangButDefaultLang(t *testing.T) {
	repo := setUp(t)
	defer tearDown(t)

	manager := NewDefaultI18nManager([]ressource.Repository{repo})
	manager.SetDefaultLanguage("en")

	const vExp = "hello"

	value, err := manager.Get("hello", "es")
	if err != nil {
		t.Fatalf("returned error despite available default language: %v", err)
	} else if value != vExp {
		t.Fatalf("result found (%s), but expected \"%s\"", value, vExp)
	}
}

func TestGetWrongLangAndWrongDefaultLang(t *testing.T) {
	repo := setUp(t)
	defer tearDown(t)

	manager := NewDefaultI18nManager([]ressource.Repository{repo})
	manager.SetDefaultLanguage("es")

	value, err := manager.Get("hello", "hu")
	if err == nil {
		t.Fatalf("result found (%s). should have returned error", value)
	}
}
