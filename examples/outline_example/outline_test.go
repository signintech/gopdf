package outline_example

import "testing"

func TestOutlineExample(t *testing.T) {
	err := OutlineWithPositionExample()
	if err != nil {
		panic(err)
	}
}

func TestOutlineWithLevelExample(t *testing.T) {
	err := OutlineWithLevelExample()
	if err != nil {
		panic(err)
	}
}
