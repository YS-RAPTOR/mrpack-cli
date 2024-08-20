package main

import (
	"flag"
	"fmt"
	"math/rand/v2"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func main() {
  mrpack := os.Args[1]
  
  downPtr := flag.Bool("download", true, "Set to false to skip downloads")

  flag.Parse()

  var tempfolder = "mrpack-cli-" + strconv.FormatInt(rand.Int64N(99999), 10) + "/"

  if runtime.GOOS == "windows" {
    tempfolder = os.Getenv("APPDATA") + tempfolder
  } else if runtime.GOOS == "linux" {
    tempfolder = "/tmp/" + tempfolder
  }
  os.MkdirAll(tempfolder, 0700)
  unzip(mrpack, tempfolder)
  var jsonf map[string]interface{}
  jsonf = openjson(tempfolder + "modrinth.index.json")

  exePath, err := os.Executable()
  if err != nil {
    fmt.Printf("Error getting executable path: %v\n", err)
  }

  exePath, err = filepath.Abs(exePath)
  if err != nil {
    fmt.Printf("Error getting absolute path: %v\n", err)
  }

  var packFolder = "" 
  packFolder = filepath.Dir(exePath) + "/" + strings.ToLower(strings.ReplaceAll(jsonf["name"].(string), " ", "-") + "/")
  os.MkdirAll(packFolder + "mods/", os.ModePerm)
  os.MkdirAll(packFolder + "resourcepacks/", os.ModePerm)

  fmt.Println("The modpack will be downloaded to: '" + packFolder +"'")
 
  if *downPtr == true {
    downloadMods(packFolder, jsonf)
    downloadResourcePacks(packFolder, jsonf)
  }

  addOverrides(packFolder, tempfolder)

  os.RemoveAll(tempfolder)
}
