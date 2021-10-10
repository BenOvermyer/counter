package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/fogleman/gg"
)

func handleCountImage(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	path := query["path"][0]

	log.Info("fetching counter image for path ", path)
	count, err := incrementCount(path)
	if err != nil {
		log.Error("failed to increment count", err)
	}

	log.Debug("generating image...")
	dc := gg.NewContext(int(config.imageWidth), int(config.imageHeight))
	dc.SetHexColor(config.backgroundColor)
	dc.Fill()
	dc.Clear()
	dc.SetHexColor(config.fontColor)
	if err := dc.LoadFontFace(config.fontPath, config.fontSize); err != nil {
		panic(err)
	}
	dc.DrawStringAnchored(strconv.Itoa(count), config.imageWidth/2, config.imageHeight/2, 0.5, 0.5)
	img := dc.Image()
	drawImage(w, &img)
}

func handleGetCountText(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	path := query["path"][0]
	log.Info("fetching count for path ", path)
	count, err := getCount(path)
	if err != nil {
		log.Error("failed to get count ", err)
	}

	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]int)
	resp["count"] = count
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Error("failed to marshal response ", err)
	}
	w.Write(jsonResp)
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	log.Info("fetching root content")
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]string)
	resp["message"] = "we're online"
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Error("failed to marshal response", err)
	}
	w.Write(jsonResp)
}
