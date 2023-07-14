package json

import (
	"encoding/json"
	"reflect"
	"testing"
)

// 数组测试
func TestUnmarshal2(t *testing.T) {
	type args struct {
		v []byte
		s any
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Array JSON",
			args: args{
				v: []byte("[1, 2, 3, 4, 5]"),
				s: new([]int),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Unmarshal(tt.args.v, tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUnmarshal(t *testing.T) {
	type TestData struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	// 创建测试数据
	data := TestData{Name: "John Doe", Age: 30}

	// 将测试数据转换为 JSON 字符串
	jsonData, err := json.Marshal(data)
	//log.Println(string(jsonData))
	if err != nil {
		t.Fatalf("Failed to marshal test data: %v", err)
	}
	// 调用 Unmarshal 函数进行测试
	var result TestData
	err = Unmarshal(jsonData, &result)
	if err != nil {
		t.Errorf("Unmarshal() error: %v", err)
	}
	// 对比结果
	if result.Name != data.Name || result.Age != data.Age {
		t.Errorf("Unmarshaled data does not match expected data")
	}
}

func TestUnmarshal1(t *testing.T) {
	type args struct {
		v []byte
		s any
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "String JSON",
			args: args{
				v: []byte(`"Hello, World!"`),
				s: new(string),
			},
			wantErr: false,
		},
		{
			name: "Boolean JSON",
			args: args{
				v: []byte(`true`),
				s: new(bool),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Unmarshal(tt.args.v, tt.args.s); (err != nil) != tt.wantErr {

				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUnmarshal3(t *testing.T) {
	type Address struct {
		Street  string
		City    string
		Country string
	}

	type Person struct {
		Name         string
		Age          int
		IsStudent    bool
		Addresses    []Address
		PhoneNumbers []string
	}

	jsonStr := `
	[
		{
			"name": "John Doe",
			"age": 30,
			"isStudent": false,
			"addresses": [
				{
					"street": "123 Main St",
					"city": "New York",
					"country": "USA"
				},
				{
					"street": "456 Elm St",
					"city": "Los Angeles",
					"country": "USA"
				}
			],
			"phoneNumbers": [
				"123-456-7890",
				"987-654-3210"
			]
		},
		{
			"name": "Jane Smith",
			"age": 28,
			"isStudent": true,
			"addresses": [
				{
					"street": "789 Oak St",
					"city": "Chicago",
					"country": "USA"
				}
			],
			"phoneNumbers": [
				"555-123-4567"
			]
		}
	]
	`
	var data []Person

	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		t.Errorf("Failed to unmarshal JSON: %v", err)
	}

	// 对每个 Person 对象进行断言验证
	// 这里使用 reflect.DeepEqual 进行对象的深度比较
	// 你也可以根据具体的测试需求编写自定义的断言逻辑
	for i, p := range data {
		switch i {
		case 0:
			expected := Person{
				Name:      "John Doe",
				Age:       30,
				IsStudent: false,
				Addresses: []Address{
					{Street: "123 Main St", City: "New York", Country: "USA"},
					{Street: "456 Elm St", City: "Los Angeles", Country: "USA"},
				},
				PhoneNumbers: []string{"123-456-7890", "987-654-3210"},
			}
			if !reflect.DeepEqual(p, expected) {
				t.Errorf("Mismatch in data at index %d", i)
			}
		case 1:
			expected := Person{
				Name:      "Jane Smith",
				Age:       28,
				IsStudent: true,
				Addresses: []Address{
					{Street: "789 Oak St", City: "Chicago", Country: "USA"},
				},
				PhoneNumbers: []string{"555-123-4567"},
			}
			if !reflect.DeepEqual(p, expected) {
				t.Errorf("Mismatch in data at index %d", i)
			}
		default:
			t.Errorf("Unexpected data at index %d", i)
		}
	}
}
