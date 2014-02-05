package i18n

import (
	"github.com/chrstphlbr/resource"
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

func setUp(t *testing.T) (repo resource.Repository) {
	os.Mkdir(filesDirectory, 0700)
	testHelpers.CreateFile(t, greetingFileName, greetingJson)

	repo = resource.NewFileRepository(filesDirectory)
	return
}

func tearDown(t *testing.T) {
	testHelpers.RemoveFile(t, filesDirectory)
}

func TestSetGetLanguage(t *testing.T) {
	const lanaguage = "en"

	manager := &DefaultManager{}
	manager.SetLanguage(lanaguage)

	if manager.GetLanguage() != "en" {
		t.Fatalf("did not set default language (\"%s\") to \"%s\"", manager.GetDefaultLanguage(), lanaguage)
	}
}

func TestSetGetDefaultLanguage(t *testing.T) {
	const defaultLanaguage = "en"

	manager := &DefaultManager{}
	manager.SetDefaultLanguage(defaultLanaguage)

	if manager.GetDefaultLanguage() != "en" {
		t.Fatalf("did not set default language (\"%s\") to \"%s\"", manager.GetDefaultLanguage(), defaultLanaguage)
	}
}

func TestConstructMapping(t *testing.T) {
	repo := setUp(t)
	defer tearDown(t)

	manager := NewDefaultManager(repo)

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

	manager := NewDefaultManager(repo)
	manager.SetLanguage("en")

	const (
		vEn = "hello"
	)

	value, err := manager.Get("hello")
	if err != nil {
		t.Fatalf("should have found value (err: %v)", err)
	} else if value != vEn {
		t.Fatalf("found value (%s), but expected \"%s\"", value, vEn)
	}
}

func TestGetByLanguageFound(t *testing.T) {
	repo := setUp(t)
	defer tearDown(t)

	manager := NewDefaultManager(repo)
	manager.SetDefaultLanguage("en")

	const (
		vEn = "hello"
		vDe = "hallo"
		vRu = "Привет"
	)

	value, err := manager.GetByLanguage("hello", "en")
	if err != nil {
		t.Fatalf("should have found value (err: %v)", err)
	} else if value != vEn {
		t.Fatalf("found value (%s), but expected \"%s\"", value, vEn)
	}
	value, err = manager.GetByLanguage("hello", "de")
	if err != nil {
		t.Fatalf("should have found value (err: %v)", err)
	} else if value != vDe {
		t.Fatalf("found value (%s), but expected \"%s\"", value, vDe)
	}
	value, err = manager.GetByLanguage("hello", "ru")
	if err != nil {
		t.Fatalf("should have found value (err: %v)", err)
	} else if value != vRu {
		t.Fatalf("found value (%s), but expected \"%s\"", value, vRu)
	}
}

func TestGetByLanguageFoundWithWrongDefaultLang(t *testing.T) {
	repo := setUp(t)
	defer tearDown(t)

	manager := NewDefaultManager(repo)
	manager.SetDefaultLanguage("es")

	const vEn = "hello"

	value, err := manager.GetByLanguage("hello", "en")
	if err != nil {
		t.Fatalf("should have found value (err: %v)", err)
	} else if value != vEn {
		t.Fatalf("found value (%s), but expected \"%s\"", value, vEn)
	}
}

func TestGetByLanguageWrongKey(t *testing.T) {
	repo := setUp(t)
	defer tearDown(t)

	manager := NewDefaultManager(repo)

	value, err := manager.GetByLanguage("huhu", "en")
	if err == nil {
		t.Fatalf("result found (%s). should have returned error", value)
	}
}

func TestGetByLanguageWrongLanguage(t *testing.T) {
	repo := setUp(t)
	defer tearDown(t)

	manager := NewDefaultManager(repo)

	value, err := manager.GetByLanguage("hello", "es")
	if err == nil {
		t.Fatalf("result found (%s). should have returned error", value)
	}
}

func TestGetByLanguageWrongLangButDefaultLang(t *testing.T) {
	repo := setUp(t)
	defer tearDown(t)

	manager := NewDefaultManager(repo)
	manager.SetDefaultLanguage("en")

	const vExp = "hello"

	value, err := manager.GetByLanguage("hello", "es")
	if err != nil {
		t.Fatalf("returned error despite available default language: %v", err)
	} else if value != vExp {
		t.Fatalf("result found (%s), but expected \"%s\"", value, vExp)
	}
}

func TestGetByLanguageWrongLangAndWrongDefaultLang(t *testing.T) {
	repo := setUp(t)
	defer tearDown(t)

	manager := NewDefaultManager(repo)
	manager.SetDefaultLanguage("es")

	value, err := manager.GetByLanguage("hello", "hu")
	if err == nil {
		t.Fatalf("result found (%s). should have returned error", value)
	}
}

func TestGetAll(t *testing.T) {
	repo := setUp(t)
	defer tearDown(t)

	manager := NewDefaultManager(repo)

	values, err := manager.GetAll("hello")
	if err != nil {
		t.Fatalf("error but there should not be one: %v", err)
	}
	length := len(values)
	if length != 3 {
		t.Fatalf("should have exactly 3 values but has %d", length)
	}
}

func TestGetAllError(t *testing.T) {
	repo := setUp(t)
	defer tearDown(t)

	manager := NewDefaultManager(repo)

	values, err := manager.GetAll("huhu")
	if err == nil {
		t.Fatalf("did not return but returned values (%v)", values)
	}
}
