package json

import (
	"log"
	"reflect"
	"testing"
)

func TestMarshal(t *testing.T) {
	type args struct {
		v any
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "Test simple struct",
			args: args{
				v: struct {
					Name string `json:"name"`
					Age  int    `json:"age"`
				}{
					Name: "Alice",
					Age:  30,
				},
			},
			want:    []byte(`{"name":"Alice","age":30}`),
			wantErr: false,
		},
		{
			name: "Test nested struct",
			args: args{
				v: struct {
					Name  string `json:"name"`
					Age   int    `json:"age"`
					Email struct {
						Address string `json:"address"`
						Type    string `json:"type"`
					} `json:"email"`
				}{
					Name: "Bob",
					Age:  25,
					Email: struct {
						Address string `json:"address"`
						Type    string `json:"type"`
					}{
						Address: "bob@example.com",
						Type:    "personal",
					},
				},
			},
			want: []byte(`{"name":"Bob","age":25,"email":{"address":"bob@example.com","type":"personal"}}`),
		},
		{
			name: "array",
			args: args{
				v: struct {
					Str [5]string `json:"str"`
				}{
					Str: [5]string{"hello", "2", "3", "4", "5"},
				},
			},
			want: []byte(`{"str":["hello","2","3","4","5"]}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Marshal(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				log.Println(string(got))
				log.Println(string(tt.want))
				t.Errorf("Marshal() got = %v, want %v", got, tt.want)
			}
		})
	}
}
