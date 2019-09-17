package libhssh

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	cr "github.com/squarescale/cloudresolver"
)

var (
	hostAttributes = map[string]string{}
)

func init() {
	hostAttributes = initHostAttributes()
}

// ---

type Filter struct {
	fieldName       string
	structFieldName string
	pattern         *regexp.Regexp
}

func NewFilterFromString(s string) (*Filter, error) {
	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("Invalid filter format %q", s)
	}

	n := strings.ToLower(parts[0])

	if !validFieldName(n) {
		return nil, fmt.Errorf("Invalid field name %q", n)
	}

	reg, err := regexp.Compile(parts[1])
	if err != nil {
		return nil, err
	}

	return &Filter{
		fieldName:       n,
		structFieldName: hostAttributes[n],
		pattern:         reg,
	}, nil
}

func (f *Filter) HostMatch(h *cr.Host) bool {
	val := reflect.ValueOf(h).Elem()

	s, ok := val.FieldByName(f.structFieldName).Interface().(string)
	if !ok {
		return false
	}

	return f.pattern.MatchString(s)
}

// ---

func initHostAttributes() map[string]string {
	val := reflect.ValueOf(new(cr.Host)).Elem()

	buff := map[string]string{}

	for i := 0; i < val.NumField(); i++ {
		v := val.Type().Field(i).Name
		k := strings.ToLower(v)
		buff[k] = v
	}

	return buff
}

func validFieldName(n string) bool {
	_, found := hostAttributes[n]
	return found
}
