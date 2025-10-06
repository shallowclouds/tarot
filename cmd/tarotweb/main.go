package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"image/jpeg"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/shallowclouds/tarot"
)

type DivineRequest struct {
	Question string `json:"question"`
	Reader   string `json:"reader"`
}

type DivineResponse struct {
	Cards   []CardInfo `json:"cards"`
	Reading string     `json:"reading"`
	Image   string     `json:"image"`
	Error   string     `json:"error,omitempty"`
}

type CardInfo struct {
	Name     string `json:"name"`
	NameZh   string `json:"name_zh"`
	Position string `json:"position"`
}

var (
	gptReader tarot.GPTReader
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	apiKey := os.Getenv("DEEPSEEK_API_KEY")
	baseURL := os.Getenv("DEEPSEEK_BASE_URL")
	model := os.Getenv("DEEPSEEK_MODEL")
	readerType := os.Getenv("READER_TYPE")
	if readerType == "" {
		readerType = "dumb"
	}

	switch readerType {
	case "deepseek":
		if apiKey != "" {
			gptReader = tarot.NewDeepSeekReader(apiKey, baseURL, model)
			log.Println("Using DeepSeek reader")
		} else {
			gptReader = &tarot.DumbGPTReader{}
			log.Println("No DeepSeek API key, using DumbGPTReader")
		}
	default:
		gptReader = &tarot.DumbGPTReader{}
		log.Println("Using DumbGPTReader")
	}

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/api/divine", handleDivine)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Printf("Server starting on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, "static/index.html")
}

func handleDivine(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req DivineRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.Question == "" {
		req.Question = "我今天的运势如何？"
	}

	reader, err := tarot.NewReader(gptReader, "", "", tarot.GetDefaultAssets())
	if err != nil {
		sendError(w, "Failed to initialize reader: "+err.Error(), http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	res, err := reader.DivineWithOption(ctx, tarot.DivineOption{
		Question:  req.Question,
		Asker:     "访客",
		AskerImg:  nil,
		Reader:    "",
		ReaderImg: nil,
		Callback:  nil,
	})
	if err != nil {
		sendError(w, "Failed to divine: "+err.Error(), http.StatusInternalServerError)
		return
	}

	cards := make([]CardInfo, len(res.Cards))
	for i, card := range res.Cards {
		cards[i] = CardInfo{
			Name:     card.Name,
			NameZh:   card.ZhName,
			Position: string(card.Position),
		}
	}

	buf := make([]byte, 0)
	imgBuf := &bytesBuffer{buf: buf}
	if err := jpeg.Encode(imgBuf, res.Img, &jpeg.Options{Quality: 90}); err != nil {
		sendError(w, "Failed to encode image: "+err.Error(), http.StatusInternalServerError)
		return
	}
	imageBase64 := base64.StdEncoding.EncodeToString(imgBuf.buf)

	response := DivineResponse{
		Cards:   cards,
		Reading: res.Result,
		Image:   imageBase64,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func sendError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(DivineResponse{Error: message})
}

type bytesBuffer struct {
	buf []byte
}

func (b *bytesBuffer) Write(p []byte) (n int, err error) {
	b.buf = append(b.buf, p...)
	return len(p), nil
}
