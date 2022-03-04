package internal

import (
	"bufio"
	"os"
)

// Uniq is just like the shell command `sort | uniq` - it returns a set of unique lines
// found in a file.
func Uniq(path string) (map[string]bool, error) {
	res := map[string]bool{}
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}

		res[scanner.Text()] = true
	}
	return res, nil
}
