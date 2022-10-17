package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/signintech/gopdf"
	rpdf "rsc.io/pdf"
)

type Response []struct {
	ID     string      `json:"id"`
	Title  string      `json:"title"`
	Length interface{} `json:"length"`
	Pages  interface{} `json:"pages"`
	Pdf    string      `json:"pdf"`
}

func main() {
	//resp, err := http.Get("http://localhost:3000/api/pdfs")
	resp, err := http.Get("http://localhost:3001/pdfs")

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
	//var counter int
	pdf := gopdf.GoPdf{}

	pdf.Start(gopdf.Config{
		PageSize: *gopdf.PageSizeA4,
	})

	t1 := time.Now()

	for i := 0; i < 3; i++ {
		for x := 0; x < len(result); x++ {
			//archivo := "myfilename" + strconv.Itoa(x) + ".pdf"
			archivo := "myfilename.pdf"

			//fmt.Println(result[x].Title)
			dec, err := base64.StdEncoding.DecodeString(result[x].Pdf)
			if err != nil {
				panic(err)
			}
			f, err := os.Create(archivo)
			if err != nil {
				panic(err)
			}
			//defer f.Close()

			if _, err := f.Write(dec); err != nil {
				panic(err)
			}
			if err := f.Sync(); err != nil {
				panic(err)
			}

			defer f.Close()

			fi, err := rpdf.Open(archivo)
			if err != nil {
				panic(err)
			}

			pageNum := fi.NumPage()
			fi = nil
			f = nil

			for p := 1; p <= pageNum; p++ {
				pdf.AddPage()

				//bgpdf := gofpdi.ImportPage(pdf, filer, p, "/MediaBox")
				tpl1 := pdf.ImportPage(archivo, p, "/MediaBox")
				pdf.UseImportedTemplate(tpl1, 50, 100, 400, 0)

				//gofpdi.UseImportedTemplate(pdf, bgpdf, 0, 0, 210*MM_TO_RMPOINTS, 297

			}

		}
		pdf.WritePdf("example.pdf")

	}

	t2 := time.Now()

	diff := t2.Sub(t1)
	fmt.Println(diff)

}
