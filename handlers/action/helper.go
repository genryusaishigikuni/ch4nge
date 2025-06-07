package action

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
)

var titleCaser = cases.Title(language.English)

func formatGreenActionTitle(option string) string {
	switch strings.ToLower(option) {
	case "planted_tree":
		return "Planted a tree"
	case "lights_off":
		return "Turned off lights"
	case "used_bike":
		return "Used bicycle"
	case "recycling":
		return "Recycled items"
	case "water_conservation":
		return "Conserved water"
	case "energy_saving":
		return "Saved energy"
	case "waste_reduction":
		return "Reduced waste"
	case "composting":
		return "Composted organic waste"
	case "solar_power":
		return "Used solar power"
	case "public_transport":
		return "Used public transport"
	default:
		return titleCaser.String(strings.ReplaceAll(option, "_", " "))
	}
}

func formatTransportationActionTitle(option, vehicle string) string {
	switch strings.ToLower(option) {
	case "active commute":
		return "Commuted by " + vehicle
	case "private vehicle":
		return "Used " + vehicle
	case "public transport":
		return "Took " + vehicle
	default:
		return option + " - " + vehicle
	}
}
