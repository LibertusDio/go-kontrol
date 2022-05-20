package gokontrol

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"
)

//KontrolOption kontrol config options
type KontrolOption struct {
	DefaultTimeout int64
	SecretKey      string
}

//Default config for kontrol
var DefaultKontrolOption = KontrolOption{
	DefaultTimeout: 1800, // second
	SecretKey:      "",
}

//DefaultKontrol simple Kontrol
type DefaultKontrol struct {
	store  KontrolStore
	Option KontrolOption
}

//NewBasicKontrol simple Kontrol with default option, stores still have to be provided
func NewBasicKontrol(store KontrolStore) Kontrol {
	return &DefaultKontrol{store: store, Option: DefaultKontrolOption}
}

//ValidateToken validate the given token
func (k DefaultKontrol) ValidateToken(c context.Context, token string, serviceid string) (*Object, error) {
	object, err := k.store.GetObjectByToken(c, token, serviceid, time.Now().Unix())
	if err != nil && err != CommonError.NOT_FOUND {
		return nil, err
	}
	if err == CommonError.NOT_FOUND {
		return nil, CommonError.INVALID_TOKEN
	}

	return object, nil
}

//IssueCertForService issue cert for issed time, does not authen, must be authen-ed beforehand
func (k DefaultKontrol) IssueCertForService(ctx context.Context, objID string, serID string) (*ObjectPermission, error) {
	// check object
	obj, err := k.store.GetObjectByID(ctx, objID)
	if err != nil && err != CommonError.NOT_FOUND {
		return nil, err
	}
	if obj == nil || err == CommonError.NOT_FOUND {
		return nil, CommonError.OBJECT_NOT_FOUND
	}
	// check service/policy
	if strings.Compare(serID, obj.ServiceID) != 0 {
		return nil, CommonError.INVALID_SERVICE
	}
	service, err := k.store.GetServiceByID(ctx, serID)
	if err != nil && err != CommonError.NOT_FOUND {
		return nil, err
	}
	if service == nil || err == CommonError.NOT_FOUND {
		return nil, CommonError.INVALID_SERVICE
	}
	// generate cert
	cert, sign, err := k.CreateCert(obj, service.DefaultPolicy, service.EnforcePolicy)
	if err != nil {
		return nil, err
	}
	if strings.Compare(obj.Token, sign) != 0 {
		return nil, CommonError.INVALID_TOKEN
	}

	return &ObjectPermission{
		Object:     *obj,
		Permission: cert.Permission,
	}, nil
}

//IssueCertForClient issue cert for current time, does not authen, must be authen-ed beforehand
func (k DefaultKontrol) IssueCertForClient(ctx context.Context, externalID string, serID string) (*ObjectPermission, error) {
	// check object
	obj, err := k.store.GetObjectByExternalID(ctx, externalID, serID)
	if err != nil && err != CommonError.NOT_FOUND {
		return nil, err
	}
	if obj == nil || err == CommonError.NOT_FOUND {
		return nil, CommonError.OBJECT_NOT_FOUND
	}
	// check service/policy
	if strings.Compare(serID, obj.ServiceID) != 0 {
		return nil, CommonError.INVALID_SERVICE
	}
	service, err := k.store.GetServiceByID(ctx, serID)
	if err != nil && err != CommonError.NOT_FOUND {
		return nil, err
	}
	if service == nil || err == CommonError.NOT_FOUND {
		return nil, CommonError.INVALID_SERVICE
	}
	obj.ExpiryDate = time.Now().Unix() + k.Option.DefaultTimeout
	// generate cert
	cert, sign, err := k.CreateCert(obj, service.DefaultPolicy, service.EnforcePolicy)
	if err != nil {
		return nil, err
	}
	obj.Token = sign
	err = k.store.UpdateObject(ctx, obj)
	if err != nil {
		return nil, err
	}

	return &ObjectPermission{
		Object:     *obj,
		Permission: cert.Permission,
	}, nil
}

//AddSimpleObjectWithDefaultPolicy add object with default service schema
func (k DefaultKontrol) AddSimpleObjectWithDefaultPolicy(ctx context.Context, externalid string, serviceid string, servicekey string) (*ObjectPermission, error) {
	// check service/policy
	service, err := k.store.GetServiceByID(ctx, serviceid)
	if err != nil && err != CommonError.NOT_FOUND {
		return nil, err
	}
	if service == nil || err == CommonError.NOT_FOUND {
		return nil, CommonError.INVALID_SERVICE
	}

	// check service key
	scert := append([]byte(k.Option.SecretKey), []byte(servicekey)...)
	hash := sha256.Sum256(scert)
	sign := base64.URLEncoding.EncodeToString(hash[:])
	if strings.Compare(sign, service.Key) != 0 {
		return nil, CommonError.INVALID_TOKEN
	}

	testobj, err := k.store.GetObjectByExternalID(ctx, externalid, serviceid)
	if err != nil && err != CommonError.NOT_FOUND {
		return nil, err
	}
	if testobj != nil || err != CommonError.NOT_FOUND {
		return nil, CommonError.INVALID_OBJECT
	}

	obj := &Object{
		ID:          uuid.New().String(),
		ExternalID:  externalid,
		ServiceID:   serviceid,
		Status:      ObjectStatus.ENABLE,
		Attributes:  nil,
		Token:       "",
		ExpiryDate:  time.Now().Unix() + k.Option.DefaultTimeout,
		ApplyPolicy: nil,
	}

	cert, sign, err := k.CreateCert(obj, service.DefaultPolicy, service.EnforcePolicy)
	if err != nil {
		return nil, err
	}
	obj.Token = sign
	err = k.store.CreateObject(ctx, obj)
	if err != nil {
		return nil, err
	}

	return &ObjectPermission{
		Object:     *obj,
		Permission: cert.Permission,
	}, nil
}

//UpdateObject update Object info
func (k DefaultKontrol) UpdateObject(ctx context.Context, obj *Object, servicekey string) error {
	// check service
	service, err := k.store.GetServiceByID(ctx, obj.ServiceID)
	if err != nil && err != CommonError.NOT_FOUND {
		return err
	}
	if service == nil || err == CommonError.NOT_FOUND {
		return CommonError.INVALID_SERVICE
	}

	// check service key
	scert := append([]byte(k.Option.SecretKey), []byte(servicekey)...)
	hash := sha256.Sum256(scert)
	sign := base64.URLEncoding.EncodeToString(hash[:])
	if strings.Compare(sign, service.Key) != 0 {
		return CommonError.INVALID_TOKEN
	}

	// check duplicate
	old, err := k.store.GetObjectByID(ctx, obj.ID)
	if err != nil && err != CommonError.NOT_FOUND {
		return err
	}
	if old == nil || err == CommonError.NOT_FOUND {
		return CommonError.OBJECT_NOT_FOUND
	}

	return k.store.UpdateObject(ctx, obj)
}

//CreateCert create final cert then sign
func (k DefaultKontrol) CreateCert(obj *Object, policy []*Policy, enforce []*Policy) (*CertForSign, string, error) {
	tempcert := &CertForSign{
		ID:         obj.ID,
		GlobalID:   obj.GlobalID,
		ExternalID: obj.ExternalID,
		ServiceID:  obj.ServiceID,
		ExpiryDate: obj.ExpiryDate,
		Attributes: obj.Attributes,
	}

	tempperm := make(map[string]map[string]bool)
	// apply default policies
	for _, dp := range policy {
		ts, exist := tempperm[dp.ServiceID]
		if !exist {
			ts = make(map[string]bool)
		}
		for k, v := range dp.Permission {
			switch v {
			case PolicyPermission.TRUE:
				ts[k] = true
			case PolicyPermission.FALSE:
			case PolicyPermission.ANY:
			default:
				return nil, "", CommonError.MALFORM_PERMISSION
			}
		}
		tempperm[dp.ServiceID] = ts
	}
	// apply custom policies
	for _, cp := range obj.ApplyPolicy {
		ts, exist := tempperm[cp.ServiceID]
		if !exist {
			ts = make(map[string]bool)
		}
		for k, v := range cp.Permission {
			switch v {
			case PolicyPermission.TRUE:
				ts[k] = true
			case PolicyPermission.FALSE:
				delete(ts, k)
			case PolicyPermission.ANY:
			default:
				return nil, "", CommonError.MALFORM_PERMISSION
			}
		}
		tempperm[cp.ServiceID] = ts
	}

	// apply enforce policy
	for _, cp := range obj.ApplyPolicy {
		ts, exist := tempperm[cp.ServiceID]
		if !exist {
			ts = make(map[string]bool)
		}
		for k, v := range cp.Permission {
			switch v {
			case PolicyPermission.TRUE:
			case PolicyPermission.FALSE:
				delete(ts, k)
			case PolicyPermission.ANY:
			default:
				return nil, "", CommonError.MALFORM_PERMISSION
			}
		}
		tempperm[cp.ServiceID] = ts
	}

	tempcert.Permission = tempperm
	certstr, err := json.Marshal(tempcert)
	if err != nil {
		return nil, "", err
	}
	scert := append([]byte(k.Option.SecretKey), certstr...)
	hash := sha256.Sum256(scert)
	sign := base64.URLEncoding.EncodeToString(hash[:])
	return tempcert, sign, nil
}

//CreatePolicy create a policy
func (k DefaultKontrol) CreatePolicy(ctx context.Context, servicekey string, policy *Policy) error {
	// check service
	service, err := k.store.GetServiceByID(ctx, policy.ServiceID)
	if err != nil && err != CommonError.NOT_FOUND {
		return err
	}
	if service == nil || err == CommonError.NOT_FOUND {
		return CommonError.INVALID_SERVICE
	}

	// check service key
	scert := append([]byte(k.Option.SecretKey), []byte(servicekey)...)
	hash := sha256.Sum256(scert)
	sign := base64.URLEncoding.EncodeToString(hash[:])

	if strings.Compare(sign, service.Key) != 0 {
		return CommonError.INVALID_TOKEN
	}

	// check duplicate policy
	testpolicy, err := k.store.GetPolicyByID(ctx, policy.ID)
	if err != nil && err != CommonError.NOT_FOUND {
		return err
	}
	if testpolicy != nil || err != CommonError.NOT_FOUND {
		return CommonError.INVALID_POLICY
	}

	return k.store.CreatePolicy(ctx, policy)
}
