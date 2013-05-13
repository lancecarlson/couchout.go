package main

import (
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
	TotalRows int
	Offset int
	Rows []Row
}

func main() {
	url := flag.String("url", "", "for the couchdb database")
	flag.Parse()
	resp, err := http.Get(*url)
	if err != nil {
		log.Fatal(err)
	}
	r := bufio.NewReader(resp.Body)
	line, err := r.ReadString('\n')
	for err == nil {
		line = line[:len(line)-3]

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