package utils

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
)

func StringPtr(s string) *string {
	return &s
}

func BoolPtr(b bool) *bool {
	return &b
}
func PrepareRequestBody(v any) (
	*bytes.Buffer,
	error,
) {
	jsonBuffer, err := json.Marshal(v)

	if err != nil {
		return nil, err
	}

	requestBody := bytes.NewBuffer(jsonBuffer)

	return requestBody, nil
}

func ExecuteSQLFile(
	ctx context.Context,
	db *sql.DB,
	filename string,
) error {
	// Read the embedded SQL file
	sqlBytes, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf(
			"failed to read embedded SQL file %s: %w",
			filename,
			err,
		)
	}

	// Execute the SQL
	_, err = db.ExecContext(
		ctx,
		string(sqlBytes),
	)
	if err != nil {
		return fmt.Errorf(
			"failed to execute SQL from file %s: %w",
			filename,
			err,
		)
	}

	return nil
}

func IntPtr(i int) *int {
	return &i
}
