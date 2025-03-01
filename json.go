package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func openjson(tmpfolder string) ModPack {
	data, err := os.ReadFile(tmpfolder)
	if err != nil {
		fmt.Println("Could not open JSON:", err)
	}
	result := ModPack{}
	json.Unmarshal(data, &result)
	return result
}
