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

func installQuilt(tempfolder, gameversion, loaderver string) error {
	var java_home = "java"
	if runtime.GOOS == "windows" {
		java_home = os.Getenv("JAVA_HOME") + "\\bin\\java.exe"
	}

	out, err := os.Create(tempfolder + "/quilt-installer.jar")
	if err != nil {
		color.Set(color.FgRed, color.Bold)
		fmt.Println("ERROR: Could not create Quilt installer file:", err)
		color.Unset()
		return err
	}

	color.Set(color.FgCyan)
	fmt.Println("Downloading Quilt installer...")
	color.Unset()

	resp, err := http.Get("https://quiltmc.org/api/v1/download-latest-installer/java-universal")
	if err != nil {
		color.Set(color.FgRed, color.Bold)
		fmt.Println("ERROR: Could not download Quilt installer from "+resp.Request.URL.RequestURI(), "due to:", err)
		color.Unset()
		return err
	}

	n, err := io.Copy(out, resp.Body)
	if err != nil {
		color.Set(color.FgRed, color.Bold)
		fmt.Println("ERROR: Could not copy Quilt installer:", err)
		color.Unset()
		return err
	}
	_ = n

	color.Set(color.FgCyan, color.Bold)
	fmt.Println("Running Quilt installer... (This might take a while)")
	color.Unset()

	fi := exec.Command(java_home, "-jar", tempfolder+"quilt-installer.jar", "install", "client", gameversion, loaderver, "--no-profile")

	output, err := fi.CombinedOutput()
	if err != nil {
		fmt.Printf("Quilt installer failed:", err)
	}

	fmt.Printf("%s\n", output)

	return nil
}
