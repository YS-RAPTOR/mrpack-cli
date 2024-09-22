package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/fatih/color"
)

type MineLauncher struct {
	Profiles []map[string]Profile
}

type Profile struct {
	Created       string `json:"created"`
	Name          string `json:"name"`
	Type          string `json:"type"`
	Icon          string `json:"icon"`
	LastUsed      string `json:"lastUsed"`
	GameDirectory string `json:"gameDir"`
	LastVersionID string `json:"lastVersionId"`
}

func addEntry(packfolder, packname, fancypackname, gamever, loaderver, loader string) error {
	var launcherfolder string
	var versionId string

	userhome, err := os.UserHomeDir()
	if err != nil {
		color.Set(color.FgRed)
		fmt.Println("ERROR: Could not get home directory")
		color.Unset()
		return err
	}
	switch runtime.GOOS {
	case "windows":
		launcherfolder = os.Getenv("APPDATA") + "\\.minecraft\\"
	case "linux":
		launcherfolder = userhome + "/.minecraft/"
	}

	color.Set(color.FgGreen)
	fmt.Println("Getting icon from Modrinth...")
	color.Unset()

	n, err := http.Get("https://api.modrinth.com/v2/project/" + packname)
	if err != nil {
		color.Set(color.FgRed)
		fmt.Println("ERROR: Could not communicate with Modrinth API:", err)
		color.Unset()
		return err
	}

	api, err := io.ReadAll(n.Body)
	if err != nil {
		color.Set(color.FgRed)
		fmt.Println("ERROR: Could not read Modrinth API response body:", err)
		color.Unset()
		return err
	}

	var iconURI string

	var js = openMRjson(string(api))
	uri := js.IconURL
	n, err = http.Get(uri)
	if err != nil {
		color.Set(color.FgRed)
		fmt.Println("ERROR: Could not get icon:", err)
		color.Unset()
	}

	d, err := io.ReadAll(n.Body)
	if err != nil {
		color.Set(color.FgRed)
		fmt.Println("ERROR: Could not read icon:", err)
		color.Unset()
	}

	img := base64.StdEncoding.EncodeToString(d)

	iconURI = fmt.Sprintf("data:image/png;base64," + img)

	color.Set(color.FgGreen)
	fmt.Println("Adding entry to Minecraft launcher")
	color.Unset()

	switch loader {
	case "neoforge":
		versionId = loader + "-" + loaderver
	case "fabric-loader":
		versionId = loader + "-" + loaderver + "-" + gamever
	case "forge":
		versionId = gamever + "-" + loader + "-" + loaderver
	case "quilt-loader":
		versionId = loader + "-" + loaderver + "-" + gamever
	}

	var ljs = openMCjson(launcherfolder + "launcher_profiles.json")
	ljs = MineLauncher{
		Profiles: []map[string]Profile{
			{
				packname: {
					Name:          fancypackname,
					Type:          "custom",
					Created:       time.Now().Format(time.RFC3339),
					LastUsed:      time.Time{}.String(),
					Icon:          iconURI,
					GameDirectory: packfolder,
					LastVersionID: versionId,
				},
			},
		},
	}

	ujs, err := json.MarshalIndent(ljs, "", "  ")
	if err != nil {
		color.Set(color.FgRed)
		fmt.Println("ERROR: marshalling JSON:", err)
		color.Unset()
		return err
	}
	_ = ujs

	err = os.WriteFile(launcherfolder+"launcher_profiles.json", ujs, 0664)
	if err != nil {
		color.Set(color.FgRed)
		fmt.Println("ERROR: writing file:", err)
		color.Unset()
		return err
	}

	return nil
}
