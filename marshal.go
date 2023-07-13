package json

import (
	"fmt"
	"reflect"
	"strings"
)

const tagName = "json"

// Marshal 序列化json
// 输入any类型，返回[]byte与err
func Marshal(v any) ([]byte, error) {
	if v == nil {
		return nil, nil
	}
	jsonStr, err := handelInterfaceInMarshal(v)
	if err != nil {
		return nil, err
	}
	return []byte(jsonStr), nil
}

func isNeedSkip(kind reflect.Kind) bool {
	switch kind {
	case reflect.Complex64, reflect.Complex128, reflect.Chan, reflect.Func, reflect.Invalid:
		return true
	default:
		return false
	}
}

// 处理interface
func handelInterfaceInMarshal(v any) (string, error) {
	vType := reflect.TypeOf(v)
	vKind := vType.Kind()
	if isNeedSkip(vKind) {
		return "", nil
	}
	switch vKind {
	//数值直接返回
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%v", reflect.ValueOf(v)), nil
	//string加引号
	case reflect.String:
		return fmt.Sprintf("\"%v\"", reflect.ValueOf(v)), nil
	//
	case reflect.Map:
		var mapStr []string
		for _, val := range reflect.ValueOf(v).MapKeys() {
			key := val.Interface()
			value := reflect.ValueOf(v).MapIndex(val)
			if !value.CanInterface() || isNeedSkip(reflect.TypeOf(value.Interface()).Kind()) {
				continue
			}
			tempStr, err := handelInterfaceInMarshal(val)
			if err != nil {
				return "", err
			}
			mapStr = append(mapStr, fmt.Sprintf("\"%v\":%v", key, tempStr))
		}
		return fmt.Sprintf("{%v}", strings.Join(mapStr, ",")), nil
	case reflect.Array, reflect.Slice:
		var arrayStr []string
		val := reflect.ValueOf(v)
		vLen := val.Len()
		for i := 0; i < vLen; i++ {
			value := val.Index(i).Interface()
			tempStr, err := handelInterfaceInMarshal(value)
			if err != nil {
				return "", err
			}
			arrayStr = append(arrayStr, tempStr)
		}
		return fmt.Sprintf("[%v]", strings.Join(arrayStr, ",")), nil
	case reflect.Struct:
		var structStr []string
		val := reflect.ValueOf(v)
		fieLdNum := val.NumField()
		for i := 0; i < fieLdNum; i++ {
			field := val.Field(i)
			vTy := val.Type()
			if !field.CanInterface() || isNeedSkip(reflect.TypeOf(field).Kind()) {
				continue
			}
			key := vTy.Field(i).Tag.Get(tagName)
			//没有设置key的时候自动以结构体字段名称为json名称
			if key == "" {
				key = vType.Field(i).Name
			}
			tempStr, err := handelInterfaceInMarshal(field.Interface())
			if err != nil {
				return "", err
			}
			structStr = append(structStr, fmt.Sprintf("\"%v\":%v", key, tempStr))
		}
		return fmt.Sprintf("{%v}", strings.Join(structStr, ",")), nil
	default:
		return "", nil
	}
}
