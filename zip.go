package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func unzip(source, dest string) error {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("tar", "-xf", source, "-C", dest)
		cmd.Run()
	case "linux":
		if _, err := os.Stat("/bin/unzip"); err == nil {
			//fmt.Println("Running unzip")
			cmd := exec.Command("unzip", source, "-d", dest)
			cmd.Run()
			//fmt.Println(cmd.Output())
		} else {
			return gounzip(source, dest)
		}
	}

	return nil
}

func gounzip(source, dest string) error {
	fmt.Println("unzip not found, using Go unzip which will NOT extract overrides")
	r, err := zip.OpenReader(source)
	if err != nil {
		return fmt.Errorf("failed to open ZIP file: %w", err)
	}
	defer r.Close()

	if err := os.MkdirAll(dest, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return fmt.Errorf("failed to open file in ZIP: %w", err)
		}

		fPath := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(fPath, f.Mode()); err != nil {
				rc.Close()
				return fmt.Errorf("failed to create directory: %w", err)
			}
			rc.Close()
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fPath), 0755); err != nil {
			rc.Close()
			return fmt.Errorf("failed to create directory for file: %w", err)
		}

		outFile, err := os.Create(fPath)
		if err != nil {
			rc.Close()
			return fmt.Errorf("failed to create file: %w", err)
		}

		if _, err := io.Copy(outFile, rc); err != nil {
			outFile.Close()
			rc.Close()
			return fmt.Errorf("failed to copy file contents: %w", err)
		}

		outFile.Close()
		rc.Close()
	}
	return nil
}
