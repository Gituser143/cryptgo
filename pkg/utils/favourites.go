package utils

import (
	"encoding/json"
	"os"
)

func GetFavourites() map[string]bool {
	favourites := make(map[string]bool)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return favourites
	}

	configPath := homeDir + "/.cryptgo-favourites.json"

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return favourites
	}

	configFile, err := os.Open(configPath)
	if err != nil {
		return favourites
	}

	err = json.NewDecoder(configFile).Decode(&favourites)
	if err != nil {
		return map[string]bool{}
	}

	return favourites
}

func SaveFavourites(favourites map[string]bool) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configPath := homeDir + "/cryptgo-favourites.json"
	hiddenPath := homeDir + "/.cryptgo-favourites.json"

	data, err := json.MarshalIndent(favourites, "", "\t")
	if err != nil {
		return err
	}

	err = os.WriteFile(configPath, data, 0666)
	if err != nil {
		return err
	}

	err = os.Rename(configPath, hiddenPath)
	if err != nil {
		return err
	}

	return nil
}
