package tarot

import (
	"context"
	"image"
	"image/color"
	"math/rand"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type Reader struct {
	gptReader             GPTReader
	assets                Assets
	sysPrompt, userPrompt string
}

type GPTReader interface {
	Chat(ctx context.Context, systemMsg, userMsg string) (string, error)
}

func NewReader(gptReader GPTReader, userPrompt, systemPrompt string, assets Assets) (*Reader, error) {
	reader := &Reader{
		gptReader:  gptReader,
		assets:     assets,
		sysPrompt:  systemPrompt,
		userPrompt: userPrompt,
	}

	if len(reader.sysPrompt) == 0 {
		reader.sysPrompt = defaultSystemPrompt
	}
	if len(reader.userPrompt) == 0 {
		reader.userPrompt = defaultUserPrompt
	}

	return reader, nil
}

func (r *Reader) Choose() ([3]Card, error) {
	cards := [3]Card{}
	ps := rand.Perm(len(r.assets.Cards))
	for idx := 0; idx < 3; idx++ {
		cards[idx] = r.assets.Cards[ps[idx]]
		cards[idx].Position = PositionUpright
		if rand.Int()%100 > 50 {
			cards[idx].Position = PositionReversed
		}
	}

	return cards, nil
}

const (
	defaultSystemPrompt = `你是一位神秘的塔罗牌占卜师，请根据三张牌面和这件具体事情进行客观解读，语言简练精辟客观，不准使用“虽然...但是...”这样模棱两可的话，不用提醒用户占卜的局限性或者意义。`
	defaultUserPrompt   = `我抽到的三张塔罗牌分别是：{{card1}}，{{card2}}，{{card3}}。我想占卜的事情是：“{{thing}}”，请解读。`
	defaultPrompt       = `
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

func (r *Reader) Prompt(cards [3]Card, thing, template string) string {
	p := template
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
	sysPrompt := r.Prompt(cards, thing, r.sysPrompt)
	userPrompt := r.Prompt(cards, thing, r.userPrompt)
	resp, err := r.gptReader.Chat(ctx, sysPrompt, userPrompt)
	if err != nil {
		return "", errors.WithMessage(err, "failed to read from gpt")
	}

	return resp, nil
}

func (r *Reader) DivineSync(ctx context.Context, thing string) ([3]Card, image.Image, string, error) {
	cards, _ := r.Choose()
	logger := logrus.StandardLogger()
	logger.Infof("chosen cards: %s, %s, %s", cards[0].ZhString(), cards[1].ZhString(), cards[2].ZhString())

	logger.Infof("start to call chat gpt reader")
	start := time.Now()
	res, err := r.Read(ctx, cards, thing)
	logger.Infof("call gpt reader cost %.2f s", time.Since(start).Seconds())
	if err != nil {
		logger.WithError(err).Warn("failed to call chat gpt")
		return cards, nil, "", errors.WithMessage(err, "failed to read cards")
	}

	img, err := r.Render(cards, thing, res)
	if err != nil {
		return cards, nil, "", errors.WithMessage(err, "failed to render img")
	}

	return cards, img, res, nil
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

	img, err := r.Render(cards, thing, "")
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

func wrapText(text string, maxWidth int, face font.Face) []string {
	var lines []string
	var line string
	var lineWidth fixed.Int26_6
	spaceRealWidth := font.MeasureString(face, " ")
	spaceWidth := font.MeasureString(face, "永")
	spaceTime := spaceWidth / spaceRealWidth
	tabWidth := spaceWidth * 2
	tabTime := spaceTime * 2
	space := strings.Repeat(" ", int(spaceTime))
	tab := strings.Repeat(" ", int(tabTime))

	for _, word := range strings.Split(text, "\n") {
		for _, r := range word {
			advance := font.MeasureString(face, string(r))
			if r == '\n' {
				lines = append(lines, line)
				line = ""
				lineWidth = 0
				continue
			} else if r == ' ' {
				if lineWidth+spaceWidth > fixed.I(maxWidth) {
					// start new line
					lines = append(lines, line)
					line = ""
					lineWidth = 0
				} else {
					line += space
					lineWidth += spaceWidth
				}
				continue
			} else if r == '\t' {
				if lineWidth+tabWidth > fixed.I(maxWidth) {
					// start new line
					lines = append(lines, line)
					line = ""
					lineWidth = 0
				} else {
					line += tab
					lineWidth += tabWidth
				}
				continue
			}
			rWidth := font.MeasureString(face, string(r))
			if lineWidth+rWidth > fixed.I(maxWidth) {
				// start new line
				lines = append(lines, line)
				line = ""
				lineWidth = 0
			}
			line += string(r)
			lineWidth += advance
		}
		lines = append(lines, line)
		line = ""
		lineWidth = 0
	}

	return lines
}

// DrawStringWrapped word-wraps the specified string to the given max width
// and then draws it at the specified anchor point using the given line
// spacing and text alignment.
func DrawStringWrapped(dc *gg.Context, ff font.Face, s string, x, y, ax, ay, width, lineSpacing float64, align gg.Align) {
	lines := wrapText(s, int(width), ff)

	// sync h formula with MeasureMultilineString
	h := float64(len(lines)) * dc.FontHeight() * lineSpacing
	h -= (lineSpacing - 1) * dc.FontHeight()

	x -= ax * width
	y -= ay * h
	switch align {
	case gg.AlignLeft:
		ax = 0
	case gg.AlignCenter:
		ax = 0.5
		x += width / 2
	case gg.AlignRight:
		ax = 1
		x += width
	}
	ay = 1
	for _, line := range lines {
		dc.DrawStringAnchored(line, x, y, ax, ay)
		y += dc.FontHeight() * lineSpacing
	}
}

func (r *Reader) Render(cards [3]Card, Q, A string) (image.Image, error) {
	img := image.NewNRGBA64(image.Rect(0, 0, defaultImageWidth, defaultImageHeight))

	dc := gg.NewContextForImage(img)
	dc.SetColor(image.Black)
	dc.Clear()
	dc.DrawImageAnchored(r.assets.BackgroundImg, defaultImageWidth*1/3, defaultImageHeight/2, 0.5, 0.5)

	ff := truetype.NewFace(r.assets.Font, &truetype.Options{
		Size: 20,
	})
	dc.SetFontFace(ff)
	dc.SetColor(color.White)

	// aW := defaultImageWidth * 2 / 3
	// // aH := defaultImageHeight
	// startH := 200
	// var card Card
	// for idx := 0; idx < 3; idx++ {
	// 	card = cards[idx]
	// 	pic := card.Pic
	// 	if card.Position == PositionReversed {
	// 		pic = imaging.Rotate(pic, 180, color.NRGBA{0, 0, 0, 0})
	// 	}
	// 	center := aW * (idx + 1) / 4
	// 	dc.DrawImageAnchored(pic, center, startH, 0.5, 0)

	// 	s := card.ZhString()
	// 	dc.DrawStringAnchored(s, float64(center), float64(startH-20), 0.5, 0.5)
	// }

	aW := defaultImageWidth * 2 / 3
	// aH := defaultImageHeight
	startH := 100
	var card Card
	for idx := 0; idx < 1; idx++ {
		card = cards[idx]
		pic := card.Pic
		if card.Position == PositionReversed {
			pic = imaging.Rotate(pic, 180, color.NRGBA{0, 0, 0, 0})
		}
		center := aW * (idx + 1) / 2
		dc.DrawImageAnchored(pic, center, startH, 0.5, 0)

		s := card.ZhString()
		dc.DrawStringAnchored(s, float64(center), float64(startH-20), 0.5, 0.5)
	}
	startH = 370
	for idx := 1; idx < 3; idx++ {
		card = cards[idx]
		pic := card.Pic
		if card.Position == PositionReversed {
			pic = imaging.Rotate(pic, 180, color.NRGBA{0, 0, 0, 0})
		}
		center := aW * (idx) / 3
		dc.DrawImageAnchored(pic, center, startH, 0.5, 0)

		s := card.ZhString()
		dc.DrawStringAnchored(s, float64(center), float64(startH-20), 0.5, 0.5)
	}

	ff = truetype.NewFace(r.assets.Font, &truetype.Options{
		Size: 18,
	})
	dc.SetFontFace(ff)
	DrawStringWrapped(dc, ff, Q+"\n\n"+A,
		float64(defaultImageWidth*5/6)-75, float64(defaultImageHeight/2), 0.5, 0.5, float64(defaultImageWidth/3), 1, gg.AlignLeft)

	return dc.Image(), nil
}
