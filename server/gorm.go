package main

import (
	"context"
	"time"

	gokontrol "github.com/LibertusDio/go-kontrol"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type gormSession struct {
	cfg     *MySQL
	session *gorm.DB
}

func ConnectMySQL(cfg *MySQL) (Database, error) {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       cfg.ConnectionString(), // data source name
		DontSupportRenameIndex:    true,                   // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,                   // `change` when rename column, rename column not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false,                  // auto configure based on currently MySQL version

	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(cfg.ConnectionIdleMax)
	sqlDB.SetMaxOpenConns(cfg.ConnectionMax)
	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(cfg.ConnectionTime))
	sqlDB.SetConnMaxIdleTime(time.Second * time.Duration(cfg.ConnectionIdleTime))

	if cfg.Log {
		db = db.Debug()
	}

	return &gormSession{
		cfg:     cfg,
		session: db,
	}, nil
}

func (db *gormSession) Session() (interface{}, error) {
	return db.session, nil
}

func (db *gormSession) Transaction() (interface{}, error) {
	return db.session.Begin(), nil
}

type gormStorage struct {
}

func NewGormStorage() Storage {
	return &gormStorage{}
}

type kontrolStorage struct{}

func NewKontrolStorage() gokontrol.KontrolStore {
	return &kontrolStorage{}
}
func (k *kontrolStorage) GetObjectByToken(c context.Context, token string, serviceid string, timestamp int64) (*gokontrol.Object, error)
func (k *kontrolStorage) CreateObject(c context.Context, obj *gokontrol.Object) error
func (k *kontrolStorage) UpdateObject(c context.Context, obj *gokontrol.Object) error
func (k *kontrolStorage) GetObjectByID(c context.Context, id string) (*gokontrol.Object, error)
func (k *kontrolStorage) GetPolicyByID(c context.Context, id string) (*gokontrol.Policy, error)
func (k *kontrolStorage) GetServiceByID(c context.Context, id string) (*gokontrol.Service, error)
