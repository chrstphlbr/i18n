package i18n

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type AcceptLanguage struct {
	headerValue string
	// languages ordered by descending quality; q=0 included
	values *acceptLanguageParameters
	// maps languages to qualities
	valuesToQuality map[string]float32
}

func (al AcceptLanguage) String() string {
	return al.headerValue
}

// AcceptLanguages returns a read-only channel that provides all accepted languages of that AcceptHeader in order from most accpeted
// to least accepted and leaves out those with quality 0 (q=0)
func (al AcceptLanguage) AcceptedLanguages() <-chan string {
	c := make(chan string)
	go func() {
		for _, v := range *al.values.params {
			c <- v.language
		}
		close(c)
	}()
	return c
}

func NewAcceptLanguage(headerValue string) (acceptLanguage *AcceptLanguage, err error) {
	if headerValue == "" {
		err = fmt.Errorf("headerValue can not be empty")
		return
	}

	al := &AcceptLanguage{headerValue: headerValue}
	al.valuesToQuality = map[string]float32{}
	params := make([]acceptLanguageParameter, 0, 5)

	add := func(alp *acceptLanguageParameter) {
		params = append(params, *alp)
		al.valuesToQuality[alp.language] = alp.quality
	}

	valueArray := strings.Split(headerValue, ",")
	var alp *acceptLanguageParameter
	// multiple languages
	for _, entry := range valueArray {
		alp, err = newAcceptLanguageParameter(entry)
		if err != nil {
			err = fmt.Errorf("headerValue malformed. error somewhere in \"%s\"", entry)
			return
		}
		add(alp)
	}

	al.values = &acceptLanguageParameters{&params}
	// sort descending
	sort.Sort(sort.Reverse(al.values))
	acceptLanguage = al
	return
}

type acceptLanguageParameter struct {
	language string
	quality  float32
}

func newAcceptLanguageParameter(value string) (alp *acceptLanguageParameter, err error) {
	// trims all whitespaces and colons
	value = strings.Trim(value, " ,")
	valueArray := strings.Split(value, ";")
	alp = &acceptLanguageParameter{}

	switch len(valueArray) {
	case 1:
		// no split
		alp.language = value
		alp.quality = 1
	case 2:
		// splitted with 2 elements
		alp.language = strings.TrimSpace(valueArray[0])
		vTrimmed := strings.TrimSpace(valueArray[1])
		// vTrimmed is of form "q=0.7"
		vRemovedSuffix := strings.TrimPrefix(vTrimmed, "q=")
		q, err1 := strconv.ParseFloat(vRemovedSuffix, 32)
		if err1 != nil {
			err = fmt.Errorf("language entry malformed. could not parse: ", vTrimmed)
			return
		}
		alp.quality = float32(q)
	default:
		// splitted and invalid amount of parameters
		// malformed language parameter
		err = fmt.Errorf("language entry malformed (%s)", value)
	}
	return
}

type acceptLanguageParameters struct {
	params *[]acceptLanguageParameter
}

func (p acceptLanguageParameters) String() string {
	return fmt.Sprint(p.params)
}

func (p acceptLanguageParameters) Len() int {
	return len(*p.params)
}

func (p acceptLanguageParameters) Less(i, j int) bool {
	elemI := (*p.params)[i]
	elemJ := (*p.params)[j]
	if elemI.quality <= elemJ.quality {
		return true
	}
	return false
}

func (p *acceptLanguageParameters) Swap(i, j int) {
	(*p.params)[i], (*p.params)[j] = (*p.params)[j], (*p.params)[i]
}
