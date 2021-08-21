package config

import (
	"go-chi-ddd/constant"
)

var (
	modeName = constant.DebugMode
)

// SetMode sets mode according to input string.
func SetMode(value string) {
	if value == "" {
		value = constant.DebugMode
	}

	switch value {
	case constant.DebugMode:
		modeName = constant.DebugMode
	case constant.ReleaseMode:
		modeName = constant.ReleaseMode
	case constant.TestMode:
		modeName = constant.TestMode
	default:
		panic("mode unknown: " + value + " (available mode: debug release test)")
	}
}

// func Mode() string {
// 	return modeName
// }

func IsDebugging() bool {
	return modeName == constant.DebugMode
}
