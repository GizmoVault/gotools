package main

import (
	"fmt"
	"os"
	"time"

	"github.com/GizmoVault/gotools/printerx"
)

func main() {
	printer := printerx.NewScrollPrinter(os.Stdout, 3)

	for i := 1; i <= 10; i++ {
		if i%2 == 0 {
			printer.PrintLines(fmt.Sprintf("第%d行", i), fmt.Sprintf("第%d——行", i+1))
			i++
		} else {
			printer.PrintLines(fmt.Sprintf("第%d行", i))
		}

		time.Sleep(time.Second)
	}
}
