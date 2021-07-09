// Code generated by go generate; DO NOT EDIT.
// This file was generated by robots.

package update

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

var dereferencableKindsProjectConfig = map[reflect.Kind]struct{}{
	reflect.Array: {}, reflect.Chan: {}, reflect.Map: {}, reflect.Ptr: {}, reflect.Slice: {},
}

// Checks if t is a kind that can be dereferenced to get its underlying type.
func canGetElementProjectConfig(t reflect.Kind) bool {
	_, exists := dereferencableKindsProjectConfig[t]
	return exists
}

// This decoder hook tests types for json unmarshaling capability. If implemented, it uses json unmarshal to build the
// object. Otherwise, it'll just pass on the original data.
func jsonUnmarshalerHookProjectConfig(_, to reflect.Type, data interface{}) (interface{}, error) {
	unmarshalerType := reflect.TypeOf((*json.Unmarshaler)(nil)).Elem()
	if to.Implements(unmarshalerType) || reflect.PtrTo(to).Implements(unmarshalerType) ||
		(canGetElementProjectConfig(to.Kind()) && to.Elem().Implements(unmarshalerType)) {

		raw, err := json.Marshal(data)
		if err != nil {
			fmt.Printf("Failed to marshal Data: %v. Error: %v. Skipping jsonUnmarshalHook", data, err)
			return data, nil
		}

		res := reflect.New(to).Interface()
		err = json.Unmarshal(raw, &res)
		if err != nil {
			fmt.Printf("Failed to umarshal Data: %v. Error: %v. Skipping jsonUnmarshalHook", data, err)
			return data, nil
		}

		return res, nil
	}

	return data, nil
}

func decode_ProjectConfig(input, result interface{}) error {
	config := &mapstructure.DecoderConfig{
		TagName:          "json",
		WeaklyTypedInput: true,
		Result:           result,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
			jsonUnmarshalerHookProjectConfig,
		),
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	return decoder.Decode(input)
}

func join_ProjectConfig(arr interface{}, sep string) string {
	listValue := reflect.ValueOf(arr)
	strs := make([]string, 0, listValue.Len())
	for i := 0; i < listValue.Len(); i++ {
		strs = append(strs, fmt.Sprintf("%v", listValue.Index(i)))
	}

	return strings.Join(strs, sep)
}

func testDecodeJson_ProjectConfig(t *testing.T, val, result interface{}) {
	assert.NoError(t, decode_ProjectConfig(val, result))
}

func testDecodeRaw_ProjectConfig(t *testing.T, vStringSlice, result interface{}) {
	assert.NoError(t, decode_ProjectConfig(vStringSlice, result))
}

func TestProjectConfig_GetPFlagSet(t *testing.T) {
	val := ProjectConfig{}
	cmdFlags := val.GetPFlagSet("")
	assert.True(t, cmdFlags.HasFlags())
}

func TestProjectConfig_SetFlags(t *testing.T) {
	actual := ProjectConfig{}
	cmdFlags := actual.GetPFlagSet("")
	assert.True(t, cmdFlags.HasFlags())

	t.Run("Test_activateProject", func(t *testing.T) {

		t.Run("Override", func(t *testing.T) {
			testValue := "1"

			cmdFlags.Set("activateProject", testValue)
			if vBool, err := cmdFlags.GetBool("activateProject"); err == nil {
				testDecodeJson_ProjectConfig(t, fmt.Sprintf("%v", vBool), &actual.ActivateProject)

			} else {
				assert.FailNow(t, err.Error())
			}
		})
	})
	t.Run("Test_archiveProject", func(t *testing.T) {

		t.Run("Override", func(t *testing.T) {
			testValue := "1"

			cmdFlags.Set("archiveProject", testValue)
			if vBool, err := cmdFlags.GetBool("archiveProject"); err == nil {
				testDecodeJson_ProjectConfig(t, fmt.Sprintf("%v", vBool), &actual.ArchiveProject)

			} else {
				assert.FailNow(t, err.Error())
			}
		})
	})
}
