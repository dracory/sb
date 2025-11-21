package sb

import "strings"

// rawColumnProcess parses a SQL column type string to extract the base type, length, and decimal places.
// This function handles column type definitions that may include parentheses for size specifications,
// such as VARCHAR(100), DECIMAL(10,2), or simpler types like INT.
//
// Parameters:
//   - columnType: A string representing the SQL column type, e.g., "VARCHAR(100)", "DECIMAL(10,2)", "INT".
//
// Returns:
//   - scolumnType: The base column type without size specifications, e.g., "VARCHAR", "DECIMAL", "INT".
//   - length: The length specification as a string, e.g., "100" for VARCHAR(100), or empty string if not specified.
//   - decimals: The decimal places specification as a string, e.g., "2" for DECIMAL(10,2), or empty string if not specified.
//
// Behavior:
//   - If the columnType does not contain "(", it returns the columnType as scolumnType with empty length and decimals.
//   - For types like "VARCHAR(100)", it returns "VARCHAR", "100", "".
//   - For types like "DECIMAL(10,2)", it returns "DECIMAL", "10", "2".
//   - If parsing fails (e.g., malformed input), it falls back to returning the original columnType with empty length and decimals.
//
// Examples:
//   rawColumnProcess("VARCHAR(100)") → ("VARCHAR", "100", "")
//   rawColumnProcess("DECIMAL(10,2)") → ("DECIMAL", "10", "2")
//   rawColumnProcess("INT") → ("INT", "", "")
//   rawColumnProcess("VARCHAR(100, extra)") → ("VARCHAR", "100", " extra")  // Note: handles unexpected formats gracefully
//
// Note:
//   - The function is designed to be robust against malformed input, returning the original columnType with empty length and decimals if parsing fails.
//   - It gracefully handles unexpected formats by returning the original columnType with empty length and decimals.
func rawColumnProcess(columnType string) (scolumnType, length, decimals string) {
	// If the columnType does not contain "(", it returns the columnType as scolumnType with empty length and decimals.
	if !strings.Contains(columnType, "(") {
		return columnType, "", ""
	}

	// Split the columnType by "(" and return the first part as scolumnType.
	splitByParen := strings.Split(columnType, "(")
	if len(splitByParen) < 2 {
		return columnType, "", ""
	}

	// Remove the first part of the split result and trim it.
	columnType = strings.TrimSpace(splitByParen[0])

	// Remove the last part of the split result and trim it.
	properties := strings.TrimRight(splitByParen[1], ")")

	// If the properties does not contain "," then return the properties as length.
	if !strings.Contains(properties, ",") {
		length = strings.TrimSpace(properties)
		return columnType, length, decimals
	}

	// Split the properties by "," and return the first part as length and the second part as decimals.
	splitByComma := strings.Split(properties, ",")
	if len(splitByComma) < 2 {
		return columnType, strings.TrimSpace(properties), ""
	}
	length = strings.TrimSpace(splitByComma[0])
	decimals = strings.TrimSpace(splitByComma[1])

	return columnType, length, decimals
}
