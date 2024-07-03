package main

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// downloadFile downloads a file from the given URL and saves it to the specified path.
func downloadFile(filepath string, url string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

// extractZIP extracts files from a ZIP archive to the specified destination directory.
func extractZIP(zipPath, destDir string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	// Create destination directory if it doesn't exist
	err = os.MkdirAll(destDir, os.ModePerm)
	if err != nil {
		return err
	}

	// Iterate through each file in the ZIP archive
	for _, f := range r.File {
		// Check if the file is within the cheklistMakerFiles-main folder
		if strings.HasPrefix(f.Name, "cheklistMakerFiles-main/") {
			// Trim the prefix to get relative path within cheklistMakerFiles-main
			relPath := strings.TrimPrefix(f.Name, "cheklistMakerFiles-main/")
			path := filepath.Join(destDir, relPath)

			if f.FileInfo().IsDir() {
				// Create directory
				os.MkdirAll(path, os.ModePerm)
			} else {
				// Create file's directory
				dir := filepath.Dir(path)
				os.MkdirAll(dir, os.ModePerm)

				// Open file from ZIP
				rc, err := f.Open()
				if err != nil {
					return err
				}
				defer rc.Close()

				// Create destination file
				w, err := os.Create(path)
				if err != nil {
					return err
				}
				defer w.Close()

				// Copy contents from ZIP file to destination file
				_, err = io.Copy(w, rc)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func main() {
	zipURL := "https://github.com/markallenarchviz/cheklistMakerFiles/archive/refs/heads/main.zip" // Replace with your ZIP file URL
	zipPath := "downloaded.zip"
	destDir := "C:\\checklistMaker" // Replace with your destination directory

	// Download the ZIP file
	fmt.Println("Downloading ZIP file...")
	err := downloadFile(zipPath, zipURL)
	if err != nil {
		fmt.Println("Error downloading file:", err)
		return
	}
	fmt.Println("Downloaded ZIP file:", zipPath)

	// Extract the ZIP file
	fmt.Println("Extracting ZIP file...")
	err = extractZIP(zipPath, destDir)
	if err != nil {
		fmt.Println("Error extracting ZIP file:", err)
		return
	}
	fmt.Println("Extracted files to:", destDir)

	// Delete the ZIP file
	fmt.Println("Deleting ZIP file...")
	err = os.Remove(zipPath)
	if err != nil {
		fmt.Println("Error deleting ZIP file:", err)
		return
	}
	fmt.Println("Deleted ZIP file:", zipPath)

	// Copy checklistMakerGo.lnk to desktop
	src := filepath.Join(destDir, "checklistMakerFiles-main", "checklistMakerGo.lnk")
	dst := filepath.Join(os.Getenv("USERPROFILE"), "Desktop", "checklistMakerGo.lnk")

	err = copyFile(src, dst)
	if err != nil {
		fmt.Println("Error copying file to desktop:", err)
		return
	}
	fmt.Println("Copied checklistMakerGo.lnk to desktop.")
}

// copyFile copies a file from src to dst.
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
