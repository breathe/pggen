// Code generated by pggen. DO NOT EDIT.

package custom_types

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/jschaf/pggen/example/custom_types/mytype"
)

// Querier is a typesafe Go interface backed by SQL queries.
//
// Methods ending with Batch enqueue a query to run later in a pgx.Batch. After
// calling SendBatch on pgx.Conn, pgxpool.Pool, or pgx.Tx, use the Scan methods
// to parse the results.
type Querier interface {
	CustomTypes(ctx context.Context) (CustomTypesRow, error)
	// CustomTypesBatch enqueues a CustomTypes query into batch to be executed
	// later by the batch.
	CustomTypesBatch(batch genericBatch)
	// CustomTypesScan scans the result of an executed CustomTypesBatch query.
	CustomTypesScan(results pgx.BatchResults) (CustomTypesRow, error)

	CustomMyInt(ctx context.Context) (int, error)
	// CustomMyIntBatch enqueues a CustomMyInt query into batch to be executed
	// later by the batch.
	CustomMyIntBatch(batch genericBatch)
	// CustomMyIntScan scans the result of an executed CustomMyIntBatch query.
	CustomMyIntScan(results pgx.BatchResults) (int, error)

	IntArray(ctx context.Context) ([][]int32, error)
	// IntArrayBatch enqueues a IntArray query into batch to be executed
	// later by the batch.
	IntArrayBatch(batch genericBatch)
	// IntArrayScan scans the result of an executed IntArrayBatch query.
	IntArrayScan(results pgx.BatchResults) ([][]int32, error)
}

type DBQuerier struct {
	conn  genericConn   // underlying Postgres transport to use
	types *typeResolver // resolve types by name
}

var _ Querier = &DBQuerier{}

// genericConn is a connection to a Postgres database. This is usually backed by
// *pgx.Conn, pgx.Tx, or *pgxpool.Pool.
type genericConn interface {
	// Query executes sql with args. If there is an error the returned Rows will
	// be returned in an error state. So it is allowed to ignore the error
	// returned from Query and handle it in Rows.
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)

	// QueryRow is a convenience wrapper over Query. Any error that occurs while
	// querying is deferred until calling Scan on the returned Row. That Row will
	// error with pgx.ErrNoRows if no rows are returned.
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row

	// Exec executes sql. sql can be either a prepared statement name or an SQL
	// string. arguments should be referenced positionally from the sql string
	// as $1, $2, etc.
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
}

// genericBatch batches queries to send in a single network request to a
// Postgres server. This is usually backed by *pgx.Batch.
type genericBatch interface {
	// Queue queues a query to batch b. query can be an SQL query or the name of a
	// prepared statement. See Queue on *pgx.Batch.
	Queue(query string, arguments ...interface{})
}

// NewQuerier creates a DBQuerier that implements Querier. conn is typically
// *pgx.Conn, pgx.Tx, or *pgxpool.Pool.
func NewQuerier(conn genericConn) *DBQuerier {
	return NewQuerierConfig(conn, QuerierConfig{})
}

type QuerierConfig struct {
	// DataTypes contains pgtype.Value to use for encoding and decoding instead
	// of pggen-generated pgtype.ValueTranscoder.
	//
	// If OIDs are available for an input parameter type and all of its
	// transitive dependencies, pggen will use the binary encoding format for
	// the input parameter.
	DataTypes []pgtype.DataType
}

// NewQuerierConfig creates a DBQuerier that implements Querier with the given
// config. conn is typically *pgx.Conn, pgx.Tx, or *pgxpool.Pool.
func NewQuerierConfig(conn genericConn, cfg QuerierConfig) *DBQuerier {
	return &DBQuerier{conn: conn, types: newTypeResolver(cfg.DataTypes)}
}

// WithTx creates a new DBQuerier that uses the transaction to run all queries.
func (q *DBQuerier) WithTx(tx pgx.Tx) (*DBQuerier, error) {
	return &DBQuerier{conn: tx}, nil
}

// preparer is any Postgres connection transport that provides a way to prepare
// a statement, most commonly *pgx.Conn.
type preparer interface {
	Prepare(ctx context.Context, name, sql string) (sd *pgconn.StatementDescription, err error)
}

// PrepareAllQueries executes a PREPARE statement for all pggen generated SQL
// queries in querier files. Typical usage is as the AfterConnect callback
// for pgxpool.Config
//
// pgx will use the prepared statement if available. Calling PrepareAllQueries
// is an optional optimization to avoid a network round-trip the first time pgx
// runs a query if pgx statement caching is enabled.
func PrepareAllQueries(ctx context.Context, p preparer) error {
	if _, err := p.Prepare(ctx, customTypesSQL, customTypesSQL); err != nil {
		return fmt.Errorf("prepare query 'CustomTypes': %w", err)
	}
	if _, err := p.Prepare(ctx, customMyIntSQL, customMyIntSQL); err != nil {
		return fmt.Errorf("prepare query 'CustomMyInt': %w", err)
	}
	if _, err := p.Prepare(ctx, intArraySQL, intArraySQL); err != nil {
		return fmt.Errorf("prepare query 'IntArray': %w", err)
	}
	return nil
}

// typeResolver looks up the pgtype.ValueTranscoder by Postgres type name.
type typeResolver struct {
	connInfo *pgtype.ConnInfo // types by Postgres type name
}

func newTypeResolver(types []pgtype.DataType) *typeResolver {
	ci := pgtype.NewConnInfo()
	for _, typ := range types {
		if txt, ok := typ.Value.(textPreferrer); ok && typ.OID != unknownOID {
			typ.Value = txt.ValueTranscoder
		}
		ci.RegisterDataType(typ)
	}
	return &typeResolver{connInfo: ci}
}

// findValue find the OID, and pgtype.ValueTranscoder for a Postgres type name.
func (tr *typeResolver) findValue(name string) (uint32, pgtype.ValueTranscoder, bool) {
	typ, ok := tr.connInfo.DataTypeForName(name)
	if !ok {
		return 0, nil, false
	}
	v := pgtype.NewValue(typ.Value)
	return typ.OID, v.(pgtype.ValueTranscoder), true
}

// setValue sets the value of a ValueTranscoder to a value that should always
// work and panics if it fails.
func (tr *typeResolver) setValue(vt pgtype.ValueTranscoder, val interface{}) pgtype.ValueTranscoder {
	if err := vt.Set(val); err != nil {
		panic(fmt.Sprintf("set ValueTranscoder %T to %+v: %s", vt, val, err))
	}
	return vt
}

const customTypesSQL = `SELECT 'some_text', 1::bigint;`

type CustomTypesRow struct {
	Column mytype.String `json:"?column?"`
	Int8   CustomInt     `json:"int8"`
}

// CustomTypes implements Querier.CustomTypes.
func (q *DBQuerier) CustomTypes(ctx context.Context) (CustomTypesRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "CustomTypes")
	row := q.conn.QueryRow(ctx, customTypesSQL)
	var item CustomTypesRow
	if err := row.Scan(&item.Column, &item.Int8); err != nil {
		return item, fmt.Errorf("query CustomTypes: %w", err)
	}
	return item, nil
}

// CustomTypesBatch implements Querier.CustomTypesBatch.
func (q *DBQuerier) CustomTypesBatch(batch genericBatch) {
	batch.Queue(customTypesSQL)
}

// CustomTypesScan implements Querier.CustomTypesScan.
func (q *DBQuerier) CustomTypesScan(results pgx.BatchResults) (CustomTypesRow, error) {
	row := results.QueryRow()
	var item CustomTypesRow
	if err := row.Scan(&item.Column, &item.Int8); err != nil {
		return item, fmt.Errorf("scan CustomTypesBatch row: %w", err)
	}
	return item, nil
}

const customMyIntSQL = `SELECT '5'::my_int as int5;`

// CustomMyInt implements Querier.CustomMyInt.
func (q *DBQuerier) CustomMyInt(ctx context.Context) (int, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "CustomMyInt")
	row := q.conn.QueryRow(ctx, customMyIntSQL)
	var item int
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("query CustomMyInt: %w", err)
	}
	return item, nil
}

// CustomMyIntBatch implements Querier.CustomMyIntBatch.
func (q *DBQuerier) CustomMyIntBatch(batch genericBatch) {
	batch.Queue(customMyIntSQL)
}

// CustomMyIntScan implements Querier.CustomMyIntScan.
func (q *DBQuerier) CustomMyIntScan(results pgx.BatchResults) (int, error) {
	row := results.QueryRow()
	var item int
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("scan CustomMyIntBatch row: %w", err)
	}
	return item, nil
}

const intArraySQL = `SELECT ARRAY ['5', '6', '7']::int[] as ints;`

// IntArray implements Querier.IntArray.
func (q *DBQuerier) IntArray(ctx context.Context) ([][]int32, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "IntArray")
	rows, err := q.conn.Query(ctx, intArraySQL)
	if err != nil {
		return nil, fmt.Errorf("query IntArray: %w", err)
	}
	defer rows.Close()
	items := [][]int32{}
	for rows.Next() {
		var item []int32
		if err := rows.Scan(&item); err != nil {
			return nil, fmt.Errorf("scan IntArray row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close IntArray rows: %w", err)
	}
	return items, err
}

// IntArrayBatch implements Querier.IntArrayBatch.
func (q *DBQuerier) IntArrayBatch(batch genericBatch) {
	batch.Queue(intArraySQL)
}

// IntArrayScan implements Querier.IntArrayScan.
func (q *DBQuerier) IntArrayScan(results pgx.BatchResults) ([][]int32, error) {
	rows, err := results.Query()
	if err != nil {
		return nil, fmt.Errorf("query IntArrayBatch: %w", err)
	}
	defer rows.Close()
	items := [][]int32{}
	for rows.Next() {
		var item []int32
		if err := rows.Scan(&item); err != nil {
			return nil, fmt.Errorf("scan IntArrayBatch row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close IntArrayBatch rows: %w", err)
	}
	return items, err
}

// textPreferrer wraps a pgtype.ValueTranscoder and sets the preferred encoding
// format to text instead binary (the default). pggen uses the text format
// when the OID is unknownOID because the binary format requires the OID.
// Typically occurs if the results from QueryAllDataTypes aren't passed to
// NewQuerierConfig.
type textPreferrer struct {
	pgtype.ValueTranscoder
	typeName string
}

// PreferredParamFormat implements pgtype.ParamFormatPreferrer.
func (t textPreferrer) PreferredParamFormat() int16 { return pgtype.TextFormatCode }

func (t textPreferrer) NewTypeValue() pgtype.Value {
	return textPreferrer{ValueTranscoder: pgtype.NewValue(t.ValueTranscoder).(pgtype.ValueTranscoder), typeName: t.typeName}
}

func (t textPreferrer) TypeName() string {
	return t.typeName
}

// unknownOID means we don't know the OID for a type. This is okay for decoding
// because pgx call DecodeText or DecodeBinary without requiring the OID. For
// encoding parameters, pggen uses textPreferrer if the OID is unknown.
const unknownOID = 0
