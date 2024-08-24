package util

import (
	"github.com/go-ego/gse"
)

var seg gse.Segmenter

func Fenci(src string) []string {
	return seg.Slice(src, true)
}

// TODO: add assets path
func InitFenci() error {
	if err := seg.LoadStop(); err != nil {
		return err
	}
	if err := seg.LoadDict(); err != nil {
		return err
	}
	return nil
}
