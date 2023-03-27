package tarot

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math/rand"
	"time"

	"github.com/disintegration/imaging"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/pkg/errors"
	"github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
	"golang.org/x/image/font"
)

type Reader struct {
	gptCli *openai.Client
	assets Assets
}

func NewReader(gptCli *openai.Client, assets Assets) (*Reader, error) {
	reader := &Reader{
		gptCli: gptCli,
		assets: assets,
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
	defaultPrompt = `我想要用塔罗牌占卜 “%s”，随后我抽中了三张牌，他们分别是：
	1. %s，
	2. %s，
	3. %s。
请模仿一位寡言少语高深莫测的塔罗牌占卜师，简单解读一下这个牌面，请从我要占卜的事情角度解读，话语精简，不用提示占卜的局限性。`
)

func (r *Reader) Prompt(cards [3]Card, thing string) string {
	p := fmt.Sprintf(defaultPrompt, thing,
		cards[0].ZhString(), cards[1].ZhString(), cards[2].ZhString())

	return p
}

func (r *Reader) Read(ctx context.Context, cards [3]Card, thing string) (string, error) {
	prompt := r.Prompt(cards, thing)
	resp, err := r.gptCli.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo0301,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)
	if err != nil {
		return "", errors.WithMessage(err, "failed to read from gpt")
	}

	return resp.Choices[0].Message.Content, nil
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
		callback(err, res)
	}()

	img, err := r.Render(cards)
	if err != nil {
		return cards, nil, errors.WithMessage(err, "failed to render img")
	}

	return cards, img, nil
}

const (
	fontSize = 22.0
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
	w := (defaultImageWidth - span*4) / 3
	var card Card
	for idx := 0; idx < 3; idx++ {
		card = cards[idx]
		pic := card.Pic
		if card.Position == PositionReversed {
			pic = imaging.Rotate(pic, 180, color.NRGBA{0, 0, 0, 0})
		}
		cen := span + (span+w)*idx + w/2
		picP := image.Pt(cen-pic.Bounds().Dx()/2, 120)

		pt := pic.Bounds().Add(picP)
		draw.Draw(img, pt, pic, image.Point{}, draw.Src)

		s := card.ZhString()
		fW, fH := r.getTextSize(s)
		wordP := freetype.Pt(cen-fW/2, 120+pic.Bounds().Dy()+fH+10)
		_, _ = c.DrawString(s, wordP)
	}

	// utils.SavePng(img, "test.png")
	return img, nil
}
