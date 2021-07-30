package usecase

import (
	"errors"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/fs"
	"math"
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
	HueHistogram50               [51]float64
	ChromaHistogram50            [51]float64
	LuminanceHistogram50         [51]float64
}

type PageAnalyzed struct {
	model.Page
	AspectRatio              float64
	AspectRatioStandardScore float64
	HueHistogram50           [51]float64
	ChromaHistogram50        [51]float64
	LuminanceHistogram50     [51]float64
}

type HCL struct {
	Hue       float64 // [0, 360]
	Chroma    float64 // [-1, 1]
	Luminance float64 // [0, 1]
}

func AnalyzeComic(src source.Source, dir string) (*ComicAnalyzed, error) {
	var (
		pages                []PageAnalyzed
		aspectRatioSum       float64
		hueHistogram50       [51]float64
		chromaHistogram50    [51]float64
		luminanceHistogram50 [51]float64

		area float64
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
		var (
			pageHueHistogram50       [51]float64
			pageChromaHistogram50    [51]float64
			pageLuminanceHistogram50 [51]float64
			pageArea                 float64
		)
		for y := bounds.Min.Y; y <= bounds.Max.Y; y++ {
			for x := bounds.Min.X; x <= bounds.Max.X; x++ {
				col := img.At(x, y)
				ful, ok := colorful.MakeColor(col)
				if !ok {
					return errors.New("cannot make colorful color")
				}
				h, c, l := ful.Hcl()
				if c > 0.0001 { // 彩度のないビットにおけるHueはブレが大きいのでHistogramに入れない
					hi := int(math.Round(h * 50 / 360))
					hueHistogram50[hi]++
					pageHueHistogram50[hi]++
				}
				ci := int(math.Min(50, math.Round(c*50)))
				chromaHistogram50[ci]++
				pageChromaHistogram50[ci]++
				li := int(math.Round(l * 50))
				luminanceHistogram50[li]++
				pageLuminanceHistogram50[li]++
				area++
				pageArea++
			}
		}
		aspectRatio := float64(config.Width) / float64(config.Height)
		aspectRatioSum += aspectRatio
		for i, v := range pageHueHistogram50 {
			pageHueHistogram50[i] = v * 10 / pageArea
		}
		for i, v := range pageChromaHistogram50 {
			pageChromaHistogram50[i] = v * 10 / pageArea
		}
		for i, v := range pageLuminanceHistogram50 {
			pageLuminanceHistogram50[i] = v * 10 / pageArea
		}
		pages = append(pages, PageAnalyzed{
			Page: model.Page{
				Rel:    strings.TrimLeft(strings.TrimPrefix(p, dir), src.Separator()),
				Config: config,
				Format: format,
			},
			AspectRatio:          aspectRatio,
			HueHistogram50:       pageHueHistogram50,
			ChromaHistogram50:    pageChromaHistogram50,
			LuminanceHistogram50: pageLuminanceHistogram50,
		})
		return nil
	}); err != nil {
		return nil, err
	}

	var (
		aspectRatioAverage  = aspectRatioSum / float64(len(pages))
		aspectRatioVariance float64
	)
	for _, p := range pages {
		d := p.AspectRatio - aspectRatioAverage
		aspectRatioVariance += d * d
	}
	aspectRatioVariance /= float64(len(pages))
	aspectRatioStandardDeviation := math.Sqrt(aspectRatioVariance)

	for i, v := range hueHistogram50 {
		hueHistogram50[i] = v * 10 / area
	}
	for i, v := range chromaHistogram50 {
		chromaHistogram50[i] = v * 10 / area
	}
	for i, v := range luminanceHistogram50 {
		luminanceHistogram50[i] = v * 10 / area
	}
	return &ComicAnalyzed{
		Pages:                        pages,
		AspectRatioAverage:           aspectRatioAverage,
		AspectRatioVariance:          aspectRatioVariance,
		AspectRatioStandardDeviation: aspectRatioStandardDeviation,
		HueHistogram50:               hueHistogram50,
		ChromaHistogram50:            chromaHistogram50,
		LuminanceHistogram50:         luminanceHistogram50,
	}, nil
}
