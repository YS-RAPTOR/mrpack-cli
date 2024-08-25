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

func installNeoforge(tempfolder, loaderver string) error {
	var java_home = "java"
	if runtime.GOOS == "windows" {
		java_home = os.Getenv("JAVA_HOME") + "\\bin\\java.exe"
	}

	out, err := os.Create(tempfolder + "/neoforge-installer.jar")
	if err != nil {
		color.Set(color.FgRed, color.Bold)
		fmt.Println("ERROR: Could not create Neoforge installer file:", err)
		color.Unset()
		return err
	}

	color.Set(color.FgCyan)
	fmt.Println("Downloading Neoforge installer...")
	color.Unset()

	resp, err := http.Get("https://maven.neoforged.net/releases/net/neoforged/neoforge/" + loaderver + "/neoforge-" + loaderver + "-installer.jar")
	if err != nil {
		color.Set(color.FgRed, color.Bold)
		fmt.Println("ERROR: Could not download Neoforge installer from "+resp.Request.URL.RequestURI(), "due to:", err)
		color.Unset()
		return err
	}

	n, err := io.Copy(out, resp.Body)
	if err != nil {
		color.Set(color.FgRed, color.Bold)
		fmt.Println("ERROR: Could not copy Neoforge installer:", err)
		color.Unset()
		return err
	}
	_ = n

	color.Set(color.FgCyan, color.Bold)
	fmt.Println("Running Neoforge installer... (This might take a while)")
	color.Unset()

	fi := exec.Command(java_home, "-jar", tempfolder+"neoforge-installer.jar", "--install-client")

	output, err := fi.CombinedOutput()
	if err != nil {
		fmt.Printf("Neoforge installer failed:", err)
	}

	fmt.Printf("%s\n", output)

	return nil
}
