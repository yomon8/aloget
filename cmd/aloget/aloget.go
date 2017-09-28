package main

import (
	"fmt"
	"os"
	"sort"
	"time"

	humanize "github.com/dustin/go-humanize"

	"github.com/yomon8/aloget/config"
	"github.com/yomon8/aloget/downloader"
	"github.com/yomon8/aloget/list"
)

const (
	timeFormatInput = "2006-01-02 15:04:05 MST"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	zone := "UTC"
	if !config.IsUTC {
		zone, _ = time.Now().In(time.Local).Zone()
	}

	start, _ := time.Parse(
		timeFormatInput,
		fmt.Sprintf("%s %s", config.StartTime, zone),
	)

	end, _ := time.Parse(
		timeFormatInput,
		fmt.Sprintf("%s %s", config.EndTime, zone),
	)

	list, err := list.GetObjectList(start, end, config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if list.Len() == 0 {
		fmt.Println("No S3 objects selected, maybe invalid values in parameters")
		os.Exit(1)
	}

	totalSizeBytes := list.GetTotalByte()
	sort.Sort(list)

	// wait for user prompt
	if !config.ForceMode {
		var key string
		var ok bool
		for !ok {
			fmt.Printf("%s %s  -  %s\n",
				"From-To      \t:",
				list.GetOldestTime().In(time.Local).Format(timeFormatInput),
				list.GetLatestTime().In(time.Local).Format(timeFormatInput),
			)
			fmt.Printf("%s %s  -  %s\n",
				"From-To(UTC) \t:",
				list.GetOldestTime().Format(timeFormatInput),
				list.GetLatestTime().Format(timeFormatInput),
			)
			fmt.Printf("%s %s\n",
				"Donwload Size\t:",
				humanize.Bytes(uint64(totalSizeBytes)),
			)
			fmt.Printf("%s %d objects\n",
				"S3 Objects   \t:",
				list.Len(),
			)
			fmt.Print("Start/Cancel>")
			fmt.Scanf("%s", &key)
			switch key {
			case "S", "s", "Start", "start":
				ok = true
			case "C", "c", "Cancel", "cancel":
				fmt.Println("canceled.")
				os.Exit(1)
			default:
				continue
			}
		}
	}

	err = downloader.NewDownloader(config).Download(list)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Download Completed.\n")
}