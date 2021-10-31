package main

import (
	"github.com/cloudflare/cloudflare-go"
)

func RuleProcess(r Rule, zones []string) {
	log.Debugf("Zones: %v", zones)
	for _, zone := range zones {
		log.Debugf("Zone: %s", zone)
		rules, err := APIClient.FirewallRules(ctx, zone, cloudflare.PaginationOptions{})
		if err != nil {
			log.WithError(err).Errorf("Error getting rules for zone: %s", zone)
			continue
		}
		var filter, ruleID string
		for _, rule := range rules {
			if rule.Description == r.Name {
				log.Debugf("Found rule by name for zone: %s", zone)
				filter = rule.Filter.ID
				ruleID = rule.ID
				break
			}
			// This should only be found if the name changes or on first run
			if rule.Filter.Expression == r.Expression {
				log.Debugf("Found rule by expression for zone: %s", zone)
				filter = rule.Filter.ID
				ruleID = rule.ID
				break
			}

		}
		if filter != "" {
			log.Debugf("Updating existing rule for zone %s", zone)
			UpdateRule(r, zone, filter, ruleID)
		} else {
			log.Debugf("Creating rule for zone %s", zone)
			CreateRule(r, zone)
		}
	}
}


//UpdateRule updates and existing firewall rule with a new filter
func UpdateRule(r Rule, ZoneID, Filter, RuleID string) {

	filter := cloudflare.Filter{
		ID:          Filter,
		Expression:  r.Expression,
		Description: r.Name,
	}
	newFilter, err := APIClient.UpdateFilter(ctx, ZoneID, filter)
	if err != nil {
		log.WithError(err).Error("Error updating filter")
	}
	updatedRule := cloudflare.FirewallRule{
		ID: RuleID,
		Filter: cloudflare.Filter{
			ID:          newFilter.ID,
			Expression:  r.Expression,
			Description: r.Name,
		},
		Action:      "block",
		Description: r.Name,
		Priority:    900,
	}
	_, err = APIClient.UpdateFirewallRule(ctx, ZoneID, updatedRule)
	if err != nil {
		log.WithError(err).Errorf("Error updating rule: %s in zone %s", RuleID, ZoneID)
		return
	}
	log.Infof("RuleID: %s in Zone %s updated", RuleID, ZoneID)
}

//CreateRule creates a new rule if one with the same is not detected
func CreateRule(rule Rule, ZoneID string) {
	newRule := []cloudflare.FirewallRule{{
		Filter: cloudflare.Filter{
			Expression:  rule.Expression,
			Description: rule.Name,
		},
		Action:      "block",
		Description: rule.Name,
		Priority:    900,
	}}
	_, err := APIClient.CreateFirewallRules(ctx, ZoneID, newRule)
	if err != nil {
		log.WithError(err).Error("Error creating new rule")
		return
	}
	log.Infof("Created new firewall rule in zone %s", ZoneID)
}
