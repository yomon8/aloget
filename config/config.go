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
	AccountID    string
	S3Prefix     string
	S3Bucket     string
	Region       string
	IsUTC        bool
	ForceMode    bool
	PreserveGzip bool
	MaxKeyCount  int64
	StartTime    string
	EndTime      string
	IsELB        bool
}

var (
	version = "0"
)

const (
	maxkey          = 10240
	timeFormatInput = "2006-01-02 15:04:05"
)

func (c *Config) fetchAccountID() error {
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

	c.AccountID = *result.Account
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
		"[Required] S3 ALB AccessLog Prefix",
	)

	flag.StringVar(
		&c.LogPrefix,
		"o",
		"",
		"[Required] Output file prefix. (ex \"/tmp/alb\")",
	)

	flag.StringVar(
		&c.StartTime,
		"s",
		defaultStartTime.Format(timeFormatInput),
		"Start Time. default 10 minutes ago",
	)

	flag.StringVar(
		&c.EndTime,
		"e",
		defaultEndTime.Format(timeFormatInput),
		"End Time. defalut now ",
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
		&useDefaultCredensial,
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

	flag.Parse()

	if isVersion {
		fmt.Println("version :", version)
		os.Exit(0)
	}

	if len(os.Args) == 1 || isHelp || c.S3Prefix == "" || c.S3Bucket == "" || c.LogPrefix == "" {
		fmt.Println("Command Line:")
		fmt.Println("aloget -o <OutputFilePrefix> -b <S3Bucket> -p <ALBAccessLogPrefix> [options]\n")
		flag.Usage()
		os.Exit(255)
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
			return nil, fmt.Errorf("No AWS Region set, use -r option or OS variable AWS_REGION")
		}
		validRegion := ""
		for key := range endpoints.AwsPartition().Regions() {
			validRegion += fmt.Sprintf("%s\n", key)
		}
		return nil, fmt.Errorf("Invalid Region set (%s),it shoud be one of follow.\n%s", c.Region, validRegion)
	}

	var err error
	if useDefaultCredensial {
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
		return nil, err
	}

	err = c.fetchAccountID()
	if err != nil {
		return nil, err
	}

	return c, nil
}
