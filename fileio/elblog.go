package fileio

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/yomon8/aloget/config"
)

type ELBLog struct {
	cfg *config.Config
}

func NewELBLog(c *config.Config) *ELBLog {
	return &ELBLog{
		cfg: c,
	}
}

func (l *ELBLog) GetReadStream(file *os.File) (*io.Reader, error) {
	var rs io.Reader
	rs, err := os.OpenFile(file.Name(), os.O_RDONLY, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to read tmpfile ,%s, %v", file.Name(), err)

	}
	return &rs, nil
}

func (l *ELBLog) GetWriteFile(key *string) (*os.File, error) {
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
