package main

import "strings"

type EndpointValidator struct {
	parser *ContractParser
}

func NewEndpointValidator(parser *ContractParser) *EndpointValidator {
	return &EndpointValidator{parser: parser}
}

func (v *EndpointValidator) Validate(endpoints []endpoint) {
	if v.parser == nil {
		panic("public API validator requires a contract parser")
	}
	operations := map[string]bool{}
	routes := map[string]bool{}
	for _, route := range endpoints {
		key := route.Method + " " + route.Path
		if route.Path != strings.ToLower(route.Path) || strings.Contains(route.Path, "_") {
			panic("route path must use lowercase kebab-case: " + key)
		}
		for _, parameter := range pathNames(route.Path) {
			if kebabCase(parameter) != parameter {
				panic("path parameter must use kebab-case: " + key)
			}
		}
		if routes[key] {
			panic("duplicate route: " + key)
		}
		routes[key] = true
		if operations[route.OperationID] {
			panic("duplicate operationId: " + route.OperationID)
		}
		operations[route.OperationID] = true
		contractless := route.Tag == "graphql" || route.Tag == "static"
		if !contractless && route.RequestType == "" {
			panic("missing request DTO: " + key)
		}
		if !contractless && route.ResponseType == "" {
			panic("missing response DTO: " + key)
		}
		if route.RequestType != "" {
			if _, ok := v.parser.FindType(route.RequestType); !ok {
				panic("request DTO not found: " + route.RequestType)
			}
		}
		if route.ResponseType != "" {
			if _, ok := v.parser.FindType(route.ResponseType); !ok {
				panic("response DTO not found: " + route.ResponseType)
			}
		}
	}
}

func authenticated(operationID string) bool {
	anonymous := map[string]bool{
		"register": true, "registerViaGoogle": true, "login": true,
		"loginViaGoogle": true, "sendAuthCode": true, "forgetPassword": true,
	}
	return !anonymous[operationID]
}

func csrfRequired(operationID string) bool {
	switch operationID {
	case "validateEmail", "resetEmail", "resetMe", "deleteMe":
		return true
	default:
		return false
	}
}
