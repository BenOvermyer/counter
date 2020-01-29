package main

import (
	"bytes"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"strconv"

	"github.com/fogleman/gg"
	"github.com/go-chi/chi"
)

var counters = map[string]int{}

func drawImage(w http.ResponseWriter, img *image.Image) {
	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, *img, nil); err != nil {
		log.Println("unable to encode image.")
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Println("unable to write image.")
	}
}

func countHandler(w http.ResponseWriter, r *http.Request) {
	domain := chi.URLParam(r, "domain")

	counters[domain]++

	count := strconv.Itoa(counters[domain])

	dc := gg.NewContext(200, 50)
	dc.SetRGB(0, 0, 0)
	dc.Fill()
	dc.Clear()
	dc.SetRGB(1, 1, 1)
	if err := dc.LoadFontFace("Berylium.ttf", 32); err != nil {
		panic(err)
	}
	dc.DrawStringAnchored(count, 200/2, 50/2, 0.5, 0.5)
	img := dc.Image()
	drawImage(w, &img)
}

func main() {
	r := chi.NewRouter()
	r.Get("/{domain}/counter.jpg", countHandler)
	log.Println("listening on port 9776")
	err := http.ListenAndServe(":9776", r)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
