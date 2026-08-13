package main

import (
	"go/ast"
	"reflect"
	"strconv"
	"strings"
)

func objectSchema() map[string]any {
	return map[string]any{"type": "object", "properties": map[string]any{}}
}

func schemaFor(expr ast.Expr, request bool, stack map[string]bool) map[string]any {
	switch value := expr.(type) {
	case *ast.StarExpr:
		schema := schemaFor(value.X, request, stack)
		if schemaType, ok := schema["type"].(string); ok {
			schema["type"] = []string{schemaType, "null"}
		}
		return schema
	case *ast.ArrayType:
		return map[string]any{"type": "array", "items": schemaFor(value.Elt, request, stack)}
	case *ast.MapType:
		return map[string]any{"type": "object", "additionalProperties": schemaFor(value.Value, request, stack)}
	case *ast.InterfaceType:
		return map[string]any{}
	case *ast.SelectorExpr:
		packageName := ""
		if identifier, ok := value.X.(*ast.Ident); ok {
			packageName = identifier.Name
		}
		if packageName == "time" && value.Sel.Name == "Time" {
			return map[string]any{"type": "string", "format": "date-time"}
		}
		if packageName == "uuid" && value.Sel.Name == "UUID" {
			return map[string]any{"type": "string", "format": "uuid"}
		}
		if packageName == "json" && value.Sel.Name == "RawMessage" {
			return map[string]any{}
		}
		if value.Sel.Name == "JSON" {
			return map[string]any{"type": "object", "additionalProperties": true}
		}
		return schemaFor(&ast.Ident{Name: value.Sel.Name}, request, stack)
	case *ast.Ident:
		switch value.Name {
		case "string":
			return map[string]any{"type": "string"}
		case "bool":
			return map[string]any{"type": "boolean"}
		case "int", "int8", "int16", "int32", "uint", "uint8", "uint16", "uint32", "byte":
			return map[string]any{"type": "integer", "format": "int32"}
		case "int64", "uint64":
			return map[string]any{"type": "integer", "format": "int64"}
		case "float32":
			return map[string]any{"type": "number", "format": "float"}
		case "float64":
			return map[string]any{"type": "number", "format": "double"}
		case "any", "error":
			return map[string]any{}
		}
		if stack[value.Name] {
			return map[string]any{"type": "object"}
		}
		definition, ok := findType(value.Name)
		if !ok {
			return map[string]any{"type": "string", "x-go-type": value.Name}
		}
		next := cloneBoolMap(stack)
		next[value.Name] = true
		schema := schemaFor(definition.expr, request, next)
		if values := enumValues[value.Name]; len(values) > 0 {
			schema["enum"] = values
		}
		return schema
	case *ast.StructType:
		properties := map[string]any{}
		required := []string{}
		for _, field := range value.Fields.List {
			if len(field.Names) == 0 {
				embedded := schemaFor(field.Type, request, stack)
				if embeddedProperties, ok := embedded["properties"].(map[string]any); ok {
					for name, property := range embeddedProperties {
						properties[name] = property
					}
				}
				continue
			}
			tag := ""
			if field.Tag != nil {
				tag, _ = strconv.Unquote(field.Tag.Value)
			}
			jsonTag := reflect.StructTag(tag).Get("json")
			name := strings.Split(jsonTag, ",")[0]
			if name == "-" {
				continue
			}
			if name == "" {
				name = lowerFirst(field.Names[0].Name)
			}
			property := schemaFor(field.Type, request, stack)
			validation := reflect.StructTag(tag).Get("validate")
			applyValidation(property, validation)
			properties[name] = property
			isPointer := false
			_, isPointer = field.Type.(*ast.StarExpr)
			if request {
				if containsValidation(validation, "required") {
					required = append(required, name)
				}
			} else if !strings.Contains(jsonTag, "omitempty") && !isPointer {
				required = append(required, name)
			}
		}
		schema := map[string]any{"type": "object", "properties": properties}
		if len(required) > 0 {
			schema["required"] = required
		}
		return schema
	case *ast.IndexExpr:
		return map[string]any{"type": "array", "items": schemaFor(value.Index, request, stack)}
	case *ast.IndexListExpr:
		return objectSchema()
	default:
		return map[string]any{}
	}
}

func cloneBoolMap(source map[string]bool) map[string]bool {
	result := make(map[string]bool, len(source))
	for key, value := range source {
		result[key] = value
	}
	return result
}

func containsValidation(validation, wanted string) bool {
	for _, item := range strings.Split(validation, ",") {
		if item == wanted {
			return true
		}
	}
	return false
}

func applyValidation(schema map[string]any, validation string) {
	for _, rule := range strings.Split(validation, ",") {
		parts := strings.SplitN(rule, "=", 2)
		switch parts[0] {
		case "email":
			schema["format"] = "email"
		case "alphaandnum":
			schema["pattern"] = `^[A-Za-z0-9]+$`
		case "isnumberstring":
			schema["pattern"] = `^[0-9]+$`
		case "isimageurl", "isurl":
			schema["format"] = "uri"
			schema["maxLength"] = 2048
		case "ishexcodecolor":
			schema["pattern"] = `^#[0-9A-Fa-f]{6}$`
		case "isuseragent":
			schema["minLength"] = 3
			schema["maxLength"] = 2048
		case "isaccesscontrolpermission":
			schema["enum"] = []string{"Read", "Write", "Admin", "Owner"}
		case "isaccount":
			schema["description"] = "Email address, or an account identifier containing at least one letter and one digit."
		case "isstrongpassword":
			schema["pattern"] = `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[^\w\s]).+$`
		case "notfuture":
			schema["description"] = "Must not be later than the current server time."
		case "oneof":
			if len(parts) == 2 {
				schema["enum"] = strings.Fields(parts[1])
			}
		case "len":
			if len(parts) == 2 && schema["type"] == "string" {
				length, _ := strconv.Atoi(parts[1])
				schema["minLength"], schema["maxLength"] = length, length
			}
		case "min":
			if len(parts) == 2 {
				value, _ := strconv.Atoi(parts[1])
				switch schema["type"] {
				case "string":
					schema["minLength"] = value
				case "array":
					schema["minItems"] = value
				case "integer", "number":
					schema["minimum"] = value
				}
			}
		case "max":
			if len(parts) == 2 {
				value, _ := strconv.Atoi(parts[1])
				switch schema["type"] {
				case "string":
					schema["maxLength"] = value
				case "array":
					schema["maxItems"] = value
				case "integer", "number":
					schema["maximum"] = value
				}
			}
		}
	}
}

func properties(schema map[string]any) map[string]any {
	result, _ := schema["properties"].(map[string]any)
	return result
}

func schemaExample(schema map[string]any, fieldName string) any {
	if enum, ok := schema["enum"].([]string); ok && len(enum) > 0 {
		return enum[0]
	}
	if format := schema["format"]; format == "uuid" {
		return "00000000-0000-4000-8000-000000000001"
	} else if format == "date-time" {
		return "2026-01-01T00:00:00Z"
	} else if format == "email" {
		return "developer@example.com"
	}
	schemaType := schema["type"]
	if types, ok := schemaType.([]string); ok {
		for _, candidate := range types {
			if candidate != "null" {
				schemaType = candidate
				break
			}
		}
	}
	switch schemaType {
	case "object":
		result := map[string]any{}
		for name, raw := range properties(schema) {
			if child, ok := raw.(map[string]any); ok {
				result[name] = schemaExample(child, name)
			}
		}
		return result
	case "array":
		if item, ok := schema["items"].(map[string]any); ok {
			return []any{schemaExample(item, fieldName)}
		}
		return []any{}
	case "integer":
		return 1
	case "number":
		return 1.0
	case "boolean":
		return true
	default:
		name := strings.ToLower(fieldName)
		switch {
		case strings.Contains(name, "password"):
			return "Example-Password-123!"
		case strings.Contains(name, "authcode"):
			return "123456"
		case strings.HasSuffix(name, "id"):
			return "00000000-0000-4000-8000-000000000001"
		case strings.Contains(name, "url") || strings.Contains(name, "endpoint"):
			return "https://example.com"
		default:
			return "example"
		}
	}
}
