package constant

import (
	"go-chi-ddd/config"
)

var (
	modeName = DebugMode
)

const (
	// DebugMode indicates mode is debug.
	DebugMode = "debug"
	// ReleaseMode indicates mode is release.
	ReleaseMode = "release"
	// TestMode indicates mode is test.
	TestMode = "test"
)

func init() {
	SetMode(config.Env.Mode)
}

// SetMode sets mode according to input string.
func SetMode(value string) {
	if value == "" {
		value = DebugMode
	}

	switch value {
	case DebugMode:
		modeName = DebugMode
	case ReleaseMode:
		modeName = ReleaseMode
	case TestMode:
		modeName = TestMode
	default:
		panic("mode unknown: " + value + " (available mode: debug release test)")
	}
}

// func Mode() string {
// 	return modeName
// }

func IsDebugging() bool {
	return modeName == DebugMode
}
