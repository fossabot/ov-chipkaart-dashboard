package main

import (
	"bufio"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// FileSystemCSVReader implements the CSVReader interface
type FileSystemCSVReader struct {
	basePath string
}

// NewFileSystemCSVReader initializes a new CSV file reader
func NewFileSystemCSVReader(basePath string) FileSystemCSVReader {
	return FileSystemCSVReader{
		basePath: basePath,
	}
}

// ReadAll returns the contents of a CSV file as struct
func (reader FileSystemCSVReader) ReadAll(fileID string) (records [][]string, err error) {
	CSVFile, err := os.Open(filepath.Join(reader.basePath, fileID))
	if err != nil {
		return records, errors.Wrapf(err, "could not open the csv fileID with id '%s'", fileID)
	}

	r := bufio.NewReader(CSVFile)

	for {
		line, err := r.ReadString('\n')
		if err == io.EOF {
			break
		}

		if err != nil {
			return records, errors.Wrapf(err, "could not read line %d from the CSV fileID %s", len(records)+1, fileID)
		}

		records = append(records, reader.fetchRecordsFromLine(line))
	}

	return records, err
}

// basic CSV file parser
func (reader FileSystemCSVReader) fetchRecordsFromLine(line string) (result []string) {
	prefix := ""
	prefixIsOpen := false
	for _, char := range line {
		if char == ';' || char == '\n' || char == '\r' {
			continue
		}
		if char == '"' && !prefixIsOpen {
			prefixIsOpen = true
			continue
		}

		if char == '"' && prefixIsOpen {
			result = append(result, prefix)
			prefixIsOpen = false
			prefix = ""
			continue
		}

		prefix += string(char)
	}

	return result
}
