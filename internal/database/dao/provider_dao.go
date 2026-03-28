package dao

import (
	"database/sql"
	"encoding/json"
)

type ProviderRow struct {
	ID             int
	Name           string
	SettingsConfig string
	IsCurrent      bool
	SortIndex      int
}

type ProviderDAO struct {
	db      *sql.DB
	appType string
}

func NewProviderDAO(db *sql.DB, appType string) *ProviderDAO {
	return &ProviderDAO{db: db, appType: appType}
}

func (d *ProviderDAO) ListAll() ([]ProviderRow, error) {
	rows, err := d.db.Query(
		`SELECT id, name, settings_config, is_current, sort_index
		 FROM providers WHERE app_type=? ORDER BY sort_index`, d.appType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []ProviderRow
	for rows.Next() {
		var r ProviderRow
		var isCurrent int
		if err := rows.Scan(&r.ID, &r.Name, &r.SettingsConfig, &isCurrent, &r.SortIndex); err != nil {
			return nil, err
		}
		r.IsCurrent = isCurrent == 1
		result = append(result, r)
	}
	return result, rows.Err()
}

func (d *ProviderDAO) GetCurrent() (*ProviderRow, error) {
	var r ProviderRow
	var isCurrent int
	err := d.db.QueryRow(
		`SELECT id, name, settings_config, is_current, sort_index
		 FROM providers WHERE app_type=? AND is_current=1 LIMIT 1`, d.appType).
		Scan(&r.ID, &r.Name, &r.SettingsConfig, &isCurrent, &r.SortIndex)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	r.IsCurrent = true
	return &r, nil
}

func (d *ProviderDAO) GetByName(name string) (*ProviderRow, error) {
	var r ProviderRow
	var isCurrent int
	err := d.db.QueryRow(
		`SELECT id, name, settings_config, is_current, sort_index
		 FROM providers WHERE app_type=? AND name=? LIMIT 1`, d.appType, name).
		Scan(&r.ID, &r.Name, &r.SettingsConfig, &isCurrent, &r.SortIndex)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	r.IsCurrent = isCurrent == 1
	return &r, nil
}

func (d *ProviderDAO) SetCurrent(name string) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.Exec(`UPDATE providers SET is_current=0 WHERE app_type=?`, d.appType); err != nil {
		return err
	}
	if _, err := tx.Exec(`UPDATE providers SET is_current=1 WHERE app_type=? AND name=?`, d.appType, name); err != nil {
		return err
	}
	return tx.Commit()
}

func ParseSettingsConfig(raw string) (env map[string]string, model string, err error) {
	var cfg struct {
		Env   map[string]string `json:"env"`
		Model string            `json:"model"`
	}
	if err := json.Unmarshal([]byte(raw), &cfg); err != nil {
		return nil, "", err
	}
	return cfg.Env, cfg.Model, nil
}
