package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/pltanton/autobrowser/macos/internal/macevents"
)

func main() {
	pid, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}
	appInfo := macevents.GetRunningAppInfo(pid)
	fmt.Printf("pid: %d, %v", pid, appInfo)
}
