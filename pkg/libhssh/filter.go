package libhssh

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strings"

	cr "github.com/squarescale/cloudresolver"
)

type Filter struct {
	matchableFields []string
	pattern         *regexp.Regexp
}

func NewFilterFromString(s string) (*Filter, error) {
	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("Invalid filter format %q", s)
	}

	mf, err := matchableFields(parts[0])
	if err != nil {
		return nil, err
	}

	reg, err := regexp.Compile(parts[1])
	if err != nil {
		return nil, err
	}

	return &Filter{
		matchableFields: mf,
		pattern:         reg,
	}, nil
}

func (f *Filter) HostMatch(h *cr.Host) bool {
	val := reflect.ValueOf(h).Elem()

	for _, mf := range f.matchableFields {
		s, ok := val.FieldByName(mf).Interface().(string)
		if !ok {
			return false
		}

		matched := f.pattern.MatchString(s)
		if matched {
			return true
		}

	}
	return false
}

// ---

func matchableFields(pattern string) ([]string, error) {
	reg, err := regexp.Compile(pattern)
	if err != nil {
		return []string{}, err
	}

	val := reflect.ValueOf(new(cr.Host)).Elem()

	buff := []string{}

	for i := 0; i < val.NumField(); i++ {
		n := val.Type().Field(i).Name

		if reg.MatchString(strings.ToLower(n)) {
			buff = append(buff, n)
		}
	}

	if len(buff) == 0 {
		return buff, errors.New("No matchable fields")
	}

	sort.Strings(buff)

	return buff, nil
}
