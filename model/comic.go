package model

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type ComicType int

const (
	ComicTypeUnknown ComicType = iota
	ComicTypeRAR
	ComicTypeDir
)

type Comic struct {
	Author string
	Title  string
	Volume string
	Append string
	Number float64
}

var (
	fuzzyAuthorPattern = `[\[［](?P<Author>[^\]］]+)[\]］]`
	fuzzyTitlePattern  = `(?: *(?P<Title>[^\[［\(（]+))`
	fuzzyNumberPattern = `(?: 第?(?P<Number>\d+(?:\.\d+)?)[巻集])`
	fuzzyVolumePattern = `(?: (?P<Volume>[^\[［\(（]+))`
	fuzzyAppendPattern = ` *(?:[\(（【](?P<Append>[^\)）】]+)[\)）】])* *`
)

func ParseComicFuzzy(name string) (comic *Comic, _ error) {
	for _, pattern := range []string{
		strings.Join([]string{
			"^",
			fuzzyAppendPattern,
			fuzzyAuthorPattern,
			fuzzyAppendPattern,
			fuzzyTitlePattern,
			fuzzyNumberPattern,
			fuzzyVolumePattern + "?",
			fuzzyAppendPattern,
			"$",
		}, ""),
		strings.Join([]string{
			"^",
			fuzzyAppendPattern,
			fuzzyTitlePattern,
			fuzzyNumberPattern,
			fuzzyVolumePattern + "?",
			fuzzyAppendPattern,
			fuzzyAuthorPattern,
			fuzzyAppendPattern,
			"$",
		}, ""),
		strings.Join([]string{
			"^",
			fuzzyAppendPattern,
			fuzzyAuthorPattern,
			fuzzyAppendPattern,
			fuzzyTitlePattern,
			fuzzyAppendPattern,
			"$",
		}, ""),
		strings.Join([]string{
			"^",
			fuzzyAppendPattern,
			fuzzyTitlePattern,
			fuzzyAppendPattern,
			fuzzyAuthorPattern,
			fuzzyAppendPattern,
			"$",
		}, ""),
	} {
		reg := regexp.MustCompile(pattern)
		subs := reg.FindStringSubmatch(name)
		if len(subs) == 0 {
			continue
		}

		comic = new(Comic)
		for i, name := range reg.SubexpNames() {
			if subs[i] == "" {
				continue
			}
			switch name {
			case "Author":
				comic.Author = strings.TrimSpace(subs[i])
			case "Title":
				comic.Title = strings.TrimSpace(subs[i])
			case "Volume":
				comic.Volume = strings.TrimSpace(subs[i])
			case "Append":
				if comic.Append == "" {
					comic.Append = subs[i]
				} else {
					comic.Append = comic.Append + " " + subs[i]
				}
			case "Number":
				num, err := strconv.ParseFloat(subs[i], 64)
				if err != nil {
					return nil, err
				}
				comic.Number = num
			}
		}
		break
	}
	if comic == nil {
		return nil, fmt.Errorf("invalid comic: %s", name)
	}
	return
}

var (
	titlePattern  = `(?P<Title>.+)`
	numberPattern = ` (?P<Number>\d+(?:\.\d+)?)巻`
	volumePattern = `(?: (?P<Volume>.+))?`
	authorPattern = ` \[(?P<Author>[^\]]+)\]`
	appendPattern = `(?: \((?P<Append>[^\)]+)\))?`
)

func ParseComicName(name string) (comic *Comic, _ error) {
	for _, pattern := range [][]string{
		{titlePattern, numberPattern, volumePattern, authorPattern, appendPattern},
		{titlePattern, authorPattern, appendPattern},
	} {
		reg := regexp.MustCompile("^" + strings.Join(pattern, "") + "$")
		subs := reg.FindStringSubmatch(name)
		if len(subs) == 0 {
			continue
		}

		comic = new(Comic)
		for i, name := range reg.SubexpNames() {
			if subs[i] == "" {
				continue
			}
			switch name {
			case "Author":
				comic.Author = strings.TrimSpace(subs[i])
			case "Title":
				comic.Title = strings.TrimSpace(subs[i])
			case "Volume":
				comic.Volume = strings.TrimSpace(subs[i])
			case "Append":
				comic.Append = subs[i]
			case "Number":
				num, err := strconv.ParseFloat(subs[i], 64)
				if err != nil {
					return nil, err
				}
				comic.Number = num
			}
		}
		break
	}
	if comic.Author == "" {
		return nil, fmt.Errorf("invalid comic name %q: no author", name)
	}
	if comic.Title == "" {
		return nil, fmt.Errorf("invalid comic name %q: no title", name)
	}
	return
}

func (c Comic) String() string {
	parts := []string{
		c.Title,
	}
	if c.Number != 0 {
		parts = append(parts, strconv.FormatFloat(c.Number, 'f', -1, 64)+"巻")
	}
	if c.Volume != "" {
		parts = append(parts, c.Volume)
	}
	parts = append(parts, "["+c.Author+"]")
	if c.Append != "" {
		parts = append(parts, "("+c.Append+")")
	}
	return strings.Join(parts, " ")
}
