/*
Created on Jan 21 2023
@author: Hue (MohammadHossein) Salari
@email:hue.salari@gmail.com
*/

package util

import (
	"encoding/csv"
	"fmt"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/cavaliergopher/grab/v3"
	"github.com/karmdip-mi/go-fitz"
)

// Function to download a file from a given URL and save it to a specified file path
// using the grab library
func DownloadFile(fileOutPath string, url string) string {
	// Create a new client
	client := grab.NewClient()
	// Create a new request to download the file from the given URL
	req, err := grab.NewRequest(fileOutPath, url)
	if err != nil {
		// Log the error and exit if there is a problem creating the request
		log.Fatal(err)
	}
	// Set the User-Agent header to mimic a browser request
	req.HTTPRequest.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:108.0) Gecko/20100101 Firefox/108.0")

	// Perform the download
	resp := client.Do(req)

	if err := resp.Err(); err != nil {
		// Log the error and exit if there is a problem during the download
		log.Fatal("Download failed: ", err)
	}

	// Return the file name of the downloaded file
	return resp.Filename
}

// Function to convert a PDF to images and save them in a specified directory
func PdfToImage(pdfPath string, imagesOutPath string) {
	// Open the PDF at the given path
	doc, err := fitz.New(pdfPath)
	if err != nil {
		panic(err)
	}
	defer doc.Close()

	// Extract pages as images
	for pageNum := 0; pageNum < doc.NumPage(); pageNum++ {
		// Only convert the first 5 pages
		if pageNum >= 5 {
			break
		}

		// Get image of the page
		img, err := doc.Image(pageNum)
		if err != nil {
			panic(err)
		}

		// Create a new file with the page number as the file name
		f, err := os.Create(filepath.Join(imagesOutPath, fmt.Sprintf("%d.jpg", pageNum)))
		if err != nil {
			panic(err)
		}
		defer f.Close()

		// Encode the image as a JPEG with a quality of 50
		err = jpeg.Encode(f, img, &jpeg.Options{Quality: 50})
		if err != nil {
			panic(err)
		}
	}
}

// Function to search for keywords in a given string and return any matches
func SearchKeywords(str string, keywords []string) []string {
	var matches []string
	for _, keyword := range keywords {
		// Check if the keyword exists in the string, ignoring case
		if strings.Contains(strings.ToLower(str), strings.ToLower(keyword)) {
			matches = append(matches, keyword)
		}
	}
	return matches
}

// Function to read a list of keywords from a CSV file
func ReadKeywords(filePath string) []string {
	// Open the CSV file at the given path
	f, err := os.Open(filePath)
	if err != nil {
		// Log the error and exit if there is a problem opening the file
		log.Fatal("Unable to read input file "+filePath, "\n", err)
	}
	// Close the file when the function exits
	defer f.Close()

	// Create a new CSV reader
	csvReader := csv.NewReader(f)
	// Read all records from the CSV file
	records, err := csvReader.ReadAll()
	if err != nil {
		// Log the error and exit if there is a problem parsing the file as a CSV
		log.Fatal("Unable to parse file as CSV for "+filePath, "\n", err)
	}

	var keywords []string
	// Iterate through the records and add the first field of each record to the keywords slice
	for _, record := range records {
		keywords = append(keywords, record[0])
	}

	return keywords

}

// func makeThumbnail(pdfPath string) {
// 	var err error

// 	tempPath := filepath.Join(".", "tmp")
// 	err = os.MkdirAll(tempPath, os.ModePerm)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	err = pdfcpu.ExtractImagesFile(pdfPath, tempPath, []string{"1-"}, nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	files, err := ioutil.ReadDir(tempPath)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	var pngImage []string

// 	for _, file := range files {
// 		if filepath.Ext(file.Name()) == ".png" {
// 			pngImage = append(pngImage, file.Name())
// 		}
// 	}
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	}

// 	err = os.RemoveAll(tempPath)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }	thumbnailImagePath := filepath.Join(filepath.Dir(pdfPath), "thumbnail.png")
// 		fmt.Println(thumbnailImagePath)
// 		err = os.Rename(filepath.Join(tempPath, thumbnailImageName), thumbnailImagePath)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	}

// 	err = os.RemoveAll(tempPath)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }
