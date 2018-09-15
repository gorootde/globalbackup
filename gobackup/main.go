package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"./internal/backuprecord"
)

const datadirectory = "../../"
const chunksize = 1024 * 1024

func main() {

	var entries []backuprecord.Backuprecord

	err := filepath.Walk(datadirectory,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			isSymlink := info.Mode()&os.ModeSymlink != 0
			if !info.IsDir() && !isSymlink {

				// abspath, _ := filepath.Abs(path)
				// record := backuprecord.New(abspath, info, chunksize)
				// fmt.Println(path, record.Checksum, len(record.ChunkSums))
				// entries = append(entries, record)
			}

			return nil
		})
	if err != nil {
		log.Println(err)
	}

	fmt.Printf("Backup contains %d files", len(entries))
	/*enc := json.NewEncoder(os.Stdout)
	if err := enc.Encode(entries); err != nil {
		log.Println(err)
	}
	*/
}
