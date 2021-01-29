// Code generated by pggen. DO NOT EDIT.

package author

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

// Querier is a typesafe Go interface backed by SQL queries.
//
// Methods ending with Batch enqueue a query to run later in a pgx.Batch. After
// calling SendBatch on pgx.Conn, pgxpool.Pool, or pgx.Tx, use the Scan methods
// to parse the results.
type Querier interface {
	// FindAuthorById finds one (or zero) authors by ID.
	FindAuthorByID(ctx context.Context, authorID int32) (FindAuthorByIDRow, error)
	// FindAuthorByIDBatch enqueues a FindAuthorByID query into batch to be executed
	// later by the batch.
	FindAuthorByIDBatch(ctx context.Context, batch *pgx.Batch, authorID int32)
	// FindAuthorByIDScan scans the result of an executed FindAuthorByIDBatch query.
	FindAuthorByIDScan(results pgx.BatchResults) (FindAuthorByIDRow, error)

	// FindAuthors finds authors by first name.
	FindAuthors(ctx context.Context, firstName string) ([]FindAuthorsRow, error)
	// FindAuthorsBatch enqueues a FindAuthors query into batch to be executed
	// later by the batch.
	FindAuthorsBatch(ctx context.Context, batch *pgx.Batch, firstName string)
	// FindAuthorsScan scans the result of an executed FindAuthorsBatch query.
	FindAuthorsScan(results pgx.BatchResults) ([]FindAuthorsRow, error)

	// DeleteAuthors deletes authors with a first name of "joe".
	DeleteAuthors(ctx context.Context) (pgconn.CommandTag, error)
	// DeleteAuthorsBatch enqueues a DeleteAuthors query into batch to be executed
	// later by the batch.
	DeleteAuthorsBatch(ctx context.Context, batch *pgx.Batch)
	// DeleteAuthorsScan scans the result of an executed DeleteAuthorsBatch query.
	DeleteAuthorsScan(results pgx.BatchResults) (pgconn.CommandTag, error)

	// InsertAuthor inserts an author by name and returns the ID.
	InsertAuthor(ctx context.Context, firstName string, lastName string) (int32, error)
	// InsertAuthorBatch enqueues a InsertAuthor query into batch to be executed
	// later by the batch.
	InsertAuthorBatch(ctx context.Context, batch *pgx.Batch, firstName string, lastName string)
	// InsertAuthorScan scans the result of an executed InsertAuthorBatch query.
	InsertAuthorScan(results pgx.BatchResults) (int32, error)
}

type DBQuerier struct {
	conn genericConn
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

// NewQuerier creates a DBQuerier that implements Querier. conn is typically
// *pgx.Conn, pgx.Tx, or *pgxpool.Pool.
func NewQuerier(conn genericConn) *DBQuerier {
	return &DBQuerier{
		conn: conn,
	}
}

// WithTx creates a new DBQuerier that uses the transaction to run all queries.
func (q *DBQuerier) WithTx(tx pgx.Tx) (*DBQuerier, error) {
	return &DBQuerier{conn: tx}, nil
}

const findAuthorByIDSQL = `SELECT * FROM author WHERE author_id = $1;`

type FindAuthorByIDRow struct {
	AuthorID  int32
	FirstName string
	LastName  string
	Suffix    pgtype.Text
}

// FindAuthorByID implements Querier.FindAuthorByID.
func (q *DBQuerier) FindAuthorByID(ctx context.Context, authorID int32) (FindAuthorByIDRow, error) {
	row := q.conn.QueryRow(ctx, findAuthorByIDSQL, authorID)
	var item FindAuthorByIDRow
	if err := row.Scan(&item.AuthorID, &item.FirstName, &item.LastName, &item.Suffix); err != nil {
		return item, fmt.Errorf("query FindAuthorByID: %w", err)
	}
	return item, nil
}

// FindAuthorByIDBatch implements Querier.FindAuthorByIDBatch.
func (q *DBQuerier) FindAuthorByIDBatch(ctx context.Context, batch *pgx.Batch, authorID int32) {
	batch.Queue(findAuthorByIDSQL, authorID)
}

// FindAuthorByIDScan implements Querier.FindAuthorByIDScan.
func (q *DBQuerier) FindAuthorByIDScan(results pgx.BatchResults) (FindAuthorByIDRow, error) {
	row := results.QueryRow()
	var item FindAuthorByIDRow
	if err := row.Scan(&item.AuthorID, &item.FirstName, &item.LastName, &item.Suffix); err != nil {
		return item, fmt.Errorf("scan FindAuthorByIDBatch row: %w", err)
	}
	return item, nil
}

const findAuthorsSQL = `SELECT * FROM author WHERE first_name = $1;`

type FindAuthorsRow struct {
	AuthorID  int32
	FirstName string
	LastName  string
	Suffix    pgtype.Text
}

// FindAuthors implements Querier.FindAuthors.
func (q *DBQuerier) FindAuthors(ctx context.Context, firstName string) ([]FindAuthorsRow, error) {
	rows, err := q.conn.Query(ctx, findAuthorsSQL, firstName)
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return nil, fmt.Errorf("query FindAuthors: %w", err)
	}
	var items []FindAuthorsRow
	for rows.Next() {
		var item FindAuthorsRow
		if err := rows.Scan(&item.AuthorID, &item.FirstName, &item.LastName, &item.Suffix); err != nil {
			return nil, fmt.Errorf("scan FindAuthors row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, err
}

// FindAuthorsBatch implements Querier.FindAuthorsBatch.
func (q *DBQuerier) FindAuthorsBatch(ctx context.Context, batch *pgx.Batch, firstName string) {
	batch.Queue(findAuthorsSQL, firstName)
}

// FindAuthorsScan implements Querier.FindAuthorsScan.
func (q *DBQuerier) FindAuthorsScan(results pgx.BatchResults) ([]FindAuthorsRow, error) {
	rows, err := results.Query()
	if rows != nil {
		defer rows.Close()
	}
	if err != nil {
		return nil, err
	}
	var items []FindAuthorsRow
	for rows.Next() {
		var item FindAuthorsRow
		if err := rows.Scan(&item.AuthorID, &item.FirstName, &item.LastName, &item.Suffix); err != nil {
			return nil, fmt.Errorf("scan FindAuthorsBatch row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, err
}

const deleteAuthorsSQL = `DELETE FROM author WHERE first_name = 'joe';`

// DeleteAuthors implements Querier.DeleteAuthors.
func (q *DBQuerier) DeleteAuthors(ctx context.Context) (pgconn.CommandTag, error) {
	cmdTag, err := q.conn.Exec(ctx, deleteAuthorsSQL)
	return cmdTag, err
}

// DeleteAuthorsBatch implements Querier.DeleteAuthorsBatch.
func (q *DBQuerier) DeleteAuthorsBatch(ctx context.Context, batch *pgx.Batch) {
	batch.Queue(deleteAuthorsSQL)
}

// DeleteAuthorsScan implements Querier.DeleteAuthorsScan.
func (q *DBQuerier) DeleteAuthorsScan(results pgx.BatchResults) (pgconn.CommandTag, error) {
	cmdTag, err := results.Exec()
	return cmdTag, err
}

const insertAuthorSQL = `INSERT INTO author (first_name, last_name)
VALUES ($1, $2)
RETURNING author_id;`

// InsertAuthor implements Querier.InsertAuthor.
func (q *DBQuerier) InsertAuthor(ctx context.Context, firstName string, lastName string) (int32, error) {
	row := q.conn.QueryRow(ctx, insertAuthorSQL, firstName, lastName)
	var item int32
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("query InsertAuthor: %w", err)
	}
	return item, nil
}

// InsertAuthorBatch implements Querier.InsertAuthorBatch.
func (q *DBQuerier) InsertAuthorBatch(ctx context.Context, batch *pgx.Batch, firstName string, lastName string) {
	batch.Queue(insertAuthorSQL, firstName, lastName)
}

// InsertAuthorScan implements Querier.InsertAuthorScan.
func (q *DBQuerier) InsertAuthorScan(results pgx.BatchResults) (int32, error) {
	row := results.QueryRow()
	var item int32
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("scan InsertAuthorBatch row: %w", err)
	}
	return item, nil
}
