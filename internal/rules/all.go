package rules

import "github.com/takeuchi-shogo/md-linker/internal/engine"

// All は v0.1 の全 Rule を登録順で返す。
func All() []engine.Rule {
	return []engine.Rule{
		MDL001{},
		MDL002{},
		MDL003{},
		MDL004{},
		MDL006{},
		MDL007{},
	}
}
