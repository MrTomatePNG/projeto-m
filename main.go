package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/nfnt/resize"
)

type ResquestBody struct {
	User    string `json:"user"`
	Comment string `json:"comment"`
	Image   []byte `json:"image"`
}

func main() {
	r := chi.NewRouter()

	r.Post("/upload", uploadHandler)

	http.ListenAndServe(":8080", r)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody ResquestBody

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading body", http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &reqBody); err != nil {
		http.Error(w, "invalid json format", http.StatusBadRequest)
		return
	}
	imagePath, err := saveAndResizeImage(reqBody.Image)
	if err != nil {
		http.Error(w, "unprocessable image", http.StatusUnsupportedMediaType)
		return
	}

	response := map[string]string{"message": "image upload sucessfull", "path": imagePath}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func saveAndResizeImage(b []byte) (string, error) {
	img, _, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		return "", fmt.Errorf("error decoding image: %w", err)
	}
	resizedImg := resize.Resize(800, 0, img, resize.Lanczos3)
	outFile, err := os.Create("images/resized_image.jpg")
	if err != nil {
		return "", fmt.Errorf("error creating file:%w", err)
	}
	defer outFile.Close()
	jpeg.Encode(outFile, resizedImg, nil)
	return "images/resized_image.jpg", nil
}
