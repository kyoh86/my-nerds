package usecase

import (
	"errors"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/fs"
	"strings"

	"github.com/kyoh86/my-nerds/driver/source"
	"github.com/kyoh86/my-nerds/model"
	"github.com/lucasb-eyer/go-colorful"
)

type ComicAnalyzed struct {
	Pages                        []PageAnalyzed
	AspectRatioAverage           float64
	AspectRatioVariance          float64
	AspectRatioStandardDeviation float64
}

type PageAnalyzed struct {
	model.Page
	HCLPoints                []HCL
	AspectRatio              float64
	AspectRatioStandardScore float64
}

type HCL struct {
	Hue       float64 // [0, 360]
	Chroma    float64 // [-1, 1]
	Luminance float64 // [0, 1]
}

func AnalyzeComic(src source.Source, dir string) (*ComicAnalyzed, error) {
	var (
		pages                        []PageAnalyzed
		aspectRatioSum               float64
		aspectRatioAverage           float64
		aspectRatioVariance          float64
		aspectRatioStandardDeviation float64
	)
	if err := src.Walk(dir, func(p string, info fs.FileInfo) (retErr error) {
		if info.IsDir() {
			return nil
		}
		if info.Mode()&fs.ModeSymlink == fs.ModeSymlink {
			return nil
		}
		entry, err := src.Open(p)
		if err != nil {
			return err
		}
		defer func() {
			if closeErr := entry.Close(); closeErr != nil && retErr == nil {
				retErr = closeErr
			}
		}()
		img, format, err := image.Decode(entry)
		if err != nil {
			return err
		}
		bounds := img.Bounds()
		config := image.Config{
			Width:      bounds.Dx(),
			Height:     bounds.Dy(),
			ColorModel: img.ColorModel(),
		}
		var points []HCL
		for y := bounds.Min.Y; y <= bounds.Max.Y; y++ {
			for x := bounds.Min.X; x <= bounds.Max.X; x++ {
				col := img.At(x, y)
				ful, ok := colorful.MakeColor(col)
				if !ok {
					return errors.New("cannot make colorful color")
				}
				h, c, l := ful.Hcl()
				points = append(points, HCL{
					Hue:       h,
					Chroma:    c,
					Luminance: l,
				})
			}
		}
		aspectRatio := float64(config.Width) / float64(config.Height)
		aspectRatioSum += aspectRatio
		pages = append(pages, PageAnalyzed{
			Page: model.Page{
				Rel:    strings.TrimLeft(strings.TrimPrefix(p, dir), src.Separator()),
				Config: config,
				Format: format,
			},
			AspectRatio: aspectRatio,
			HCLPoints:   points,
		})
		return nil
	}); err != nil {
		return nil, err
	}
	return &ComicAnalyzed{
		Pages: pages,
	}, nil
}
