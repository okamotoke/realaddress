package realaddress

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

const (
	addressFile    = "../realaddress/data/KEN_ALL_ROME.CSV"
	totalLineCount = 124523
)

// exclude towns if contains following words
// because it might not address
var excludeTown = []string{
	"以下に掲載がない場合",
	"（",
}

type Address struct {
	postalCode string
	prefecture string
	city       string
	town       string
}

func (a *Address) GetPostalCode() string {
	return a.postalCode
}

func (a *Address) GetPrefecture() string {
	return a.prefecture
}

func (a *Address) GetCity() string {
	return a.city
}

func (a *Address) GetTown() string {
	return a.town
}

func GetRandomAddress() (Address, error) {
	return getRandomAddress(addressFile, totalLineCount)
}

func getRandomAddress(filePath string, lineCount int) (Address, error) {
	var retryCount int
	for true {
		randomLineNum := rand.Intn(lineCount) + 1

		line, err := readLine(filePath, randomLineNum)
		if err != nil {
			return Address{}, err
		}

		address := strings.Split(line, ",")
		if len(address) != 7 {
			// 本来起らないはずだが念の為
			retryCount++
			if retryCount <= 3 {
				continue
			}
			return Address{}, fmt.Errorf("unexpceted line")
		}

		if skipAddress(address[3]) {
			continue
		}

		return Address{
			postalCode: strings.Trim(address[0], `"`),
			prefecture: strings.Trim(address[1], `"`),
			city:       strings.Trim(address[2], `"`),
			town:       strings.Trim(address[3], `"`),
		}, nil
	}
	return Address{}, nil
}

func skipAddress(town string) bool {
	for _, t := range excludeTown {
		if strings.Contains(town, t) {
			return true
		}
	}
	return false
}

func countLine(filePath string) (int, error) {
	p, err := filepath.Abs(filePath)
	if err != nil {
		return 0, err
	}

	f, err := os.OpenFile(p, os.O_RDONLY, os.ModeDir)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	var (
		count int
		buf   = make([]byte, bufio.MaxScanTokenSize)
	)

	for {
		c, err := f.Read(buf)
		if err != nil && err != io.EOF {
			return 0, fmt.Errorf("failed to count line: %v", err)
		}
		if err == io.EOF {
			return count, nil
		}
		count += bytes.Count(buf[:c], []byte{'\n'})
	}
}

func readLine(filePath string, lineNumber int) (string, error) {
	p, err := filepath.Abs(filePath)
	if err != nil {
		return "", err
	}

	f, err := os.OpenFile(p, os.O_RDONLY, os.ModeDir)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var (
		c       int
		s       = transform.NewReader(f, japanese.ShiftJIS.NewDecoder())
		scanner = bufio.NewScanner(s)
	)

	for scanner.Scan() {
		c++
		if c == lineNumber {
			return scanner.Text(), nil
		}
	}
	return "", fmt.Errorf("unexpected line number %d", lineNumber)
}
