package string

import "fmt"

func Example_findRelativePath() {
	fmt.Println(findRelativePath(
		"/root/config/fish/settings.json",
		"/root/config/cat/food/menu.json",
	))

	// Output:
	// 1
}
