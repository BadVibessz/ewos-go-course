package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

const (
	fileTag        = "file"
	defaultFileTag = ""
)

// todo: encapsulate
func parseJson(j string) []map[string]any {

	var results []map[string]any

	err := json.Unmarshal([]byte(j), &results)
	if err != nil {
		return nil // todo:
	}

	for key, result := range results {
		address := result["operation"].(map[string]interface{})

		fmt.Println("Reading Value for Key :", key)

		fmt.Println("company :", result["company"],
			"- type :", result["type"],
			"- id :", result["id"])

		fmt.Println("operation :", address["created_at"])
	}
	return results // todo: maybe map resulting object to Billing{} struct?
}

func main() {
	var filename string
	flag.StringVar(&filename, fileTag, defaultFileTag, "file to read.")
	flag.Parse()

	if filename == defaultFileTag {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalf("Some error occured. Err: %s", err) // todo: do i need to exit here? or maybe just fallthrough to stdin?
		}

		val, exist := os.LookupEnv("FILE_NAME")
		if !exist || val == "" {
			fmt.Println("Enter path to file")

			reader := bufio.NewReader(os.Stdin)

			str, readErr := reader.ReadString('\n')
			if readErr != nil {
				log.Fatal(readErr) // todo:
			}

			filename = str
		} else {
			filename = val
		}
	}

	println(filename)

	json, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err) // todo:
	}

	_ = parseJson(string(json))
}
