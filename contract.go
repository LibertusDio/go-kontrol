package gokontrol

import "context"

type kontrol interface{}

type kontrolstore interface {
	GetObjectByToken(c context.Context, token string, serviceid string, timestamp int64) (*Object, error)
	CreateObject(c context.Context, obj *Object) error
	UpdateObject(c context.Context, obj *Object) error
	GetObjectByID(c context.Context, id string) (*Object, error)
	GetPolicyByID(c context.Context, id string) (*Policy, error)
	GetServiceByID(c context.Context, id string) (*Service, error)
}
