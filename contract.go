package gokontrol

import "context"

type Kontrol interface {
	ValidateToken(c context.Context, token string, serviceid string) (*Object, error)                                                        // validate if token existed, for tighter check, use IssueCertForService
	IssueCertForService(ctx context.Context, objID string, serID string) (*ObjectPermission, error)                                          // get client cert for service to store
	AddSimpleObjectWithDefaultPolicy(ctx context.Context, externalid string, serviceid string, servicekey string) (*ObjectPermission, error) //service create new object
	UpdateObject(ctx context.Context, obj *Object, servicekey string) error                                                                  //service update object
	CreateCert(obj *Object, policy []*Policy, enforce []*Policy) (*CertForSign, string, string, error)                                       // internal use, centralise function to issue permission
	CreatePolicy(ctx context.Context, servicekey string, policy *Policy) error                                                               // service create policy
	IssueCertForClient(ctx context.Context, externalID string, serID string) (*ObjectPermission, error)                                      // issue cert for client when login success
}

type KontrolStore interface {
	GetObjectByToken(c context.Context, token string, serviceid string, timestamp int64) (*Object, error)
	CreateObject(c context.Context, obj *Object) error
	UpdateObject(c context.Context, obj *Object) error
	GetObjectByID(c context.Context, id string) (*Object, error)
	GetObjectByExternalID(c context.Context, extid string, serviceid string) (*Object, error)
	GetPolicyByID(c context.Context, id string) (*Policy, error)
	CreatePolicy(c context.Context, policy *Policy) error
	GetServiceByID(c context.Context, id string) (*Service, error)
}
