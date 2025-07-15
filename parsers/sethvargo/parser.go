package sethvargo

import (
	"cmp"
	"reflect"
	"strings"

	"github.com/therenotomorrow/enw"
)

const (
	tagKeyPrefix   = "prefix="
	tagKeyDefault  = "default="
	tagKeyRequired = "required"

	defaultTagKey = "env"
)

type (
	Config struct {
		TagKey string
	}

	Parser struct {
		config Config
	}
)

func New() *Parser {
	return &Parser{config: Config{TagKey: defaultTagKey}}
}

func NewWithConfig(config Config) *Parser {
	if config.TagKey == "" {
		config.TagKey = defaultTagKey
	}

	return &Parser{config: config}
}

func (p *Parser) Config() Config {
	return p.config
}

func (p *Parser) Parse(field reflect.StructField, path string, pkg string) (*enw.Env, string) {
	tagVal := field.Tag.Get(p.config.TagKey)

	if tagVal == "" || tagVal == "-" {
		return nil, ""
	}

	var (
		tag    enw.Tag
		prefix string
		parts  = strings.Split(tagVal, ",")
	)

	for _, part := range parts[1:] {
		trimmedPart := strings.TrimSpace(part)

		if val, found := strings.CutPrefix(trimmedPart, tagKeyPrefix); found {
			prefix = val
		}

		if val, found := strings.CutPrefix(trimmedPart, tagKeyDefault); found {
			tag.Default = val
		}

		if trimmedPart == tagKeyRequired {
			tag.Required = true
		}
	}

	tag.Empty = prefix == "" && tag.Default == "" && !tag.Required

	value := strings.TrimSpace(parts[0])
	if value == "" {
		return nil, prefix
	}

	rType := field.Type
	fieldPkg := rType.PkgPath()
	fieldType := cmp.Or(rType.Name(), rType.String())

	if fieldPkg != "" {
		fieldType = fieldPkg + "." + fieldType
	}

	return &enw.Env{
		Value:   value,
		Field:   field.Name,
		Type:    fieldType,
		Path:    path,
		Package: pkg,
		Tag:     tag,
	}, prefix
}
