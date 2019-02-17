package main

import (
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"log"
)

var storePath = "/tmp/recordstore"

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	dat, err := ioutil.ReadFile(storePath)
	check(err)

	record := &Record{}
	err = proto.Unmarshal(dat, record)
	check(err)

	log.Printf("Read from disk: %s\n", record)

	r := &Record{Word: proto.String("abc"), Count: proto.Int32(1), Resolved: proto.Bool(true)}
	marshalled, err := proto.Marshal(r)
	check(err)

	ioutil.WriteFile(storePath, marshalled, 0644)
	log.Printf("Written to disk: %s\n", r)
}
