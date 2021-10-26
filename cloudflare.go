package main

import (
	"github.com/cloudflare/cloudflare-go"
)

func UpdateRule(ZoneId, Filter, Rule string) {
	_, err := ApiClient.Filter(ctx, ZoneId, Filter)
	if err != nil {
		log.WithError(err).Error("Error checking for existing filter")
		return
	}
	filter := cloudflare.Filter{
		ID: Filter,
		Expression: Conf.FilterExpression,
		Description: Conf.RuleName,
	}
	newFilter, err := ApiClient.UpdateFilter(ctx, ZoneId, filter)
	if err != nil {
		log.WithError(err).Error("Error updating filter")
	}
	updatedRule := cloudflare.FirewallRule{
		ID: Rule,
		Filter: cloudflare.Filter{
			ID:          newFilter.ID,
			Expression:  Conf.FilterExpression,
			Description: Conf.RuleName,
		},
		Action:      "block",
		Description: Conf.RuleName,
		Priority:    900,
	}
	_, err = ApiClient.UpdateFirewallRule(ctx, ZoneId, updatedRule)
	if err != nil {
		log.WithError(err).Errorf("Error updating rule: %s in zone %s", Rule, ZoneId)
		return
	}
	log.Infof("Rule: %s in Zone %s updated", Rule, ZoneId)
}

func CreateRule(ZoneId string) {

	newRule := []cloudflare.FirewallRule{{
		Filter: cloudflare.Filter{
			Expression:  Conf.FilterExpression,
			Description: Conf.RuleName,
		},
		Action:      "block",
		Description: Conf.RuleName,
		Priority:    900,
	}}
	_, err := ApiClient.CreateFirewallRules(ctx, ZoneId, newRule)
	if err != nil {
		log.WithError(err).Error("Error creating new rule")
		return
	}
	log.Infof("Created new firewall rule in zone %s", ZoneId)
}