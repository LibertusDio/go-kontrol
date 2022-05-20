package main

import (
	"context"
	"encoding/json"
	"time"

	gokontrol "github.com/LibertusDio/go-kontrol"
	"github.com/google/uuid"
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

//system storage layer
type gormStorage struct {
}

func NewGormStorage() Storage {
	return &gormStorage{}
}

//kontrol support storage
type kontrolStorage struct{}

func NewKontrolStorage() gokontrol.KontrolStore {
	return &kontrolStorage{}
}

type serviceStore struct {
	ID         string
	ServiceID  string
	Name       string
	Key        string
	Status     string
	ExpiryDate int64
}

type servicepolicymesh struct {
	ID        string
	ServiceID string
	PolicyID  string
	Type      string
}

type objectStore struct {
	ID         string
	GlobalID   string
	ExternalID string
	ServiceID  string
	Status     string
	Token      string
	ExpiryDate int64
}

type objectpolicymesh struct {
	ID       string
	ObjectID string
	PolicyID string
}

type policystore struct {
	ID         string
	Name       string
	ServiceID  string
	Permission string
	Status     string
	ApplyFrom  int64
	ApplyTo    int64
}

func (k *kontrolStorage) GetObjectByToken(c context.Context, token string, serviceid string, timestamp int64) (*gokontrol.Object, error) {
	tx := c.Value(ContextKeyTransaction).(*gorm.DB)
	var objectstore objectStore
	err := tx.WithContext(c).Table(DBTableName.TB_OBJECTS).Where("token = ? AND service_id = ? AND expiry_date >= ?", token, serviceid, timestamp).First(&objectstore).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == gorm.ErrRecordNotFound {
		return nil, gokontrol.CommonError.NOT_FOUND
	}

	var mesh []*objectpolicymesh
	err = tx.WithContext(c).Table(DBTableName.TB_OBJECT_POLICY_MESH).Where("object_id = ? ", objectstore.ID).Scan(&mesh).Error
	if err != nil {
		return nil, err
	}
	defaultpolicy := make([]*gokontrol.Policy, 0)
	for _, m := range mesh {
		policy, err := k.GetPolicyByID(c, m.PolicyID)
		if err != nil {
			return nil, err
		}
		defaultpolicy = append(defaultpolicy, policy)
	}
	return &gokontrol.Object{
		ID:          objectstore.ID,
		GlobalID:    objectstore.GlobalID,
		ExternalID:  objectstore.ExternalID,
		ServiceID:   objectstore.ServiceID,
		Status:      objectstore.Status,
		Attributes:  nil,
		Token:       objectstore.Token,
		ExpiryDate:  objectstore.ExpiryDate,
		ApplyPolicy: defaultpolicy,
	}, nil
}

func (k *kontrolStorage) CreateObject(c context.Context, obj *gokontrol.Object) error {
	tx := c.Value(ContextKeyTransaction).(*gorm.DB)
	object := objectStore{
		ID:         obj.ID,
		GlobalID:   obj.GlobalID,
		ExternalID: obj.ExternalID,
		ServiceID:  obj.ServiceID,
		Status:     obj.Status,
		Token:      obj.Token,
		ExpiryDate: obj.ExpiryDate,
	}
	err := tx.WithContext(c).Table(DBTableName.TB_OBJECTS).Create(&object).Error
	if err != nil {
		return err
	}

	// assuming policies are validated
	for _, p := range obj.ApplyPolicy {
		opm := objectpolicymesh{
			ID:       uuid.NewString(),
			ObjectID: obj.ID,
			PolicyID: p.ID,
		}
		err := tx.WithContext(c).Table(DBTableName.TB_OBJECT_POLICY_MESH).Create(&opm).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (k *kontrolStorage) UpdateObject(c context.Context, obj *gokontrol.Object) error {
	tx := c.Value(ContextKeyTransaction).(*gorm.DB)
	object := objectStore{
		ID:         obj.ID,
		GlobalID:   obj.GlobalID,
		ExternalID: obj.ExternalID,
		ServiceID:  obj.ServiceID,
		Status:     obj.Status,
		Token:      obj.Token,
		ExpiryDate: obj.ExpiryDate,
	}
	err := tx.WithContext(c).Table(DBTableName.TB_OBJECTS).Updates(&object).Where("id = ?", obj.ID).Error
	if err != nil {
		return err
	}

	// clean old policy
	err = tx.WithContext(c).Table(DBTableName.TB_OBJECT_POLICY_MESH).Delete(&objectpolicymesh{}, "object_id = ? ", obj.ID).Error
	if err != nil {
		return err
	}

	// assuming policies are validated
	for _, p := range obj.ApplyPolicy {
		opm := objectpolicymesh{
			ID:       uuid.NewString(),
			ObjectID: obj.ID,
			PolicyID: p.ID,
		}
		err = tx.WithContext(c).Table(DBTableName.TB_OBJECT_POLICY_MESH).Save(&opm).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (k *kontrolStorage) GetObjectByID(c context.Context, id string) (*gokontrol.Object, error) {
	tx := c.Value(ContextKeyTransaction).(*gorm.DB)
	var objectstore objectStore
	err := tx.WithContext(c).Table(DBTableName.TB_OBJECTS).Where("id = ? ", id).First(&objectstore).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == gorm.ErrRecordNotFound {
		return nil, gokontrol.CommonError.NOT_FOUND
	}

	var mesh []*objectpolicymesh
	err = tx.WithContext(c).Table(DBTableName.TB_OBJECT_POLICY_MESH).Where("object_id = ? ", id).Scan(&mesh).Error
	if err != nil {
		return nil, err
	}
	defaultpolicy := make([]*gokontrol.Policy, 0)
	for _, m := range mesh {
		policy, err := k.GetPolicyByID(c, m.PolicyID)
		if err != nil {
			return nil, err
		}
		defaultpolicy = append(defaultpolicy, policy)
	}
	return &gokontrol.Object{
		ID:          objectstore.ID,
		GlobalID:    objectstore.GlobalID,
		ExternalID:  objectstore.ExternalID,
		ServiceID:   objectstore.ServiceID,
		Status:      objectstore.Status,
		Attributes:  nil,
		Token:       objectstore.Token,
		ExpiryDate:  objectstore.ExpiryDate,
		ApplyPolicy: defaultpolicy,
	}, nil
}

func (k *kontrolStorage) GetObjectByExternalID(c context.Context, extid string, serviceid string) (*gokontrol.Object, error) {
	tx := c.Value(ContextKeyTransaction).(*gorm.DB)
	var objectstore objectStore
	err := tx.WithContext(c).Table(DBTableName.TB_OBJECTS).Where("external_id = ? AND service_id = ? ", extid, serviceid).First(&objectstore).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == gorm.ErrRecordNotFound {
		return nil, gokontrol.CommonError.NOT_FOUND
	}

	var mesh []*objectpolicymesh
	err = tx.WithContext(c).Table(DBTableName.TB_OBJECT_POLICY_MESH).Where("object_id = ? ", objectstore.ID).Scan(&mesh).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	defaultpolicy := make([]*gokontrol.Policy, 0)
	for _, m := range mesh {
		policy, err := k.GetPolicyByID(c, m.PolicyID)
		if err != nil {
			return nil, err
		}
		defaultpolicy = append(defaultpolicy, policy)
	}
	return &gokontrol.Object{
		ID:          objectstore.ID,
		GlobalID:    objectstore.GlobalID,
		ExternalID:  objectstore.ExternalID,
		ServiceID:   objectstore.ServiceID,
		Status:      objectstore.Status,
		Attributes:  nil,
		Token:       objectstore.Token,
		ExpiryDate:  objectstore.ExpiryDate,
		ApplyPolicy: defaultpolicy,
	}, nil
}

func (k *kontrolStorage) GetPolicyByID(c context.Context, id string) (*gokontrol.Policy, error) {
	tx := c.Value(ContextKeyTransaction).(*gorm.DB)
	var policystore policystore
	err := tx.WithContext(c).Table(DBTableName.TB_POLICIES).Where("id = ?", id).First(&policystore).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == gorm.ErrRecordNotFound {
		return nil, gokontrol.CommonError.NOT_FOUND
	}

	perm := make(map[string]int)
	err = json.Unmarshal([]byte(policystore.Permission), &perm)
	if err != nil {
		return nil, err
	}

	return &gokontrol.Policy{
		ID:         policystore.ID,
		Name:       policystore.Name,
		ServiceID:  policystore.ServiceID,
		Permission: perm,
		Status:     policystore.Status,
		ApplyFrom:  policystore.ApplyFrom,
		ApplyTo:    policystore.ApplyTo,
	}, nil
}

func (k *kontrolStorage) CreatePolicy(c context.Context, policy *gokontrol.Policy) error {
	tx := c.Value(ContextKeyTransaction).(*gorm.DB)
	// convert perm
	perm, err := json.Marshal(policy.Permission)
	if err != nil {
		return err
	}

	// save DB
	policystore := policystore{
		ID:         policy.ID,
		Name:       policy.Name,
		ServiceID:  policy.ServiceID,
		Permission: string(perm),
		Status:     policy.Status,
		ApplyFrom:  policy.ApplyFrom,
		ApplyTo:    policy.ApplyTo,
	}
	err = tx.WithContext(c).Table(DBTableName.TB_POLICIES).Create(&policystore).Error
	if err != nil {
		return err
	}

	return nil
}

func (k *kontrolStorage) GetServiceByID(c context.Context, id string) (*gokontrol.Service, error) {
	tx := c.Value(ContextKeyTransaction).(*gorm.DB)
	var servicestore serviceStore
	err := tx.WithContext(c).Table(DBTableName.TB_SERVICES).Where("id = ? ", id).First(&servicestore).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == gorm.ErrRecordNotFound {
		return nil, gokontrol.CommonError.NOT_FOUND
	}
	service := &gokontrol.Service{
		ID:         servicestore.ID,
		ServiceID:  servicestore.ServiceID,
		Name:       servicestore.Name,
		Key:        servicestore.Key,
		Status:     servicestore.Status,
		ExpiryDate: servicestore.ExpiryDate,
	}

	var defaultmesh []*servicepolicymesh
	err = tx.WithContext(c).Table(DBTableName.TB_SERVICE_POLICY_MESH).Where("service_id = ? AND `type` = ? ", id, ServicePolicyType.DEFAULT).Scan(&defaultmesh).Error
	if err != nil {
		return nil, err
	}
	defaultpolicy := make([]*gokontrol.Policy, 0)
	for _, m := range defaultmesh {
		policy, err := k.GetPolicyByID(c, m.PolicyID)
		if err != nil {
			return nil, err
		}
		defaultpolicy = append(defaultpolicy, policy)
	}
	service.DefaultPolicy = defaultpolicy

	var enforcemesh []*servicepolicymesh
	err = tx.WithContext(c).Table(DBTableName.TB_SERVICE_POLICY_MESH).Where("service_id = ? AND `type` = ? ", id, ServicePolicyType.ENFORCE).Scan(&enforcemesh).Error
	if err != nil {
		return nil, err
	}
	enforcepolicy := make([]*gokontrol.Policy, 0)
	for _, m := range enforcemesh {
		policy, err := k.GetPolicyByID(c, m.PolicyID)
		if err != nil {
			return nil, err
		}
		enforcepolicy = append(enforcepolicy, policy)
	}
	service.EnforcePolicy = enforcepolicy
	return service, nil
}
