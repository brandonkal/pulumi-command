//
// conv.go
//
// Copyright (c) 2016 Junpei Kawamoto
//
// This software is released under the MIT License.
//
// http://opensource.org/licenses/mit-license.php
//

package structpbconv

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/golang/protobuf/ptypes/struct"
)

// tagKey defines a structure tag name for ConvertStructPB.
const tagKey = "structpb"

// Convert converts a structpb.Struct object to a concrete object.
func Convert(src *structpb.Struct, dst interface{}) error {
	return convertStruct(src, reflect.ValueOf(dst))
}

func toPrimitive(src *structpb.Value) (reflect.Value, bool) {
	switch t := src.GetKind().(type) {
	case *structpb.Value_BoolValue:
		return reflect.ValueOf(t.BoolValue), true
	case *structpb.Value_NullValue:
		return reflect.ValueOf(nil), true
	case *structpb.Value_NumberValue:
		return reflect.ValueOf(t.NumberValue), true
	case *structpb.Value_StringValue:
		return reflect.ValueOf(t.StringValue), true
	default:
		return reflect.Value{}, false
	}
}

func convertValue(src *structpb.Value, dest reflect.Value) error {
	dst := reflect.Indirect(dest)
	if v, ok := toPrimitive(src); ok {
		if !v.Type().AssignableTo(dst.Type()) {
			if !v.Type().ConvertibleTo(dst.Type()) {
				return fmt.Errorf("cannot assign %T to %s", src.GetKind(), dst.Type())
			}
			v = v.Convert(dst.Type())
		}
		dst.Set(v)
		return nil
	}
	switch t := src.GetKind().(type) {
	case *structpb.Value_ListValue:
		return convertList(t.ListValue, dst)
	case *structpb.Value_StructValue:
		return convertStruct(t.StructValue, dst)
	default:
		return fmt.Errorf("unsuported value: %T", src.GetKind())
	}
}

func convertList(src *structpb.ListValue, dest reflect.Value) error {
	dst := reflect.Indirect(dest)
	if dst.Kind() != reflect.Slice {
		return fmt.Errorf("cannot convert %T to %s", src, dst.Type())
	}
	values := src.GetValues()
	elemType := dst.Type().Elem()
	converted := make([]reflect.Value, len(values))
	for i, value := range values {
		element := reflect.New(elemType).Elem()
		if err := convertValue(value, element); err != nil {
			return err
		}
		converted[i] = element
	}
	dst.Set(reflect.Append(dst, converted...))
	return nil
}

func convertStruct(src *structpb.Struct, dest reflect.Value) error {
	dst := reflect.Indirect(dest)
	if dst.Kind() == reflect.Struct {
		fields := src.GetFields()
		for i := 0; i < dst.NumField(); i++ {
			target := dst.Field(i)
			field := dst.Type().Field(i)
			name := field.Tag.Get(tagKey)
			if name == "" {
				name = strings.ToLower(field.Name)
			}
			if v, ok := fields[name]; ok {
				if err := convertValue(v, target); err != nil {
					return err
				}
			}
		}
		return nil
	} else if dst.Kind() == reflect.Map {
		elemType := dst.Type().Elem()
		mapType := reflect.MapOf(reflect.TypeOf(string("")), elemType)
		aMap := reflect.MakeMap(mapType)
		fields := src.GetFields()
		for key, value := range fields {
			element := reflect.New(elemType).Elem()
			if err := convertValue(value, element); err != nil {
				return err
			}
			aMap.SetMapIndex(reflect.ValueOf(key), element)
		}
		dst.Set(aMap)
		return nil
	}

	return fmt.Errorf("cannot convert %T to %s", src, dst.Type())
}
