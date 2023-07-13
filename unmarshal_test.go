package json

import (
	"encoding/json"
	"log"
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
	log.Println(string(jsonData))
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
		log.Println(result)
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
