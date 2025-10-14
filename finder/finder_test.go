package finder

import (
	"reflect"
	"testing"
)

func TestHeading(t *testing.T) {
	input := `# H1

## H2

### H3

TLB`
	index := 3
	want := "H1/H2/H3"

	candidates := Find(input)
	if candidates[index].Path != want {
		t.Errorf("got %q, want %q", candidates[index].Path, want)
	}
}

func TestFind(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{"happy path", `## 缘起

> BTW I use Arch

所以你可以猜到我会选择 latest。作为 edge-runner 除了冒着相信上游的风险外，也是背上了时刻更新的负担。那看都看了，不妨留下记录，未来就可以直接过而无需在 déjà vu 中纠结是否曾经考虑过这个问题。

另外我的一贯作风是把 latest 作为常态基准，在特定情况下难以更新的话，就是一个 TODO 甩上去表示以后会如何使用新特性改善。了解更新历史为兼容旧版本提供可能，尽管我有在尽力避免不能滚动更新的情况。

`, []string{"latest", "edge", "runner", "déjà vu", "latest", "TODO"}},
		{"table", `| name! |
| ---- |
| ABS?  |
| SLB.  |`, []string{"name", "ABS", "SLB"}},
		{"stand alone", "TLB", []string{"TLB"}},
		{"space before", "基于 latest", []string{"latest"}},
		{"code span", "换而言之，无视需要开 `GOEXPERIMENT=xxx` 才能使用的内容。", nil},
		{"date", "2014-12-10", nil},
		{"version number", "1.10", nil},
		{"abbreviation with dots", "茶汤颜色也从唐宋时的绿色（a.k.a. 抹茶）变成了明清时的黄色", []string{"a.k.a."}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			candidates := Find(tt.input)
			var words []string
			for _, candidate := range candidates {
				words = append(words, candidate.Word)
			}
			if got := words; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Find() = %v, want %v", got, tt.want)
			}
		})
	}
}
