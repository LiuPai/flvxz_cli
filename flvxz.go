package main

import (
	"net/http"
	//	"net/url"
	"encoding/base64"
	"encoding/json"
	"strings"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

var (
	baseURL = "http://api.flvxz.com/jsonp/purejson/url/"
	quality = "2"
	dir = "/tmp/"
)

type FILE struct {
	Furl string
	Ftype string
	Seconds uint64
	Bytes uint64
	Time string
	Size string
}

type Video struct {
	Title string
	Files []FILE
	Site string
	Quality string
}

func encodeURL(url string) (result string){
	r := strings.NewReplacer("://", ":##")
	result = r.Replace(url)
	result = base64.StdEncoding.EncodeToString([]byte(result))
	return
}

// just url get method, no time do other way
func getFLV(url string){
	eURL := encodeURL(url)
	apiURL := baseURL + eURL
	resp, err := http.Get(apiURL)
	if err != nil {
		log.Fatalln("http.Get err = %T(%v)\n", err, err)
	}
	defer resp.Body.Close() 
	j, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln("http.Get err = %T(%v)\n", err, err)
	}
	
	var videos []Video
	err = json.Unmarshal(j, &videos)
	if err != nil {
		log.Fatalln("error:", err)
	}
	for index, video := range videos {
		fmt.Printf("Index : %02d Title : %s\n  Quality : %s\n", index, video.Title, video.Quality)
		for i, file := range video.Files {
			fmt.Printf("%02d :\n    Type : %s\n    Time : %s\n    Size : %s\n", i, file.Ftype, file.Time, file.Size)
		}
	}
	var sltidx uint64 = 1
	fmt.Print("Select Video Index Number Which You Want: ") 
	if len(videos) > 0 {
		if len(videos) > 1{
			fmt.Fscanf(os.Stdin, "%d\n", &sltidx)	
		}
		for i, file := range videos[sltidx].Files {
			fileIdx := fmt.Sprintf("%02d", i)
			r := strings.NewReplacer(" ", "_")
			
			outFile := dir + r.Replace(videos[sltidx].Title) + "_" + fileIdx + "." + file.Ftype
			resp, err := http.Get(file.Furl)
			if err != nil {
				log.Fatalln("http.Get err = %T(%v)\n", err, err)
			}
			defer resp.Body.Close() 
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatalln("http.Get err = %T(%v)\n", err, err)
			}
			ioutil.WriteFile(outFile, body, 0644)
		}
	}
	
}
func main(){
	getFLV(os.Args[1])
}
