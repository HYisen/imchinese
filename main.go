package main

import (
	"fmt"
	"imchinese/finder"
)

func main() {
	text := `## 缘起

> BTW I use Arch

所以你可以猜到我会选择 latest。作为 edge-runner 除了冒着相信上游的风险外，也是背上了时刻更新的负担。那看都看了，不妨留下记录，未来就可以直接过而无需在 déjà vu 中纠结是否曾经考虑过这个问题。

另外我的一贯作风是把 latest 作为常态基准，在特定情况下难以更新的话，就是一个 TODO 甩上去表示以后会如何使用新特性改善。了解更新历史为兼容旧版本提供可能，尽管我有在尽力避免不能滚动更新的情况。

`

	lines := finder.FilterText(text)

	var words []string
	for _, line := range lines {
		words = append(words, finder.FilterWord(line)...)
	}

	prettyPrint(words)
}

func prettyPrint(words []string) {
	for i, word := range words {
		fmt.Printf("%4d %s\n", i, word)
	}
}
