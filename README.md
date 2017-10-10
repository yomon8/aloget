aloget
====
[![Build Status](https://travis-ci.org/yomon8/aloget.svg?branch=master)](https://travis-ci.org/yomon8/aloget)
[![Latest Version](http://img.shields.io/github/release/yomon8/aloget.svg?style=flat-square)](https://github.com/yomon8/aloget/releases)

AWS ALB(Application Load Balancer)/ELB(Classic Load Balancer) Access Log Downloader

## Usage

Set below if necessary or use `-cred` and `-r` option.

```
$ export AWS_ACCESS_KEY_ID='yourkey'
$ export AWS_SECRET_ACCESS_KEY='yoursecretkey'
$ export AWS_REGION='us-east-1'
```


```
#-- output to stdout
$ aloget -b <S3Bucket> -p <ALBAccessLogPrefix> -stdout [options] 

#-- output to file
$ aloget -o <OutputFilePrefix> -b <S3Bucket> -p <ALBAccessLogPrefix> [options] 
```

## Install


```
go get github.com/yomon8/aloget/...
go install github.com/yomon8/aloget/...
```

or 
 
Download from [released file](https://github.com/yomon8/aloget/releases)

## Arguments

```
Usage:
  aloget -b <S3Bucket> -p <ALBAccessLogPrefix> {-o <OutputFilePrefix>|-stdout}
         [-s yyyy-MM-ddTHH:mm:ss] [-e yyyy-MM-ddTHH:mm:ss]
         [-r aws-region]
         [-cred] [-gz|-elb] [-utc] [-force] [-debug] [-version]
```

|Arguments|Description|Example|
|:--|:--|:--|
|-b|S3 Bucket name| -b yourbucket|
|-p|S3 ALB AccessLog Prefix| -p alb-log/alb-name|
|-o|Output file prefix,if provided no value,set output to STDOUT|-l /tmp/alblog|
|-stdout|Write access log to stdout|-stdout|
|-r|Required to set AWS Region or set env variable AWS_REGION| -r us-west-1|
|-s|Download files newer than start time (default 10 minutes ago)| -s 2017-09-28T11:59:54|
|-e|Download files older than end time (defalut now)| -e 2017-09-28T12:59:54 |
|-cred|Use default profile of credential file (Usually ~/.aws/credentials)| -cred|
|-gz|Don't decompress gzip file | -gz |
|-version|Show Version|-version|
|-utc|Recognize the datetime value of -s and -e as UTC| -utc|
|-elb|ELB(Classic Load Balancer) mode| -elb|
|-force|Don't prompt before start of downloading|-force|
|-debug|Print debug message|-debug|
