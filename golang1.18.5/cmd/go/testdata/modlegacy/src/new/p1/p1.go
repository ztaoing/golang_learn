package p1

import (
	_ "new/v2"
	_ "old/p2"

	_ "new/v2/p2"

	_ "new/sub/v2/x/v1/y"

	// v2 is module, v1 is directory in module
	_ "new/sub/inner/x"
)

// new/sub/inner/go.mod overrides new/sub/go.mod
