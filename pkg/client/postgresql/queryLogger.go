package postgresql

import "github.com/VrMolodyakov/vote-service/pkg/logging"

func QueryLogger(sql, table string, logger logging.Logger, args []interface{}) *logging.Logger {
	return logger.ExtraFields(map[string]interface{}{
		"sql":   sql,
		"table": table,
		"args":  args,
	})
}
