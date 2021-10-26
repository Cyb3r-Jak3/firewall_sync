package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
)

func ParseConfig(filename string, outputInterface interface{}) (err error) {
	filename = filepath.Clean(filename)
	file, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(file, outputInterface)
}
func GetZones() {
	for i := 1; i <= 10; i++ {
		zoneId := os.Getenv(fmt.Sprintf("CF_ZONE_%d", i))
		if zoneId == "" {
			return
		}
		Conf.ZoneIDs = append(Conf.ZoneIDs, zoneId)
	}

}