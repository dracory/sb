package sb

import "strings"

func (b *Builder) quote(s string, quoteType string) string {
	var quoteChar string
	switch b.Dialect {
	case DIALECT_MYSQL:
		quoteChar = "`"
	case DIALECT_POSTGRES:
		quoteChar = `"`
	case DIALECT_SQLITE:
		quoteChar = `"`
	}

	if quoteType == "value" {
		switch b.Dialect {
		case DIALECT_MYSQL:
			return `"` + b.escapeMysql(s) + `"`
		case DIALECT_POSTGRES:
			return `"` + b.escapePostgres(s) + `"`
		case DIALECT_SQLITE:
			return `'` + b.escapeSqlite(s) + `'`
		}
		return s
	}

	if quoteChar != "" {
		return quoteChar + s + quoteChar
	}

	return s
}

func (b *Builder) quoteValue(value string) string {
	if b.Dialect == DIALECT_MYSQL {
		value = `"` + b.escapeMysql(value) + `"`
	}

	if b.Dialect == DIALECT_POSTGRES {
		value = `"` + b.escapePostgres(value) + `"`
	}

	if b.Dialect == DIALECT_SQLITE {
		value = `'` + b.escapeSqlite(value) + `'`
	}

	return value
}

func (b *Builder) escapeMysql(value string) string {
	// escapeRegexp       = regexp.MustCompile(`[\0\t\x1a\n\r\"\'\\]`)
	// characterEscapeMap = map[string]string{
	// 	"\\0":  `\\0`,  //ASCII NULL
	// 	"\b":   `\\b`,  //backspace
	// 	"\t":   `\\t`,  //tab
	// 	"\x1a": `\\Z`,  //ASCII 26 (Control+Z);
	// 	"\n":   `\\n`,  //newline character
	// 	"\r":   `\\r`,  //return character
	// 	"\"":   `\\"`,  //quote (")
	// 	"'":    `\'`,   //quote (')
	// 	"\\":   `\\\\`, //backslash (\)
	// 	// "\\%":  `\\%`,  //% character
	// 	// "\\_":  `\\_`,  //_ character
	// }
	// return escapeRegexp.ReplaceAllStringFunc(val, func(s string) string {

	// 	mVal, ok := characterEscapeMap[s]
	// 	if ok {
	// 		return mVal
	// 	}
	// 	return s
	// })

	escapedStr := strings.ReplaceAll(value, `"`, `""`)
	return escapedStr
}

func (b *Builder) escapePostgres(value string) string {
	escapedStr := strings.ReplaceAll(value, `"`, `""`)
	return escapedStr
}

func (b *Builder) escapeSqlite(value string) string {
	escapedStr := strings.ReplaceAll(value, "'", "''")
	return escapedStr
}
