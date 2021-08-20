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

// Metadata holds persistant information to be stored to disk
type Metadata struct {
	Favourites map[string]bool    `json:"favourites"`
	Currency   string             `json:"currency"`
	Portfolio  map[string]float64 `json:"portfolio"`
}

// Currency holds currency data when fetched from CoinCap
type Currency struct {
	ID             string `json:"id"`
	Symbol         string `json:"symbol"`
	CurrencySymbol string `json:"currencySymbol"`
	Type           string `json:"type"`
	RateUSD        string `json:"rateUsd"`
}

// AllCurrencyData holds data of a group of currencies fetched from CoinCap
type AllCurrencyData struct {
	Data      []Currency `json:"data"`
	Timestamp uint       `json:"timestamp"`
}

// GetFavourites reads stored favourite coin details from
// ~/.cryptgo-data.json and returns a map.
func GetFavourites() map[string]bool {
	metadata := Metadata{}

	// Get home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return map[string]bool{}
	}

	// Check if metadata file exists
	configPath := homeDir + "/.cryptgo-data.json"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return map[string]bool{}
	}

	// Open file
	configFile, err := os.Open(configPath)
	if err != nil {
		return map[string]bool{}
	}

	// Read content
	err = json.NewDecoder(configFile).Decode(&metadata)
	if err != nil {
		return map[string]bool{}
	}

	if len(metadata.Favourites) > 0 {
		return metadata.Favourites
	}

	return map[string]bool{}
}

// GetPortfolio reads stored portfolio details from
// ~/.cryptgo-data.json and returns a map.
func GetPortfolio() map[string]float64 {
	metadata := Metadata{}

	// Get home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return map[string]float64{}
	}

	// Check if metadta file exists
	configPath := homeDir + "/.cryptgo-data.json"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return map[string]float64{}
	}

	// Open file
	configFile, err := os.Open(configPath)
	if err != nil {
		return map[string]float64{}
	}

	// Read content
	err = json.NewDecoder(configFile).Decode(&metadata)
	if err != nil {
		return map[string]float64{}
	}

	if len(metadata.Portfolio) > 0 {
		return metadata.Portfolio
	}

	return map[string]float64{}
}

// GetCurrencyID returns the currencyID stored from metadata
func GetCurrencyID() string {
	metadata := Metadata{}

	// Get home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "united-states-dollar"
	}

	// Check if metadta file exists
	configPath := homeDir + "/.cryptgo-data.json"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return "united-states-dollar"
	}

	// Open file
	configFile, err := os.Open(configPath)
	if err != nil {
		return "united-states-dollar"
	}

	// Read content
	err = json.NewDecoder(configFile).Decode(&metadata)
	if err != nil {
		return "united-states-dollar"
	}

	return metadata.Currency
}

// SaveMetadata exports favourites, currency and portfolio to disk.
// Data is saved on ~/.cryptgo-data.json
func SaveMetadata(favourites map[string]bool, currency string, portfolio map[string]float64) error {
	// Get Home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// configPath and hidden path are used explicitly because we
	// get a permission denied error on trying to write/create
	// to a hidden file
	configPath := homeDir + "/cryptgo-data.json"
	hiddenPath := homeDir + "/.cryptgo-data.json"

	// Create data
	metadata := Metadata{
		Favourites: favourites,
		Currency:   currency,
		Portfolio:  portfolio,
	}

	data, err := json.MarshalIndent(metadata, "", "\t")
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
