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

	localzone, _ := time.Now().In(time.Local).Zone()
	start, _ := time.Parse(
		timeFormatInput,
		fmt.Sprintf("%s %s", config.StartTime, localzone),
	)

	end, _ := time.Parse(
		timeFormatInput,
		fmt.Sprintf("%s %s", config.EndTime, localzone),
	)

	list, err := list.GetObjectList(start, end, config)
	totalSizeBytes := list.GetTotalByte()
	sort.Sort(list)

	// wait for user prompt
	var key string
	var ok bool
	for !ok {
		fmt.Printf("%s %s  -  %s\n",
			"From-To      \t:",
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

	err = downloader.NewDownloader(config).Download(list)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Download Completed.\n")

}
