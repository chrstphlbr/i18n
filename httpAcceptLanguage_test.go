package i18n

import (
	"fmt"
	"testing"
)

func TestStringer(t *testing.T) {
	header, err := NewAcceptLanguage("en")
	var _ fmt.Stringer = header
	if err != nil {
		// should not have error
		t.Fatalf("unexpected error from NewAcceptLanguage returned: %v", err)
	}
	s := header.String()
	if "en" != s {
		t.Fatalf("header.String (%s) did not return \"en\"", s)
	}
}

func TestAcceptedLanguages(t *testing.T) {
	header, err := NewAcceptLanguage("en")
	if err != nil {
		// should not have error
		t.Fatalf("unexpected error from NewAcceptLanguage returned: %v", err)
	}

	languages := header.AcceptedLanguages()
	elem := <-languages
	if elem != "en" {
		t.Fatalf("did not return correct element (%s)", elem)
	}

	elem, ok := <-languages
	if ok {
		t.Fatalf("receive on channel was ok with value (%s)", elem)
	}
}

func TestALHttpHeader(t *testing.T) {
	testAL(t, "da, en-gb;q=0.8, en;q=0.7")
}

func TestALHttpHeaderUnordered(t *testing.T) {
	testAL(t, "en-gb;q=0.8, da;q=1, en;q=0.7")
}

func TestALHttpHeaderQ0Included(t *testing.T) {
	testAL(t, "en-gb;q=0.8, ru;q=0 ,en;q=0.7, da;q=.9")
}

func testAL(t *testing.T, acceptLanguageHeader string) {
	header, err := NewAcceptLanguage(acceptLanguageHeader)
	if err != nil {
		// should not have error
		t.Fatalf("unexpected error from NewAcceptLanguage returned: %v", err)
	}

	languages := header.AcceptedLanguages()

	elem, ok := <-languages
	if !ok || elem != "da" {
		t.Fatalf("first element was not \"da\", but \"%s\" and channel returned ok=%t", elem, ok)
	}

	elem, ok = <-languages
	if !ok || elem != "en-gb" {
		t.Fatalf("first element was not \"da\", but \"%s\" and channel returned ok=%t", elem, ok)
	}

	elem, ok = <-languages
	if !ok || elem != "en" {
		t.Fatalf("first element was not \"da\", but \"%s\" and channel returned ok=%t", elem, ok)
	}

	elem, ok = <-languages
	if ok {
		t.Fatalf("receive on channel was ok with value (%s)", elem)
	}
}
