package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Response []struct {
	ID     string      `json:"id"`
	Title  string      `json:"title"`
	Length interface{} `json:"length"`
	Pages  interface{} `json:"pages"`
	Pdf    string      `json:"pdf"`
}

func main() {
	resp, err := http.Get("http://localhost:3000/api/pdfs")

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode > 299 {
		log.Fatal("Response failed with status code: %d and\nbody: %s\n", resp.StatusCode, body)
	}

	if err != nil {
		log.Fatal(err)
	}
	var result Response

	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to the go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}
	var counter int
	//pdf := gopdf.GoPdf{}
	counter = 0
	for i := 0; i < 50; i++ {
		for _, rec := range result {
			dec, err := base64.StdEncoding.DecodeString(rec.Pdf)
			if err != nil {
				panic(err)
			}
			f, err := os.Create("myfilename" + strconv.Itoa(counter) + ".pdf")
			if err != nil {
				panic(err)
			}
			defer f.Close()

			if _, err := f.Write(dec); err != nil {
				panic(err)
			}
			if err := f.Sync(); err != nil {
				panic(err)
			}
			output, err := os.OpenFile("output.pdf", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Fatal(err)
			}
			for i := 0; i < 2; i++ {
				pdfComoBytes, err := ioutil.ReadFile("myfilename" + strconv.Itoa(i) + ".pdf")
				if err != nil {
					panic(err)
				}
				if _, err := output.Write([]byte(pdfComoBytes)); err != nil {
					log.Fatal(err)
				}
				if err := output.Close(); err != nil {
					log.Fatal(err)
				}

			}

			//fmt.Println(pdfComoBytes)
			counter++
		}

	}

}
