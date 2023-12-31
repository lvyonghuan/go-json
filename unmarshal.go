package json

import (
	"errors"
	"log"
	"reflect"
	"strconv"
)

/*
作用解释：
解析字符串，将字符串按类存入map当中，再将map里的内容填入struct和slice
*/

func Unmarshal(v []byte, s any) error {
	//当JSON字符串为NULL时，直接赋值nil
	if v == nil {
		s = nil
		return nil
	}
	val := string(v)
	//当JSON字符串为布尔值的时候，直接赋值布尔类型
	switch val {
	case "true":
		s = true
		return nil
	case "false":
		s = false
		return nil
	}
	//当JSON字符串为字符串的时候，返回本身，同时去除引号
	if val[0] == '"' {
		s = val[1 : len(val)-1]
		return nil
	}
	//当JSON字符串为数字类型的时候
	if temp, err := strconv.Atoi(val); err == nil {
		s = temp
		return nil
	}
	if temp, err := strconv.ParseFloat(val, 64); err == nil {
		s = temp
		return nil
	}
	//当JSON字符串为数组或object时，进行进一步处理
	return handelInterfaceInUnmarshal(val, s)
}

//处理对应类型的JSON串

func handelInterfaceInUnmarshal(v string, s any) error {
	//定基调，整个结构是数组还是对象
	var first string
	for i := 0; i < len(v); i++ {
		temp := string(v[i])
		if temp == " " || temp == "\n" || temp == "\r" || temp == "\t" {
			continue
		} else {
			first = temp
			break
		}
	}
	typ := checkType(first)
	var index int
	var handel interface{}
	if typ == array {
		arrMap, err := handelArray(v, &index)
		if err != nil {
			return err
		}
		handel = arrMap
		elemVal := reflect.ValueOf(s).Elem()
		if elemVal.Kind() != reflect.Slice {
			return errors.New("无法解析的切片类型")
		}
		s, err = madeSlice(handel, elemVal)
		if err != nil {
			return err
		}
		return nil
	} else if typ == object {
		objMap, err := handelObject(v, &index)
		if err != nil {
			return err
		}
		handel = objMap
		val := reflect.TypeOf(s).Elem()
		if val.Kind() != reflect.Struct {
			return errors.New("无法解析的结构体类型")
		}
		temp, err := madeObject(handel, val)
		value := reflect.ValueOf(s).Elem()
		value.Set(reflect.ValueOf(temp))
		log.Println(s)
		if err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("预处理逻辑的问题")
	}
}

// 检测到"["时，按照数组逻辑进行处理
func handelArray(v string, index *int) (arrayMap, error) {
	var arrayM = make(arrayMap)
	for ; string(v[*index]) != "["; *index++ {
		temp := string(v[*index])
		if temp == " " || temp == "\n" || temp == "\r" || temp == "\t" {
			continue
		} else {
			break
		}
	}
	for ; v[*index] != ']'; *index++ {
		skip(v, index)
		typ := checkType(string(v[*index]))
		switch typ {
		case str:
			*index++ //略过起始的"
			arrayM[str] = handelStr(v, index)
		case object:
			*index++ //略过起始的{
			objMap, err := handelObject(v, index)
			if err != nil {
				return nil, err
			}
			arrayM[object] = objMap
		case startCount:
			var TempStr string
			var isFloat = false
			//进行一个for循环，拼接字符串
			for j := index; v[*j] != ',' && v[*j] != ']'; *j++ {
				skip(v, j)
				TempStr += string(v[*j])
				if v[*j] == '.' {
					isFloat = true
				}
			}
			switch TempStr {
			case "true":
				arrayM[boolean] = true
			case "false":
				arrayM[boolean] = false
			default:
				if isFloat {
					num, err := strconv.ParseFloat(TempStr, 64)
					if err != nil {
						return nil, err
					}
					arrayM[numberFloat] = num
				} else {
					num, err := strconv.Atoi(TempStr)
					if err != nil {
						return nil, err
					}
					arrayM[numberInt] = num
				}
			}
			//鉴于在上面的情况下有可能直接循环到]，故这里做一次判断
			if v[*index] == ']' {
				return arrayM, nil
			}
		}
	}
	return arrayM, nil
}

// 检测到"{"时，按照结构体逻辑进行处理
func handelObject(v string, index *int) (objectMap, error) {
	var objMap = make(objectMap)
	//跳过某些特殊符号
	for ; string(v[*index]) != "\""; *index++ {
	}
	for ; v[*index] != '}'; *index++ {
		skip(v, index)
		//提取字符串中的tag
		var tag string //存储tag，作为objectMap的key
		for i := index; string(v[*i]) != "\""; *i++ {
			tag += string(v[*i])
		}
		*index += 2 //略过:
		typ := checkType(string(v[*index]))
		switch typ {
		case str:
			temp := make(map[int]interface{})
			temp[str] = handelStr(v, index)
			objMap[tag] = temp
		case array:
			arrMap, err := handelArray(v, index)
			if err != nil {
				return nil, err
			}
			temp := make(map[int]interface{})
			temp[array] = arrMap
			objMap[tag] = temp
		case startCount:
			var TempStr string
			var isFloat = false
			//进行一个for循环，拼接字符串
			for j := index; v[*j] != ',' && v[*j] != '}'; *j++ {
				skip(v, j)
				TempStr += string(v[*j])
				if v[*j] == '.' {
					isFloat = true
				}
			}
			switch TempStr {
			case "true":
				temp := make(map[int]interface{})
				temp[boolean] = true
				objMap[tag] = temp
			case "false":
				temp := make(map[int]interface{})
				temp[boolean] = false
				objMap[tag] = temp
			default:
				if isFloat {
					num, err := strconv.ParseFloat(TempStr, 64)
					if err != nil {
						return nil, err
					}
					temp := make(map[int]interface{})
					temp[numberFloat] = num
					objMap[tag] = temp
				} else {
					num, err := strconv.Atoi(TempStr)
					if err != nil {
						return nil, err
					}
					temp := make(map[int]interface{})
					temp[numberInt] = num
					objMap[tag] = temp
				}
			}
			//鉴于在上面的情况下有可能直接循环到]，故这里做一次判断
			if v[*index] == '}' {
				return objMap, nil
			}
		}
	}
	return objMap, nil
}

// 检测到"时，按照字符串逻辑处理
func handelStr(v string, index *int) string {
	var str string
	for *index += 1; v[*index] != '"'; *index++ {
		str += string(v[*index])
	}
	*index += 2
	return str
}

// 跳过空格，换行符和tab
func skip(v string, i *int) {
	for ; *i < len(v); *i++ {
		temp := string(v[*i])
		if temp == " " || temp == "\n" || temp == "\r" || temp == "\t" {
			continue
		} else {
			break
		}
	}
}

// 将数据映射到传入结构上

// 处理字符串类型
func getStr(v interface{}) (string, error) {
	tempStr, ok := v.(string)
	if !ok {
		return "", errors.New("无法解析的字符串")
	}
	return tempStr, nil
}

// 处理bool类型
func getBool(v interface{}) (bool, error) {
	tempBool, ok := v.(bool)
	if !ok {
		return tempBool, errors.New("无法解析的布尔类型")
	}
	return tempBool, nil
}

// 处理int类型
func getInt(v interface{}) (int, error) {
	tempInt, ok := v.(int)
	if !ok {
		return tempInt, errors.New("无法解析的int类型")
	}
	return tempInt, nil
}

// 处理float64类型
func getFloat64(v interface{}) (float64, error) {
	tempFloat, ok := v.(float64)
	if !ok {
		return tempFloat, errors.New("无法解析的float64类型")
	}
	return tempFloat, nil
}

//分割字符串，定位JSON类型

// 定义类型常量
const (
	startCount  = 0 //这是一个信号标记，用于通知程序开始计数字符。鉴于特殊符号总是对称的，计数从特殊符号前半部分未出现的时候开始，直到第一个逗号
	array       = 1 //切片类型
	object      = 2 //结构体类型
	str         = 3 //字符串类型
	boolean     = 4 //布尔类型
	numberFloat = 5 //float64类型
	numberInt   = 6 //int类型
)

// 通过map存储解析的字符串
type arrayMap map[int]interface{}             //键为类型，参考上述类型常量，值为接口。用在数组类型上。
type objectMap map[string]map[int]interface{} //存储object。由object的tag索引到存储的类型map，再由类型的tag索引到存储的值。可以由断言取出。

// 确认类型
func checkType(sign string) int {
	//确认特殊字符
	switch sign {
	case "[":
		return array
	case "{":
		return object
	case "\"":
		return str
	default:
		return startCount
	}
}

//对进行映射

// 制造切片
func madeSlice(v interface{}, s reflect.Value) (interface{}, error) {
	newSlice := reflect.MakeSlice(s.Type(), 0, 0)
	tempMap, ok := v.(arrayMap)
	if !ok {
		return nil, errors.New("不是切片类型")
	}
	for key, value := range tempMap {
		switch key {
		case str:
			tempStr, err := getStr(value)
			if err != nil {
				return nil, err
			}
			newSlice = reflect.Append(newSlice, reflect.ValueOf(tempStr))
		case boolean:
			tempBool, err := getBool(value)
			if err != nil {
				return nil, err
			}
			newSlice = reflect.Append(newSlice, reflect.ValueOf(tempBool))
		case numberFloat:
			tempFloat, err := getFloat64(value)
			if err != nil {
				return nil, err
			}
			newSlice = reflect.Append(newSlice, reflect.ValueOf(tempFloat))
		case numberInt:
			tempInt, err := getInt(value)
			if err != nil {
				return nil, err
			}
			newSlice = reflect.Append(newSlice, reflect.ValueOf(tempInt))
		case array:
			tempSlice, err := madeSlice(value, s)
			if err != nil {
				return nil, err
			}
			newSlice = reflect.Append(newSlice, reflect.ValueOf(tempSlice))
		case object:
			elem := s.Type().Elem().Elem()
			tempObj, err := madeObject(value, elem)
			if err != nil {
				return nil, err
			}
			newSlice = reflect.Append(newSlice, reflect.ValueOf(tempObj))
		}
	}
	return newSlice.Interface(), nil
}

// 制造结构体
func madeObject(v interface{}, s reflect.Type) (interface{}, error) {
	newStruct := reflect.New(s).Elem()
	tempMap := v.(objectMap)
	for i := 0; i < s.NumField(); i++ {
		field := s.Field(i)
		//根据tag取值
		tag := field.Tag.Get(tagName) //tagName在marshal.go下面
		handelMap := tempMap[tag]
		for key, value := range handelMap {
			switch key {
			case str:
				tempStr, err := getStr(value)
				if err != nil {
					return nil, err
				}
				newStruct.Field(i).SetString(tempStr)
			case boolean:
				tempBool, err := getBool(value)
				if err != nil {
					return nil, err
				}
				newStruct.Field(i).SetBool(tempBool)
			case numberFloat:
				tempFloat, err := getFloat64(value)
				if err != nil {
					return nil, err
				}
				newStruct.Field(i).SetFloat(tempFloat)
			case numberInt:
				tempInt, err := getInt(value)
				if err != nil {
					return nil, err
				}
				newStruct.Field(i).SetInt(int64(tempInt))
			case array:
				tempSlice, err := madeSlice(value, reflect.ValueOf(newStruct.Field(i)))
				if err != nil {
					return nil, err
				}
				newStruct.Field(i).Set(reflect.ValueOf(tempSlice))
			case object:
				tempObj, err := madeObject(value, newStruct.Field(i).Type())
				if err != nil {
					return nil, err
				}
				newStruct.Field(i).Set(reflect.ValueOf(tempObj))
			}
		}
	}
	return newStruct.Interface(), nil
}
