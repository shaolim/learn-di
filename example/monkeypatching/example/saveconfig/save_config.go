package saveconfig

import (
	"encoding/json"
	"io/fs"
	"log"
	"os"
)

func SaveConfig(filename string, cfg *Config) error {
	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	// save file
	err = writeFile(filename, data, 0666)
	if err != nil {
		log.Printf("failed to save file '%s' with err: %s", filename, err)
		return err
	}

	return nil
}

type fileWriter func(filename string, data []byte, perm fs.FileMode) error

var writeFile fileWriter = os.WriteFile

type Config struct {
	Host string
	Port string
}
