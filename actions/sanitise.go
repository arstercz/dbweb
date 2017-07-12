package actions

// This function sanitises SQL for use by this program. The reason for this is
// to prevent SQL being executed that would have un-desirable side effects or
// cause errors. For example, when dealing with SQL updates many people craft
// SQL for use on the command line. Such command line clients support SQL
// features that this program does not. These features are removed from any SQL
// executed by this program.
//
// All carriage returns are removed from the passed SQL string to help parsing.
// All SQL lines containing the following statements are removed:
//
// 1. merge the space char
// 2. CREATE DATABASE/SCHEMA, DROP DATABASE/SCHEMA
// 3. TRUNCATE/DROP TABLE/GRANT../REVOKE../SET..
// 4. DELETE or UPDATE with no WHERE
// 5. USE <..>
import (
	"regexp"
	"strings"
)

func SanitiseSQL(sql string) string {
	sql = MergeSpace(sql)
	sql = ConvertToUnixLineEndings(sql)
	sql = removeAlterDrop(sql)
	sql = removeCreateDatabaseStatements(sql)
	sql = removeNotPermitStatements(sql)
	sql = removeChangeNoWhere(sql)
	sql = removeSelectNoWhere(sql)
	sql = removeUseStatements(sql)
	return sql
}

func MergeSpace(sql string) string {
	pattern := regexp.MustCompile("(\\s+)")
	if pattern.MatchString(sql) {
		sql = pattern.ReplaceAllString(sql, " ")
	}
	return sql
}

// Convert the SQL string line endings to unix format.
func ConvertToUnixLineEndings(sql string) string {
	sql = strings.TrimSpace(sql)
	sql = strings.Replace(sql, "\r\n", "\n", -1)
	sql = strings.Replace(sql, "\r", "\n", -1)
	return sql
}

// remove any CREATE DATABASE or CREATE SCHEMA statements in the passed SQL.
func removeCreateDatabaseStatements(sql string) string {
	pattern := regexp.MustCompile("(?i)(^CREATE\\s+DATABASE\\s+|^CREATE SCHEMA\\s+|^DROP\\s+DATABASE\\s|^DROP\\s+SCHEMA\\s+)")
	lines := strings.Split(sql, "\n")
	output := make([]string, 0)
	for _, line := range lines {
		if !pattern.MatchString(line) {
			output = append(output, line)
		}
	}
	return strings.Join(output, "\n")
}

func removeNotPermitStatements(sql string) string {
	pattern := regexp.MustCompile("(?i)(^TRUNCATE\\s+|^DROP\\s+TABLE\\s+|^GRANT\\s+|^REVOKE\\s+|^SET\\s+)")
	lines := strings.Split(sql, "\n")
	output := make([]string, 0)
	for _, line := range lines {
		if !pattern.MatchString(line) {
			output = append(output, line)
		}
	}
	return strings.Join(output, "\n")
}

func removeChangeNoWhere(sql string) string {
	pattern1 := regexp.MustCompile("(?i)(^DELETE\\s+|^UPDATE\\s+)")
	pattern2 := regexp.MustCompile("(?i)(\\s+WHERE\\s+)")
	lines := strings.Split(sql, "\n")
	output := make([]string, 0)
	for _, line := range lines {
		// remove statement which change file with no where condition.
		if !(pattern1.MatchString(line) && !pattern2.MatchString(line)) {
			output = append(output, line)
		}
	}
	return strings.Join(output, "\n")
}

func removeAlterDrop(sql string) string {
	pattern1 := regexp.MustCompile("(?i)(^ALTER\\s+)")
	pattern2 := regexp.MustCompile("(?i)(\\s+DROP\\s+)")
	lines := strings.Split(sql, "\n")
	output := make([]string, 0)
	for _, line := range lines {
		// remove statement which change file with no where condition.
		if !(pattern1.MatchString(line) && !pattern2.MatchString(line)) {
			output = append(output, line)
		}
	}
	return strings.Join(output, "\n")
}

func removeSelectNoWhere(sql string) string {
	pattern1 := regexp.MustCompile("(?i)(^SELECT\\s+)")
	pattern2 := regexp.MustCompile("(?i)(\\s+WHERE\\s+|\\s+LIMIT\\s+)")
	lines := strings.Split(sql, "\n")
	output := make([]string, 0)
	for _, line := range lines {
		// remove statement which change file with no where condition.
		if !(pattern1.MatchString(line) && !pattern2.MatchString(line)) {
			output = append(output, line)
		}
	}
	return strings.Join(output, "\n")
}

func removeUseStatements(sql string) string {
	pattern := regexp.MustCompile("(?i)(^USE\\s+)")
	lines := strings.Split(sql, "\n")
	output := make([]string, 0)
	for _, line := range lines {
		if !pattern.MatchString(line) {
			output = append(output, line)
		}
	}
	return strings.Join(output, "\n")
}
