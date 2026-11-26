package service

import "github.com/jackc/pgx/v5/pgtype"

func numericToString(n pgtype.Numeric) string {
	if !n.Valid {
			return ""
	}
	val, err := n.Value()
	if err != nil || val == nil {
			return ""
	}
	return val.(string)
}