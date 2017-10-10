package migrations

import (
	"github.com/go-pg/migrations"
	log "github.com/sirupsen/logrus"
)

func init() {
	migrations.Register(func(db migrations.DB) error {
		log.Info("migrate 20171010114357_create-locales")
		_, err := db.Exec(`
CREATE TABLE locales (
  id BIGSERIAL PRIMARY KEY,
  code VARCHAR(255) NOT NULL,
  lang VARCHAR(8) NOT NULL DEFAULT 'en-US',
  message TEXT NOT NULL,
  created_at TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
  updated_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);
CREATE UNIQUE INDEX idx_locales_code_lang ON locales (code, lang);
CREATE INDEX idx_locales_code ON locales (code);
CREATE INDEX idx_locales_lang ON locales (lang);
			`)
		return err
	}, func(db migrations.DB) error {
		log.Info("rollback 20171010114357_create-locales")
		_, err := db.Exec(`DROP TABLE locales;`)
		return err
	})
}
