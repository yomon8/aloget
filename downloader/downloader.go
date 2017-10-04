package downloader

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/yomon8/aloget/config"
	"github.com/yomon8/aloget/entry"
	"github.com/yomon8/aloget/fileio"
	"github.com/yomon8/aloget/objects"
)

type Downloader struct {
	instance *s3manager.Downloader
	cfg      *config.Config
	files    []string
	fileio   fileio.FileIO
	buffer   []*entry.Entry
	parse    func(string) (*entry.Entry, error)
}

func NewDownloader(cfg *config.Config) *Downloader {
	d := &Downloader{
		instance: s3manager.NewDownloader(cfg.Session),
		cfg:      cfg,
		files:    make([]string, 0),
	}
	if cfg.IsELB {
		// ELB Text
		d.fileio = fileio.NewELBLog(cfg)
		d.parse = entry.ParseELBLog
	} else if cfg.PreserveGzip {
		// ALB Gzip
		d.fileio = fileio.NewALBLogGzip(cfg)
	} else {
		// ALB Text
		d.fileio = fileio.NewALBLog(cfg)
		d.parse = entry.ParseALBLog
	}
	if cfg.Stdout {
		d.buffer = make([]*entry.Entry, 0)
	}

	return d
}

func (dl *Downloader) Download(list *objects.List) error {
	for _, key := range list.GetAllKeys() {
		err := dl.downloadObject(key)
		if err != nil {
			return err
		}
	}
	if dl.cfg.Stdout {
		dl.printBuffer()
	}
	return nil
}

func (dl *Downloader) debugLog(text string) {
	if dl.cfg.Debug {
		fmt.Println(text)
	}
}

func (dl *Downloader) printBuffer() {
	var entries entry.Entries = dl.buffer
	sort.Sort(entries)
	for _, entry := range entries.GetAllEntries() {
		fmt.Println(entry.Line)
	}

}

func (dl *Downloader) addBuffer(rs *io.Reader) error {
	scanner := bufio.NewScanner(*rs)
	for scanner.Scan() {
		e, err := dl.parse(scanner.Text())
		if err != nil {
			return err
		}
		dl.buffer = append(dl.buffer, e)
	}
	if err := scanner.Err(); err != nil {
		return err
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
	dl.debugLog(fmt.Sprintf("download size\t: %d", d))

	rs, err := dl.fileio.GetReadStream(tmpfile)
	if err != nil {
		return fmt.Errorf("failed to read tmpfile, %v", err)
	}

	if dl.cfg.Stdout {
		// Print to STDOUT
		err := dl.addBuffer(rs)
		if err != nil {
			return fmt.Errorf("failed to write buffer, %v", err)
		}
		return nil
	} else {
		// Write Output
		wf, err := dl.fileio.GetWriteFile(key)
		if err != nil {
			return fmt.Errorf("failed to open logfile, %v", err)
		}
		defer wf.Close()
		dl.debugLog(fmt.Sprintf("outfile      \t: %s", wf.Name()))

		n, err := io.Copy(wf, *rs)
		if err != nil {
			return fmt.Errorf("failed to write file, %v", err)
		}
		dl.debugLog(fmt.Sprintf("write size   \t: %d", n))
		dl.debugLog(fmt.Sprintf("s3obj        \t: %s", *key))
	}
	return nil
}
