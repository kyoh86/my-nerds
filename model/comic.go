package model

import (
	"strconv"
	"strings"
)

type ComicType int

const (
	ComicTypeUnknown ComicType = iota
	ComicTypeRAR
	ComicTypeDir
)

type ComicInfo struct {
	Author string
	Title  string
	Number float64
	Volume string
	Append string
}

type Comic struct {
	Type ComicType
	Path string
	Info ComicInfo
}

func (c Comic) Rename() string {
	parts := []string{
		c.Info.Title,
	}
	if c.Info.Number != 0 {
		parts = append(parts, strconv.FormatFloat(c.Info.Number, 'f', -1, 64)+"巻")
	}
	if c.Info.Volume != "" {
		parts = append(parts, c.Info.Volume)
	}
	parts = append(parts, "["+c.Info.Author+"]")
	if c.Info.Append != "" {
		parts = append(parts, "("+c.Info.Append+")")
	}
	return strings.Join(parts, " ")
}
