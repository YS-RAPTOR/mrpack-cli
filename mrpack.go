package main

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

var downloaded = 1

func downloadMods(packFolder string, jsonf map[string]interface{}) {
	if mods, ok := jsonf["files"].([]interface{}); ok {
		for _, mod := range mods {
			if modMap, ok := mod.(map[string]interface{}); ok {
				path := modMap["path"].(string)

				if strings.Contains(path, "mods") {
					color.Set(color.FgGreen)
					fmt.Print("Downloading mod: ")
					color.Set(color.Bold)
					fmt.Print(strings.Split(path, "/")[1])
					color.Set(color.ResetBold)
					fmt.Println(" ("+strconv.FormatInt(int64(downloaded), 10)+"/"+strconv.FormatInt(int64(len(mods)), 10)+")")
					color.Unset()
					out, err := os.Create(packFolder + "mods/" + strings.Split(path, "/")[1])
					if err != nil {
						panic(err)
					}
					defer out.Close()
					for _, dwn := range modMap["downloads"].([]interface{}) {
						//fmt.Println(dwn.(string))
						resp, err := http.Get(dwn.(string))
						if err != nil {
							panic(err)
						}
						n, err := io.Copy(out, resp.Body)
						if err != nil {
							panic(err)
						}
						defer resp.Body.Close()
						_ = n // Just to make the compiler shut up

						// Read SHA 512 hash
						f, err := os.Open(packFolder + "mods/" + strings.Split(path, "/")[1])
						if err != nil {
							panic(err)
						}

						has := sha512.New()
						if _, err := io.Copy(has, f); err != nil {
							panic(err)
						}
						f.Close()

						if hashes, ok := modMap["hashes"].(map[string]interface{}); ok {
							fhas := hashes["sha512"].(string)
							if fhas != hex.EncodeToString(has.Sum(nil)) {
								color.Set(color.FgRed)
								fmt.Println("Warning: Potential Fake Mod")
								fmt.Println("The mod hash doesn't match what is recorded in the .mrpack, which may indicate a fake or modified version. Please verify the mod’s source and ensure it’s from a trusted provider. (e.g., Modrinth)")
								fmt.Println("File deleted.")
								color.Unset()
								os.Remove(packFolder + "mods/" + strings.Split(path, "/")[1])
							}
						}
					}
					downloaded++
				}
			}
		}
	} else {
		fmt.Println("Expected 'files' to be a slice, but got something else.")
	}
}

func downloadResourcePacks(packFolder string, jsonf map[string]interface{}) {
	if mods, ok := jsonf["files"].([]interface{}); ok {
		for _, mod := range mods {
			if modMap, ok := mod.(map[string]interface{}); ok {
				path := modMap["path"].(string)

				if strings.Contains(path, "resourcepack") {
					color.Set(color.FgBlue)
					fmt.Print("Downloading resourcepack:") 
					color.Set(color.Bold)
					fmt.Print(strings.Split(path, "/")[1])
					color.Set(color.ResetBold)
					fmt.Println("("+strconv.FormatInt(int64(downloaded), 10)+"/"+strconv.FormatInt(int64(len(mods)), 10)+")")
					color.Unset()
					out, err := os.Create(packFolder + "resourcepacks/" + strings.Split(path, "/")[1])
					if err != nil {
						panic(err)
					}
					defer out.Close()
					for _, dwn := range modMap["downloads"].([]interface{}) {
						//fmt.Println(dwn.(string))
						resp, err := http.Get(dwn.(string))
						if err != nil {
							panic(err)
						}
						n, err := io.Copy(out, resp.Body)
						if err != nil {
							panic(err)
						}
						defer resp.Body.Close()
						_ = n // Just to make the compiler shut up

						// Read SHA 512 hash
						f, err := os.Open(packFolder + "resourcepacks/" + strings.Split(path, "/")[1])
						if err != nil {
							panic(err)
						}

						has := sha512.New()
						if _, err := io.Copy(has, f); err != nil {
							panic(err)
						}
						f.Close()

						if hashes, ok := modMap["hashes"].(map[string]interface{}); ok {
							fhas := hashes["sha512"].(string)
							if fhas != hex.EncodeToString(has.Sum(nil)) {
								color.Set(color.FgRed)
								fmt.Println("Warning: Potential Fake Resource Pack")
								fmt.Println("The resource pack hash doesn't match what is recorded in the .mrpack, which may indicate a fake or modified version. Please verify the resource pack’s source and ensure it’s from a trusted provider. (e.g., Modrinth)")
								fmt.Println("File deleted.")
								color.Unset()
								os.Remove(packFolder + "resourcepacks/" + strings.Split(path, "/")[1])
							}
						}
					}
					downloaded++
				}
			}
		}
	} else {
		fmt.Println("Expected 'files' to be a slice, but got ")
		os.Exit(1)
	}
}

func addOverrides(packFolder string, tempFolder string) {
	fmt.Println("Copy: "+tempFolder+"overrides", "->", packFolder)
	switch runtime.GOOS {
	case "linux":
		cmd, err := exec.Command("/bin/sh", "-c", "cp -r "+tempFolder+"overrides/* "+packFolder).Output()
		if err != nil {
			color.Set(color.FgRed)
			fmt.Println("Could not copy overrides: %v", err)
			color.Unset()
		}
		_ = cmd
	case "windows":
		cmd, err := exec.Command("xcopy", tempFolder+"overrides\\*", packFolder, "/E").Output()
		if err != nil {
			color.Set(color.FgRed)
			fmt.Println("Could not copy overrides: %v", err)
			color.Unset()
		}
		_ = cmd
	}
}
