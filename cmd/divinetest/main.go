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
	"github.com/sashabaranov/go-openai"
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
		thingArg   = flag.String("thing", "我这周运势怎么样？", "Thing you want to divine")
		readerType = flag.String("reader", "dumb", "Reader type: dumb, openai, deepseek")
		apiKey     = flag.String("api-key", "", "API key for OpenAI or DeepSeek")
		baseURL    = flag.String("base-url", "", "Base URL for API (optional, for DeepSeek)")
		model      = flag.String("model", "", "Model name (optional)")
	)
	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	if *apiKey == "" {
		if *readerType == "deepseek" {
			*apiKey = os.Getenv("DEEPSEEK_API_KEY")
		} else if *readerType == "openai" {
			*apiKey = os.Getenv("OPENAI_API_KEY")
		}
	}

	if *baseURL == "" && *readerType == "deepseek" {
		*baseURL = os.Getenv("DEEPSEEK_BASE_URL")
	}

	if *model == "" {
		if *readerType == "deepseek" {
			*model = os.Getenv("DEEPSEEK_MODEL")
		} else if *readerType == "openai" {
			*model = os.Getenv("OPENAI_MODEL")
		}
	}

	var gptReader tarot.GPTReader
	switch *readerType {
	case "deepseek":
		if *apiKey == "" {
			fmt.Println("Warning: No DeepSeek API key provided, using DumbGPTReader instead")
			gptReader = &tarot.DumbGPTReader{}
		} else {
			gptReader = tarot.NewDeepSeekReader(*apiKey, *baseURL, *model)
			if *model == "" {
				*model = "deepseek-chat"
			}
			fmt.Printf("Using DeepSeek reader with model: %s\n", *model)
		}
	case "openai":
		if *apiKey == "" {
			fmt.Println("Warning: No OpenAI API key provided, using DumbGPTReader instead")
			gptReader = &tarot.DumbGPTReader{}
		} else {
			client := openai.NewClient(*apiKey)
			gptReader = tarot.NewChatGPTReader(client)
			fmt.Println("Using OpenAI ChatGPT reader")
		}
	default:
		gptReader = &tarot.DumbGPTReader{}
		fmt.Println("Using DumbGPTReader (default)")
	}

	reader, err := tarot.NewReader(gptReader, "", "", tarot.GetDefaultAssets())
	if err != nil {
		panic(err)
	}

	res, err := reader.DivineWithOption(context.Background(), tarot.DivineOption{
		Question:  *thingArg,
		Asker:     "岳云翎",
		AskerImg:  nil,
		Reader:    "",
		ReaderImg: nil,
		Callback:  nil,
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", res)

	err = SaveJpg(res.Img, "assets/divine_results.jpg")
	if err != nil {
		panic(err)
	}
}
