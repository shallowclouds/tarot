package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/png"
	"math/rand"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/shallowclouds/tarot"
)

func SavePng(img image.Image, p string) error {
	_ = os.Remove(p)
	data := bytes.Buffer{}
	if err := png.Encode(&data, img); err != nil {
		return errors.WithMessage(err, "failed to encode png")
	}

	return os.WriteFile(p, data.Bytes(), os.ModePerm)
}

func main() {
	var (
		thingArg = flag.String("thing", "我这周运势怎么样？", "Thing you want to divine")
	)
	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	reader, err := tarot.NewReader(nil, tarot.GetDefaultAssets())
	if err != nil {
		panic(err)
	}

	cards, err := reader.Choose()

	fmt.Println(reader.Prompt(cards, *thingArg))

	img, err := reader.Render(cards)
	if err != nil {
		panic(err)
	}

	err = SavePng(img, "divine_results.png")
	if err != nil {
		panic(err)
	}
}
