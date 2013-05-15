package main

import (
	"os"
	"fmt"
	"log"
	"net/http"
	"encoding/json"
	"encoding/base64"
	"bufio"
	"strings"
	"flag"
)

type Val struct {
	Ref string
}

type Row struct {
	Id string
	Key string
	Value Val
	Doc json.RawMessage
}

type Response struct {
	TotalRows int `json:"total_rows"`
	Offset int
//	Rows []Row
}

func ParseResponse(line string) (response *Response, err error) {
	if line[len(line)-1:len(line)] == ":" {
		line = line + "["
	}
	line = line + "]}"
	response = &Response{}
	err = json.Unmarshal([]byte(line), response)
	return 
}

func main() {
	flag.Usage = func() {
                fmt.Fprintf(os.Stderr, "Usage of %s [options] [url]:\n", os.Args[0])
		flag.PrintDefaults()
        }
	flag.Parse()
	url := flag.Arg(0)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	r := bufio.NewReader(resp.Body)
	line, err := r.ReadString('\n')
	firstLine := true
	for err == nil {
		line = line[:len(line)-3]
		if firstLine {
			firstLine = false
			response, err := ParseResponse(line)
			if err != nil { log.Fatal(err) }
			if response.TotalRows == 0 {
				log.Fatal("0 rows to parse")
			}
		}

		row := &Row{}
		err = json.Unmarshal([]byte(line), row)
		if err != nil {
			if strings.TrimSpace(line) != "" {
				line = line + "}"
				err = json.Unmarshal([]byte(line), row)
				if err != nil {
					fmt.Println(err)
				}
			}

		}

		doc := row.Doc
		encodedDoc := base64.StdEncoding.EncodeToString(doc)
		if row.Id != "" {
			fmt.Print("SET ")
			fmt.Print(row.Id, " ")
			fmt.Print(encodedDoc, "\n")
		}

		line, err = r.ReadString('\n')
	}
}