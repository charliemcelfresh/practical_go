package practical_go

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type database interface {
	sqlx.ExtContext
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type Item struct {
	ItemID    string `db:"item_id"`
	Name      string `db:"name"`
	CreatedAt string `db:"created_at"`
}

type repository struct {
	db database
}

func NewRepository(db *sqlx.DB) repository {
	return repository{
		db,
	}
}

func (r repository) CreateItem(ctx context.Context, name string) error {
	statement := `
		INSERT INTO item (name) values ($1);
	`
	_, err := r.db.ExecContext(ctx, statement, name)
	return err
}

func (r repository) GetItem(ctx context.Context, itemID string) (Item, error) {
	itemToReturn := Item{}
	statement := `
		SELECT
			item_id, name, created_at
		FROM
			item
		WHERE
			item_id = $1;
		;
	`
	err := r.db.GetContext(ctx, &itemToReturn, statement, itemID)
	return itemToReturn, err
}
