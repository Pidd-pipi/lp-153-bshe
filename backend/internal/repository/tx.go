package repository

import "gorm.io/gorm"

// TxManager 事务执行器，由 service 层编排多步写操作。
type TxManager interface {
	Transaction(fn func(tx *gorm.DB) error) error
}

type gormTxManager struct {
	db *gorm.DB
}

// NewTxManager 构造事务执行器。
func NewTxManager(db *gorm.DB) TxManager {
	return &gormTxManager{db: db}
}

func (m *gormTxManager) Transaction(fn func(tx *gorm.DB) error) error {
	return m.db.Transaction(fn)
}
