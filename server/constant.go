package main

import "errors"

const ContextKeyTransaction string = "Tx"

type servicepolicytype struct {
	INIT    string
	DEFAULT string
	ENFORCE string
}

var ServicePolicyType = servicepolicytype{
	INIT:    "",
	DEFAULT: "default",
	ENFORCE: "enforce",
}

type dbtablename struct {
	TB_SERVICES            string
	TB_SERVICE_POLICY_MESH string
	TB_OBJECTS             string
	TB_OBJECT_POLICY_MESH  string
	TB_POLICIES            string
}

var DBTableName = dbtablename{
	TB_SERVICES:            "services",
	TB_SERVICE_POLICY_MESH: "service_policy_mesh",
	TB_OBJECTS:             "objects",
	TB_OBJECT_POLICY_MESH:  "object_policy_mesh",
	TB_POLICIES:            "policies",
}

type commonerror struct {
	INVALID_PARAM error
	FORBIDDEN     error
}

var CommonError = commonerror{
	INVALID_PARAM: errors.New("invalid params"),
	FORBIDDEN:     errors.New("forbidden"),
}
