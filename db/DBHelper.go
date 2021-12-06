package db

import "database/sql"

// Update Add, delete and change operations
func Update(connection *sql.DB, sql string, params []interface{}) (sql.Result, error) {
	result, err := connection.Exec(sql, params)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// UpdateByTx Add, delete and change operations
func UpdateByTx(tx *sql.Tx, sql string, params []interface{}) (sql.Result, error) {
	result, err := tx.Exec(sql, params)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// SelectList Query list based on sql and parameters
func SelectList(connection *sql.DB, sql string, params []interface{}) ([]map[string]string, error) {
	rows, err := connection.Query(sql, params)

	if err != nil {
		return nil, err
	}

	return toMap(rows)
}

// SelectListByTx Query list based on sql and parameters
func SelectListByTx(tx *sql.Tx, sql string, params []interface{}) ([]map[string]string, error) {
	rows, err := tx.Query(sql, params)

	if err != nil {
		return nil, err
	}

	return toMap(rows)
}

// toMap Converting query results to map
func toMap(rows *sql.Rows) ([]map[string]string, error) {
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	values := make([][]byte, len(cols))
	scans := make([]interface{}, len(cols))

	for i := range values {
		scans[i] = &values[i]
	}

	var rowLine = make([]map[string]string, 0, 5)
	for rows.Next() {

		err = rows.Scan(scans...)
		if err != nil {
			return nil, err
		}
		row := make(map[string]string, 5)

		for k, v := range values {
			key := cols[k]
			row[key] = string(v)
		}

		rowLine = append(rowLine, row)
	}

	return rowLine, nil
}
