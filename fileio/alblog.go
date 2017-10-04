package fileio

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/yomon8/aloget/config"
)

type ALBLog struct {
	cfg *config.Config
}

func NewALBLog(c *config.Config) *ALBLog {
	return &ALBLog{
		cfg: c,
	}
}

func (l *ALBLog) GetReadStream(file *os.File) (*io.Reader, error) {
	rgz, err := os.OpenFile(file.Name(), os.O_RDONLY, 0666)
	var rs io.Reader
	rs, err = gzip.NewReader(rgz)
	if err != nil {
		return nil, fmt.Errorf("failed to extract gzip, if donwload elb log, use -elb option %v", err)
	}
	return &rs, nil
}

func (l *ALBLog) GetWriteFile(key *string) (*os.File, error) {
	splitedKey := strings.Split(*key, "_")
	suffix := splitedKey[len(splitedKey)-2]
	outfile := fmt.Sprintf("%s_%s.log", l.cfg.LogPrefix, suffix)
	writeFlg := os.O_CREATE | os.O_WRONLY | os.O_APPEND
	wf, err := os.OpenFile(
		outfile,
		writeFlg,
		0666,
	)
	if err != nil {
		return nil, err
	}
	return wf, nil
}
