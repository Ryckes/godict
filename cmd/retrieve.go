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
	"math/rand"
	"os"
	"strings"
	"time"
)

func init() {
	rootCmd.AddCommand(retrieveCommand)
}

var retrieveCommand = &cobra.Command{
	Use:   "retrieve",
	Short: "Start new retrieval session",
	Long: `Start new retrieval session.

Unresolved records will be retrieved in random order, and you can
enter 'r' to mark them as resolved, or 's' to skip the last record
retrieved without resolving.

Enter 'quit' or press Ctrl-D to exit.`,
	Run: func(cmd *cobra.Command, args []string) {
		config := config.ReadAndMaybeCreateConfig()
		storePath, err := homedir.Expand(config.GetStorePath())
		if err != nil {
			log.Fatalf("Invalid store path: %s.\n", err)
		}

		store := &storage.RecordStore{}

		storage.InitializeStore(storePath, store)

		unresolvedKeys := make([]string, 0, len(store.Record))
		for key, value := range store.Record {
			if !value.GetResolved() {
				unresolvedKeys = append(unresolvedKeys, key)
			}
		}
		log.Printf("There are %d unresolved records.\n", len(unresolvedKeys))

		// Shuffle the keys.
		rand.Seed(time.Now().UnixNano())
		for i := len(unresolvedKeys) - 1; i > 0; i-- {
			j := rand.Intn(i + 1)
			unresolvedKeys[i], unresolvedKeys[j] = unresolvedKeys[j], unresolvedKeys[i]
		}

		reader := bufio.NewReader(os.Stdin)
		var nonpersistedChanges int32 = 0
		current := 0
		for {
			if len(unresolvedKeys) == 0 {
				log.Println("There are no unresolved records left.")
				break
			}

			if current >= len(unresolvedKeys) {
				current = 0
			}
			currentKey := unresolvedKeys[current]
			log.Println(currentKey)
			log.Printf("Count: %d\n", store.Record[currentKey].GetCount())

			text, err := reader.ReadString('\n')
			if err == io.EOF {
				log.Println("Quitting...")
				break
			}

			text = strings.TrimSpace(text)

			switch text {
			case "quit":
				log.Println("Quitting...")
				break
			case "r":
				store.Record[currentKey].Resolved = proto.Bool(true)
				swap(unresolvedKeys, current, len(unresolvedKeys)-1)
				unresolvedKeys = unresolvedKeys[:len(unresolvedKeys)-1]
				nonpersistedChanges += 1
				// Don't increment current, since current now points
				// to an unresolved key.
			case "s":
				current++
			}

			if nonpersistedChanges > config.GetMaxNonpersistedChanges() {
				storage.WriteStore(storePath, store)
				nonpersistedChanges = 0
				log.Printf("There are %d unresolved records.\n", len(unresolvedKeys))
			}
		}

		if nonpersistedChanges > 0 {
			storage.WriteStore(storePath, store)
		}
	},
}

func swap(a []string, i, j int) {
	a[i], a[j] = a[j], a[i]
}
