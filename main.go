package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type SummaryRequest struct {
	VideoURL string `json:"videoUrl"`
	Lang     string `json:"lang"`
}

type OpenRouterRequest struct {
	Model    string               `json:"model"`
	Messages []OpenRouterMessage `json:"messages"`
}

type OpenRouterMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenRouterResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env dosyası yüklenemedi:", err)
	}

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.POST("/summarize", func(c *gin.Context) {
		var request SummaryRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			log.Println("JSON parse hatası:", err)
			c.JSON(http.StatusBadRequest, gin.H{"summary": "Geçersiz istek."})
			return
		}

		transcript := getTranscript(request.VideoURL)
		if transcript == "" {
			log.Println("Transkript alınamadı")
			c.JSON(http.StatusInternalServerError, gin.H{"summary": "Videodan transkript alınamadı."})
			return
		}

		log.Println("GÖNDERİLEN TRANSKRİPT:")
		log.Println(transcript)

		summary := getSummaryFromOpenRouter(transcript, request.Lang)
		c.JSON(http.StatusOK, gin.H{"summary": summary})
	})

	r.Run(":8080")
}

func getVideoID(videoURL string) string {
	if strings.Contains(videoURL, "youtu.be/") {
		parts := strings.Split(videoURL, "youtu.be/")
		return strings.Split(parts[1], "?")[0]
	}
	if strings.Contains(videoURL, "v=") {
		parts := strings.Split(videoURL, "v=")
		return strings.Split(parts[1], "&")[0]
	}
	return ""
}

func getTranscript(videoURL string) string {
	videoID := getVideoID(videoURL)
	if videoID == "" {
		log.Println("Video ID bulunamadı")
		return ""
	}

	cmd := exec.Command("./yt-dlp", "--skip-download", "--write-auto-sub", "--sub-lang", "en", "-o", videoID+".%(ext)s", videoURL)
	err := cmd.Run()
	if err != nil {
		log.Println("yt-dlp çalıştırılamadı:", err)
		return ""
	}

	vttFile := videoID + ".en.vtt"
	content, err := ioutil.ReadFile(vttFile)
	if err != nil {
		log.Println("VTT dosyası okunamadı:", err)
		return ""
	}
	defer os.Remove(vttFile)

	lines := strings.Split(string(content), "\n")
	var transcript strings.Builder
	for _, line := range lines {
		if len(line) > 0 && !strings.Contains(line, "-->") && !isNumber(line) {
			transcript.WriteString(line + " ")
		}
	}

	cleaned := regexp.MustCompile(`<.*?>`).ReplaceAllString(transcript.String(), "")
	return cleaned
}

func isNumber(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func getSummaryFromOpenRouter(transcript string, lang string) string {
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		log.Println("OpenRouter API anahtarı bulunamadı.")
		return "API anahtarı eksik."
	}

	prompt := map[string]string{
		"en": "Can you summarize the video below shortly and clearly?",
		"tr": "Aşağıdaki videoyu sade ve kısa şekilde özetler misin?",
		"ru": "Можешь кратко и ясно резюмировать это видео?",
		"ar": "هل يمكنك تلخيص هذا الفيديو بإيجاز ووضوح؟",
	}[lang]

	reqData := OpenRouterRequest{
		Model: "anthropic/claude-3-haiku",
		Messages: []OpenRouterMessage{
			{Role: "user", Content: prompt + "\n" + transcript},
		},
	}

	jsonData, err := json.Marshal(reqData)
	if err != nil {
		log.Println("JSON marshal hatası:", err)
		return "İstek hazırlanırken hata oluştu."
	}

	req, err := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println("HTTP isteği oluşturulamadı:", err)
		return "İstek gönderilemedi."
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HTTP-Referer", "http://localhost")
	req.Header.Set("X-Title", "YouTube Summary")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("OpenRouter API çağrısı başarısız:", err)
		return "API isteği başarısız."
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Yanıt okunamadı:", err)
		return "Yanıt işlenemedi."
	}

	log.Println("OPENROUTER YANITI:")
	log.Println(string(body))

	var orResp OpenRouterResponse
	err = json.Unmarshal(body, &orResp)
	if err != nil {
		log.Println("JSON çözümleme hatası:", err)
		return "Yanıt çözümlenemedi."
	}

	if len(orResp.Choices) > 0 {
		return orResp.Choices[0].Message.Content
	}
	return "Özet alınamadı."
}
