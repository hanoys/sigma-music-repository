package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/hanoys/sigma-music-core/domain"
	"github.com/hanoys/sigma-music-core/ports"
	"github.com/hanoys/sigma-music-core/util"
	entity2 "github.com/hanoys/sigma-music-repository/repository/postgres/entity"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

const (
	subscriptionGetByID = "SELECT * FROM subscriptions WHERE id = $1"
)

type PostgresSubscriptionRepository struct {
	db *sqlx.DB
}

func NewPostgresSubscriptionRepository(db *sqlx.DB) *PostgresSubscriptionRepository {
	return &PostgresSubscriptionRepository{db: db}
}

func (sr *PostgresSubscriptionRepository) Create(ctx context.Context, sub domain.Subscription) (domain.Subscription, error) {
	pgSubscription := entity2.NewPgSuscription(sub)
	queryString := entity2.InsertQueryString(pgSubscription, "subscriptions")
	_, err := sr.db.NamedExecContext(ctx, queryString, pgSubscription)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return domain.Subscription{}, util.WrapError(ports.ErrSubDuplicate, err)
			}
		}
		return domain.Subscription{}, util.WrapError(ports.ErrInternalSubRepo, err)
	}

	var createdSubscription entity2.PgSubscription
	err = sr.db.GetContext(ctx, &createdSubscription, subscriptionGetByID, pgSubscription.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Subscription{}, util.WrapError(ports.ErrSubIDNotFound, err)
		}
		return domain.Subscription{}, util.WrapError(ports.ErrInternalSubRepo, err)
	}

	return createdSubscription.ToDomain(), nil
}
