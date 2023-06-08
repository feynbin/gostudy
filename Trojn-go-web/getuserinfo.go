package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

//用法
//go run getuserinfo.go username&password

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: demo3 <cmd>")
		return
	}

	name := os.Args[1]
	cmd := exec.Command("trojan-go", "-api-addr", "web.feyncode.top:10000", "-api", "get", "-target-password", name)
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
		return
	}
	jsonString := string(out)
	var response Response
	err = json.Unmarshal([]byte(jsonString), &response)
	if err != nil {
		panic(err)
	}
	//response.Status.TrafficTotal.DownloadTraffic /= 1024*1024
	response.Status.TrafficTotal.DownloadTraffic /= 1024 * 1024
	response.Status.TrafficTotal.UploadTraffic /= 1024 * 1024
	result, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}
	formattedJson, err := jsonformat(string(result))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(formattedJson)

}

type User struct {
	Hash     string `json:"hash"`
	Password string `json:"password"`
}

type Status struct {
	SpeedCurrent interface{} `json:"speed_current"`
	SpeedLimit   interface{} `json:"speed_limit"`
	TrafficTotal Traffic     `json:"traffic_total"`
	User         User        `json:"user"`
}

type Traffic struct {
	DownloadTraffic int `json:"download_traffic"`
	UploadTraffic   int `json:"upload_traffic"`
}

type Response struct {
	Status  Status `json:"status"`
	Success bool   `json:"success"`
}

//格式化json

func jsonformat(data string) (string, error) {
	var str bytes.Buffer
	_ = json.Indent(&str, []byte(data), "", "  ")
	return str.String(), nil
}
