package database

import (
	"fmt"

	"github.com/joaopandolfi/core"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// GormMemoryBackend implements core.MemoryBackend
// Uses GORM connection to retrieve memory from Postgres database
type GormMemoryBackend struct {
	userID string
	db     *gorm.DB
}

// NewGormMemoryBackend returns a new GormMemoryBackend
func NewGormMemoryBackend(db *gorm.DB, userID string) *GormMemoryBackend {
	return &GormMemoryBackend{
		userID: userID,
		db:     db,
	}
}

// Migrate database
// add the infraestruct tables and relations
func (a *GormMemoryBackend) Migrate() error {
	return a.db.AutoMigrate(
		&core.ToolCall{},
		&core.ToolResult{},
		&core.Metadata{},
		&core.Image{},
		&core.Message{},
	)
}

// Add adds messages to the GormMemoryBackend using "append"
func (a *GormMemoryBackend) Add(m ...*core.Message) error {
	tx := a.db.Create(m)

	if tx.Error != nil {
		return fmt.Errorf("saving data on bd %w", tx.Error)
	}
	return nil
}

// GetMaxN returns the last N number of messages
func (a *GormMemoryBackend) GetMaxN(n int) ([]*core.Message, error) {
	var messages []*core.Message
	tx := a.db.Preload(clause.Associations).Where("user_id = ?", a.userID)

	if n > 0 {
		tx = tx.Limit(n)
	}

	tx = tx.Find(&messages)

	if tx.Error != nil {
		return nil, fmt.Errorf("getting last messages on database: %w", tx.Error)
	}

	return messages, nil
}

// Dump returns the whole GormMemoryBackend array
func (a *GormMemoryBackend) Dump() ([]*core.Message, error) {
	return a.GetMaxN(0)
}

// Prune resets the array in the GormMemoryBackend
func (a *GormMemoryBackend) Prune() {
	var m *core.Message
	tx := a.db.Delete(&m, "user_id = ?", a.userID)
	if tx.Error != nil {
		fmt.Println("[!][GormMemoryBackend] Error on prunning database: ", tx.Error)
	}
}
