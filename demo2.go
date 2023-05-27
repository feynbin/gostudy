package main

import (
    "encoding/json"
    "fmt"
)

func main() {
    jsonString := `[{"status":{"user":{"hash":"b963218925b4d41d3e64dd6b1b447d6012105374328954842123e95e"},"traffic_total":{"upload_traffic":468220,"download_traffic":34515611},"speed_current":{},"speed_limit":{}}},{"status":{"user":{"hash":"870fbc96689ae85d5de3db94b017c305306c539fd9e7d8bd5b8b8fe6"},"traffic_total":{},"speed_current":{},"speed_limit":{}}}]`

    var data interface{}
    err := json.Unmarshal([]byte(jsonString), &data)
    if err != nil {
        fmt.Println("Error decoding JSON: ", err)
        return
    }

    formattedJSON, err := json.MarshalIndent(data, "", "    ")
    if err != nil {
        fmt.Println("Error encoding JSON: ", err)
        return
    }

    fmt.Println(string(formattedJSON))
}





// #	列出所有用户信息
// #   获取一个用户信息
// #   添加一个用户信息
// #	删除一个用户信息
// #	修改一个用户信息

