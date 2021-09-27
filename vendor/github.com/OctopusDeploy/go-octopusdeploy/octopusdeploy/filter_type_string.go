// Code generated by "enumer -type=FilterType -json -output filter_type_string.go"; DO NOT EDIT.

//
package octopusdeploy

import (
	"encoding/json"
	"fmt"
)

const _FilterTypeName = "ContinuousDailyScheduleCronExpressionScheduleDailyScheduleDaysPerMonthScheduleDaysPerWeekScheduleMachineFilterOnceDailySchedule"

var _FilterTypeIndex = [...]uint8{0, 23, 45, 58, 78, 97, 110, 127}

func (i FilterType) String() string {
	if i < 0 || i >= FilterType(len(_FilterTypeIndex)-1) {
		return fmt.Sprintf("FilterType(%d)", i)
	}
	return _FilterTypeName[_FilterTypeIndex[i]:_FilterTypeIndex[i+1]]
}

var _FilterTypeValues = []FilterType{0, 1, 2, 3, 4, 5, 6}

var _FilterTypeNameToValueMap = map[string]FilterType{
	_FilterTypeName[0:23]:    0,
	_FilterTypeName[23:45]:   1,
	_FilterTypeName[45:58]:   2,
	_FilterTypeName[58:78]:   3,
	_FilterTypeName[78:97]:   4,
	_FilterTypeName[97:110]:  5,
	_FilterTypeName[110:127]: 6,
}

// FilterTypeString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func FilterTypeString(s string) (FilterType, error) {
	if val, ok := _FilterTypeNameToValueMap[s]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to FilterType values", s)
}

// FilterTypeValues returns all values of the enum
func FilterTypeValues() []FilterType {
	return _FilterTypeValues
}

// IsAFilterType returns "true" if the value is listed in the enum definition. "false" otherwise
func (i FilterType) IsAFilterType() bool {
	for _, v := range _FilterTypeValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for FilterType
func (i FilterType) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for FilterType
func (i *FilterType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("FilterType should be a string, got %s", data)
	}

	var err error
	*i, err = FilterTypeString(s)
	return err
}
