package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/wzshiming/ctc"
)

type EarthViewItem struct {
	Id          string `json:"id"`
	Slug        string `json:"slug"`
	Title       string `json:"title"`
	Region      string `json:"region"`
	Country     string `json:"country"`
	Lat         float64 `json:"lat"`
	Lng         float64 `json:"lng"`
	PhotoURL    string `json:"photoUrl"`
	Attribution string `json:"attribution"`
	MapsLink    string `json:"mapsLink"`
	EarthLink   string `json:"earthLink"`
}

const Usage = `Usage: earth [options...] EARTH_VIEW_URL
Options:
	-o name of output file
`

var (
	o = flag.String("o", "", "")
)

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, Usage)
	}

	flag.Parse()

	earthURL := ""
	if flag.NArg() > 0 {
		earthURL = flag.Args()[0]
	}
	output := *o

	if earthURL == "" {
		fmt.Fprint(os.Stderr, Usage)
		return
	}

	resp, err := http.Get(earthURL)
	if err != nil {
		printError("error requesting url, ", earthURL)
		return
	}

	expandedURL := resp.Request.URL.String()
	apiURL, err := url.ParseRequestURI(expandedURL)
	if err != nil || apiURL.Host != "earthview.withgoogle.com" {
		printError("error resolving url, expecting host earthview.withgoogle.com")
		return
	}

	apiURL.Path = path.Join("_api", apiURL.Path+".json")
	resp, err = http.Get(apiURL.String())
	if err != nil {
		printError("error requesting api, ", apiURL.String())
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		printError("error reading api request body")
		return
	}

	var item EarthViewItem
	err = json.Unmarshal(body, &item)
	if err != nil {
		printError("error decoding json")
		return
	}
	fmt.Print(ctc.ForegroundGreen, item.Region, ", ", item.Country, ctc.Reset, "\n")
	fmt.Print("Lat: ", ctc.ForegroundBlue, item.Lat, "°", ctc.Reset, ", Lng: ", ctc.ForegroundBlue, item.Lng, "°", ctc.Reset, ", ", item.Attribution, "\n")
	fmt.Print("Link: ", ctc.ForegroundBlue, item.EarthLink, ctc.Reset, "\n")

	if output == "" {
		output = item.Slug + ".jpg"
	}
	out, err := os.Create(output)
	if err != nil {
		printError("error writing output, ", output)
		return
	}
	defer out.Close()

	resp, err = http.Get(item.PhotoURL)
	if err != nil {
		printError("error downloading file, ", item.PhotoURL)
		return
	}
	defer resp.Body.Close()

	n, err := io.Copy(out, resp.Body)
	if err != nil {
		printError("error downloading file, ", item.PhotoURL)
		return
	}
	fmt.Print("Downloaded ", ctc.ForegroundBlue, item.PhotoURL, ctc.Reset, " to ", ctc.ForegroundBlue, output, ctc.ForegroundYellow, " (1 file, ", byteCountDecimal(n), ")", ctc.Reset, "\n")
}

func printError(a ...interface{}) {
	var message []interface{}
	message = append(message, ctc.ForegroundRed)
	message = append(message, a...)
	message = append(message, ctc.Reset, "\n")
	fmt.Fprint(os.Stderr, message...)
}

// programming.guide/go/formatting-byte-size-to-human-readable-format.html
func byteCountDecimal(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}
