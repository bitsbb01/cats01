package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"

	app "iptvcat-scraper/pkg"

	"github.com/gocolly/colly"
)

const aHref = "a[href]"

func downloadFile(filepath string, url string) (err error) {
	fmt.Println("downloadFile from ", url, "to ", filepath)

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func getUrlFromFile(filepath string, origUrl string) (string, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	// Splits on newlines by default.
	scanner := bufio.NewScanner(f)

	line := 1
	// https://golang.org/pkg/bufio/#Scanner.Scan
	for scanner.Scan() {
		if strings.HasPrefix(strings.ToLower(scanner.Text()), "http") {
			return scanner.Text(), nil
		}
		line++
	}

	if err := scanner.Err(); err != nil {
		// Handle the error
	}

	return origUrl, err
}

func checkNestedUrls() {
	fmt.Println("checkNestedUrls()")

	converted_urls := map[string]string{}
	ignored := 0
	processed := 0

	for _, stream := range app.Streams.All {
		url_lower := strings.ToLower(stream.Link)

		if strings.Contains(url_lower, "list.iptvcat.com") {
			if _, ok := converted_urls[url_lower]; ok {
				// stream.Link = converted_urls[url_lower]
				ignored++
				fmt.Println(">>> SKIP DUPLICATE: ", ignored)
				continue
			}

			const tmpFile = "tmp.m3u8"
			// Download the file
			downloadFile(tmpFile, stream.Link)

			// Get the Url
			newUrl, err := getUrlFromFile(tmpFile, stream.Link)
			if err != nil {
				fmt.Println(err)
				//return
			}
			//fmt.Println("newUrl found in link: ", newUrl)
			stream.Link = newUrl
			converted_urls[url_lower] = newUrl

			processed++

			// Delete the file
			err2 := os.Remove(tmpFile)
			if err2 != nil {
				fmt.Println(err2)
				return
			}

		} else {
			fmt.Println("no m3u8 found in link: ", stream.Link)
		}
	}

	fmt.Println("### MAP ", converted_urls)
	fmt.Println("### ignored ", ignored)
	fmt.Println("### processed ", processed)

}

func writeToFile() {
	streamsAll, err := json.MarshalIndent(app.Streams.All, "", "    ")
	streamsCountry, err := json.MarshalIndent(app.Streams.ByCountry, "", "    ")
	if err != nil {
		fmt.Println("error:", err)
	}

	os.MkdirAll("data/countries", os.ModePerm)

	ioutil.WriteFile("data/all-streams.json", streamsAll, 0644)
	ioutil.WriteFile("data/all-by-country.json", streamsCountry, 0644)
	for key, val := range app.Streams.ByCountry {
		// streamsCountry, err := json.Marshal(val)
		streamsCountry, err := json.MarshalIndent(val, "", "    ")
		if err != nil {
			fmt.Println("error:", err)
		}
		ioutil.WriteFile("data/countries/"+key+".json", streamsCountry, 0644)
	}
}

func processUrl(url string, domain string) {
	urlFilters := regexp.MustCompile(url + ".*")
	c := colly.NewCollector(
		colly.AllowedDomains(domain),
		colly.URLFilters(urlFilters),
	)

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)
	})

	c.OnHTML(aHref, app.HandleFollowLinks(c))
	c.OnHTML(app.GetStreamTableSelector(), app.HandleStreamTable(c))

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error: %d %s\n", r.StatusCode, r.Request.URL)
	})

	c.Visit(url)
	c.Wait()
	checkNestedUrls()
	writeToFile()
}

func main() {
	const iptvCatDomain = "iptvcat.com"

	urlList := [...]string{

		
		"https://iptvcat.com/algeria",
		
		"https://iptvcat.com/bahrain",
		
		"https://iptvcat.com/djibouti",
		
		"https://iptvcat.com/iraq",
		
		"https://iptvcat.com/jordan",
		
		"https://iptvcat.com/kuwait",
		
		"https://iptvcat.com/lebanon",
		
		"https://iptvcat.com/libya",
		
		"https://iptvcat.com/mauritania",
		
		"https://iptvcat.com/morocco",
		
		"https://iptvcat.com/oman",
		
		"https://iptvcat.com/qatar",
		
		"https://iptvcat.com/saudi_arabia",
		
		"https://iptvcat.com/somalia",
		
		"https://iptvcat.com/sudan",
		
		"https://iptvcat.com/syria",
		
		"https://iptvcat.com/tunisia",
		
		"https://iptvcat.com/united_arab_emirates",
		
		"https://iptvcat.com/yemen",
		
		"https://iptvcat.com/brazil",
		
		"https://iptvcat.com/mexico",
		
		"https://iptvcat.com/puerto_rico",
		
		"https://iptvcat.com/egypt",
		
		"https://iptvcat.com/seychelles",
		
		"https://iptvcat.com/new_zealand",
		
		"https://iptvcat.com/afghanistan",
		
		"https://iptvcat.com/azerbaijan",
		
		"https://iptvcat.com/china",
		
		"https://iptvcat.com/hong_kong",
		
		"https://iptvcat.com/india",
		
		"https://iptvcat.com/indonesia",
		
		"https://iptvcat.com/iran",
		
		"https://iptvcat.com/israel",
		
		"https://iptvcat.com/japan",
		
		"https://iptvcat.com/malaysia",
		
		"https://iptvcat.com/pakistan",
		
		"https://iptvcat.com/palestine",
		
		"https://iptvcat.com/south_korea",
		
		"https://iptvcat.com/taiwan",
		
		"https://iptvcat.com/thailand",
		
		"https://iptvcat.com/turkmenistan",
		
		"https://iptvcat.com/ex_yugoslavia",
		
		"https://iptvcat.com/scandinavia",
		
		"https://iptvcat.com/albania",
		
		"https://iptvcat.com/austria",
		
		"https://iptvcat.com/belgium",
		
		"https://iptvcat.com/bulgaria",
		
		"https://iptvcat.com/czech_republic",
		
		"https://iptvcat.com/france",
		
		"https://iptvcat.com/germany",
		
		"https://iptvcat.com/greece",
		
		"https://iptvcat.com/hungary",
		
		"https://iptvcat.com/ireland",
		
		"https://iptvcat.com/italy",
		
		"https://iptvcat.com/netherlands",
		
		"https://iptvcat.com/poland",
		
		"https://iptvcat.com/portugal",
		
		"https://iptvcat.com/russian_federation",
		
		"https://iptvcat.com/spain",
		
		"https://iptvcat.com/switzerland",
		
		"https://iptvcat.com/turkey",
		
		"https://iptvcat.com/united_kingdom",
		
		"https://iptvcat.com/canada",
		
		"https://iptvcat.com/united_states_of_america",
		
		"https://iptvcat.com/australia",
		
		"https://iptvcat.com/new_channels-26",
		
		"https://iptvcat.com/sport",
		
		"https://iptvcat.com/s/movie/mark/movie",
		
		"https://iptvcat.com/s/star/mark/star",
		
		"https://iptvcat.com/s/hbo/mark/hbo",
		}

	for _, element := range urlList {
		processUrl(element, iptvCatDomain)
	}

}
