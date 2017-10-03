package objects

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

func (objs List) Len() int {
	return len(objs)
}

func (objs List) Swap(i, j int) {
	objs[i], objs[j] = objs[j], objs[i]
}

func (objs List) Less(i, j int) bool {
	return objs[i].LastModified.Before(*objs[j].LastModified)
}

func (objs List) GetTotalByte() int64 {
	var total int64
	for _, o := range objs {
		total += *o.Size
	}
	return total
}

func (objs List) GetOldestTime() time.Time {
	return *objs[0].LastModified
}

func (objs List) GetLatestTime() time.Time {
	return *objs[objs.Len()-1].LastModified
}

func (objs List) GetAllKeys() []*string {
	keys := make([]*string, objs.Len())
	for i, obj := range objs {
		keys[i] = obj.Key
	}
	return keys
}

func getTargetPaths(config *config.Config) map[string]string {
	target := config.StartTime
	end := config.EndTime
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

func GetObjectList(config *config.Config) (*List, error) {
	s3Objects := make([]*s3.Object, 0)
	for _, path := range getTargetPaths(config) {
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
			if o.LastModified.Sub(config.StartTime) > 0 &&
				o.LastModified.Sub(config.EndTime) < 0 {
				s3Objects = append(s3Objects, o)
			}
		}
	}
	var objs List = s3Objects
	return &objs, nil
}
