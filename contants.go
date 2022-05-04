package gokontrol

import "errors"

type commonerror struct {
	NOT_FOUND            error
	INVALID_TOKEN        error
	INVALID_SERVICE      error
	INVALID_POLICY       error
	OBJECT_NOT_FOUND     error
	PERMISSION_NOT_FOUND error
	POLICY_NOT_FOUND     error
	SERVICE_NOT_FOUND    error
	MALFORM_PERMISSION   error
}

var CommonError = commonerror{
	NOT_FOUND:            errors.New("not found"),
	OBJECT_NOT_FOUND:     errors.New("object not found"),
	PERMISSION_NOT_FOUND: errors.New("permission not found"),
	POLICY_NOT_FOUND:     errors.New("policy not found"),
	SERVICE_NOT_FOUND:    errors.New("service not found"),
	INVALID_TOKEN:        errors.New("invalid or expired token"),
	INVALID_SERVICE:      errors.New("invalid service"),
	INVALID_POLICY:       errors.New("invalid policy"),
	MALFORM_PERMISSION:   errors.New("policy permission malform"),
}

type objectstatus struct {
	INIT    string
	ENABLE  string
	DISABLE string
}

var ObjectStatus = objectstatus{
	INIT:    "",
	ENABLE:  "enable",
	DISABLE: "disable",
}

type objectpolicystatus struct {
	INIT    string
	ENABLE  string
	DISABLE string
	DEFAULT string
}

var ObjectPolicyStatus = objectpolicystatus{
	INIT:    "",
	ENABLE:  "enable",
	DISABLE: "disable",
	DEFAULT: "default", // set as default settings
}

type policypermission struct {
	ANY   int
	TRUE  int
	FALSE int
}

var PolicyPermission = policypermission{
	ANY:   0,
	TRUE:  1,
	FALSE: 2,
}
