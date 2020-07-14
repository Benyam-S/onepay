package tools

import "os"

// RemoveFile is a method that removes a given file path from the assets folder.
func RemoveFile(filePath string) error {

	if err := os.Remove(filePath); err != nil {
		return err
	}

	return nil

}
