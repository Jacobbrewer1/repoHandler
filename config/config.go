package config

import (
	"encoding/json"
	"errors"
	"github.com/Jacobbrewer1/repoHandler/api"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var (
	override *overrideStruct
)

type overrideStruct struct {
	GithubApiToken *string `json:"GithubApiToken"`
}

func ReadConfig() error {
	if abs, exists := findFile("./config/override.json"); exists {
		log.Println("Override detected - Reading file")

		file, err := ioutil.ReadFile(abs)
		if err != nil {
			return err
		}

		log.Println(string(file))

		err = json.Unmarshal(file, &override)
		if err != nil {
			return err
		}

		api.GithubApiToken = *override.GithubApiToken
	} else {
		log.Println("No override detected. Using production config")
		api.GithubApiToken, exists = os.LookupEnv("GITHUBAPITOKEN")
		if !exists {
			return errors.New("github api token is nil, make sure this is set")
		}
	}

	if abs, exists := findFile("./config/labelConfig.json"); exists {
		log.Println("Config detected - Reading file")

		c, err := ioutil.ReadFile(abs)
		if err != nil {
			return err
		}

		log.Println(string(c))
		err = json.Unmarshal(c, &api.ConfiguredLabels)
		if err != nil {
			log.Println(err)
			return err
		}

	}
	return nil
}

func findFile(path string) (string, bool) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", false
	}
	log.Println(abs)

	file, err := os.Open(abs)
	if err != nil {
		return "", false
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	return abs, true
}
