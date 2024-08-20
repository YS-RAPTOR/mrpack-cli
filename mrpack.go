package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func downloadMods(packFolder string, jsonf map[string]interface{}) {
  if mods, ok := jsonf["files"].([]interface{}); ok {
    for i, mod := range mods {
  		if modMap, ok := mod.(map[string]interface{}); ok {
  				path := modMap["path"].(string)
           
          if strings.Contains(path, "mods") {
          fmt.Println("Downloading mod:", strings.Split(path, "/")[1], "(" + strconv.FormatInt(int64(i), 10) + "/" + strconv.FormatInt(int64(len(mods)), 10) + ")")
          out, err := os.Create(packFolder + "mods/" + strings.Split(path, "/")[1])
          if err != nil {
            panic(err)
          }
          defer out.Close()
          for _, dwn := range modMap["downloads"].([]interface{}) {
              fmt.Println(dwn.(string))
              resp, err := http.Get(dwn.(string))
              n, err := io.Copy(out, resp.Body)
              if err != nil {
                panic(err)
              }
              defer resp.Body.Close()
              _ = n // Just to make the compiler shut up
            }
          }
  			}
  		}
  } else {
  		fmt.Println("Expected 'files' to be a slice, but got something else.")
  }
}

func addOverrides(packFolder string, tempFolder string) {
  fmt.Println("Copy: " + tempFolder + "overrides", "->", packFolder)
  cmd, err := exec.Command("/bin/sh", "-c", "cp -r " + tempFolder + "overrides/* " + packFolder).Output()
  if err != nil {
    fmt.Errorf("%v", err)
  }
  fmt.Println(cmd)
}
