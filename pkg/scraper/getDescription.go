package scraper

import (
	"fmt"
	"strings"
)

func GetDescription(html *string) (string, error) {
	// Split html code into lines
	lines, err := StringToLines(html)
	if err != nil {
		return "", err
	}

	// Search the line that contains "<div class="description dont-break-out ">"
	var rowData string
	for _, line := range lines {
		if strings.Contains(line, "<div class=\"description dont-break-out \">") {
			// Keep from <div> to </div>
			index := strings.Index(line, ">")
			rowData = line[index+1 : len(line)-6]
			break
		}
	}
	if len(rowData) == 0 {
		return "", fmt.Errorf("description not found")
	}

	return rowData, nil
}
