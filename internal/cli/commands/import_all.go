package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/farion1231/ccmux/internal/config"
	"github.com/farion1231/ccmux/internal/database"
	"github.com/farion1231/ccmux/internal/database/dao"
)

// RunImportAll imports providers from cc-switch.db into JSON profiles.
// Dedup by ANTHROPIC_BASE_URL + ANTHROPIC_MODEL.
// Returns (imported count, skipped count, error).
func RunImportAll(dbPath, profilesDir string) (int, int, error) {
	db, err := database.Open(dbPath)
	if err != nil {
		return 0, 0, fmt.Errorf("cannot open cc-switch database: %w", err)
	}
	defer db.Close()

	providerDAO := dao.NewProviderDAO(db, config.AppType)
	rows, err := providerDAO.ListAll()
	if err != nil {
		return 0, 0, fmt.Errorf("cannot read providers: %w", err)
	}

	// Load existing profiles to build dedup set (base_url|model)
	existingKeys := make(map[string]bool)
	existingFiles, _ := filepath.Glob(filepath.Join(profilesDir, "*.json"))
	for _, f := range existingFiles {
		data, err := os.ReadFile(f)
		if err != nil {
			continue
		}
		var profile struct {
			Env map[string]string `json:"env"`
		}
		if json.Unmarshal(data, &profile) != nil {
			continue
		}
		key := profile.Env["ANTHROPIC_BASE_URL"] + "|" + profile.Env["ANTHROPIC_MODEL"]
		existingKeys[key] = true
	}

	os.MkdirAll(profilesDir, 0755)

	imported, skipped := 0, 0
	for _, row := range rows {
		env, model, err := dao.ParseSettingsConfig(row.SettingsConfig)
		if err != nil {
			continue
		}

		key := env["ANTHROPIC_BASE_URL"] + "|" + env["ANTHROPIC_MODEL"]
		if existingKeys[key] {
			fmt.Printf("  skipped '%s' (base_url+model already exists)\n", row.Name)
			skipped++
			continue
		}

		profile := map[string]interface{}{
			"name": row.Name,
			"env":  env,
		}
		if model != "" {
			profile["model"] = model
		}

		data, err := json.MarshalIndent(profile, "", "  ")
		if err != nil {
			continue
		}

		filename := sanitizeName(row.Name) + ".json"
		path := filepath.Join(profilesDir, filename)
		// Avoid overwriting existing file with same filename
		if _, err := os.Stat(path); err == nil {
			path = filepath.Join(profilesDir, sanitizeName(row.Name)+"_imported.json")
		}

		if err := os.WriteFile(path, data, 0644); err != nil {
			continue
		}

		existingKeys[key] = true
		fmt.Printf("  imported '%s'\n", row.Name)
		imported++
	}

	return imported, skipped, nil
}

func sanitizeName(name string) string {
	result := make([]byte, 0, len(name))
	for i := 0; i < len(name); i++ {
		c := name[i]
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' || c == '_' || c == '.' {
			result = append(result, c)
		} else {
			result = append(result, '-')
		}
	}
	if len(result) == 0 {
		return "unnamed"
	}
	return string(result)
}

// RunImportAllCLI is the CLI entry point
func RunImportAllCLI() error {
	dbPath := config.CCswitchDB
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return fmt.Errorf("cc-switch database not found at %s", dbPath)
	}
	fmt.Println("Importing providers from cc-switch...")
	imported, skipped, err := RunImportAll(dbPath, config.CCCProfiles)
	if err != nil {
		return err
	}
	fmt.Printf("✓ imported %d providers, skipped %d\n", imported, skipped)
	return nil
}

// CountCCSwitchProviders returns the number of providers in cc-switch.db
func CountCCSwitchProviders(dbPath string) int {
	db, err := database.Open(dbPath)
	if err != nil {
		return 0
	}
	defer db.Close()
	providerDAO := dao.NewProviderDAO(db, config.AppType)
	rows, err := providerDAO.ListAll()
	if err != nil {
		return 0
	}
	return len(rows)
}
