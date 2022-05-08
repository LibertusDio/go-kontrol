package gokontrol

import "context"

type Kontrol interface {
	ValidateToken(c context.Context, token string, serviceid string) (*Object, error)
	IssueCertForService(ctx context.Context, objID string, serID string) (*ObjectPermission, error)
	AddSimpleObjectWithDefaultPolicy(ctx context.Context, externalid string, serviceid string) (*ObjectPermission, error)
	UpdateObject(ctx context.Context, obj *Object) error
	CreateCert(obj *Object, policy []Policy, enforce []Policy) (*CertForSign, string, error)
}

type KontrolStore interface {
	GetObjectByToken(c context.Context, token string, serviceid string, timestamp int64) (*Object, error)
	CreateObject(c context.Context, obj *Object) error
	UpdateObject(c context.Context, obj *Object) error
	GetObjectByID(c context.Context, id string) (*Object, error)
	GetPolicyByID(c context.Context, id string) (*Policy, error)
	GetServiceByID(c context.Context, id string) (*Service, error)
}
