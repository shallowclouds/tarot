package tarot

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io/fs"
	"math"
	"sync"

	"github.com/disintegration/imaging"
	"github.com/golang/freetype/truetype"
)

const (
	defaultImageWidth  = 1200
	defaultImageHeight = 720
)

type Position string

const (
	PositionUpright  = "Upright"
	PositionReversed = "Reversed"
)

func (p *Position) ZhName() string {
	if *p == PositionUpright {
		return "正位"
	}
	return "逆位"
}

type Card struct {
	Name     string      `json:"name"`
	ZhName   string      `json:"zh_name"`
	Position Position    `json:"-"`
	Pic      image.Image `json:"-"`
}

func (c *Card) ZhString() string {
	return fmt.Sprintf("%s（%s）", c.ZhName, c.Position.ZhName())
}

type Assets struct {
	Cards         []Card
	BackgroundImg image.Image
	Font          *truetype.Font
}

var (
	//go:embed assets/*
	assetsFS embed.FS

	initAssetsOnce sync.Once
	assets         Assets
)

func mustReadImg(p string) image.Image {
	data, err := fs.ReadFile(assetsFS, p)
	if err != nil {
		panic(err)
	}

	pic, err := jpeg.Decode(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}

	return pic
}

func imageTypeToRGBA64(m image.Image) *image.RGBA64 {
	bounds := m.Bounds()
	dx := bounds.Dx()
	dy := bounds.Dy()
	img := image.NewRGBA64(bounds)
	for x := 0; x < dx; x++ {
		for y := 0; y < dy; y++ {
			colorRgb := m.At(x, y)
			r, g, b, a := colorRgb.RGBA()
			rr := uint16(r)
			gg := uint16(g)
			bb := uint16(b)
			aa := uint16(a)
			img.SetRGBA64(x, y, color.RGBA64{
				R: rr,
				G: gg,
				B: bb,
				A: aa,
			})
		}
	}
	return img
}

func initCards() {
	data, err := fs.ReadFile(assetsFS, "assets/cards.json")
	if err != nil {
		panic(err)
	}

	var cards []Card

	if err := json.Unmarshal(data, &cards); err != nil {
		panic(err)
	}

	for idx, card := range cards {
		card.Pic = mustReadImg(fmt.Sprintf("assets/%d.jpg", idx))
		assets.Cards = append(assets.Cards, card)
	}
}

func processReaderImg(pic image.Image) image.Image {
	pic = imaging.Resize(pic, 0, defaultImageHeight, imaging.Lanczos)
	rgb64 := imageTypeToRGBA64(pic)
	bounds := rgb64.Bounds()
	dx := bounds.Dx()
	dy := bounds.Dy()
	img := image.NewNRGBA64(bounds)
	for x := 0; x < dx; x++ {
		for y := 0; y < dy; y++ {
			colorRgb := rgb64.At(x, y)
			r, g, b, a := colorRgb.RGBA()
			percent := 1.0 - math.Abs(float64(x)-float64(dx)/2.0)/(float64(dx)/2.0)
			percent = percent * percent

			rr, gg, bb, aa := img.ColorModel().Convert(color.NRGBA64{
				R: uint16(r),
				G: uint16(g),
				B: uint16(b),
				A: uint16(float64(a) * percent),
			}).RGBA()

			img.SetRGBA64(x, y, color.RGBA64{
				R: uint16(rr),
				G: uint16(gg),
				B: uint16(bb),
				A: uint16(aa),
			})
		}
	}

	return img
}

func initBgImg() {
	img := image.NewNRGBA(image.Rect(0, 0, 1200, 720))
	draw.Draw(img, img.Bounds(), image.Black, image.Point{}, draw.Src)

	readerImg := mustReadImg("assets/reader.jpg")
	readerImg = processReaderImg(readerImg)
	draw.Draw(img, readerImg.Bounds().Add(image.Pt((defaultImageWidth-readerImg.Bounds().Dx())/2, 0)), readerImg, image.Point{}, draw.Over)

	assets.BackgroundImg = img
}

func initFonts() {
	data, err := fs.ReadFile(assetsFS, "assets/arplmingu20lt.ttf")
	if err != nil {
		panic(err)
	}

	font, err := truetype.Parse(data)
	if err != nil {
		panic(err)
	}

	assets.Font = font
}

func GetDefaultAssets() Assets {
	initAssetsOnce.Do(func() {
		initCards()
		initBgImg()
		initFonts()
	})

	return assets
}
