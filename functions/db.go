package functions

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

// Connect simplifies the connection functionality
// and returns a connection and an error.
func Connect(connString string) (*pgx.Conn, error) {

	conn, err := pgx.Connect(context.Background(), connString)
	return conn, err
}

func ConvertReturn(ret []interface{}) ([]string, error) {
	var err error
	stringValues := make([]string, 0)
	for _, val := range ret {
		switch v := val.(type) {
		case bool:
			s := strconv.FormatBool(v)
			stringValues = append(stringValues, s)
		case int64:
			s := strconv.FormatInt(v, 10)
			stringValues = append(stringValues, s)
		case int32:
			s := strconv.FormatInt(int64(v), 10)
			stringValues = append(stringValues, s)
		case float64:
			s := strconv.FormatFloat(v, 'f', -1, 32)
			stringValues = append(stringValues, s)
		case time.Time:
			stringValues = append(stringValues, v.String())
		case pgtype.Numeric:
			var fl float64
			v.AssignTo(&fl)
			s := strconv.FormatFloat(fl, 'f', -1, 32)
			stringValues = append(stringValues, s)
		default:
			stringValues = append(stringValues, val.(string))
		}
	}

	return stringValues, err
}

func InterfaceSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("InterfaceSlice() given a non-slice type")
	}

	// Keep the distinction between nil and empty slice input
	if s.IsNil() {
		return nil
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}

// Truncate deletes all the data in the provided table
func Truncate(c *pgx.Conn, table string) (pgconn.CommandTag, error) {
	return c.Exec(context.Background(), fmt.Sprintf("TRUNCATE %s", table))
}

// Importer takes an opened CSV file
// and attempts to COPY it into a table
func Importer(conn *pgx.Conn, f *os.File, table string) (pgconn.CommandTag, error) {
	res, err := conn.PgConn().CopyFrom(context.Background(), f, fmt.Sprintf("COPY %s FROM STDIN DELIMITER ',' CSV HEADER", table))

	return res, err
}
