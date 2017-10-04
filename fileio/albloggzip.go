package fileio

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/yomon8/aloget/config"
)

type ALBLogGzip struct {
	cfg *config.Config
}

func NewALBLogGzip(c *config.Config) *ALBLogGzip {
	return &ALBLogGzip{
		cfg: c,
	}
}

func (l *ALBLogGzip) GetReadStream(file *os.File) (*io.Reader, error) {
	var rs io.Reader
	rs, err := os.OpenFile(file.Name(), os.O_RDONLY, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to read tmpfile ,%s, %v", file.Name(), err)
	}
	return &rs, nil

}

func (l *ALBLogGzip) GetWriteFile(key *string) (*os.File, error) {
	splitedKey := strings.Split(*key, "_")
	suffix := strings.Join(splitedKey[len(splitedKey)-4:], "_")
	outfile := fmt.Sprintf("%s_%s", l.cfg.LogPrefix, suffix)
	writeFlg := os.O_CREATE | os.O_WRONLY | os.O_TRUNC
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
