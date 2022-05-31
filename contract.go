package gokontrol

import "context"

type Kontrol interface {
	ValidateToken(c context.Context, token string, reqPath string, reqMethod string) (*Object, error)                                        // validate if token existed, for tighter check, use IssueCertForService
	IssueCertForService(ctx context.Context, objID string, externalid string) (*ObjectPermission, error)                                     // get client cert for service to store
	AddSimpleObjectWithDefaultPolicy(ctx context.Context, externalid string, serviceid string, servicekey string) (*ObjectPermission, error) //service create new object
	UpdateObject(ctx context.Context, obj *Object, servicekey string) error                                                                  //service update object
	CreateCert(obj *Object, policy []*Policy, enforce []*Policy, objectExtendServiceIds []string) (*CertForSign, string, string, error)      // internal use, centralise function to issue permission
	CreatePolicy(ctx context.Context, servicekey string, policy *Policy) error
	UpdatePolicy(ctx context.Context, servicekey string, policy *Policy) error
	IssueCertForClient(ctx context.Context, externalID string, serID string) (*ObjectPermission, error) // issue cert for client when login success
	GetObjectExtendServiceIds(ctx context.Context, objId string) ([]string, error)                      // GET LIST EXTEND SERVICE THAT OBJECT CAN ACCESS
}

type KontrolStore interface {
	GetObjectByToken(c context.Context, token string, timestamp int64) (*Object, error)
	CreateObject(c context.Context, obj *Object) error
	UpdateObject(c context.Context, obj *Object) error
	GetObjectByID(c context.Context, id string) (*Object, error)
	GetObjectByExternalID(c context.Context, extid string, serviceid string) (*Object, error)
	GetPolicyByID(c context.Context, id string) (*Policy, error)
	CreatePolicy(c context.Context, policy *Policy) error
	UpdatePolicy(c context.Context, policy *Policy) error
	ExpiredObjectsByPolicy(c context.Context, policyId string) error
	GetServiceByID(c context.Context, id string) (*Service, error)
	GetServiceByExternalId(c context.Context, externalId string) (*Service, error)
	GetObjectServiceMesh(c context.Context, objectId string) ([]*ObjectServiceMess, error)
}
