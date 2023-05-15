package tarot

import (
	"context"
	"image"
	"image/color"
	"image/draw"
	"math/rand"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/image/font"
)

type Reader struct {
	gptReader GPTReader
	assets    Assets
}

type GPTReader interface {
	Chat(ctx context.Context, systemMsg, userMsg string) (string, error)
}

func NewReader(gptReader GPTReader, assets Assets) (*Reader, error) {
	reader := &Reader{
		gptReader: gptReader,
		assets:    assets,
	}

	return reader, nil
}

func (r *Reader) Choose() ([3]Card, error) {
	cards := [3]Card{}
	ps := rand.Perm(len(r.assets.Cards))
	for idx := 0; idx < 3; idx++ {
		cards[idx] = r.assets.Cards[ps[idx]]
		cards[idx].Position = PositionUpright
		if rand.Int()%100 > 60 {
			cards[idx].Position = PositionReversed
		}
	}

	return cards, nil
}

const (
	defaultPrompt = `
假如你是一位神秘的塔罗牌占卜师，我想要占卜事情是 “{{thing}}”，我抽到的三张牌分别是：{{card1}}，{{card2}}，{{card3}}。
请根据三张牌面和这件具体事情进行解读，语言简练精辟客观，不准使用“虽然...但是...”这样模棱两可的话，千万不要建议我或者安慰我，不要提醒我占卜的局限性或者意义。`
)

func (r *Reader) sanitizeThing(thing string) string {
	// Replace puncuation to space.
	for _, remove := range []string{"“", "”", "\""} {
		thing = strings.ReplaceAll(thing, remove, " ")
	}

	return thing
}

func (r *Reader) Prompt(cards [3]Card, thing string) string {
	p := defaultPrompt
	fills := map[string]string{
		"thing": r.sanitizeThing(thing),
		"card1": cards[0].ZhString(),
		"card2": cards[1].ZhString(),
		"card3": cards[2].ZhString(),
	}
	for k, v := range fills {
		p = strings.ReplaceAll(p, "{{"+k+"}}", v)
	}

	return p
}

func (r *Reader) Read(ctx context.Context, cards [3]Card, thing string) (string, error) {
	prompt := r.Prompt(cards, thing)
	resp, err := r.gptReader.Chat(ctx, "", prompt)
	if err != nil {
		return "", errors.WithMessage(err, "failed to read from gpt")
	}

	return resp, nil
}

func (r *Reader) Divine(ctx context.Context, thing string, callback func(err error, res string)) ([3]Card, image.Image, error) {
	cards, _ := r.Choose()
	logger := logrus.StandardLogger()
	logger.Infof("chosen cards: %s, %s, %s", cards[0].ZhString(), cards[1].ZhString(), cards[2].ZhString())

	go func() {
		logger.Infof("start to call chat gpt reader")
		start := time.Now()
		res, err := r.Read(ctx, cards, thing)
		logger.Infof("call gpt reader cost %.2f s", time.Since(start).Seconds())
		if err != nil {
			logger.WithError(err).Warn("failed to call chat gpt")
		}
		if callback != nil {
			callback(err, res)
		}
	}()

	img, err := r.Render(cards)
	if err != nil {
		return cards, nil, errors.WithMessage(err, "failed to render img")
	}

	return cards, img, nil
}

const (
	fontSize = 26.0
)

func (r *Reader) getTextSize(s string) (int, int) {
	face := truetype.NewFace(r.assets.Font, &truetype.Options{
		Size: fontSize,
	})
	width := font.MeasureString(face, s).Ceil()
	height := face.Metrics().Ascent.Ceil()

	return width, height
}

func (r *Reader) Render(cards [3]Card) (image.Image, error) {
	img := image.NewNRGBA64(image.Rect(0, 0, defaultImageWidth, defaultImageHeight))
	draw.Draw(img, r.assets.BackgroundImg.Bounds(), r.assets.BackgroundImg, image.Point{}, draw.Src)

	c := freetype.NewContext()
	c.SetFont(r.assets.Font)
	c.SetDst(img)
	c.SetClip(img.Bounds())
	c.SetSrc(image.White)
	c.SetFontSize(fontSize)

	span := 100
	startH := 120
	w := (defaultImageWidth - span*4) / 3
	var card Card
	for idx := 0; idx < 3; idx++ {
		card = cards[idx]
		pic := card.Pic
		if card.Position == PositionReversed {
			pic = imaging.Rotate(pic, 180, color.NRGBA{0, 0, 0, 0})
		}
		cen := span + (span+w)*idx + w/2
		picP := image.Pt(cen-pic.Bounds().Dx()/2, startH)

		pt := pic.Bounds().Add(picP)
		draw.Draw(img, pt, pic, image.Point{}, draw.Src)

		s := card.ZhString()
		fW, fH := r.getTextSize(s)
		wordP := freetype.Pt(cen-fW/2, startH-fH)
		_, _ = c.DrawString(s, wordP)
	}

	return img, nil
}
