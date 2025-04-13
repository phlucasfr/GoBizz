// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type Querier interface {
	CheckCustomSlugExists(ctx context.Context, customSlug pgtype.Text) (bool, error)
	CheckShortURLExists(ctx context.Context, shortUrl string) (bool, error)
	CountLinksByCustomer(ctx context.Context, customerID pgtype.UUID) (int64, error)
	CreateLink(ctx context.Context, arg CreateLinkParams) (Link, error)
	DeleteExpiredLinks(ctx context.Context) error
	DeleteLink(ctx context.Context, id pgtype.UUID) error
	GetExpiredLinks(ctx context.Context) ([]Link, error)
	GetLinkByCustomSlug(ctx context.Context, customSlug pgtype.Text) (Link, error)
	GetLinkByID(ctx context.Context, id pgtype.UUID) (Link, error)
	GetLinkByShortURL(ctx context.Context, shortUrl string) (Link, error)
	GetLinksByCustomer(ctx context.Context, arg GetLinksByCustomerParams) ([]GetLinksByCustomerRow, error)
	UpdateLink(ctx context.Context, arg UpdateLinkParams) (Link, error)
	UpdateLinkClicks(ctx context.Context, id pgtype.UUID) (Link, error)
}

var _ Querier = (*Queries)(nil)
