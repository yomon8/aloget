# aloget
AWS ALB(Application Load Balancer)/ELB(Classic Load Balancer) Access Log Downloader

## Usage

Set below if necessary or use `-cred` and `-r` option.

```
$ export AWS_ACCESS_KEY_ID='yourkey'
$ export AWS_SECRET_ACCESS_KEY='yoursecretkey'
$ export AWS_REGION='us-east-1'
```


```
$ aloget -o <OutputFilePrefix> -b <S3Bucket> -p <ALBAccessLogPrefix> [options]
```

## Install


```
go get github.com/yomon8/aloget/...
go install github.com/yomon8/aloget/...
```

or 
 
Download from [released file](https://github.com/yomon8/aloget/releases)

## OPTIONS

|Option|Description|Example|
|:--|:--|:--|
|-o(Required)|Output file prefix|-l /tmp/alblog|
|-b(Required)|S3 Bucket name| -b yourbucket|
|-p(Required)|S3 ALB AccessLog Prefix| -p alb-log/alb-name|
|-r|Required to set AWS Region or set env variable AWS_REGION| -r us-west-1|
|-s|Download files newer than start time (default 10 minutes ago)| -s 2017-09-28T11:59:54|
|-e|Download files older than end time (defalut now)| -e 2017-09-28T12:59:54 |
|-cred|Use default profile of credential file (Usually ~/.aws/credentials)| -cred|
|-gz|Don't decompress gzip file | -gz |
|-version|Show Version|-version|
|-utc|Recognize the datetime value of -s and -e as UTC| -utc|
|-elb|ELB(Classic Load Balancer) mode| -elb|
|-force|Don't prompt before start of downloading|-force|
