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

func storeExists(storePath string) bool {
	_, err := os.Stat(storePath)
	return !os.IsNotExist(err)
}

func initializeStore(storePath string, store *RecordStore) {
	log.Printf("Store path: %s.\n", storePath)
	if !storeExists(storePath) {
		log.Println("No existing store found. Will initialize a new one.")
	} else {
		dat, err := ioutil.ReadFile(storePath)
		if err != nil {
			log.Fatalf("Couldn't open store: %s. Exiting.\n", err)
		}
		err = proto.Unmarshal(dat, store)
		if err != nil {
			log.Fatalf("Couldn't deserialize store: %s. Exiting.\n", err)
		}
		log.Printf("Store successfully retrieved from disk. Entries in store: %d\n", len(store.Record))
	}

	if store.Record == nil {
		store.Record = make(map[string]*Record)
	}
}

func writeStore(storePath string, store *RecordStore) {
	// TODO: error handling here is pretty bad, attempt to save to a
	// couple other places (/tmp?, $HOME?) to avoid data loss.
	marshalled, err := proto.Marshal(store)
	if err != nil {
		log.Fatalf("Couldn't serialize new store: %s. Exiting.\n", err)
	}

	err = ioutil.WriteFile(storePath, marshalled, 0644)
	if err != nil {
		log.Fatalf("Couldn't write new store to disk: %s. Exiting.\n", err)
	}
	log.Printf("Store successfully written to disk. Entries in store: %d\n", len(store.Record))
}

func readConfig(configPath string) *Configuration {
	config := &Configuration{}
	data, _ := ioutil.ReadFile(configPath)

	proto.UnmarshalText(string(data), config)
	return config
}

func main() {
	config := readConfig("./config.textproto")
	storePath := config.GetStorePath()
	
	store := &RecordStore{}

	initializeStore(storePath, store)

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

	writeStore(storePath, store)
}
