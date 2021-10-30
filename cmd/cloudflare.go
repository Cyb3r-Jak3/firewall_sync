package main

import (
	"github.com/cloudflare/cloudflare-go"
)

//UpdateRule updates and existing firewall rule with a new filter
func UpdateRule(ZoneID, Filter, Rule string) {
	_, err := APIClient.Filter(ctx, ZoneID, Filter)
	if err != nil {
		log.WithError(err).Error("Error checking for existing filter")
		return
	}
	filter := cloudflare.Filter{
		ID:          Filter,
		Expression:  conf.FilterExpression,
		Description: conf.RuleName,
	}
	newFilter, err := APIClient.UpdateFilter(ctx, ZoneID, filter)
	if err != nil {
		log.WithError(err).Error("Error updating filter")
	}
	updatedRule := cloudflare.FirewallRule{
		ID: Rule,
		Filter: cloudflare.Filter{
			ID:          newFilter.ID,
			Expression:  conf.FilterExpression,
			Description: conf.RuleName,
		},
		Action:      "block",
		Description: conf.RuleName,
		Priority:    900,
	}
	_, err = APIClient.UpdateFirewallRule(ctx, ZoneID, updatedRule)
	if err != nil {
		log.WithError(err).Errorf("Error updating rule: %s in zone %s", Rule, ZoneID)
		return
	}
	log.Infof("Rule: %s in Zone %s updated", Rule, ZoneID)
}

//CreateRule creates a new rule if one with the same is not detected
func CreateRule(ZoneID string) {
	newRule := []cloudflare.FirewallRule{{
		Filter: cloudflare.Filter{
			Expression:  conf.FilterExpression,
			Description: conf.RuleName,
		},
		Action:      "block",
		Description: conf.RuleName,
		Priority:    900,
	}}
	_, err := APIClient.CreateFirewallRules(ctx, ZoneID, newRule)
	if err != nil {
		log.WithError(err).Error("Error creating new rule")
		return
	}
	log.Infof("Created new firewall rule in zone %s", ZoneID)
}
