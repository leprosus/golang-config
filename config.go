package tcjconfig

import (
	"path/filepath"
	"io/ioutil"
	"github.com/tidwall/gjson"
	"os"
)

type config struct {
	FilePath string
	json     string
}

var cfg config = config{}

func Init(filePath string) {
	if len(cfg.FilePath) == 0 {
		cfg = config{
			FilePath: filePath,
		}

		cfg.loadJson()
	}
}

func getResult(path string) gjson.Result {
	return gjson.Get(cfg.json, path)
}

func String(path string) string {
	return getResult(path).String()
}

func Int(path string) int64 {
	return getResult(path).Int()
}

func Array(path string) (slice []string) {
	result := getResult(path)

	for _, el := range result.Array() {
		slice = append(slice, el.String())
	}

	return slice
}

func (config *config) loadJson() {
	if len(config.json) == 0 {
		fileName, _ := filepath.Abs(config.FilePath)
		json, err := ioutil.ReadFile(fileName)
		if err != nil {
			os.Exit(1)
		}

		config.json = string(json)
	}
}
