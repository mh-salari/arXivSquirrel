/*
Created on Jan 21 2023
@author: Hue (MohammadHossein) Salari
@email:hue.salari@gmail.com
*/

package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"arxiv.ai-hue.ir/arxiv/src/util"
	"github.com/mmcdole/gofeed"
)

type Channel struct {
	XMLName     xml.Name  `xml:"channel"`
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	PubDate     time.Time `xml:"pubDate"`
	Image       Image     `xml:"image"`
	ItemList    []Item    `xml:"item"`
}

type Image struct {
	Title string `xml:"title"`
	Url   string `xml:"url"`
	Link  string `xml:"link"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description Cdata  `xml:"description"`
	Date        string `xml:"published"`
	// Thumbnail   string `xml:"media:thumbnail"`
}

type Cdata struct {
	Value string `xml:",cdata"`
}

const basePath = "/var/www/arxiv.ai-hue.ir"
const sitePath = "https://arxiv.ai-hue.ir"

func generateFeed(feed *gofeed.Feed, timePeriod time.Time, keywords []string) *Channel {

	//Making the file path of our mascots/logo arXivSquirrel
	imageUrl := filepath.Join(sitePath, "resources", "arXivSquirrel.png")

	// Making the image data for RSS channel
	image := Image{Title: "Hue personalized cs.CV updates on arXiv.org", Url: imageUrl, Link: imageUrl}

	// Filling out the RSS chanel metadata
	channel := &Channel{
		Title:       "Hue personalized cs.CV updates on arXiv.org",
		Link:        "arXiv.ai.hue.ir",
		Description: "Selection of latest Computer Science updates on the arXiv.org base on Hue keywords.", //+ strings.Join(keywords, ", "),
		PubDate:     time.Now(),
		Image:       image,
		ItemList:    []Item{},
	}

	// Generating the Items of rss feed
	for _, item := range feed.Items {

		updated, err := time.Parse(time.RFC3339, item.Updated)
		if err != nil {
			log.Fatal(err)
		}

		// Ignore papers that have not updated on the past 24 hours
		if updated.Before(timePeriod) {
			continue
		}

		//Searching for the user target keywords in the paper title and description
		matches := util.SearchKeywords(item.Title+" "+item.Description, keywords)

		// if there is a mach, add that paper to RSS items
		if len(matches) > 0 {
			// Generate the link to pdf file of the paper
			link := strings.Replace(item.Link, "abs", "pdf", 1)

			// Get the paper name and extension
			_, paperName := filepath.Split(link)

			//Make a output path for this paper files
			paperOutPath := filepath.Join(basePath, "papers", paperName)

			// Generate dirs for the paper output dirs is doesn't exist
			if _, err := os.Stat(paperOutPath); os.IsNotExist(err) {
				err := os.MkdirAll(paperOutPath, os.ModePerm)
				if err != nil {
					log.Fatal(err)
				}
				//Downloading  the PDF file of the paper
				log.Println("Downloading the PDF file from", link, "🦥...")
				pdfPath := util.DownloadFile(paperOutPath, link)

				// Converting the paper pages to images fo max 5 pages
				log.Println("Converting the first five pages of pdf file to images 🐨...")
				util.PdfToImage(pdfPath, paperOutPath)

				// log.Println("Extracting images of the paper to create a thumbnail ...")
				// makeThumbnail(pdfPath)

				// Remove Downloaded pdf file to free up space!
				log.Println("Removing downloaded PDF file 🦈.")
				err = os.Remove(pdfPath)
				if err != nil {
					log.Fatal(err)
				}
			}

			// Get the tile of the paper, remove (arXiv __some_numbers__) from it
			reg := regexp.MustCompile(`\(arXiv:.+\)`)
			title := strings.TrimSpace(reg.ReplaceAllString(item.Title, "${1}"))

			//  Add title and list of the matched keywords to the description
			content := fmt.Sprintf("<b> <a href='%s'>%s</a></b><br><br>", item.Link, title)
			content += "<b>Keywords: </b>" + strings.Join(matches, ", ") + "<br><br>"
			content += item.Description

			// Add max 5 images generated from each page of the PDF file to the description
			imagesTable := "<table>\n<tr>\n"
			for i := 0; i < 5; i++ {
				// make the local address of the image on the dir to check if we have it or not
				imagePath := filepath.Join("papers", paperName, fmt.Sprintf("%d.jpg", i))
				// Just add the images exist
				if _, err := os.Stat(filepath.Join(basePath, imagePath)); err == nil {
					imageURL := filepath.Join(sitePath, imagePath)
					imagesTable += fmt.Sprintf("<td><a href='%s'><img src='%s' width='212' height='275'></a></td>\n", imageURL, imageURL)
				}
			}
			imagesTable += "\n</tr></table>"

			// Add all of the data of each item to its positions in the channel.itemList
			channel.ItemList = append(channel.ItemList, Item{
				Title:       title,
				Link:        item.Link,
				Description: Cdata{content + imagesTable},
				Date:        item.Published,
				// Thumbnail:   "url='" + filepath.Join("papers", paperName, "thumbnail.png") + "' width='75' height='50'",
			})
		}

	}
	return channel
}

func main() {

	// Create outputs path if not exist
	err := os.MkdirAll(basePath, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	// Get the list of user Keywords from a local CSV file
	keywords := util.ReadKeywords(filepath.Join(".", "keywords.csv"))
	log.Println(keywords)
	
	// Download and parsing the RSS feed for cs.CV list from arXiv
	log.Println("Getting the list of latest Computer Science and Pattern Recognition papers from arxiv.og 🗃️")

	fp := gofeed.NewParser()

	// Define the categories you want to search for
	categories := []string{"cs.HC", "cs.CV", "cs.AI", "cs.RO", "cs.CL", "cs.SE", "cs.LG"}

	var allItems []*gofeed.Item

	// Loop through each category and fetch items
	for idx, category := range categories {
		searchQuery := "cat:" + category
		url := fmt.Sprintf("http://export.arxiv.org/api/query?search_query=%s&sortBy=lastUpdatedDate&sortOrder=descending&max_results=250", searchQuery)
		
		log.Printf("Fetching articles for category %s (%d/%d)...", category, idx+1, len(categories))

		feed, err := fp.ParseURL(url)
		if err != nil {
			log.Panicf("Error in getting the RSS feed from arXiv for category %s\n", category, err)
		}

		allItems = append(allItems, feed.Items...)

		log.Printf("Fetched %d articles for category %s.", len(feed.Items), category)
	}


	// Construct a new feed from the fetched items
	constructedFeed := &gofeed.Feed{
	    Items: allItems,
	}

	// ---------------------- Building the RSS feed ----------------------//
	// Calculate yesterday
	previousDay := time.Now().UTC().AddDate(0, 0, -1).Truncate(24 * time.Hour)

	// We now pass allItems instead of feed to the generateFeed function. If generateFeed expects a feed, you will need to adjust it.
	channel := generateFeed(constructedFeed, previousDay, keywords)

	// Don't generate an empty RSS file If we don't have any papers
	// It's probably the weekend and the arXiv feed hasn't updated
	// So lets search for the papers published since 3 days ago
	if len(channel.ItemList) == 0 {
		log.Printf("No new articles have been published in the last 24 hours 🙀!")
		timePeriod := time.Now().UTC().AddDate(0, 0, -3).Truncate(24 * time.Hour)
		log.Printf("Searching for the articles published in the last 72 hours!")
		channel = generateFeed(constructedFeed, timePeriod, keywords)
	}


	// Convert out channel struct to a xml file
	newFeed, err := xml.MarshalIndent(channel, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	// Create and empty XML file
	rssOutPath := filepath.Join(basePath, "arxiv.xml")

	file, err := os.Create(rssOutPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// write out XML data with RSS header to the local file
	_, err = file.WriteString(xml.Header + "<rss version='2.0'>" + string(newFeed) + "</rss>")
	if err != nil {
		log.Fatal(err)
	}

	// Done!
	fmt.Println("New RSS feed saved to", rssOutPath, "🦫.")

}
