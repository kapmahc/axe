package migrations

import (
	"github.com/go-pg/migrations"
	log "github.com/sirupsen/logrus"
)

func init() {
	migrations.Register(func(db migrations.DB) error {
		log.Info("migrate 20171010114403_create-settings")
		_, err := db.Exec(`
CREATE TABLE settings (
  id BIGSERIAL PRIMARY KEY,
  key VARCHAR(255) NOT NULL,
  value BYTEA NOT NULL,
  encode BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);
CREATE UNIQUE INDEX idx_settings_key ON settings (key);
			`)
		return err
	}, func(db migrations.DB) error {
		log.Info("rollback 20171010114403_create-settings")
		_, err := db.Exec(`DROP TABLE settings;`)
		return err
	})
}
