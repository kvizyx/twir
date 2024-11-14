package command

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/satont/twir/apps/parser-new/internal/entity"
	"go.uber.org/fx"
)

const (
	tableCommandsResponses = "channels_commands_responses"
)

type Repository struct {
	db *pgxpool.Pool
}

var _ RepositoryContract = (*Repository)(nil)

type Params struct {
	fx.In

	DB *pgxpool.Pool
}

func NewRepository(params Params) Repository {
	return Repository{
		db: params.DB,
	}
}

func (r *Repository) GetResponsesByID(
	ctx context.Context,
	commandID uuid.UUID,
) ([]entity.CommandResponse, error) {
	query, args, err := squirrel.
		Select("*").
		From(tableCommandsResponses).
		Where(
			squirrel.Eq{
				"commandId": commandID,
			},
		).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	responses, err := pgx.CollectRows(
		rows,
		pgx.RowToStructByName[ResponseDTO],
	)
	if err != nil {
		return nil, err
	}

	return fromCommandResponses(responses), nil
}
