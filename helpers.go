package moka

import (
	"fmt"
	"strings"
)

func FormatMethodCall(methodName string, args []interface{}) string {
	stringArgs := []string{}
	for _, arg := range args {
		stringArgs = append(stringArgs, fmt.Sprintf("%#v", arg))
	}

	return fmt.Sprintf("%s(%s)", methodName, strings.Join(stringArgs, ", "))
}
