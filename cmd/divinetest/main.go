package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
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

func SaveJpg(img image.Image, p string) error {
	_ = os.Remove(p)
	data := bytes.Buffer{}
	if err := jpeg.Encode(&data, img, &jpeg.Options{Quality: 100}); err != nil {
		return errors.WithMessage(err, "failed to encode jpg")
	}

	return os.WriteFile(p, data.Bytes(), os.ModePerm)
}

func main() {
	var (
		thingArg = flag.String("thing", "我这周运势怎么样？", "Thing you want to divine")
	)
	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	reader, err := tarot.NewReader(&tarot.DumbGPTReader{}, "", "", tarot.GetDefaultAssets())
	if err != nil {
		panic(err)
	}

	_, img, res, err := reader.DivineSync(context.Background(), *thingArg)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", res)

	// fmt.Println(reader.Prompt(cards, *thingArg, ""))

	// err = SavePng(img, "divine_results.png")
	err = SaveJpg(img, "dev/divine_results.jpg")
	if err != nil {
		panic(err)
	}
}
