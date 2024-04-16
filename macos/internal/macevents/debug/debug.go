package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/pltanton/autobrowser/macos/internal/macevents"
)

func main() {
	pid, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}
	appInfo := macevents.GetRunningAppInfo(pid)
	fmt.Printf("pid: %d, %v\n", pid, appInfo)

	macevents.StartListenNCEvents()
	macevents.StopListenNCEvents()
	_, err = macevents.WaitForURL(3 * time.Second)
	fmt.Println(err)
}
