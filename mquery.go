package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dghubble/sling"
)

const baseURL = "https://openwhisk.ng.bluemix.net/api/v1/web/estesp%40us.ibm.com_dev/default/archList.json"

// QueryParams defines the parameters sent; in our case only "image" is needed
type QueryParams struct {
	Image string `url:"image"`
}

// ImageDataResponse holds the payload response on success
type ImageDataResponse struct {
	ImageData Payload `json:"payload,omitempty"`
	Error     string  `json:"error,omitempty"`
}

// Payload contains the JSON struct we get from the web action
type Payload struct {
	ManifestList string   `json:"manifestList,omitempty"`
	Tag          string   `json:"tag,omitempty"`
	ID           string   `json:"_id,omitempty"`
	RepoTags     []string `json:"repoTags,omitempty"`
	ArchList     []string `json:"archList,omitempty"`
	Platform     string   `json:"platform,omitempty"`
}

func stripSpaces(inStr string) string {
	//split on a new line as stdin appends \n
	parts := strings.Split(inStr, "\n")
	//return te part before the first newline with all spaces removed
	return strings.Replace(parts[0], " ", "", -1)
}

func main() {

	input, err := ioutil.ReadAll(os.Stdin)

	if err != nil {
		log.Fatal("Unable to read standard input:", err)
	}

	inputStr := stripSpaces(string(input))

	if len(inputStr) == 0 {
		log.Fatalln("An image name is required.\n")
	}

	qparam := &QueryParams{
		Image: inputStr,
	}

	response := new(ImageDataResponse)
	resp, err := sling.New().Base(baseURL).QueryStruct(qparam).ReceiveSuccess(response)
	if err != nil {
		fmt.Printf("ERROR: failed to query backend: %v\n", err)
	}
	os.Exit(processResponse(resp, response))
}

func processResponse(resp *http.Response, response *ImageDataResponse) int {
	if resp.StatusCode != 200 {
		// non-success RC from our http request
		fmt.Printf("ERROR: Failure code from our HTTP request: %d\n", resp.StatusCode)
		return 1
	}
	if response.Error != "" {
		// print out error
		fmt.Printf("ERROR: %s\n", response.Error)
		return 1
	}
	printManifestInfo(response)
	return 0
}

func printManifestInfo(response *ImageDataResponse) {
	fmt.Printf("Manifest List: %s\n", response.ImageData.ManifestList)
	if strings.Compare(response.ImageData.ManifestList, "Yes") == 0 {
		fmt.Println("Supported platforms:")
		for _, archosPair := range response.ImageData.ArchList {
			fmt.Printf(" - %s\n", archosPair)
		}
	} else {
		fmt.Printf(" Supports: %s\n", response.ImageData.Platform)
	}
	fmt.Println("")
}
