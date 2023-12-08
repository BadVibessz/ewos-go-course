package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/ew0s/ewos-to-go-hw/objectModel2/homework/billing"
	"github.com/joho/godotenv"
)

const (
	fileFlag        = "file"
	defaultFileFlag = ""
)

const filenameEnvVar = "FILE_NAME"

const outputFilename = "out.json"

func getFilename() (string, error) {
	var filename string

	flag.StringVar(&filename, fileFlag, defaultFileFlag, "file to read.")
	flag.Parse()

	if filename == defaultFileFlag {
		err := godotenv.Load(".env")
		if err == nil {
			val, exist := os.LookupEnv(filenameEnvVar)
			if exist && val != "" {
				return val, nil
			}
		}

		fmt.Println("Enter path to file")

		reader := bufio.NewReader(os.Stdin)

		str, readErr := reader.ReadString('\n')
		if readErr != nil {
			return "", err
		}

		return str, nil
	}

	return filename, nil
}

func main() {
	filename, err := getFilename()
	if err != nil {
		log.Fatalf("cannot load filename, err: %s", err)
	}

	js, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	b, err := billing.ParseJson(string(js))
	if err != nil {
		log.Fatal("error occurred while parsing input: ", err)
	}

	infos := billing.CalculateBalances(b)

	sort.Slice(infos, func(i, j int) bool { return infos[i].Company < infos[j].Company })

	bytes, err := json.MarshalIndent(infos, "", "    ")
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(outputFilename, bytes, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
