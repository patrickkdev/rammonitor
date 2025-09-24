package rammonitor

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Monitor only works on UNIX based systems
//
// Run this function as a goroutine.
// Ex: go Monitor(80, onHighUsage)
func Monitor(maxUsagePercent float64, onHighUsagePercent func(usagePercent float64)) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			out, err := exec.Command("free", "-m").Output()
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			lines := strings.Split(string(out), "\n")
			if len(lines) < 2 {
				fmt.Println("Error: Unexpected output format")
				continue
			}

			fields := strings.Fields(lines[1])
			if len(fields) < 7 {
				fmt.Println("Error: Unexpected output format")
				continue
			}

			total, err := strconv.Atoi(fields[1])
			if err != nil {
				fmt.Println("Error: Unable to parse total memory")
				continue
			}

			used, err := strconv.Atoi(fields[2])
			if err != nil {
				fmt.Println("Error: Unable to parse used memory")
				continue
			}

			memUsagePercent := float64(used) / float64(total) * 100

			if memUsagePercent >= maxUsagePercent {
				onHighUsagePercent(memUsagePercent)
			}
		}
	}
}
