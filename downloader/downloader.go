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
	splitedKey := strings.Split(*key, "_")
	suffix := strings.Join(splitedKey[len(splitedKey)-4:], "_")
	tmpfile, err := ioutil.TempFile("", fmt.Sprintf("tmp_%s", suffix))
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
		rf       io.Reader
		outfile  string
		writeFlg int
	)

	if dl.cfg.IsELB {
		// ELB Text
		elbIP := splitedKey[len(splitedKey)-2]
		outfile = fmt.Sprintf("%s_%s.log", dl.cfg.LogPrefix, elbIP)
		rf, err = os.OpenFile(tmpfile.Name(), os.O_RDONLY, 0666)
		if err != nil {
			return fmt.Errorf("failed to read tmpfile ,%s, %v", tmpfile.Name(), err)
		}
		writeFlg = os.O_CREATE | os.O_WRONLY | os.O_APPEND
	} else if dl.cfg.PreserveGzip {
		// ALB Gzip
		outfile = fmt.Sprintf("%s_%s", dl.cfg.LogPrefix, suffix)
		rf, err = os.OpenFile(tmpfile.Name(), os.O_RDONLY, 0666)
		if err != nil {
			return fmt.Errorf("failed to read tmpfile ,%s, %v", tmpfile.Name(), err)
		}
		writeFlg = os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	} else {
		// ALB Text
		albIP := splitedKey[len(splitedKey)-2]
		outfile = fmt.Sprintf("%s_%s.log", dl.cfg.LogPrefix, albIP)
		rgz, err := os.OpenFile(tmpfile.Name(), os.O_RDONLY, 0666)
		rf, err = gzip.NewReader(rgz)
		if err != nil {
			return fmt.Errorf("failed to extract gzip, %v", err)
		}
		writeFlg = os.O_CREATE | os.O_WRONLY | os.O_APPEND
	}

	wf, err := os.OpenFile(
		outfile,
		writeFlg,
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
	fmt.Println("write  \t:", n)
	fmt.Println("s3obj  \t:", *key)
	fmt.Println("output \t:", outfile)
	return nil
}
