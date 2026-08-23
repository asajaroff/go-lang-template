package items_test

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	itemsdb "github.com/asajaroff/go-lang-template/internal/items/db"
)

// fakeQuerier is an in-memory stand-in for itemsdb.Querier — the same
// dependency-substitution pattern the Python sibling template uses to
// override its DB session for tests, just expressed as a Go interface
// instead of a FastAPI dependency override.
type fakeQuerier struct {
	items map[uuid.UUID]itemsdb.Item
}

func newFakeQuerier() *fakeQuerier {
	return &fakeQuerier{items: make(map[uuid.UUID]itemsdb.Item)}
}

func (f *fakeQuerier) CreateItem(_ context.Context, arg itemsdb.CreateItemParams) (itemsdb.Item, error) {
	now := pgtype.Timestamptz{Time: time.Now(), Valid: true}
	item := itemsdb.Item{
		ID:          uuid.New(),
		Name:        arg.Name,
		Description: arg.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	f.items[item.ID] = item
	return item, nil
}

func (f *fakeQuerier) GetItem(_ context.Context, id uuid.UUID) (itemsdb.Item, error) {
	item, ok := f.items[id]
	if !ok {
		return itemsdb.Item{}, pgx.ErrNoRows
	}
	return item, nil
}

func (f *fakeQuerier) ListItems(_ context.Context) ([]itemsdb.Item, error) {
	out := make([]itemsdb.Item, 0, len(f.items))
	for _, item := range f.items {
		out = append(out, item)
	}
	return out, nil
}

func (f *fakeQuerier) UpdateItem(_ context.Context, arg itemsdb.UpdateItemParams) (itemsdb.Item, error) {
	item, ok := f.items[arg.ID]
	if !ok {
		return itemsdb.Item{}, pgx.ErrNoRows
	}
	if arg.Name.Valid {
		item.Name = arg.Name.String
	}
	if arg.Description.Valid {
		item.Description = arg.Description
	}
	item.UpdatedAt = pgtype.Timestamptz{Time: time.Now(), Valid: true}
	f.items[item.ID] = item
	return item, nil
}

func (f *fakeQuerier) DeleteItem(_ context.Context, id uuid.UUID) (int64, error) {
	if _, ok := f.items[id]; !ok {
		return 0, nil
	}
	delete(f.items, id)
	return 1, nil
}
