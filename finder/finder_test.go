package finder

import (
	"reflect"
	"testing"
)

func TestFind(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{"table", `| name! |
| ---- |
| ABS?  |
| SLB.  |`, []string{"name", "ABS", "SLB"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Find(tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Find() = %v, want %v", got, tt.want)
			}
		})
	}
}
