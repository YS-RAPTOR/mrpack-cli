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

	"github.com/fatih/color"
)

func main() {
	mrpack := os.Args[1]

	downPtr := flag.Bool("download", true, "Set to false to skip downloads (default: true)")

	flag.Parse()

	var tempfolder = "mrpack-cli-" + strconv.FormatInt(rand.Int64N(99999), 10) + "/"

	if runtime.GOOS == "windows" {
		tempfolder = strings.ReplaceAll(tempfolder, "/", "\\")
		tempfolder = os.Getenv("APPDATA") + "\\" + tempfolder
	} else if runtime.GOOS == "linux" {
		tempfolder = "/tmp/" + tempfolder
	}
	os.MkdirAll(tempfolder, 0700)
	unzip(mrpack, tempfolder)
	//var jsonf
	var jsonf map[string]interface{} = openjson(tempfolder + "modrinth.index.json")

	if jsonf["game"] != "minecraft" {
		color.Set(color.FgRed)
		fmt.Println("ERROR: Game not supported.")
		color.Unset()
		os.Exit(1)
	}
	if jsonf["formatVersion"] != "1" {
		color.Set(color.FgYellow)
		fmt.Println("WARNING: formatVersion not '1'.")
		color.Unset()
	}

	exePath, err := os.Executable()
	if err != nil {
		color.Set(color.FgRed)
		fmt.Println("Error getting executable path: %v\n", err)
		color.Unset()
		os.Exit(1)
	}

	exePath, err = filepath.Abs(exePath)
	if err != nil {
		color.Set(color.FgRed)
		fmt.Println("Error getting absolute path: %v\n", err)
		color.Unset()
		os.Exit(1)
	}

	var packFolder = ""
	packFolder = filepath.Dir(exePath) + "/" + strings.ToLower(strings.ReplaceAll(jsonf["name"].(string), " ", "-")+"/")
	os.MkdirAll(packFolder+"mods/", os.ModePerm)
	os.MkdirAll(packFolder+"resourcepacks/", os.ModePerm)

	color.Set(color.FgGreen)
	fmt.Println("The modpack will be downloaded to: '" + packFolder + "'")
	color.Unset()

	if *downPtr {
		downloadMods(packFolder, jsonf)
		downloadResourcePacks(packFolder, jsonf)
	}

	addOverrides(packFolder, tempfolder)

	os.RemoveAll(tempfolder)
}
