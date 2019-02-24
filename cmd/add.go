
package cmd

import (
	"bufio"
	"github.com/Ryckes/godict/config"
	"github.com/Ryckes/godict/storage"
	"github.com/golang/protobuf/proto"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"io"
	"log"
	"os"
	"strings"
)

func init() {
	rootCmd.AddCommand(addCommand)
}

var addCommand = &cobra.Command{
  Use:   "add",
  Short: "Start new recording session",
  Long:  `Start new recording session`,
  Run: func(cmd *cobra.Command, args []string) {
	config := config.ReadConfig()
	storePath, err := homedir.Expand(config.GetStorePath())
	if err != nil {
		log.Fatalf("Invalid store path: %s.\n", err)
	}

	store := &storage.RecordStore{}

	storage.InitializeStore(storePath, store)

	reader := bufio.NewReader(os.Stdin)
	var nonpersistedRecords int32 = 0
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
			store.Record[text] = &storage.Record{}
		}

		store.Record[text].Count = proto.Int32(store.Record[text].GetCount() + 1)

		nonpersistedRecords += 1
		if nonpersistedRecords > config.GetMaxNonpersistedRecords() {
			storage.WriteStore(storePath, store)
			nonpersistedRecords = 0
		}
	}

	if nonpersistedRecords > 0 {
		storage.WriteStore(storePath, store)
	}
  },
}
