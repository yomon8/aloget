package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

type Config struct {
	Session      *session.Session
	LogPrefix    string
	AccountId    string
	S3Prefix     string
	S3Bucket     string
	Region       string
	IsUTC        bool
	ForceMode    bool
	NoDecompress bool
	MaxKeyCount  int64
	StartTime    string
	EndTime      string
}

const (
	maxkey          = 10240
	timeFormatInput = "2006-01-02 15:04:05"
)

func (c *Config) fetchAccountId() error {
	svc := sts.New(c.Session)
	input := &sts.GetCallerIdentityInput{}
	result, err := svc.GetCallerIdentity(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				return fmt.Errorf(aerr.Error())
			}
		} else {
			return fmt.Errorf(err.Error())
		}
	}

	c.AccountId = *result.Account
	return nil
}

func LoadConfig() (*Config, error) {
	var (
		defaultEndTime       = time.Now()
		defaultStartTime     = defaultEndTime.Add(time.Duration(10) * -time.Minute)
		useDefaultCredensial = false
		isVersion            = false
		isHelp               = false
	)

	c := new(Config)
	c.MaxKeyCount = maxkey

	flag.StringVar(
		&c.S3Bucket,
		"b",
		"",
		"[Required] S3 Bucket",
	)

	flag.StringVar(
		&c.S3Prefix,
		"p",
		"",
		"[Required] S3 Prefix",
	)

	flag.StringVar(
		&c.LogPrefix,
		"f",
		"",
		"[Required] Logfile prefix. (ex \"/tmp/alb_\")",
	)

	flag.StringVar(
		&c.StartTime,
		"s",
		defaultStartTime.Format(timeFormatInput),
		"Start Time, default is 10 minutes ago",
	)

	flag.StringVar(
		&c.EndTime,
		"e",
		defaultEndTime.Format(timeFormatInput),
		"End Time defalut is now ",
	)

	flag.StringVar(
		&c.Region,
		"r",
		"",
		"AWS REGION (ex. us-west-1)",
	)

	flag.BoolVar(
		&c.IsUTC,
		"UTC",
		false,
		"Input times are UTC",
	)

	flag.BoolVar(
		&c.NoDecompress,
		"gz",
		false,
		"not decompress gzip",
	)

	flag.BoolVar(
		&useDefaultCredensial,
		"c",
		false,
		"Use credentials file (~/.aws/credentials)",
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

	flag.Parse()
	if len(os.Args) == 1 || isHelp || c.S3Prefix == "" || c.S3Bucket == "" || c.LogPrefix == "" {
		flag.Usage()
		os.Exit(255)
	}

	if isVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	if c.Region == "" {
		c.Region = os.Getenv("AWS_REGION")
	}
	isValidRegion := false
	for key, _ := range endpoints.AwsPartition().Regions() {
		if c.Region == key {
			isValidRegion = true
			break
		}
	}
	if !isValidRegion {
		if c.Region == "" {
			return nil, fmt.Errorf("No AWS Region set, use -r option or OS variable AWS_REGION")
		} else {
			validRegion := ""
			for key, _ := range endpoints.AwsPartition().Regions() {
				validRegion += fmt.Sprintf("%s\n", key)
			}
			return nil, fmt.Errorf("Invalid Region set (%s),it shoud be one of follow.\n%s", c.Region, validRegion)
		}
	}

	var err error
	if useDefaultCredensial {
		c.Session, err = session.NewSession(&aws.Config{
			Region: aws.String(c.Region),
		})
	} else {
		c.Session, err = session.NewSession(&aws.Config{
			Credentials: credentials.NewEnvCredentials(),
			Region:      aws.String(c.Region),
		})
	}
	if err != nil {
		return nil, err
	}

	err = c.fetchAccountId()
	if err != nil {
		return nil, err
	}

	return c, nil
}
