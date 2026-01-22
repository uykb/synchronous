package processor

import (
	"context"
	"crypto-sync-bot/internal/database"
	"crypto-sync-bot/internal/models"
	"log"
	"time"
)

type Reconciler struct {
	executors map[string]models.ExchangeExecutor
}

func NewReconciler(execs []models.ExchangeExecutor) *Reconciler {
	m := make(map[string]models.ExchangeExecutor)
	for _, e := range execs {
		m[e.Name()] = e
	}
	return &Reconciler{executors: m}
}

func (r *Reconciler) Start(ctx context.Context) {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	log.Println("Reconciler started")
	for {
		select {
		case <-ctx.Done():
			log.Println("Reconciler stopping")
			return
		case <-ticker.C:
			r.reconcileOrders()
		}
	}
}

func (r *Reconciler) reconcileOrders() {
	// Query non-final orders. 
	// We include 'success' because PlaceOrder might return 'success' but the order is still 'NEW' or 'PARTIALLY_FILLED' on the exchange.
	rows, err := database.DB.Query("SELECT exchange, symbol, order_id FROM orders WHERE status NOT IN ('FILLED', 'CANCELLED', 'REJECTED', 'failed') AND order_id != ''")
	if err != nil {
		log.Printf("Reconciler: failed to query orders: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var exchange, symbol, orderID string
		if err := rows.Scan(&exchange, &symbol, &orderID); err != nil {
			log.Printf("Reconciler: failed to scan row: %v", err)
			continue
		}

		exec, ok := r.executors[exchange]
		if !ok {
			log.Printf("Reconciler: executor not found for %s", exchange)
			continue
		}

		res, err := exec.GetOrder(orderID, symbol)
		if err != nil {
			log.Printf("Reconciler: failed to get order %s from %s: %v", orderID, exchange, err)
			continue
		}

		_, err = database.DB.Exec("UPDATE orders SET status = ? WHERE order_id = ?", res.Status, orderID)
		if err != nil {
			log.Printf("Reconciler: failed to update order %s in DB: %v", orderID, err)
		} else {
			log.Printf("Reconciler: updated order %s status to %s", orderID, res.Status)
		}
	}
}
