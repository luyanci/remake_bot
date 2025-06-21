package bot

import (
	"fmt"

	"github.com/imroc/req/v3"
)

type Lolicon struct {
	PID        int      `json:"pid"`
	P          int      `json:"p"`
	Title      string   `json:"title"`
	Author     string   `json:"author"`
	Tags       []string `json:"tags"`
	R18        bool     `json:"r18"`
	Width      int      `json:"width"`
	Height     int      `json:"height"`
	Aitype     int      `json:"aitype"`
	Ext        string   `json:"ext"`
	Uploaddate int64    `json:"uploadDate"`
	Urls       struct {
		Original string `json:"original"`
		//Small     string `json:"small"`
		//Medium    string `json:"medium"`
		//Thumbnail string `json:"thumbnail"`
	} `json:"urls"`
}

type loliconResponse struct {
	Error string    `json:"error"`
	Data  []Lolicon `json:"data"`
}

func GetLoliconImage(data string) (string, error) {
	var loliconResponse loliconResponse
	url := fmt.Sprintf("https://api.lolicon.app/setu/v2?tag=%s&num=1", data)
	client := req.C()
	client.ImpersonateFirefox()
	resp, err := client.R().SetSuccessResult(loliconResponse).Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch lolicon image: %w", err)
	}
	if resp.IsError() {
		return "", fmt.Errorf("error response from lolicon API: %s", resp.String())
	}
	if len(loliconResponse.Data) == 0 {
		fmt.Println("loliconResponse.Data is empty")
		fmt.Println("loliconResponse:", loliconResponse)
		return "", fmt.Errorf("no data found for tag: %s", data)
	}
	return loliconResponse.Data[0].Urls.Original, nil

}
