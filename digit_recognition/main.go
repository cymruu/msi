package main

import (
	"crypto/md5"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type Digit struct {
	digit       uint8
	imageBytes  []uint8
	imageWidth  int
	imageHeight int
}
type DigitRating struct {
	Digit
	rating float64
}

func (d *Digit) writeToFile() error {
	imageB := image.Gray{Pix: d.imageBytes, Stride: d.imageWidth, Rect: image.Rect(0, 0, d.imageWidth, d.imageHeight)}
	output, err := os.Create(fmt.Sprintf("./images/%x.png", md5.Sum(d.imageBytes)))
	if err != nil {
		return err
	}
	err = png.Encode(output, &imageB)
	return err
}
func (d *Digit) compare(compareList []Digit) []DigitRating {
	bestInCollection := .0
	potential := make([]DigitRating, 0)
	for i := 0; i < len(compareList); i++ {
		score := .0
		matches := 0
		for k := 0; k < d.imageHeight*d.imageWidth; k++ {
			if d.imageBytes[k] != 0 && compareList[i].imageBytes[k] != 0 {
				matches++
			}
		}
		score = float64(matches) / float64(d.imageHeight*d.imageWidth)
		if score > bestInCollection {
			bestInCollection = score
			potential = append(potential, DigitRating{Digit: compareList[i], rating: score})
		}
	}
	return potential
}

var Digits []Digit

func readAllImages(imgW, imgH int) {
	for i := uint8('0'); i < '9'+1; i++ {
		filename := fmt.Sprintf("./digits/data%c", i)
		file, err := os.Open(filename)
		digitString, _ := strconv.Atoi(string(i))
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		for {
			imageBytes := make([]byte, imgW*imgH)
			_, err := file.Read(imageBytes)
			if err != nil {
				fmt.Printf("Finished loaded from %s\n", filename)
				break
			}
			Digits = append(Digits, Digit{digit: uint8(digitString), imageWidth: imgW, imageHeight: imgH, imageBytes: imageBytes})
		}
	}
}
func LoadImagesToRecognize(path string, imgW, imgH int) ([]Digit, error) {
	digits := make([]Digit, 0)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println(err.Error())
	}
	for _, fileI := range files {
		file, err := os.Open(fmt.Sprintf("%s/%s", path, fileI.Name()))
		if err != nil {
			return nil, err
		}
		img, err := png.Decode(file)
		if err != nil {
			return nil, err
		}
		imageBytes := make([]byte, imgW*imgH)
		for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
			for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
				c := color.GrayModel.Convert(img.At(x, y)).(color.Gray)
				pixel := c.Y
				imageBytes[y*imgW+x] = pixel
			}
		}
		filenameString := strings.Split(fileI.Name(), ".")[0]
		digitString, _ := strconv.Atoi(filenameString)
		digits = append(digits, Digit{digit: uint8(digitString), imageBytes: imageBytes, imageWidth: imgW, imageHeight: imgH})
	}
	return digits, nil
}

//digit images from http://cis.jhu.edu/~sachin/digit/digit.html
func main() {
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	imageSize := 28
	timeStarted := time.Now()
	readAllImages(imageSize, imageSize)
	imagesLoadFinish := time.Since(timeStarted)
	randomKnownImage := Digits[rand.Intn(len(Digits))]
	loaded, err := LoadImagesToRecognize("./to_recognize", imageSize, imageSize)
	loaded = append(loaded, randomKnownImage)
	if err != nil {
		fmt.Print(err)
		return
	}
	for _, digitToRecognize := range loaded {
		potentials := digitToRecognize.compare(Digits)
		fmt.Printf("Best matches for %d: \n", digitToRecognize.digit)
		for _, d := range potentials {
			fmt.Printf("%d %f\n", d.digit, d.rating)
		}
	}
	totaltime := time.Since(timeStarted)
	fmt.Printf("Loaded %d images in %s totaltime %s\n", len(Digits), imagesLoadFinish, totaltime)
}
