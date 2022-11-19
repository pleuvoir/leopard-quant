package base

import "testing"

func TestGet(t *testing.T) {

	loader := &FileConfigLoader{AbsConfigLoader: AbsConfigLoader{}}
	_ = loader.load()

	getInt, err := loader.getInt("int")

	t.Log(getInt)
	t.Log(err)

}
