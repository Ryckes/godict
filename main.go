package main

import (
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"log"
	"os"
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
			log.Fatalf("Couldn't open store: %s. Exiting.\n", err)
			return
		}
		err = proto.Unmarshal(dat, store)
		if err != nil {
			log.Fatalf("Couldn't deserialize store: %s. Exiting.\n", err)
			return
		}
		log.Printf("Read from disk: %s\n", store)
	}
}

func writeStore(store *RecordStore) {
	// TODO: error handling here is pretty bad, attempt to save to a
	// couple other places (/tmp?, $HOME?) to avoid data loss.
	// TODO: remove double logging on errors.
	marshalled, err := proto.Marshal(store)
	if err != nil {
		// TODO: this is pretty bad, attempt to save to a couple other
		// places to avoid data loss.
		log.Fatalf("Couldn't serialize new store: %s. Exiting.\n", err)
		panic(err)
	}

	err = ioutil.WriteFile(storePath, marshalled, 0644)
	if err != nil {
		log.Fatalf("Couldn't write new store to disk: %s. Exiting.\n", err)
		panic(err)
	}
	log.Printf("Written to disk: %s\n", store)
}

func main() {
	store := &RecordStore{}

	initializeStore(store)

	r := &RecordStore{Record: []*Record{{Word: proto.String("abc"), Count: proto.Int32(1)}}}
	writeStore(r)
}
