package config

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

const (
	TimeFormatInput = "2006-01-02T15:04:05"
)

type Config struct {
	Session              *session.Session
	LogPrefix            string
	AccountID            string
	S3Prefix             string
	S3Bucket             string
	Region               string
	IsUTC                bool
	ForceMode            bool
	PreserveGzip         bool
	MaxKeyCount          int64
	StartTime            time.Time
	EndTime              time.Time
	IsELB                bool
	UseDefaultCredensial bool
	Debug                bool
	Stdout               bool
}

func (c *Config) FetchAccountID() error {
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
