package util

import (
	"github.com/go-ego/gse"
)

var seg gse.Segmenter

func SegWord(src string) []string {
	return seg.Slice(src, true)
}

// TODO: add assets path
func InitSegWord() error {
	if err := seg.LoadStop(); err != nil {
		return err
	}
	if err := seg.LoadDict(); err != nil {
		return err
	}
	return nil
}
