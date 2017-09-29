package list

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/yomon8/aloget/config"
)

const (
	timeFormatObjectPath = "2006/01/02"
)

type List []*s3.Object

func (list List) Len() int {
	return len(list)
}

func (list List) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func (list List) Less(i, j int) bool {
	return list[i].LastModified.Before(*list[j].LastModified)
}

func (list List) GetTotalByte() int64 {
	var total int64
	for _, o := range list {
		total += *o.Size
	}
	return total
}

func (list List) GetOldestTime() time.Time {
	return *list[0].LastModified
}

func (list List) GetLatestTime() time.Time {
	return *list[list.Len()-1].LastModified
}

func (list List) GetAllKeys() []*string {
	keys := make([]*string, 0)
	for _, obj := range list {
		keys = append(keys, obj.Key)
	}
	return keys
}

func getTargetPaths(target, end time.Time, config *config.Config) map[string]string {
	targetPaths := make(map[string]string, 10)
	for i := 0; end.Sub(target) > 0; i++ {
		datekey := target.In(time.UTC).Format(timeFormatObjectPath)
		if targetPaths[datekey] == "" {
			targetPaths[datekey] = fmt.Sprintf("%s/AWSLogs/%s/elasticloadbalancing/%s/%s",
				config.S3Prefix,
				config.AccountID,
				config.Region,
				datekey,
			)
		}
		target = target.Add(time.Duration(i) * time.Hour)
	}
	return targetPaths
}

func GetObjectList(start, end time.Time, config *config.Config) (*List, error) {
	s3Objects := make([]*s3.Object, 0)
	for _, path := range getTargetPaths(start, end, config) {
		input := &s3.ListObjectsInput{
			Bucket:  aws.String(config.S3Bucket),
			MaxKeys: aws.Int64(config.MaxKeyCount),
			Prefix:  aws.String(path),
		}
		result, err := s3.New(config.Session).ListObjects(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case s3.ErrCodeNoSuchBucket:
					return nil, fmt.Errorf(
						"s3 bucket not found. %#v %#v\n",
						s3.ErrCodeNoSuchBucket,
						aerr.Error(),
					)
				default:
					return nil, fmt.Errorf(
						"get s3 list error. %#v\n",
						aerr.Error(),
					)
				}
			} else {
				return nil, fmt.Errorf(
					"get s3 list error. %#v\n",
					err,
				)
			}
		}
		for _, o := range result.Contents {
			if o.LastModified.Sub(start) > 0 && o.LastModified.Sub(end) < 0 {
				s3Objects = append(s3Objects, o)
			}
		}
	}
	var list List = s3Objects
	return &list, nil
}
