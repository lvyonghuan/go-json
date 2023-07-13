package json

import (
	"encoding/json"
	"log"
	"testing"
)

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
