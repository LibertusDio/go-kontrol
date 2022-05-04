package gokontrol

type kontrol interface{}

type kontrolstore interface {
	GetObjectByToken(token string, serviceid string, timestamp int64) (*Object, error)
	CreateObjectWithPermission(objperm ObjectPermission) error
	CreateObject(obj *Object) error
	CreateObjectWithPolicy(obj *Object) error
	UpdateObject(obj *Object) error
	GetObjectByID(id string) (*Object, error)
	GetObjectPermissionByObject(obj *Object) (*ObjectPermission, error)
	UpdateObjectPermission(perm *ObjectPermission) error
	GetPolicyByID(id string) (*Policy, error)
	GetServiceByID(id string) (*Service, error)
	GetListServiceByID(ids []string) ([]Service, error)
	GetListPolicyByID(ids []string) ([]Policy, error)
	SetService()
	GetObjectPermission()
	SetObjectPermission()
}
