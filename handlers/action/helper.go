package action

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"log"
	"strings"
)

var titleCaser = cases.Title(language.English)

func formatGreenActionTitle(option string) string {
	// Log the option that is being processed
	log.Printf("Processing green action option: %s", option)

	// Determine the formatted action title based on the option
	var actionTitle string
	switch strings.ToLower(option) {
	case "planted_tree":
		actionTitle = "Planted a tree"
	case "lights_off":
		actionTitle = "Turned off lights"
	case "used_bike":
		actionTitle = "Used bicycle"
	case "recycling":
		actionTitle = "Recycled items"
	case "water_conservation":
		actionTitle = "Conserved water"
	case "energy_saving":
		actionTitle = "Saved energy"
	case "waste_reduction":
		actionTitle = "Reduced waste"
	case "composting":
		actionTitle = "Composted organic waste"
	case "solar_power":
		actionTitle = "Used solar power"
	case "public_transport":
		actionTitle = "Used public transport"
	default:
		actionTitle = titleCaser.String(strings.ReplaceAll(option, "_", " "))
	}

	// Log the formatted action title
	log.Printf("Formatted green action title: %s", actionTitle)

	// Return the formatted title
	return actionTitle
}
