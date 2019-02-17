package main

import (
	"bufio"
	"github.com/golang/protobuf/proto"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var storePath = "/tmp/recordstore"

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func storeExists() bool {
	_, err := os.Stat(storePath)
	return !os.IsNotExist(err)
}

func initializeStore(store *RecordStore) {
	if !storeExists() {
		log.Println("No existing store found. Will initialize a new one.")
	} else {
		dat, err := ioutil.ReadFile(storePath)
		if err != nil {
			log.Fatalln("Couldn't open store. Exiting.")
			panic(err)
		}
		err = proto.Unmarshal(dat, store)
		if err != nil {
			log.Fatalln("Couldn't deserialize store. Exiting.")
			panic(err)
		}
		log.Printf("Store successfully retrieved from disk. Entries in store: %d\n", len(store.Record))
	}

	if store.Record == nil {
		store.Record = make(map[string]*Record)
	}
}

func writeStore(store *RecordStore) {
	// TODO: error handling here is pretty bad, attempt to save to a
	// couple other places (/tmp?, $HOME?) to avoid data loss.
	marshalled, err := proto.Marshal(store)
	if err != nil {
		log.Fatalln("Couldn't serialize new store. Exiting.")
		panic(err)
	}

	err = ioutil.WriteFile(storePath, marshalled, 0644)
	if err != nil {
		log.Fatalln("Couldn't write new store to disk. Exiting.")
		panic(err)
	}
	log.Printf("Store successfully written to disk. Entries in store: %d\n", len(store.Record))
}

func main() {
	store := &RecordStore{}

	initializeStore(store)

	reader := bufio.NewReader(os.Stdin)
	for {
		text, err := reader.ReadString('\n')
		if err == io.EOF {
			log.Println("Quitting...")
			break
		}

		text = strings.TrimSpace(text)

		if text == "quit" {
			log.Println("Quitting...")
			break
		}

		if store.Record[text].GetCount() == 0 {
			store.Record[text] = &Record{}
		}

		store.Record[text].Count = proto.Int32(store.Record[text].GetCount() + 1)
	}

	writeStore(store)
}
