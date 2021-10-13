package main

import (
	"bytes"
	"database/sql"
	"image"
	"image/jpeg"
	"net/http"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/go-chi/chi"
	_ "github.com/mattn/go-sqlite3"
)

type Config struct {
	backgroundColor string
	fontColor       string
	fontPath        string
	fontSize        float64
	imageWidth      float64
	imageHeight     float64
	logLevel        string
}

var config Config
var db *sql.DB

func drawImage(w http.ResponseWriter, img *image.Image) {
	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, *img, nil); err != nil {
		log.Error("failed to encode image", err)
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Error("failed to write image", err)
	}
}

func initialize() {
	backgroundColor := os.Getenv("COUNTER_BG_COLOR")
	if backgroundColor == "" {
		backgroundColor = "#000000"
	}

	fontColor := os.Getenv("COUNTER_FONT_COLOR")
	if fontColor == "" {
		fontColor = "#FFFFFF"
	}

	fontPath := os.Getenv("COUNTER_FONT_DIR")
	if fontPath == "" {
		fontPath = "./fonts/"
	}

	fontFace := os.Getenv("COUNTER_FONT_FILE")
	if fontFace == "" {
		fontFace = "FiraCode.ttf"
	}

	fontPath += fontFace

	imageWidth, err := strconv.Atoi(os.Getenv("COUNTER_IMAGE_WIDTH"))
	if err != nil || imageWidth == 0 {
		imageWidth = 200
	}

	imageHeight, err := strconv.Atoi(os.Getenv("COUNTER_IMAGE_HEIGHT"))
	if err != nil || imageHeight == 0 {
		imageHeight = 50
	}

	fontSize, err := strconv.Atoi(os.Getenv("COUNTER_FONT_SIZE"))
	if err != nil || fontSize == 0 {
		fontSize = 32
	}

	logLevel := os.Getenv("COUNTER_LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}
	
	port, err := strconv.Atoi(os.Getenv("COUNTER_PORT"))
	if err != nil || port == 0 {
		port = 9776
	}

	config = Config{
		backgroundColor: backgroundColor,
		fontColor:       fontColor,
		fontPath:        fontPath,
		fontSize:        float64(fontSize),
		imageWidth:      float64(imageWidth),
		imageHeight:     float64(imageHeight),
		logLevel:        logLevel,
		port:	         float64(port)
	}
}

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.Info("initializing...")
	initialize()
	if config.logLevel == "debug" {
		log.SetLevel(log.DebugLevel)
	}

	if _, err := os.Stat("counter.sqlite"); os.IsNotExist(err) {
		err = createDB()
		if err != nil {
			log.Fatal("failed to create database", err)
		}
	} else {
		err = connectDB()
		if err != nil {
			log.Fatal("failed to connect to database", err)
		}
	}
	defer db.Close()

	log.Debug("loading handlers...")

	r := chi.NewRouter()

	r.Get("/", handleHome)
	r.Get("/count", handleGetCountText)
	r.Get("/count/counter.jpg", handleCountImage)

	log.Info("listening on port",port)
	err := http.ListenAndServe(port, r)
	if err != nil {
		log.Fatal("failed to listen:", err)
	}
}
