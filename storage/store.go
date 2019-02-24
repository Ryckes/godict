
package storage

import (
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"log"
	"os"
)

func storeExists(storePath string) bool {
	_, err := os.Stat(storePath)
	return !os.IsNotExist(err)
}

func InitializeStore(storePath string, store *RecordStore) {
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

func WriteStore(storePath string, store *RecordStore) {
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
