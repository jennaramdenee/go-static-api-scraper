package main

import (
  "log"
  "strings"
  "os"
  "path/filepath"
  "net/http"
  "io"
  "github.com/PuerkitoBio/goquery"
  "fmt"
  )

const baseUrl = "https://apidataparliament.azure-api.net/fixed-query/"

func runScraper() {
  doc, err := goquery.NewDocument(baseUrl)
  if err != nil {
    log.Fatal(err)
  }
  // 1. Find all links
  links := doc.Find("a")

  // 2. Loop through all links
  links.Each(func(index int, link *goquery.Selection){
    linkName, _ := link.Attr("href")

    // 3.Find the name of the folder by taking everything before ?
    folderName := strings.Split(linkName, "?")[0]

    // 4. Create folder if it doesn't already exist
    if checkIfExists(folderName) == true {
      err := os.Mkdir(folderName, os.ModePerm)
      catchError(err, "Cannot create folder")
    }

    // If file does not exist..
    // 5. Create file within that folder called "index"
    fileName := filepath.Join(folderName, "index")
    var file *os.File
    if checkIfExists(fileName) == true {
      file, err = os.Create(fileName)
      catchError(err, "Cannot create file")
      defer file.Close()

      // 6. Make a call to the original link
      response, err := http.Get(baseUrl + linkName)
      catchError(err, "Cannot make call")
      defer response.Body.Close()

      // 7. Write response to index file
      io.Copy(file, response.Body)
    }
  })
}

func catchError(err error, desc string) {
  if err != nil {
    log.Fatal(err)
    fmt.Println(desc)
  }
}

// returns true if path exists
func checkIfExists(path string) bool {
  _, err := os.Stat(path)
  return os.IsNotExist(err)
}

func main() {
  runScraper()
}
