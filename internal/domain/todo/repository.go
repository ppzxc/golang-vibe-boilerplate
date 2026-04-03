package todo

import "context"

// Repository is the outbound port for persisting and retrieving Todos.
// Implementations are in adapter/postgresrepo.
type Repository interface {
	Save(ctx context.Context, todo *Todo) error
	FindByID(ctx context.Context, id string) (*Todo, error)
	FindAll(ctx context.Context, cursor string, pageSize int) (todos []*Todo, nextCursor string, err error)
	Update(ctx context.Context, id string, fn func(*Todo) (*Todo, error)) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context) (int64, error)
}
