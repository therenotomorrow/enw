package enw

import (
	"fmt"
	"reflect"
	"sync"
)

type (
	Parser interface {
		Parse(field reflect.StructField, path string, pkg string) (env *Env, prefix string)
	}

	Collector struct {
		mutex     sync.Mutex
		parser    Parser
		variables []*Env
	}
)

func New(parser Parser) *Collector {
	return &Collector{mutex: sync.Mutex{}, parser: parser, variables: make([]*Env, 0)}
}

func (c *Collector) Collect(rValue reflect.Value, currPrefix string, currPath string, currPkg string) []*Env {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.walk(rValue, currPrefix, currPath, currPkg)

	variables := c.variables
	c.variables = make([]*Env, 0)

	return variables
}

func (c *Collector) walk(rValue reflect.Value, currPrefix string, currPath string, currPkg string) {
	rValue, isStruct := extractStruct(rValue)
	if !isStruct {
		return
	}

	rType := rValue.Type()
	for i := range rType.NumField() {
		field := rType.Field(i)
		fieldValue := rValue.Field(i)

		if !field.IsExported() {
			continue
		}

		path := field.Name
		if currPath != "" {
			path = currPath + "->" + path
		}

		env, prefix := c.parser.Parse(field, path, currPkg)
		if env != nil {
			env.Value = currPrefix + env.Value

			c.variables = append(c.variables, env)
		}

		switch fieldValue.Kind() {

		case reflect.Slice, reflect.Array:
			// Итерируемся по каждому элементу в срезе/массиве
			for j := 0; j < fieldValue.Len(); j++ {
				elem := fieldValue.Index(j)

				// Проверяем, является ли элемент структурой (или указателем на нее)
				if nested, ok := extractStruct(elem); ok {
					// Формируем путь с индексом, например "Servers->0"
					elemPath := fmt.Sprintf("%s->%d", path, j)
					// Запускаем рекурсию для элемента
					c.walk(nested, currPrefix+prefix, elemPath, nested.Type().PkgPath())
				}
			}

		default:
			if nested, ok := extractStruct(fieldValue); ok {
				c.walk(nested, currPrefix+prefix, path, nested.Type().PkgPath())
			}
		}
	}
}

func extractStruct(rValue reflect.Value) (reflect.Value, bool) {
	for rValue.Kind() == reflect.Ptr {
		if rValue.IsNil() {
			return rValue, false
		}

		rValue = rValue.Elem()
	}

	return rValue, rValue.Kind() == reflect.Struct
}
