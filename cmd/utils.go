package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
)

//GenerateConfig parses the configuration
func GenerateConfig(filename string) (err error) {
	filename = filepath.Clean(filename)
	file, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	if err = yaml.Unmarshal(file, &conf); err != nil {
		return err
	}
	if len(conf.Rules) == 0 {
		return fmt.Errorf("no rules found")
	}
	GetZones()
	return nil
}

//GetZones gets ZoneIDs from environment variables
func GetZones() {
	for i := 1; i <= 10; i++ {
		zoneID := os.Getenv(fmt.Sprintf("CF_ZONE_%d", i))
		if zoneID == "" {
			return
		}
		conf.ZoneIDs = append(conf.ZoneIDs, zoneID)
	}

}
