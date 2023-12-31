package env

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type parser struct {
	keyTemplate    string
	templateValues []any
}

const (
	errParseStringEnvTemplate    = "parser '%s' should be set"
	errParseIntegerTemplate      = "Cant parse '%s' with value '%s' as integer"
	errParseTimeDurationTemplate = "Cant parse '%s' with value '%s' as time duration"
	errParseStringSliceTemplate  = "Can't split '%s' with separator '%s' to slice of strings"
	errParseStringMapTemplate    = "Can't split '%s' with separator '%s' by pair separator '%s' as map of strings"
)

func Get(keyTemplate string, templateValues ...any) *parser {
	return &parser{
		keyTemplate:    keyTemplate,
		templateValues: templateValues,
	}
}

func (p *parser) AsDuration() (time.Duration, error) {
	value, err := p.AsString()
	if err == nil {
		result, err := time.ParseDuration(value)
		if err != nil {
			return 0, errors.Errorf(errParseTimeDurationTemplate, p.keyTemplate, value)
		}

		return result, nil
	}

	return 0, err
}

func (p *parser) AsStringMapOfStrings(pairSeparator string, separator string) (map[string]string, error) {
	data, err := p.AsString()
	if err == nil {
		result := map[string]string{}
		for _, keyValues := range strings.Split(data, separator) {
			pair := strings.Split(keyValues, pairSeparator)
			if len(pair) == 2 && pair[0] != "" && pair[1] != "" {
				result[pair[0]] = pair[1]
			}
		}

		if len(result) == 0 {
			return nil, errors.Errorf(errParseStringMapTemplate, data, separator, pairSeparator)
		}

		return result, nil
	}

	return nil, err
}

func (p *parser) AsSliceOfString(separator string) ([]string, error) {
	values, err := p.AsString()
	if err == nil {
		var result []string
		for _, value := range strings.Split(values, ",") {
			result = append(result, strings.Trim(value, "\n\t\r "))
		}

		if len(result) == 1 && result[0] == values {
			return nil, errors.Errorf(errParseStringSliceTemplate, values, separator)
		}

		return result, nil
	}

	return nil, err
}

func (p *parser) AsIntegerWithDefaultValue(defaultValue int64) int64 {
	value, err := p.AsInteger()
	if err != nil {
		return defaultValue
	}

	return value
}

func (p *parser) AsInteger() (int64, error) {
	value, err := p.AsString()
	if err == nil {
		result, err := strconv.Atoi(value)
		if err != nil {
			return 0, errors.Wrap(err, errParseIntegerTemplate)
		}

		return int64(result), nil
	}

	return 0, err
}

func (p *parser) AsBoolWithDefault(defaultValue bool) bool {
	value, _ := p.AsString()

	switch strings.ToLower(value) {
	case "true":
		return true
	case "false":
		return false
	default:
		return defaultValue
	}
}

func (p *parser) AsString() (string, error) {
	envName := fmt.Sprintf(p.keyTemplate, p.templateValues...)
	result := os.Getenv(envName)
	if result != "" {
		return result, nil
	}

	return "", errors.Errorf(errParseStringEnvTemplate, envName)
}
