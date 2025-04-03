package logx

import "fmt"

type ConsoleRecorder struct{}

func (*ConsoleRecorder) Log(_ Level, a ...interface{}) {
	//nolint: forbidigo // only  debug
	fmt.Println(a...)
}
func (*ConsoleRecorder) Logf(_ Level, format string, a ...interface{}) {
	if format != "" {
		if format[len(format)-1] != '\n' {
			format += "\n"
		}
	}

	//nolint: forbidigo // only debug
	fmt.Printf(format, a...)
}
