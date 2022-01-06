package realaddress

import (
	"bufio"
	"bytes"
	"embed"
	"fmt"
	"math/rand"
	"strings"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

const (
	addressFile    = "data/KEN_ALL_ROME.CSV"
	totalLineCount = 124523
)

var (
	// exclude towns if contains following words
	// because it might not be address
	excludeTown = []string{
		"以下に掲載がない場合",
		"（",
	}
	//go:embed data/KEN_ALL_ROME.CSV data/test1.CSV
	addresses embed.FS
)

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

func readLine(filePath string, lineNumber int) (string, error) {
	f, err := addresses.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	var (
		c       int
		s       = transform.NewReader(bytes.NewReader(f), japanese.ShiftJIS.NewDecoder())
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
