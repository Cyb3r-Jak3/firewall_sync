package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/cloudflare/cloudflare-go"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

//Config is the configuration containing the zones, filter expression and rule name
type Config struct {
	FilterExpression string   `yaml:"expression,omitempty"`
	RuleName         string   `yaml:"rule,omitempty"`
	ZoneIDs          []string `yaml:"zoneIDs,omitempty"`
}

var (
	log  = logrus.New()
	ctx  = context.Background()
	conf *Config
	//APIClient is the Cloudflare Client that is used
	APIClient *cloudflare.API
	version   = "dev"     //nolint
	commit    = "none"    //nolint
	date      = "unknown" //nolint
)

func setLogLevel(c *cli.Context) {
	if c.Bool("debug") {
		log.SetLevel(logrus.DebugLevel)
	} else if c.Bool("verbose") {
		log.SetLevel(logrus.InfoLevel)
	} else {
		switch strings.ToLower(os.Getenv("LOG_LEVEL")) {
		case "trace":
			log.SetLevel(logrus.TraceLevel)
		case "debug":
			log.SetLevel(logrus.DebugLevel)
		default:
			log.SetLevel(logrus.InfoLevel)
		}
	}
	log.Debugf("Log Level set to %v", log.Level)
}

func run(c *cli.Context) error {
	setLogLevel(c)
	log.Info("Starting")
	apiToken := os.Getenv("CF_TOKEN")
	if apiToken == "" {
		return fmt.Errorf("no API token set")
	}
	err := ParseConfig(c.String("config"), &conf)
	if err != nil {
		log.WithError(err).Error("Error reading config")
		return err
	}
	APIClient, err = cloudflare.NewWithAPIToken(apiToken)
	if err != nil {
		log.WithError(err).Error("Error making API client")
		return err
	}
	if len(conf.ZoneIDs) == 0 {
		GetZones()
	}
	if len(conf.ZoneIDs) == 0 {
		return fmt.Errorf("no zone IDs found")
	}
	for _, zone := range conf.ZoneIDs {
		log.Debugf("Zone: %s", zone)
		rules, err := APIClient.FirewallRules(ctx, zone, cloudflare.PaginationOptions{})
		if err != nil {
			log.WithError(err).Errorf("Error getting rules for zone: %s", zone)
			continue
		}
		var filter, ruleId string
		for _, rule := range rules {
			if rule.Description == conf.RuleName {
				log.Debugf("Found rule by name for zone: %s", zone)
				filter = rule.Filter.ID
				ruleId = rule.ID
				break
			}
			if rule.Filter.Expression == conf.FilterExpression {
				log.Debugf("Found rule by expression for zone: %s", zone)
				filter = rule.Filter.ID
				ruleId = rule.ID
				break
			}

		}
		log.Debugf("Filter %s Rule ID %s", filter, ruleId)
		if filter != "" {
			log.Debugf("Updating existing rule for zone %s", zone)
			UpdateRule(zone, filter, ruleId)
		} else {
			log.Debugf("Creating rule for zone %s", zone)
			CreateRule(zone)
		}
	}
	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "Firewall Sync"
	app.Usage = "Sync firewall rules between zones"
	app.Version = fmt.Sprintf("%s, Commit: %s, Date: %s", version, commit, date)
	app.Authors = []*cli.Author{
		{
			Name:  "Cyb3r-Jak3",
			Email: "cyb3rjak3@pm.me",
		},
	}
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Value:   "./config.yml",
			Usage:   "Path to the configuration file",
		},
		&cli.BoolFlag{
			Name:    "verbose",
			Aliases: []string{"v"},
			EnvVars: []string{"LOG_LEVEL_VERBOSE"},
		},
		&cli.BoolFlag{
			Name:    "debug",
			Aliases: []string{"d"},
			EnvVars: []string{"LOG_LEVEL_DEBUG"},
		},
	}
	app.Action = run
	err := app.Run(os.Args)
	if err != nil {
		log.WithError(err).Fatal("Error running app")
	}
}
