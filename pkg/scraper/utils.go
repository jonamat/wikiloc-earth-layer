package scraper

import (
	"bufio"
	"strings"
)

func StringToLines(s *string) ([]string, error) {
	var lines []string

	scanner := bufio.NewScanner(strings.NewReader(*s))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
