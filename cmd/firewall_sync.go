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
	Rules   []Rule   `yaml:"rules"`
	ZoneIDs          []string `yaml:"zoneIDs,omitempty"`
}

type Rule struct {
	Name       string   `yaml:"name"`
	Expression string   `yaml:"expression"`
	ZoneIDs    []string `yaml:"zoneIDs,omitempty"`
	ZonesNames []string `yaml:"zonesNames,omitempty"`
	Zones      []string
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
	builtBy   = "unknown" //nolint
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
	err := GenerateConfig(c.String("config"))
	if err != nil {
		log.WithError(err).Error("Error generating config")
		return err
	}
	APIClient, err = cloudflare.NewWithAPIToken(apiToken)
	if err != nil {
		log.WithError(err).Error("Error making API client")
		return err
	}
	log.Debugf("There are %d rules", len(conf.Rules))
	for _, rule := range conf.Rules {
		if len(rule.Zones) == 0 {
			log.Debugf("Rule %s has no zones configured. Applying to all zones", rule.Name)
			RuleProcess(rule, conf.ZoneIDs)
		} else {
			log.Debugf("Rule %s has %d zones configured.", rule.Name, len(rule.ZoneIDs))
			RuleProcess(rule, rule.ZoneIDs)
		}
	}
	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "Firewall Sync"
	app.Usage = "Sync firewall rules between zones"
	app.Version = fmt.Sprintf("%s, Commit: %s, Date: %s, Built By: %s", version, commit, date, builtBy)
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
			Aliases: []string{"V"},
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
