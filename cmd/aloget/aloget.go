package main

import (
	"fmt"
	"os"
	"sort"
	"time"

	humanize "github.com/dustin/go-humanize"

	"github.com/yomon8/aloget/config"
	"github.com/yomon8/aloget/downloader"
	"github.com/yomon8/aloget/objects"
)

func main() {
	cfg, err := config.LoadConfig()
	if err == config.ErrOnlyPrintAndExit {
		os.Exit(255)
	} else if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	list, err := objects.GetObjectList(cfg)
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
	if !cfg.ForceMode {
		var key string
		var ok bool
		for !ok {
			fmt.Printf("%s %s  -  %s\n",
				"From-To(Local) \t:",
				list.GetOldestTime().In(time.Local).Format(config.TimeFormatParse),
				list.GetLatestTime().In(time.Local).Format(config.TimeFormatParse),
			)
			fmt.Printf("%s %s  -  %s\n",
				"From-To(UTC)   \t:",
				list.GetOldestTime().Format(config.TimeFormatParse),
				list.GetLatestTime().Format(config.TimeFormatParse),
			)
			fmt.Printf("%s %s\n",
				"Download Size  \t:",
				humanize.Bytes(uint64(totalSizeBytes)),
			)
			fmt.Printf("%s %s\n",
				"Decompress Gzip\t:",
				fmt.Sprint(!cfg.PreserveGzip),
			)
			fmt.Printf("%s %d objects\n",
				"S3 Objects    \t:",
				list.Len(),
			)
			fmt.Print("Start/Cancel?>")
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

	err = downloader.NewDownloader(cfg).Download(list)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if !cfg.Stdout {
		fmt.Printf("Download Completed.\n")
	}
}
