package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/fatih/color"
)

func installForge(tempfolder, launcherfolder, gameversion, loaderver string) error {
	var java_home = "java"
	if runtime.GOOS == "windows" {
		java_home = os.Getenv("JAVA_HOME") + "\\bin\\java.exe"
	}

	out, err := os.Create(tempfolder + "/forge-installer.jar")
	if err != nil {
		color.Set(color.FgRed, color.Bold)
		fmt.Println("ERROR: Could not create Forge installer file:", err)
		color.Unset()
		return err
	}

	color.Set(color.FgCyan)
	fmt.Println("Downloading Forge installer...")
	color.Unset()

	resp, err := http.Get("https://maven.minecraftforge.net/net/minecraftforge/forge/" + gameversion + "-" + loaderver + "/forge-" + gameversion + "-" + loaderver + "-installer.jar")
	if err != nil {
		color.Set(color.FgRed, color.Bold)
		fmt.Println("ERROR: Could not download Forge installer from "+resp.Request.URL.RequestURI(), "due to:", err)
		color.Unset()
		return err
	}

	n, err := io.Copy(out, resp.Body)
	if err != nil {
		color.Set(color.FgRed, color.Bold)
		fmt.Println("ERROR: Could not copy Forge installer:", err)
		color.Unset()
		return err
	}
	_ = n

	color.Set(color.FgCyan, color.Bold)
	fmt.Println("Running Forge installer... (This might take a while)")
	color.Unset()

	fi := exec.Command(java_home, "-jar", tempfolder+"forge-installer.jar", "--installClient", launcherfolder)

	output, err := fi.CombinedOutput()
	if err != nil {
		fmt.Printf("Forge installer failed:", err)
	}

	fmt.Printf("%s\n", output)

	return nil
}
