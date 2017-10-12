package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/yomon8/aloget/config"
)

const (
	usage = `
Usage:
  aloget -b <S3Bucket> -p <ALBAccessLogPrefix> {-o <OutputFilePrefix>|-stdout}
         [-r aws-region]
         [-s yyyy-MM-ddTHH:mm:ss] [-e yyyy-MM-ddTHH:mm:ss]
         [-duration <Minutes>]
         [-cred] [-gz|-elb] [-utc] [-force] [-debug] [-version]
`

	maxkey          = 10240
	TimeFormatParse = "2006-01-02T15:04:05 MST"
)

var (
	ErrOnlyPrintAndExit = errors.New("")
	startTimeInput      = ""
	endTimeInput        = ""
	defaultEndTime      = time.Now()
	defaultDuraion      = 10
	defaultStartTime    = defaultEndTime.Add(time.Duration(defaultDuraion) * -time.Minute)
	duration            = 0
	isVersion           = false
	isHelp              = false
	err                 error
)

func parseFlags(c *config.Config) {
	flag.StringVar(
		&c.S3Bucket,
		"b",
		"",
		"S3 Bucket",
	)

	flag.StringVar(
		&c.S3Prefix,
		"p",
		"",
		"S3 ALB AccessLog Prefix",
	)

	flag.StringVar(
		&c.LogPrefix,
		"o",
		"",
		"Output file prefix. (ex /tmp/alb)",
	)

	flag.StringVar(
		&startTimeInput,
		"s",
		defaultStartTime.Format(config.TimeFormatInput),
		"Start Time. default 10 minutes ago",
	)

	flag.StringVar(
		&endTimeInput,
		"e",
		defaultEndTime.Format(config.TimeFormatInput),
		"End Time. defalut now ",
	)

	flag.IntVar(
		&duration,
		"duration",
		defaultDuraion,
		"Duraion of log (minutes)",
	)

	flag.StringVar(
		&c.Region,
		"r",
		"",
		"AWS REGION (ex. us-west-1)",
	)

	flag.BoolVar(
		&c.IsUTC,
		"utc",
		false,
		"-s and -e as UTC",
	)

	flag.BoolVar(
		&c.PreserveGzip,
		"gz",
		false,
		"Don't decompress gzip, preserve gzip format.",
	)

	flag.BoolVar(
		&c.UseDefaultCredensial,
		"cred",
		false,
		"Use default credentials (~/.aws/credentials)",
	)

	flag.BoolVar(
		&isVersion,
		"v",
		false,
		"Show version info",
	)

	flag.BoolVar(
		&c.ForceMode,
		"force",
		false,
		"Force mode",
	)

	flag.BoolVar(
		&c.IsELB,
		"elb",
		false,
		"ELB(Classic Load Balancer) mode",
	)

	flag.BoolVar(
		&c.Debug,
		"debug",
		false,
		"Debug mode",
	)

	flag.BoolVar(
		&c.Stdout,
		"stdout",
		false,
		"Write access log to stdout.",
	)

	flag.Parse()
}

func validateOptions(c *config.Config) error {
	if isVersion {
		fmt.Println("version :", version)
		return ErrOnlyPrintAndExit
	}

	// Check Options
	if len(os.Args) == 1 || isHelp || c.S3Prefix == "" || c.S3Bucket == "" {
		fmt.Println(usage)
		flag.Usage()
		return ErrOnlyPrintAndExit
	}

	if c.LogPrefix == "" && !c.Stdout {
		fmt.Println("You should set either -o or -stdout")
		return ErrOnlyPrintAndExit
	}
	if c.LogPrefix != "" && c.Stdout {
		fmt.Println("You can only set either -o or -stdout")
		return ErrOnlyPrintAndExit
	}

	if c.IsELB && c.PreserveGzip {
		fmt.Println("-elb can't use with -gz")
		return ErrOnlyPrintAndExit
	}

	// Check Time Inputs
	zone := "UTC"
	if !c.IsUTC {
		zone, _ = time.Now().In(time.Local).Zone()
	}
	if err != nil {
		return fmt.Errorf("-s time format is %s", config.TimeFormatInput)
	}
	c.EndTime, err = time.Parse(
		TimeFormatParse,
		fmt.Sprintf("%s %s", endTimeInput, zone),
	)
	if duration == defaultDuraion {
		c.StartTime, err = time.Parse(
			TimeFormatParse,
			fmt.Sprintf("%s %s", startTimeInput, zone),
		)
	} else {
		c.StartTime = c.EndTime.Add(time.Duration(duration) * -time.Minute)
	}
	if err != nil {
		return fmt.Errorf("-e time format is %s", config.TimeFormatInput)
	}
	if c.EndTime.Sub(c.StartTime) < 0 {
		return fmt.Errorf("-s should be before -e")
	}

	if c.Stdout {
		if c.Debug {
			fmt.Println("need to set -o to use with -debug")
			return ErrOnlyPrintAndExit
		} else if c.PreserveGzip {
			fmt.Println("need to set -o to use with -gz")
			return ErrOnlyPrintAndExit
		}
		c.ForceMode = true
	}

	if c.Region == "" {
		c.Region = os.Getenv("AWS_REGION")
	}
	isValidRegion := false
	for key := range endpoints.AwsPartition().Regions() {
		if c.Region == key {
			isValidRegion = true
			break
		}
	}
	if !isValidRegion {
		if c.Region == "" {
			return fmt.Errorf("No AWS region set, use -r option or os variable AWS_REGION")
		}
		validRegion := ""
		for key := range endpoints.AwsPartition().Regions() {
			validRegion += fmt.Sprintf("%s\n", key)
		}
		return fmt.Errorf("Invalid Region set (%s),it shoud be one of follow.\n%s", c.Region, validRegion)
	}

	if c.UseDefaultCredensial {
		c.Session, err = session.NewSession(&aws.Config{
			Credentials: credentials.NewSharedCredentials("", "default"),
			Region:      aws.String(c.Region),
		})
		if err != nil {
			c.Session, err = session.NewSession(&aws.Config{
				Credentials: credentials.NewEnvCredentials(),
				Region:      aws.String(c.Region),
			})
		}
	} else {
		c.Session, err = session.NewSession(&aws.Config{
			Credentials: credentials.NewEnvCredentials(),
			Region:      aws.String(c.Region),
		})
	}
	if err != nil {
		return err
	}

	return nil
}

func loadConfig() (*config.Config, error) {
	c := new(config.Config)
	c.MaxKeyCount = maxkey
	parseFlags(c)
	err := validateOptions(c)
	if err != nil {
		return nil, err
	}
	err = c.FetchAccountID()
	if err != nil {
		return nil, err
	}
	return c, nil
}
