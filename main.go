package main

import (
	"context"
	"os"
	"strings"

	"github.com/cloudflare/cloudflare-go"
	"github.com/sirupsen/logrus"
)

type Config struct {
	FilterExpression string `yaml:"expression,omitempty"`
	RuleName         string `yaml:"rule,omitempty"`
	FilterID         string
	RuleId           string
	ZoneIDs          []string `yaml:"zoneIDs,omitempty"`
}

var (
	log       = logrus.New()
	ctx       = context.Background()
	Conf      *Config
	ApiClient *cloudflare.API
)

func setLogLevel()  {
	switch strings.ToLower(os.Getenv("LOG_LEVEL")) {
	case "trace": log.SetLevel(logrus.TraceLevel)
	case "debug": log.SetLevel(logrus.DebugLevel)
	default:
		log.SetLevel(logrus.InfoLevel)
	}
}


func main() {
	setLogLevel()
	log.Info("Starting")
	apiToken := os.Getenv("CF_TOKEN")
	if apiToken == "" {
		log.Fatal("No API Token set")
	}
	err := ParseConfig("config.yml", &Conf)
	if err != nil {
		log.WithError(err).Fatal("Error reading config")
	}
	ApiClient, err = cloudflare.NewWithAPIToken(apiToken)
	if err != nil {
		log.WithError(err).Fatal("Error making API client")
	}
	if len(Conf.ZoneIDs) == 0 {
		GetZones()
	}
	if len(Conf.ZoneIDs) == 0 {
		log.Fatal("There were no zone IDs found")
	}
	for _, zone := range Conf.ZoneIDs {
		log.Debugf("Zone: %s", zone)
		rules, err := ApiClient.FirewallRules(ctx, zone, cloudflare.PaginationOptions{})
		if err != nil {
			log.WithError(err).Errorf("Error getting rules for zone: %s", zone)
			continue
		}
		var filter, ruleid string
		for _, rule := range rules {
			if rule.Description == Conf.RuleName {
				log.Debugf("Found rule by name for zone: %s", zone)
				filter = rule.Filter.ID
				ruleid = rule.ID
				break
			}
			if rule.Filter.Expression == Conf.FilterExpression {
				log.Debugf("Found rule by expression for zone: %s", zone)
				filter = rule.Filter.ID
				ruleid = rule.ID
				break
			}

		}
		log.Debugf("Filter %s Rule ID %s", filter, ruleid)
		if filter != "" {
			log.Debugf("Updating existing rule for zone %s", zone)
			UpdateRule(zone, filter, ruleid)
		} else {
			log.Debugf("Creating rule for zone %s", zone)
			CreateRule(zone)
		}
	}
}
