package config

import (
	"errors"
	"os"
	"reflect"
	"strconv"
	"strings"
)

const (
	envTag       = "env"
	namespaceTag = "namespace"
	requiredTag  = "required"
	defaultTag   = "default"
)

var (
	ErrNotStructPointer     = errors.New("not struct pointer")
	ErrValueRequired        = errors.New("value required")
	ErrInvalidType          = errors.New("invalid type")
	ErrNamespaceTagNotFound = errors.New("namespace tag not found")

	envData = map[string]string{}

	osStat  = os.Stat
	osEnv   = os.Environ
	osDirFS = os.DirFS
)

type parserFunc func(v string) (interface{}, error)

var defaultBuiltInParsers = map[reflect.Kind]parserFunc{
	reflect.Bool: func(v string) (interface{}, error) {
		return strconv.ParseBool(v)
	},
	reflect.String: func(v string) (interface{}, error) {
		return v, nil
	},
	reflect.Int: func(v string) (interface{}, error) {
		i, err := strconv.ParseInt(v, 10, 32)
		return int(i), err
	},
	reflect.Int16: func(v string) (interface{}, error) {
		i, err := strconv.ParseInt(v, 10, 16)
		return int16(i), err
	},
	reflect.Int32: func(v string) (interface{}, error) {
		i, err := strconv.ParseInt(v, 10, 32)
		return int32(i), err
	},
	reflect.Int64: func(v string) (interface{}, error) {
		return strconv.ParseInt(v, 10, 64)
	},
	reflect.Int8: func(v string) (interface{}, error) {
		i, err := strconv.ParseInt(v, 10, 8)
		return int8(i), err
	},
	reflect.Uint: func(v string) (interface{}, error) {
		i, err := strconv.ParseUint(v, 10, 32)
		return uint(i), err
	},
	reflect.Uint16: func(v string) (interface{}, error) {
		i, err := strconv.ParseUint(v, 10, 16)
		return uint16(i), err
	},
	reflect.Uint32: func(v string) (interface{}, error) {
		i, err := strconv.ParseUint(v, 10, 32)
		return uint32(i), err
	},
	reflect.Uint64: func(v string) (interface{}, error) {
		i, err := strconv.ParseUint(v, 10, 64)
		return i, err
	},
	reflect.Uint8: func(v string) (interface{}, error) {
		i, err := strconv.ParseUint(v, 10, 8)
		return uint8(i), err
	},
	reflect.Float64: func(v string) (interface{}, error) {
		return strconv.ParseFloat(v, 64)
	},
	reflect.Float32: func(v string) (interface{}, error) {
		f, err := strconv.ParseFloat(v, 32)
		return float32(f), err
	},
}

// Ref: https://github.com/caarlos0/env
func Parse(v interface{}) error {
	// from env
	env := osEnv()
	for _, e := range env {
		p := strings.SplitN(e, "=", 2)
		envData[p[0]] = p[1]
	}

	ptrRef := reflect.ValueOf(v)
	if ptrRef.Kind() != reflect.Ptr {
		return ErrNotStructPointer
	}
	ref := ptrRef.Elem()
	if ref.Kind() != reflect.Struct {
		return ErrNotStructPointer
	}

	return parse(ref, []string{})
}

func parse(ref reflect.Value, prefix []string) error {
	refType := ref.Type()
	for i := 0; i < refType.NumField(); i++ {
		refField := ref.Field(i)
		refTypeField := refType.Field(i)

		namespace := refTypeField.Tag.Get(namespaceTag)
		required := refTypeField.Tag.Get(requiredTag) != ""
		defVal := refTypeField.Tag.Get(defaultTag)
		env := refTypeField.Tag.Get(envTag)
		if refField.Kind() == reflect.Struct {
			if namespace == "" {
				return ErrNamespaceTagNotFound
			}
			if err := parse(refField, append(prefix, namespace)); err != nil {
				return err
			}
			continue
		}

		var val string
		if env != "" {
			key := strings.Join(append(prefix, env), "_")
			tmp, ok := envData[key]
			if ok {
				val = tmp
			}
		}
		if val == "" {
			val = defVal
		}
		if required && val == "" {
			return ErrValueRequired
		}

		parser, ok := defaultBuiltInParsers[refField.Kind()]
		if ok {
			realVal, err := parser(val)
			if err != nil {
				return ErrInvalidType
			}
			refField.Set(reflect.ValueOf(realVal).Convert(refTypeField.Type))
			continue
		}

		switch refField.Kind() {
		case reflect.Slice:
			if val == "" {
				result := reflect.MakeSlice(refTypeField.Type, 0, 0)
				refField.Set(result)
			} else {
				handleSlice(refField, refTypeField, val)
			}
		default:
			return ErrInvalidType
		}
	}

	return nil
}

// TODO: supports separator
func handleSlice(refField reflect.Value, refTypeField reflect.StructField, val string) error {
	parts := strings.Split(val, ",")
	refTypeElemField := refTypeField.Type.Elem()

	parser, ok := defaultBuiltInParsers[refTypeElemField.Kind()]
	if !ok {
		return ErrInvalidType
	}
	result := reflect.MakeSlice(refTypeField.Type, 0, len(parts))
	for _, part := range parts {
		r, err := parser(strings.TrimSpace(part))
		if err != nil {
			return err
		}
		v := reflect.ValueOf(r).Convert(refTypeElemField)
		result = reflect.Append(result, v)
	}
	refField.Set(result)
	return nil
}
