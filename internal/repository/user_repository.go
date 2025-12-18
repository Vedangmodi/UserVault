package repository

import (
	"context"
	"time"

	sqlcdb "uservault/db/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
)

type UserRepository interface {
	Create(ctx context.Context, name string, dob time.Time) (sqlcdb.User, error)
	Get(ctx context.Context, id int64) (sqlcdb.User, error)
	List(ctx context.Context, limit, offset int32) ([]sqlcdb.User, error)
	Update(ctx context.Context, id int64, name string, dob time.Time) (sqlcdb.User, error)
	Delete(ctx context.Context, id int64) error
}

type userRepository struct {
	q *sqlcdb.Queries
}

// NewUserRepository accepts anything that implements sqlcdb.DBTX
// (e.g. *pgxpool.Pool or pgx.Conn).
func NewUserRepository(db sqlcdb.DBTX) UserRepository {
	return &userRepository{
		q: sqlcdb.New(db),
	}
}

func (r *userRepository) Create(ctx context.Context, name string, dob time.Time) (sqlcdb.User, error) {
	return r.q.CreateUser(ctx, sqlcdb.CreateUserParams{
		Name: name,
		Dob: pgtype.Date{
			Time:  dob,
			Valid: true,
		},
	})
}

func (r *userRepository) Get(ctx context.Context, id int64) (sqlcdb.User, error) {
	return r.q.GetUser(ctx, int32(id))
}

func (r *userRepository) List(ctx context.Context, limit, offset int32) ([]sqlcdb.User, error) {
	return r.q.ListUsers(ctx, sqlcdb.ListUsersParams{
		Limit:  limit,
		Offset: offset,
	})
}

func (r *userRepository) Update(ctx context.Context, id int64, name string, dob time.Time) (sqlcdb.User, error) {
	return r.q.UpdateUser(ctx, sqlcdb.UpdateUserParams{
		ID:   int32(id),
		Name: name,
		Dob: pgtype.Date{
			Time:  dob,
			Valid: true,
		},
	})
}

func (r *userRepository) Delete(ctx context.Context, id int64) error {
	return r.q.DeleteUser(ctx, int32(id))
}
