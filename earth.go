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
)

type EarthViewItem struct {
    Id string `json:"id"`
    Slug string `json:"slug"`
    Title string `json:"title"`
    Lat string `json:"lat"`
    Lng string `json:"lng"`
    PhotoURL string `json:"photoUrl"`
    Attribution string `json:"attribution"`
    MapsLink string `json:"mapsLink"`
    EarthLink string `json:"earthLink"`
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
		fmt.Fprintln(os.Stderr, Usage)
		return		
	}

	resp, err := http.Get(earthURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error requesting url, ", earthURL)
		return
	}

	expandedURL := resp.Request.URL.String()
	apiURL, err := url.ParseRequestURI(expandedURL)
	if err != nil || apiURL.Host != "earthview.withgoogle.com"  {
		fmt.Fprintln(os.Stderr, "error resolving url, expecting host earthview.withgoogle.com")
		return
	}

	apiURL.Path = path.Join("_api", apiURL.Path + ".json")
	resp, err = http.Get(apiURL.String())
	if err != nil {
		fmt.Fprintln(os.Stderr, "error requesting api, ", apiURL.String())
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error reading api request body")
		return
	}

	var item EarthViewItem
	err = json.Unmarshal(body, &item)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error decoding json")
	}

	if output == "" {
		output = item.Slug + ".jpg"
	}
	out, err := os.Create(output)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error writing output, ", output)
	}
	defer out.Close()

	resp, err = http.Get(item.PhotoURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error downloading file, ", item.PhotoURL)
	}
	defer resp.Body.Close()

	n, err := io.Copy(out, resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error downloading file, ", item.PhotoURL)
	}
	fmt.Println(item.Title)
	fmt.Printf("Lat: %s, Lng: %s, %s\n", item.Lat, item.Lng, item.Attribution)
	fmt.Printf("Downloaded %s to %s (1 file, %s)\n", item.PhotoURL, output, byteCountDecimal(n))
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