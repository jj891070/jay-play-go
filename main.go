package main

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/kbinani/screenshot"
)

type Webhook struct {
	Type     string     `json:"@type"`
	Text     string     `json:"text"`
	Sections []Sections `json:"sections"`
}
type Images struct {
	Image string `json:"image"`
}
type Sections struct {
	Images []Images `json:"images"`
}

func main() {
	// n := screenshot.NumActiveDisplays()
	num := 1
	// for i := 0; i < n; i++ {
	bounds := screenshot.GetDisplayBounds(num)
	// 	Large: Display as 1000 x 1000 image
	// Medium: Display as a 600 x 600 image
	// Small: Display as a 200 x 200 image
	// Tiny: Display as a 50 x 50 image
	// Text: Display as a URL
	bounds.Min.X = 100
	bounds.Max.X = 200
	bounds.Min.Y = 100
	bounds.Max.Y = 200
	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		panic(err)
	}
	fileName := fmt.Sprintf("%d_%dx%d.png", num, bounds.Dx(), bounds.Dy())
	file, _ := os.Create(fileName)
	defer file.Close()
	png.Encode(file, img)

	fmt.Printf("#%d : %v \"%s\"\n", num, bounds, fileName)
	dat, err := os.ReadFile(fileName)
	if err != nil {
		log.Println("err =>", err)
		return
	}
	// fmt.Println(b64.URLEncoding.EncodeToString(dat))

	url := "https://xsgames.webhook.office.com/webhookb2/39feab54-"
	method := "POST"

	var content Webhook
	content.Type = "MessageCard"
	content.Text = "test2"
	content.Sections = append(
		content.Sections,
		Sections{
			[]Images{
				{
					Image: "data:image/png;base64," + b64.StdEncoding.EncodeToString(dat),
				},
			},
		},
	)

	webhookJson, err := json.Marshal(content)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(webhookJson))
	payload := strings.NewReader(string(webhookJson))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))

	// }
}

func compressImageResource(data []byte) []byte {
	imgSrc, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return data
	}
	newImg := image.NewRGBA(imgSrc.Bounds())
	draw.Draw(newImg, newImg.Bounds(), &image.Uniform{C: color.White}, image.Point{}, draw.Src)
	draw.Draw(newImg, newImg.Bounds(), imgSrc, imgSrc.Bounds().Min, draw.Over)

	buf := bytes.Buffer{}
	err = jpeg.Encode(&buf, newImg, &jpeg.Options{Quality: 40})
	if err != nil {
		return data
	}
	if buf.Len() > len(data) {
		return data
	}
	return buf.Bytes()
}
