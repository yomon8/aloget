package downloader

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/yomon8/aloget/config"
	"github.com/yomon8/aloget/list"
)

type Downloader struct {
	instance *s3manager.Downloader
	cfg      *config.Config
	files    []string
}

func NewDownloader(cfg *config.Config) *Downloader {
	return &Downloader{
		instance: s3manager.NewDownloader(cfg.Session),
		cfg:      cfg,
		files:    make([]string, 0),
	}
}

func (dl *Downloader) Download(list *list.List) error {
	for _, key := range list.GetAllKeys() {
		err := dl.downloadObject(key)
		if err != nil {
			return err
		}
	}
	return nil
}

func (dl *Downloader) downloadObject(key *string) error {
	tmpfile, err := ioutil.TempFile("", fmt.Sprintf("tmp_%s", os.Executable))
	if err != nil {
		return err
	}
	defer os.Remove(tmpfile.Name())

	d, err := dl.instance.Download(
		tmpfile,
		&s3.GetObjectInput{
			Bucket: aws.String(dl.cfg.S3Bucket),
			Key:    key,
		})
	if err != nil {
		return fmt.Errorf("failed to download file, %v", err)
	}
	fmt.Println("download:", d)

	var (
		rf      io.Reader
		outfile string
	)

	splitedKey := strings.Split(*key, "_")
	if dl.cfg.NoDecompress {
		suffix := splitedKey[len(splitedKey)-4]
		outfile = fmt.Sprintf("%s_%s", dl.cfg.LogPrefix, suffix)
		rf, err = os.OpenFile(tmpfile.Name(), os.O_RDONLY, 0666)
		if err != nil {
			return fmt.Errorf("failed to read tmpfile ,%s, %v", tmpfile.Name(), err)
		}
	} else {
		albIP := splitedKey[len(splitedKey)-2]
		outfile = fmt.Sprintf("%s_%s.log", dl.cfg.LogPrefix, albIP)
		rgz, err := os.OpenFile(tmpfile.Name(), os.O_RDONLY, 0666)
		rf, err = gzip.NewReader(rgz)
		if err != nil {
			return fmt.Errorf("failed to extract gzip, %v", err)
		}
	}

	wf, err := os.OpenFile(
		outfile,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666,
	)
	if err != nil {
		return fmt.Errorf("failed to open logfile, %v", err)
	}
	defer wf.Close()

	n, err := io.Copy(wf, rf)
	if err != nil {
		return fmt.Errorf("failed to write file, %v", err)
	}
	fmt.Println("write:", n)
	fmt.Println("s3obj:", *key)
	fmt.Println("outfile:", outfile)
	return nil
}
