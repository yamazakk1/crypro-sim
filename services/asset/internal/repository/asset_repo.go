package repository

import (
	"context"
	"log"
	"time"

	"crypto-simulator/services/asset/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AssetRepo struct {
	pool *pgxpool.Pool
}

func NewAssetRepo(pool *pgxpool.Pool) *AssetRepo {
	log.Println("asset-repo: initialized")
	return &AssetRepo{pool: pool}
}

func (r *AssetRepo) List(ctx context.Context) ([]*models.Asset, error) {
	log.Println("asset-repo: List called")

	rows, err := r.pool.Query(ctx, `SELECT id, symbol, full_name, initial_price, is_active, created_at FROM assets ORDER BY created_at DESC`)
	if err != nil {
		log.Printf("asset-repo: List: query error: %v", err)
		return nil, err
	}
	defer rows.Close()

	iterationCount := 0
	var assets []*models.Asset
	for rows.Next() {
		iterationCount++
		log.Printf("asset-repo: List: iteration #%d started", iterationCount)

		a := &models.Asset{}
		err := rows.Scan(&a.ID, &a.Symbol, &a.Fullname, &a.InitialPrice, &a.IsActive, &a.CreatedAt)
		if err != nil {
			log.Printf("asset-repo: List: scan error on iteration #%d: %v", iterationCount, err)
			continue
		}
		log.Printf("asset-repo: scanned: id=%s, symbol=%s, is_active=%v", a.ID, a.Symbol, a.IsActive)
		assets = append(assets, a)
	}

	if err := rows.Err(); err != nil {
		log.Printf("asset-repo: List: rows.Err(): %v", err)
	}

	log.Printf("asset-repo: List: total iterations: %d", iterationCount)
	log.Printf("asset-repo: List: success, count=%d", len(assets))
	return assets, nil
}

func (r *AssetRepo) GetByID(ctx context.Context, id string) (*models.Asset, error) {
	log.Printf("asset-repo: GetByID called: id=%s", id)
	var a models.Asset
	err := r.pool.QueryRow(ctx, `SELECT id, symbol, full_name, initial_price, is_active, created_at, updated_at FROM assets WHERE id = $1`, id).
		Scan(&a.ID, &a.Symbol, &a.Fullname, &a.InitialPrice, &a.IsActive, &a.CreatedAt, &a.UpdatedAt)
	if err == pgx.ErrNoRows {
		log.Printf("asset-repo: GetByID: not found: id=%s", id)
		return nil, ErrAssetNotFound
	}
	if err != nil {
		log.Printf("asset-repo: GetByID: error: %v", err)
		return nil, err
	}
	log.Printf("asset-repo: GetByID: success, symbol=%s", a.Symbol)
	return &a, nil
}

func (r *AssetRepo) Create(ctx context.Context, symbol, fullname string, initPrice float64) (*models.Asset, error) {
	log.Printf("asset-repo: Create called: symbol=%s, full_name=%s", symbol, fullname)
	a := &models.Asset{
		ID:       uuid.New().String(),
		Symbol:   symbol,
		Fullname: fullname,
		InitialPrice: initPrice,
		IsActive: true,
	}
	_, err := r.pool.Exec(ctx, `INSERT INTO assets (id, symbol, full_name, initial_price, is_active) VALUES ($1, $2, $3, $4, $5)`,
		a.ID, symbol, fullname, initPrice, true)
	if err != nil {
		log.Printf("asset-repo: Create: error: %v", err)
		return nil, err
	}
	log.Printf("asset-repo: Create: success, id=%s", a.ID)
	return a, nil
}

func (r *AssetRepo) Update(ctx context.Context, id, fullname string) error {
	log.Printf("asset-repo: Update called: id=%s, full_name=%s", id, fullname)
	_, err := r.pool.Exec(ctx, `UPDATE assets SET full_name = $1, updated_at = $2 WHERE id = $3`, fullname, time.Now(), id)
	if err != nil {
		log.Printf("asset-repo: Update: error: %v", err)
		return err
	}
	log.Printf("asset-repo: Update: success, id=%s", id)
	return nil
}

func (r *AssetRepo) SetActive(ctx context.Context, id string, active bool) error {
	log.Printf("asset-repo: SetActive called: id=%s, active=%v", id, active)
	_, err := r.pool.Exec(ctx, `UPDATE assets SET is_active = $1, updated_at = $2 WHERE id = $3`, active, time.Now(), id)
	if err != nil {
		log.Printf("asset-repo: SetActive: error: %v", err)
		return err
	}
	log.Printf("asset-repo: SetActive: success, id=%s", id)
	return nil
}

func (r *AssetRepo) Delete(ctx context.Context, id string) error {
	log.Printf("asset-repo: Delete called: id=%s", id)
	_, err := r.pool.Exec(ctx, `DELETE FROM assets WHERE id = $1`, id)
	if err != nil {
		log.Printf("asset-repo: Delete: error: %v", err)
		return err
	}
	log.Printf("asset-repo: Delete: success, id=%s", id)
	return nil
}
