package httpserver

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/woragis/minecraft-campus-backend/server/internal/apperrors"
	"gorm.io/gorm"
)

func handleReady(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sqlDB, err := db.DB()
		if err != nil {
			log.Printf("ready: sql db: %v", err)
			apperrors.WriteError(w, apperrors.UnavailableCause(apperrors.CodeReadyGetHandlerSqlGetterFailed, apperrors.MsgReadyGetHandlerSqlGetterFailed, err))
			return
		}
		if err := sqlDB.PingContext(r.Context()); err != nil {
			log.Printf("ready: db ping: %v", err)
			apperrors.WriteError(w, apperrors.UnavailableCause(apperrors.CodeReadyGetHandlerDatabasePingFailed, apperrors.MsgReadyGetHandlerDatabasePingFailed, err))
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ready"})
	}
}
