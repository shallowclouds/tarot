package main

import (
	"bytes"
	"context"
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

type dumbGPTReader struct{}

func (r *dumbGPTReader) Chat(ctx context.Context, systemMsg, userMsg string) (string, error) {
	return "运气不错", nil
}

func main() {
	var (
		thingArg = flag.String("thing", "我这周运势怎么样？", "Thing you want to divine")
	)
	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	reader, err := tarot.NewReader(&dumbGPTReader{}, tarot.GetDefaultAssets())
	if err != nil {
		panic(err)
	}

	cards, img, err := reader.Divine(context.Background(), *thingArg, func (err error, res string)  {
		if err != nil {
			fmt.Printf("failed to read from tarot cards: %v\n", err)
			return
		}

		fmt.Printf("%s\n", res)
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(reader.Prompt(cards, *thingArg))

	err = SavePng(img, "divine_results.png")
	if err != nil {
		panic(err)
	}
}
