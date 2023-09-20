package quicksql

import "io"

type Insert interface {
	// Convert query into `INSERT IGNORE INTO`
	Ignore() Insert

	// Convert query into `REPLACE INTO`
	Replace() Insert

	// Start a new insert statement every `n` rows
	Every(int) Insert

	// Add another element to insert query
	Add(...any) Insert

	// Flush pending data to buffer
	Flush() Insert

	// Get last error
	Err() error
}

type insert struct {
	*spacer

	err error
}

func NewInsert(writer io.Writer, table string, columns ...any) Insert {
	var inserr error

	tableName, err := QuoteIdentifier(table)
	if err != nil {
		inserr = err
	}

	cols, err := QuoteColumnNames(columns...)
	if err != nil {
		inserr = err
	}

	return &insert{
		spacer: &spacer{
			wri:   writer,
			top:   "INSERT INTO %s (%s) VALUES\n",
			mid:   "\t(%s),\n",
			bot:   "\t(%s);\n\n",
			split: 1000,

			headerinfo: []any{tableName, cols},
		},

		err: inserr,
	}
}

func (i *insert) Ignore() Insert {
	i.top = "INSERT IGNORE INTO %s (%s) VALUES\n"
	return i
}

func (i *insert) Replace() Insert {
	i.top = "REPLACE INTO %s (%s) VALUES\n"
	return i
}

func (i *insert) Every(n int) Insert {
	i.split = n
	return i
}

func (i *insert) Add(args ...any) Insert {
	if i.err != nil {
		return i
	}

	if err := i.push(QuoteMultiple(args...)); err != nil {
		i.err = err
	}

	return i
}

func (i *insert) Flush() Insert {
	if i.err != nil {
		return i
	}

	if err := i.flush(); err != nil {
		i.err = err
	}

	return i
}

func (i *insert) Err() error {
	return i.err
}
