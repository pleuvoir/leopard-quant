package util

import (
	"fmt"
	"testing"
)

func TestRootPath(t *testing.T) {
	path, err := RootPath()
	fmt.Println(path)
	fmt.Println(err)
}
