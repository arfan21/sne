package urlshortenerrepo

import (
	"context"
	"errors"
	"fmt"

	"github.com/arfan21/backend-test/internal/entity"
	"github.com/arfan21/backend-test/pkg/constant"
	dbpostgres "github.com/arfan21/backend-test/pkg/db/postgres"
	"github.com/jackc/pgx/v5"
)

type RepositoryPostgres struct {
	db dbpostgres.Queryer
}

func NewPostgres(db dbpostgres.Queryer) *RepositoryPostgres {
	return &RepositoryPostgres{
		db: db,
	}
}

func (r RepositoryPostgres) Set(ctx context.Context, value entity.URLShortener) (err error) {
	query := `INSERT INTO url_shortener (short_url, long_url, expired_at) VALUES ($1, $2, $3)`

	_, err = r.db.Exec(ctx, query, value.ShortURL, value.LongURL, value.ExpiredAt)
	if err != nil {
		err = fmt.Errorf("urlshortener.repositoryPostgres.Create: failed to insert data %w", err)
		return
	}

	return
}

func (r RepositoryPostgres) Get(ctx context.Context, shortURL string) (res entity.URLShortener, err error) {
	query := `SELECT short_url, long_url, expired_at FROM url_shortener WHERE short_url = $1`

	row := r.db.QueryRow(ctx, query, shortURL)

	err = row.Scan(&res.ShortURL, &res.LongURL, &res.ExpiredAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = constant.ErrURLNotFound
		}
		err = fmt.Errorf("urlshortener.repositoryPostgres.Get: failed to get data %w", err)
		return
	}

	return
}
