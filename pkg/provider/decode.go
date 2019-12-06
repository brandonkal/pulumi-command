package provider

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi/pkg/resource"
	pulumirpc "github.com/pulumi/pulumi/sdk/proto/go"
)

func typeMismatch(path, expected string, actual resource.PropertyValue) *pulumirpc.CheckFailure {
	return &pulumirpc.CheckFailure{
		Property: path,
		Reason:   fmt.Sprintf("expected a %v value, received a %v", expected, actual.TypeString()),
	}
}

func missingRequiredProperty(path, key string) *pulumirpc.CheckFailure {
	return &pulumirpc.CheckFailure{
		Property: path,
		Reason:   fmt.Sprintf("missing required property %v", key),
	}
}

func failureError(f *pulumirpc.CheckFailure) error {
	return errors.Errorf("%v: %v", f.Property, f.Reason)
}

type fieldDesc struct {
	name     string
	optional bool
	forceNew bool
}

func computeName(fieldName string) string {
	runes := []rune(fieldName)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

func getFieldDesc(field reflect.StructField) (*fieldDesc, error) {
	if field.PkgPath != "" {
		return nil, nil
	}

	opts := strings.Split(field.Tag.Get("pulumi"), ",")
	if len(opts) == 0 {
		return &fieldDesc{name: computeName(field.Name)}, nil
	}
	desc := &fieldDesc{name: opts[0]}
	if desc.name == "" {
		desc.name = computeName(field.Name)
	}
	for _, opt := range opts[1:] {
		switch opt {
		case "optional":
			desc.optional = true
		case "forceNew":
			desc.forceNew = true
		default:
			return nil, errors.Errorf("unknown option '%v' in tag for struct field %v", opt, field.Name)
		}
	}
	return desc, nil
}

type checker struct {
	failures []*pulumirpc.CheckFailure
}

func (c *checker) checkProperty(path string, v resource.PropertyValue, schema reflect.Type) error {
	if v.IsComputed() {
		return nil
	}

	switch schema.Kind() {
	case reflect.Bool:
		if !v.IsBool() {
			c.failures = append(c.failures, typeMismatch(path, "bool", v))
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		if !v.IsNumber() {
			c.failures = append(c.failures, typeMismatch(path, "number", v))
		}

	case reflect.String:
		if !v.IsString() {
			c.failures = append(c.failures, typeMismatch(path, "string", v))
		}

	case reflect.Slice:
		if !v.IsArray() {
			c.failures = append(c.failures, typeMismatch(path, "[]", v))
		}
		for i, e := range v.ArrayValue() {
			if err := c.checkProperty(fmt.Sprintf("%v[%v]", path, i), e, schema.Elem()); err != nil {
				return err
			}
		}

	case reflect.Map:
		if schema.Key().Kind() != reflect.String {
			return errors.New("map schema must have string keys")
		}
		if !v.IsObject() {
			c.failures = append(c.failures, typeMismatch(path, "object", v))
		} else {
			for k, e := range v.ObjectValue() {
				if err := c.checkProperty(fmt.Sprintf("%v.%v", path, k), e, schema.Elem()); err != nil {
					return err
				}
			}
		}

	case reflect.Struct:
		if !v.IsObject() {
			c.failures = append(c.failures, typeMismatch(path, "object", v))
		} else {
			m := v.ObjectValue()
			for i := 0; i < schema.NumField(); i++ {
				f := schema.Field(i)
				desc, err := getFieldDesc(f)
				if err != nil {
					return err
				}
				if desc == nil {
					continue
				}

				e, ok := m[resource.PropertyKey(desc.name)]
				if !ok || e.IsNull() {
					if !desc.optional {
						c.failures = append(c.failures, missingRequiredProperty(path, desc.name))
					}
					continue
				}
				if err := c.checkProperty(fmt.Sprintf("%v.%v", path, desc.name), e, f.Type); err != nil {
					return err
				}
			}
		}

	case reflect.Ptr:
		if !v.IsNull() {
			return c.checkProperty(path, v, schema.Elem())
		}

	default:
		return errors.Errorf("unsupported type %v", schema.Name())
	}

	return nil
}

func decodeProperty(path string, v resource.PropertyValue, dest reflect.Value) error {
	switch dest.Kind() {
	case reflect.Bool:
		if !v.IsBool() {
			return failureError(typeMismatch(path, "bool", v))
		}
		dest.SetBool(v.BoolValue())

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if !v.IsNumber() {
			return failureError(typeMismatch(path, "number", v))
		}
		dest.SetInt(int64(v.NumberValue()))

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if !v.IsNumber() {
			return failureError(typeMismatch(path, "number", v))
		}
		dest.SetUint(uint64(v.NumberValue()))

	case reflect.Float32, reflect.Float64:
		if !v.IsNumber() {
			return failureError(typeMismatch(path, "number", v))
		}
		dest.SetFloat(v.NumberValue())

	case reflect.String:
		if !v.IsString() {
			return failureError(typeMismatch(path, "string", v))
		}
		dest.SetString(v.StringValue())

	case reflect.Slice:
		if !v.IsArray() {
			return failureError(typeMismatch(path, "[]", v))
		}
		arrayValue := v.ArrayValue()
		slice := reflect.MakeSlice(dest.Type(), len(arrayValue), len(arrayValue))
		for i, e := range arrayValue {
			if err := decodeProperty(fmt.Sprintf("%v[%v]", path, i), e, slice.Index(i)); err != nil {
				return err
			}
		}
		dest.Set(slice)

	case reflect.Map:
		if dest.Type().Key().Kind() != reflect.String {
			return errors.New("map schema must have string keys")
		}
		if !v.IsObject() {
			return failureError(typeMismatch(path, "object", v))
		}
		m := reflect.MakeMap(dest.Type())
		for k, e := range v.ObjectValue() {
			me := reflect.New(dest.Type().Elem()).Elem()
			if err := decodeProperty(fmt.Sprintf("%v.%v", path, k), e, me); err != nil {
				return err
			}
			m.SetMapIndex(reflect.ValueOf(string(k)), me)
		}
		dest.Set(m)

	case reflect.Struct:
		if !v.IsObject() {
			return failureError(typeMismatch(path, "object", v))
		}
		m := v.ObjectValue()
		for i := 0; i < dest.NumField(); i++ {
			f := dest.Field(i)
			desc, err := getFieldDesc(dest.Type().Field(i))
			if err != nil {
				return err
			}
			if desc == nil {
				continue
			}

			e, ok := m[resource.PropertyKey(desc.name)]
			if !ok || e.IsNull() {
				if !desc.optional {
					return failureError(missingRequiredProperty(path, desc.name))
				}
				f.Set(reflect.Zero(f.Type()))
				continue
			}
			if err := decodeProperty(fmt.Sprintf("%v.%v", path, desc.name), e, f); err != nil {
				return err
			}
		}

	case reflect.Ptr:
		if v.IsNull() {
			dest.Set(reflect.Zero(dest.Type()))
		} else {
			if dest.IsNil() {
				dest.Set(reflect.New(dest.Type().Elem()))
			}
			if err := decodeProperty(path, v, dest.Elem()); err != nil {
				return err
			}
		}

	default:
		return errors.Errorf("unsupported type %v", dest.Type().Name())
	}

	return nil
}
