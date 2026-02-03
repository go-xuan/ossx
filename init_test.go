package ossx

import (
	"fmt"
	"testing"
)

func TestConfig(t *testing.T) {
	fmt.Println(GetClient().GetConfig())
}
