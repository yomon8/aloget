# aloget
AWS ALB Access Log Downloader

## Usage

Set below if necessary or use `-cred` and `-r` option

```
$ export AWS_ACCESS_KEY_ID='yourkey'
$ export AWS_SECRET_ACCESS_KEY='yoursecretkey'
$ export AWS_REGION='us-east-1'
```


```
$ aloget -l <OutputFilePrefix> -b <S3Bucket> -p <ALBAccessLogPrefix> [options]
```

## OPTIONS

|Option|Description|Example|
|:--|:--|:--|
|-l(Required)|Output file prefix|-l /tmp/alblog|
|-b(Required)|S3 Bucket name| -b yourbucket|
|-p(Required)|S3 ALB AccessLog Prefix| -p alb-log/alb-name|
|-r|Required to set AWS Region or set env variable AWS_REGION| -r us-west-1|
|-s|Download files newer than start time (default 10 minutes ago)| -s "2017-09-28 11:59:54"|
|-e|Download files older than end time (defalut now)| -e "2017-09-28 12:59:54" |
|-cred|Use credential file (Usually ~/.aws/credentials)| -cred|
|-gz|Don't decompress gzip file | -gz |
|-version|Show Version|-version|
|-utc|Recognize the datetime value of -s and -e as UTC| -utc|
|-force|Don't prompt before start of downloading|-force|
