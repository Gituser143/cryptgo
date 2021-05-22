/*
Copyright © 2021 Bhargav SNV bhargavsnv100@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package utils

import (
	"encoding/json"
	"os"
)

// GetFavourites reads stored favourite coin details from ~/.cryptgo-favourites.json and returns a map.
func GetFavourites() map[string]bool {
	favourites := make(map[string]bool)

	// Get home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return favourites
	}

	// Check if favourites file exists
	configPath := homeDir + "/.cryptgo-favourites.json"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return favourites
	}

	// Open file
	configFile, err := os.Open(configPath)
	if err != nil {
		return favourites
	}

	// Read content
	err = json.NewDecoder(configFile).Decode(&favourites)
	if err != nil {
		return map[string]bool{}
	}

	return favourites
}

// SaveFavourites exports favourites to disk. Data is saved on ~/.cryptgo-favourites.json
func SaveFavourites(favourites map[string]bool) error {
	// Get Home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// configPath and hidden path are used explicitly because we
	// get a permission denied error on trying to write/create
	// to a hidden file
	configPath := homeDir + "/cryptgo-favourites.json"
	hiddenPath := homeDir + "/.cryptgo-favourites.json"

	// Create data
	data, err := json.MarshalIndent(favourites, "", "\t")
	if err != nil {
		return err
	}

	// Write to file
	err = os.WriteFile(configPath, data, 0666)
	if err != nil {
		return err
	}

	// Hide file
	err = os.Rename(configPath, hiddenPath)
	if err != nil {
		return err
	}

	return nil
}
