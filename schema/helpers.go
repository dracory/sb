package schema

import "strings"

// rawColumnProcess parses a SQL column type string to extract the base type, length, and decimals.
func rawColumnProcess(columnType string) (scolumnType, length, decimals string) {
	if !strings.Contains(columnType, "(") {
		return columnType, "", ""
	}

	splitByParen := strings.Split(columnType, "(")
	if len(splitByParen) < 2 {
		return columnType, "", ""
	}

	columnType = strings.TrimSpace(splitByParen[0])
	properties := strings.TrimRight(splitByParen[1], ")")

	if !strings.Contains(properties, ",") {
		length = strings.TrimSpace(properties)
		return columnType, length, decimals
	}

	splitByComma := strings.Split(properties, ",")
	if len(splitByComma) < 2 {
		return columnType, strings.TrimSpace(properties), ""
	}
	length = strings.TrimSpace(splitByComma[0])
	decimals = strings.TrimSpace(splitByComma[1])

	return columnType, length, decimals
}
