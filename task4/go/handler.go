package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	db  *sql.DB
	log *slog.Logger
}

func NewHandler(db *sql.DB, log *slog.Logger) *Handler {
	return &Handler{
		db:  db,
		log: log,
	}
}

// handler получает устройство по id и пишет в лог таблицу audit_log
func (h *Handler) DeviceHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		h.log.Debug("missing id")
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.log.Debug("invalid id")
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	go func() {
		select {
		case <-time.After(5 * time.Second):
			fmt.Println("long debug operation finished")
		case <-ctx.Done():
			fmt.Println("operation cancelled:", ctx.Err())
			return
		}
	}()

	query := "SELECT id, hostname, ip FROM devices WHERE id = $1"
	row := h.db.QueryRowContext(ctx, query, id)

	var d Device
	err = row.Scan(&d.ID, &d.Hostname, &d.IP)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			h.log.Debug("device not found", slog.Int64("id", id))
			http.Error(w, "device not found", http.StatusNotFound)
		default:
			h.log.Error("failed to write audit log", slog.String("error", err.Error()))
			http.Error(w, "db error", http.StatusInternalServerError)
		}
		return
	}

	_, err = h.db.ExecContext(ctx, "INSERT INTO audit_log(device_id, ts, action) VALUES ($1, now(), 'view')", d.ID)
	if err != nil {
		h.log.Error("failed to write audit log", slog.String("error", err.Error()))
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Device: %s (%s)", d.Hostname, d.IP)))
}
