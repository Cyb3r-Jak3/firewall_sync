package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
)

//ParseConfig parses the configuration
func ParseConfig(filename string, outputInterface interface{}) (err error) {
	filename = filepath.Clean(filename)
	file, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(file, outputInterface)
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