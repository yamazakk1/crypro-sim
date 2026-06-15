package repository

import (
	"context"
	"time"

	"crypto-simulator/services/trading/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionRepo struct {
	pool *pgxpool.Pool
}

func NewTransactionRepo(pool *pgxpool.Pool) *TransactionRepo {
	return &TransactionRepo{pool: pool}
}

func (r *TransactionRepo) Create(ctx context.Context, tx pgx.Tx, userID, assetID, txType string, quantity, price, total float64) (*models.Transaction, error) {
	var txID string
	var createdAt time.Time

	err := tx.QueryRow(ctx,
		`INSERT INTO transactions (id, user_id, asset_id, type, quantity, asset_price, total_usdt)
         VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6) RETURNING id, created_at`,
		userID, assetID, txType, quantity, price, total).Scan(&txID, &createdAt)
	if err != nil {
		return nil, err
	}

	return &models.Transaction{
		ID:         txID,
		UserID:     userID,
		AssetID:    assetID,
		Type:       txType,
		Quantity:   quantity,
		AssetPrice: price,
		TotalUSDT:  total,
		CreatedAt:  createdAt.Format(time.RFC3339),
	}, nil
}

func (r *TransactionRepo) GetByUserID(ctx context.Context, userID string) ([]models.Transaction, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT t.id, t.type, a.symbol, t.quantity, t.asset_price, t.total_usdt, t.created_at
         FROM transactions t JOIN assets a ON t.asset_id = a.id
         WHERE t.user_id = $1 ORDER BY t.created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txs []models.Transaction
	for rows.Next() {
		var tx models.Transaction
		var createdAt time.Time
		if err := rows.Scan(&tx.ID, &tx.Type, &tx.AssetSymbol, &tx.Quantity, &tx.AssetPrice, &tx.TotalUSDT, &createdAt); err != nil {
			return nil, err
		}
		tx.UserID = userID
		tx.CreatedAt = createdAt.Format(time.RFC3339)
		txs = append(txs, tx)
	}
	return txs, nil
}
