package config

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type ConfigManager struct {
	szConfigFilePath    string
	config 				Config
	mu                  sync.RWMutex
}

func NewConfigManager(szConfigFilePath string) *ConfigManager {
	cfgMgr := &ConfigManager{
		szConfigFilePath: szConfigFilePath,
	}

	if err := cfgMgr.loadConfig(); err != nil {
		fmt.Printf("Warning: Failed to load config: %v\n", err)
	}
	return cfgMgr
}

func (cfgMgr *ConfigManager) loadConfig() error {
	data, err := os.ReadFile(cfgMgr.szConfigFilePath)
	if err != nil {
		fmt.Printf("Failed to load config: %w\n", err)
		return nil
	}
	return json.Unmarshal(data, &cfgMgr.config)
}

func (cfgMgr *ConfigManager) GetActiveProfile() Profile {
	cfgMgr.mu.RLock()
	defer cfgMgr.mu.RUnlock()
	return cfgMgr.config.Profiles[cfgMgr.config.SzActiveProfile]
}

func (cfgMgr *ConfigManager) SwitchProfile(szName string) error {
	cfgMgr.mu.Lock()
	defer cfgMgr.mu.Unlock()
	if szName == "" {
		return fmt.Errorf("profile name cannot be empty")
	}

	if _, exists := cfgMgr.config.Profiles[szName]; !exists {
		return fmt.Errorf("profile '%s' not found", szName)
	}

	cfgMgr.config.SzActiveProfile = szName
	return cfgMgr.saveConfig()
}

func (cfgMgr *ConfigManager) saveConfig() error {
	data, err := json.MarshalIndent(cfgMgr.config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	return os.WriteFile(cfgMgr.szConfigFilePath, data, 0644)
}

func (cfgMgr *ConfigManager) ListProfile() []Profile {
	cfgMgr.mu.RLock()
	defer cfgMgr.mu.RUnlock()

	profiles := make([]Profile, 0, len(cfgMgr.config.Profiles))
	for _, profile := range cfgMgr.config.Profiles {
		profiles = append(profiles, profile)
	}

	return profiles
}

func (cfgMgr *ConfigManager) GetProfile(szName string) (Profile, error) {
	cfgMgr.mu.RLock()
	defer cfgMgr.mu.RUnlock()

	profile, exists := cfgMgr.config.Profiles[szName]
	if !exists {
		return Profile{}, fmt.Errorf("profile '%s' not found", szName)
	}

	return profile, nil
}
