package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/konrin/modesdecoder"
)

func main() {
	var msgFilePath string
	flag.StringVar(&msgFilePath, "file", "", "")
	flag.Parse()

	if msgFilePath == "" {
		log.Fatalln("Msg file not found")
	}

	decoder := modesdecoder.NewDecoder(modesdecoder.CacheTtl)

	file, err := os.Open(msgFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		m := scanner.Text()
		if len(m) == 0 {
			continue
		}

		if m[0] == '*' {
			m = string(m[1 : len(m)-1])
		}

		if m[len(m)-1] == ';' {
			m = string(m[0 : len(m)-2])
		}

		if len(m) == 14 {
			continue
		}

		msg, err := modesdecoder.NewMessage(m, time.Now())
		if err != nil {
			println(err.Error())
			continue
		}

		err = decoder.Decode(msg)
		if err != nil {
			println(err.Error())
			continue
		}

		fmt.Printf("%+v\n\n", msg)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
