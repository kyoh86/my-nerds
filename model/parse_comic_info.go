package model

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	titlePattern  = `(?P<Title>.+)`
	numberPattern = ` (?P<Number>\d+(?:\.\d+)?)巻`
	volumePattern = `(?: (?P<Volume>.+))?`
	authorPattern = ` \[(?P<Author>[^\]]+)\]`
	appendPattern = `(?: \((?P<Append>[^\)]+)\))?`
)

func ParseComicInfo(name string) (info *ComicInfo, _ error) {
	for _, pattern := range [][]string{
		{titlePattern, numberPattern, volumePattern, authorPattern, appendPattern},
		{titlePattern, authorPattern, appendPattern},
	} {
		reg := regexp.MustCompile("^" + strings.Join(pattern, "") + "$")
		subs := reg.FindStringSubmatch(name)
		if len(subs) == 0 {
			continue
		}

		info = new(ComicInfo)
		for i, name := range reg.SubexpNames() {
			if subs[i] == "" {
				continue
			}
			switch name {
			case "Author":
				info.Author = strings.TrimSpace(subs[i])
			case "Title":
				info.Title = strings.TrimSpace(subs[i])
			case "Volume":
				info.Volume = strings.TrimSpace(subs[i])
			case "Append":
				info.Append = subs[i]
			case "Number":
				num, err := strconv.ParseFloat(subs[i], 64)
				if err != nil {
					return nil, err
				}
				info.Number = num
			}
		}
		break
	}
	if info == nil {
		return nil, fmt.Errorf("invalid comic name %q", name)
	}
	if info.Author == "" {
		return nil, fmt.Errorf("invalid comic name %q: no author", name)
	}
	if info.Title == "" {
		return nil, fmt.Errorf("invalid comic name %q: no title", name)
	}
	return
}
