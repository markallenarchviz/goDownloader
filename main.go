package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"github.com/mholt/archiver/v4"
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

// extractZIP extracts a ZIP file to the specified destination directory.
func extractZIP(zipPath, destDir string) error {
	zip := archiver.Zip{}

	// Open the ZIP file
	file, err := os.Open(zipPath)
	if err != nil {
		return err
	}
	defer file.Close()

	return zip.Extract(context.Background(), file, []string{destDir}, nil)
}

// createShortcut creates a desktop shortcut for the specified executable.
func createShortcut(exePath, shortcutPath string) error {
	ole.CoInitialize(0)
	defer ole.CoUninitialize()

	unknown, err := oleutil.CreateObject("WScript.Shell")
	if err != nil {
		return err
	}
	defer unknown.Release()

	shell, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return err
	}
	defer shell.Release()

	shortcut, err := oleutil.CallMethod(shell, "CreateShortcut", shortcutPath)
	if err != nil {
		return err
	}
	shortcutObj := shortcut.ToIDispatch()
	defer shortcutObj.Release()

	_, err = oleutil.PutProperty(shortcutObj, "TargetPath", exePath)
	if err != nil {
		return err
	}

	_, err = oleutil.CallMethod(shortcutObj, "Save")
	if err != nil {
		return err
	}

	return nil
}

func main() {
	zipURL := "C:\\checklistMakerGo.zip" // Replace with your ZIP file URL
	zipPath := "downloaded.zip"
	destDir := "C:\\teste"            // Replace with your destination directory
	exeName := "checklistMakerGo.exe" // Replace with your executable name

	// Create destination directory if it doesn't exist
	err := os.MkdirAll(destDir, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating destination directory:", err)
		return
	}

	// Download the ZIP file
	fmt.Println("Downloading ZIP file...")
	err = downloadFile(zipPath, zipURL)
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

	// Create desktop shortcut
	exePath := filepath.Join(destDir, exeName)
	desktopDir := filepath.Join(os.Getenv("USERPROFILE"), "Desktop")
	shortcutPath := filepath.Join(desktopDir, exeName+".lnk")

	fmt.Println("Creating desktop shortcut...")
	err = createShortcut(exePath, shortcutPath)
	if err != nil {
		fmt.Println("Error creating shortcut:", err)
		return
	}
	fmt.Println("Created desktop shortcut:", shortcutPath)
}
