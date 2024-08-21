package main

import (
	"encoding/json"
	"io"
	"os"
)

func openjson(tmpfolder string) map[string]interface{} {
	jsonFile, err := os.Open(tmpfolder + "")
	if err != nil {
		panic(err)
	}
	byteValue, _ := io.ReadAll(jsonFile)

	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)
	defer jsonFile.Close()
	return result
}
