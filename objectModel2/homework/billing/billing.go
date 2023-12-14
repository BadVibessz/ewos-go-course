package billing

import (
	"encoding/json"
	"errors"
	"math"
	"strconv"
	"time"

	"github.com/ew0s/ewos-to-go-hw/objectModel2/homework/utils"
)

type OperationType int

const (
	Unknown = OperationType(iota)
	Outcome
	Income
)

var (
	errNoSuchVar         = errors.New("no such var")
	errInvalidType       = errors.New("invalid type")
	errInvalidValue      = errors.New("invalid type for value provided")
	errNotString         = errors.New("provided value has type different from string")
	errInvalidTimeFormat = errors.New("provided value has invalid time format")
)

type Operation struct {
	ID        any
	Type      OperationType
	Value     int
	CreatedAt time.Time
}

func (op *Operation) Validate() bool {
	if op.Type == Unknown || op.Value == 0 || op.CreatedAt.IsZero() {
		return false
	}

	return true
}

type Billings map[string][]Operation

func lookupVar(key string, data ...map[string]any) (any, error) {
	for _, v := range data {
		_, exist := v[key]
		if exist {
			return v[key], nil
		}
	}

	return nil, errNoSuchVar
}

func getOperation(data map[string]any) map[string]any {
	if opData, exist := data["operation"]; exist {
		if op, ok := opData.(map[string]any); ok {
			return op
		}
	}

	return nil
}

func (op *Operation) trySetId(id any) error {
	switch v := id.(type) {
	case int:
		*op = Operation{ID: v}

	case float64:
		*op = Operation{ID: int(v)}

	case string:
		*op = Operation{ID: v}

	default:
		return errInvalidType // invalid type
	}

	return nil
}

func (op *Operation) setType(typ string) {
	switch typ {
	case "income", "+":
		op.Type = Income

	case "outcome", "-":
		op.Type = Outcome

	default:
		op.Type = Unknown
	}
}

func (op *Operation) trySetValue(val any) error {
	switch v := val.(type) {
	case int:
		op.Value = v

	case float64:
		if math.Mod(v, 1.0) == 0 {
			op.Value = int(v)
		}

	case string:
		i, convErr := strconv.Atoi(v)
		if convErr == nil {
			op.Value = i
		}

	default:
		return errInvalidValue
	}

	return nil
}

func (op *Operation) trySetCreatedAt(createdAt any) error {
	s, ok := createdAt.(string)
	if !ok {
		return errNotString
	}

	t, parseErr := time.Parse(time.RFC3339, s)
	if parseErr != nil {
		return errInvalidTimeFormat
	}

	op.CreatedAt = t

	return nil
}

func ParseJson(j string) (Billings, error) {
	var results []map[string]any

	err := json.Unmarshal([]byte(j), &results)
	if err != nil {
		return nil, err
	}

	m := make(Billings)

	for _, result := range results {
		c, ok := result["company"].(string)
		if !ok {
			continue // no company provided
		}

		comp := c
		operation := getOperation(result)

		var op Operation

		id, idErr := lookupVar("id", result, operation)
		if idErr != nil {
			continue // ID does not exist
		}

		err = op.trySetId(id)
		if err != nil {
			continue // invalid ID
		}

		typ, typErr := lookupVar("type", result, operation)
		if typErr == nil {
			s, ok := typ.(string)
			if ok {
				op.setType(s)
			}
		}

		val, valErr := lookupVar("value", result, operation)
		if valErr == nil {
			err = op.trySetValue(val)
			if err != nil {
				op.Value = 0 // anti linter "unhandled err and empty if block"
			}
		}

		createdAt, createdErr := lookupVar("created_at", result, operation)
		if createdErr != nil {
			continue // createdAd does not exist
		}

		err = op.trySetCreatedAt(createdAt)
		if err != nil {
			continue // operation with invalid time won't be handled
		}

		ops, ok := m[comp]
		if !ok {
			m[comp] = append(make([]Operation, 0), op)
		} else {
			m[comp] = append(ops, op)
		}
	}

	return m, nil
}

type CompanyInfo struct {
	Company              string `json:"company"`
	ValidOperationsCount int    `json:"valid_operations_count"`
	Balance              int    `json:"balance"`
	InvalidOperations    []any  `json:"invalid_operations,omitempty"`
}

func CalculateBalances(b Billings) []CompanyInfo {
	infos := make([]CompanyInfo, 0, len(b))

	for c, ops := range b {
		validOps := utils.Filter(ops, func(op Operation) bool { return op.Validate() })
		invalidOps := utils.Filter(ops, func(op Operation) bool { return !op.Validate() })

		balance := 0

		for _, op := range validOps {
			if op.Type == Income {
				balance += op.Value
			} else if op.Type == Outcome {
				balance -= op.Value
			}
		}

		info := CompanyInfo{
			Company:              c,
			ValidOperationsCount: len(validOps),
			Balance:              balance,
			InvalidOperations:    utils.Map(invalidOps, func(op Operation) any { return op.ID }),
		}

		infos = append(infos, info)
	}

	return infos
}
