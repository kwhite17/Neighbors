package retrievers

import "html/template"
import "github.com/kwhite17/Neighbors/pkg/managers"

func StatusAsString(status managers.ItemStatus) string {
	switch status {
	case managers.CREATED:
		return "CREATED"
	case managers.CLAIMED:
		return "CLAIMED"
	case managers.DELIVERED:
		return "DELIVERED"
	case managers.RECEIVED:
		return "RECEIVED"
	default:
		return "UNKNOWN"
	}
}

func buildFuncMap() template.FuncMap {
	return template.FuncMap{
		"statusAsString": StatusAsString,
	}
}
