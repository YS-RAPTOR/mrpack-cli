package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/fatih/color"
)

func installfabric(tempfolder, gameversion, loaderver string) error {
	var java_home = "java"
	if runtime.GOOS == "windows" {
		java_home = os.Getenv("JAVA_HOME") + "\\bin\\java.exe"
	}

	out, err := os.Create(tempfolder + "/fabric-installer.jar")
	if err != nil {
		color.Set(color.FgRed, color.Bold)
		fmt.Println("ERROR: Could not create fabric installer file")
		color.Unset()
		return err
	}

	color.Set(color.FgCyan)
	fmt.Println("Downloading Fabric installer...")
	color.Unset()

	resp, err := http.Get("https://maven.fabricmc.net/net/fabricmc/fabric-installer/1.0.1/fabric-installer-1.0.1.jar")
	if err != nil {
		color.Set(color.FgRed, color.Bold)
		fmt.Println("ERROR: Could not download fabric installer from " + resp.Request.URL.RequestURI())
		color.Unset()
		return err
	}

	n, err := io.Copy(out, resp.Body)
	if err != nil {
		color.Set(color.FgRed, color.Bold)
		fmt.Println("ERROR: Could not copy fabric installer")
		color.Unset()
		return err
	}
	_ = n

	color.Set(color.FgCyan, color.Bold)
	fmt.Println("Running Fabric installer... (This might take a while)")
	color.Unset()

	fi := exec.Command(java_home, "-jar", tempfolder+"fabric-installer.jar", "client", "-mcversion", gameversion, "-loader", loaderver, "-noprofile")

	output, err := fi.CombinedOutput()
	if err != nil {
		fmt.Printf("Fabric installer failed: %v\n", err)
	}

	fmt.Printf("%s\n", output)

	return nil
}

func addFabricEntry(packfolder, packname, gamever, fabricver string) error {
	var launcherfolder string

	switch runtime.GOOS {
	case "windows":
		launcherfolder = os.Getenv("APPDATA") + "\\.minecraft\\"
	case "linux":
		launcherfolder = "~/.minecraft/"
	}

	color.Set(color.FgGreen)
	fmt.Println("Getting icon from Modrinth...")
	color.Unset()

	n, err := http.Get("https://api.modrinth.com/v2/project/" + packname)
	if err != nil {
		fmt.Println("Could not communicate with Modrinth API:", err)
		return err
	}

	api, err := io.ReadAll(n.Body)
	if err != nil {
		fmt.Println("Could not read Modrinth API response body:", err)
		return err
	}

	var iconURI string

	var js map[string]interface{} = openjsonfromstring(string(api))
	if uri, ok := js["icon_url"].(string); ok {
		n, err := http.Get(uri)
		if err != nil {
			fmt.Println("Could not get icon:", err)
			return err
		}

		d, err := io.ReadAll(n.Body)
		if err != nil {
			fmt.Println("Could not read icon:", err)
			return err
		}

		img := base64.StdEncoding.EncodeToString(d)

		iconURI = fmt.Sprintf("data:image/png;base64," + img)
	}

	color.Set(color.FgGreen)
	fmt.Println("Adding entry to Minecraft launcher")
	color.Unset()

	var ljs map[string]interface{} = openjson(launcherfolder + "launcher_profiles.json")
	if pf, ok := ljs["profiles"].(map[string]interface{}); ok {
		pf[packname] = map[string]interface{}{
			"name":          packname,
			"type":          "custom",
			"created":       time.Now().Format(time.RFC3339),
			"lastUsed":      time.Time{},
			"icon":          iconURI,
			"gameDir":       packfolder,
			"lastVersionId": "fabric-loader-" + fabricver + "-" + gamever,
		}
	}

	ujs, err := json.MarshalIndent(ljs, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return err
	}
	_ = ujs

	err = os.WriteFile(launcherfolder+"launcher_profiles.json", ujs, 0664)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return err
	}

	return nil
}
