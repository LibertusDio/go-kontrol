package gokontrol

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
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
	SecretKey:      "secret",
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

//Claims -- JWT claim use for specific customize
type Claims struct {
	Permission map[string]map[string]bool `json:"permission"`
	Token      string                     `json:"token"`
	jwt.StandardClaims
}

//ValidateToken validate the given token
func (k DefaultKontrol) ValidateToken(c context.Context, jwtToken string, externalServiceID string) (*Object, error) {
	customizeClaim := &Claims{}
	tkn, err := jwt.ParseWithClaims(jwtToken, customizeClaim, func(token *jwt.Token) (interface{}, error) {
		return []byte(k.Option.SecretKey), nil
	})
	if err != nil || jwtToken == "" || tkn == nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, err
		}
	}
	if !tkn.Valid {
		return nil, errors.New("Token is invalid ")
	}
	// verify service follow path
	reqService, err := k.store.GetServiceByExternalId(c, externalServiceID)
	if err != nil && err != CommonError.SERVICE_NOT_FOUND {
		return nil, err
	}
	//verify token
	object, err := k.store.GetObjectByToken(c, customizeClaim.Token, time.Now().Unix())
	if err != nil {
		return nil, CommonError.NOT_FOUND
	}
	// Verify permission access path by permission verified from JWT
	if object.ServiceID != reqService.ID {
		for jwtsid, _ := range customizeClaim.Permission {
			if jwtsid == reqService.ID {
				return object, nil
			}
		}
		return nil, CommonError.INVALID_SERVICE
	}

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
	extendSerivceIds, err := k.GetObjectExtendServiceIds(ctx, objID)
	if err != nil {
		return nil, CommonError.INVALID_POLICY
	}
	_, sign, jwtToken, err := k.CreateCert(obj, service.DefaultPolicy, service.EnforcePolicy, extendSerivceIds)
	if err != nil {
		return nil, err
	}
	if strings.Compare(obj.Token, sign) != 0 {
		return nil, CommonError.INVALID_TOKEN
	}

	return &ObjectPermission{
		ObjectId: objID,
		Token:    jwtToken,
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
	objectExtendServiceIds, err := k.GetObjectExtendServiceIds(ctx, obj.ID)
	// generate cert
	_, sign, jwtToken, err := k.CreateCert(obj, service.DefaultPolicy, service.EnforcePolicy, objectExtendServiceIds)
	if err != nil {
		return nil, err
	}
	obj.Token = sign
	err = k.store.UpdateObject(ctx, obj)
	if err != nil {
		return nil, err
	}
	return &ObjectPermission{
		ObjectId: obj.ID,
		Token:    jwtToken,
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
		GlobalID:    uuid.New().String(),
		ExpiryDate:  time.Now().Unix() + k.Option.DefaultTimeout,
		ApplyPolicy: nil,
	}

	_, sign, jwtToken, err := k.CreateCert(obj, service.DefaultPolicy, service.EnforcePolicy, []string{})
	if err != nil {
		return nil, err
	}
	obj.Token = sign
	err = k.store.CreateObject(ctx, obj)
	if err != nil {
		return nil, err
	}
	return &ObjectPermission{
		ObjectId: obj.ID,
		Token:    jwtToken,
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
func (k DefaultKontrol) GetObjectExtendServiceIds(ctx context.Context, objId string) ([]string, error) {

	objectServiceMess, err := k.store.GetObjectServiceMesh(ctx, objId)
	if err != nil && err != gorm.ErrRecordNotFound {
		return []string{}, err
	}
	rs := make([]string, len(objectServiceMess))
	for i := 0; i < len(objectServiceMess); i++ {
		rs[i] = objectServiceMess[i].ServiceID
	}
	return rs, nil
}

//CreateCert create final cert then sign
func (k DefaultKontrol) CreateCert(obj *Object, policy []*Policy, enforce []*Policy, extendServiceIds []string) (*CertForSign, string, string, error) {
	tempcert := &CertForSign{
		ID:         obj.ID,
		GlobalID:   obj.GlobalID,
		ExternalID: obj.ExternalID,
		ServiceID:  obj.ServiceID,
		ExpiryDate: obj.ExpiryDate,
		Attributes: obj.Attributes,
	}

	tempperm := make(map[string]map[string]bool)
	// apply extend serivce
	for _, v := range extendServiceIds {
		tempperm[v] = make(map[string]bool)
	}
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
				return nil, "", "", CommonError.MALFORM_PERMISSION
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
				return nil, "", "", CommonError.MALFORM_PERMISSION
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
				return nil, "", "", CommonError.MALFORM_PERMISSION
			}
		}
		tempperm[cp.ServiceID] = ts
	}

	tempcert.Permission = tempperm
	certstr, err := json.Marshal(tempcert)
	if err != nil {
		return nil, "", "", err
	}
	scert := append([]byte(k.Option.SecretKey), certstr...)
	hash := sha256.Sum256(scert)
	sign := base64.URLEncoding.EncodeToString(hash[:])
	claims := &Claims{
		Permission: tempperm,
		Token:      sign,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: obj.ExpiryDate,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign and get the complete encoded token as a string using the secret
	jwtString, err := token.SignedString([]byte(k.Option.SecretKey))
	if err != nil {
		return nil, "", "", err
	}
	return tempcert, sign, jwtString, nil
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
func (k DefaultKontrol) UpdatePolicy(ctx context.Context, servicekey string, policy *Policy) error {
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

	// check  policy exist
	_, err = k.store.GetPolicyByID(ctx, policy.ID)
	if err != nil {
		return err
	}

	if err := k.store.UpdatePolicy(ctx, policy); err != nil {
		return err
	}
	//Expired related object
	return k.store.ExpiredObjectsByPolicy(ctx, policy.ID)

}
