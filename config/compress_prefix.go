package config

import (
	"fmt"
	"strings"

	"github.com/gookit/goutil/strutil"
	"github.com/qiangyt/jog/util"
	"gopkg.in/yaml.v2"
)

// CompressPrefixAction ...
type CompressPrefixAction int

const (
	// CompressPrefixActionRemoveNonFirstLetter ...
	CompressPrefixActionRemoveNonFirstLetter CompressPrefixAction = iota

	// CompressPrefixActionRemove ...
	CompressPrefixActionRemove

	// CompressPrefixActionDefault ...
	CompressPrefixActionDefault = CompressPrefixActionRemoveNonFirstLetter
)

// Format ...
func (i CompressPrefixAction) String() string {
	if i == CompressPrefixActionRemoveNonFirstLetter {
		return "remove-non-first-letter"
	}
	if i == CompressPrefixActionRemove {
		return "remove"
	}

	return ""
}

// ParseCompressPrefixAction ...
func ParseCompressPrefixAction(text string) CompressPrefixAction {
	if "remove-non-first-letter" == text {
		return CompressPrefixActionRemoveNonFirstLetter
	}
	if "remove" == text {
		return CompressPrefixActionRemove
	}

	panic(fmt.Errorf("unknown CompressPrefixAction text: %v", text))
}

// CompressPrefixT ...
type CompressPrefixT struct {
	yaml.Unmarshaler
	yaml.Marshaler

	Enabled    bool
	Separators StringSet
	WhiteList  StringSet
	Action     CompressPrefixAction
}

// CompressPrefix ..
type CompressPrefix = *CompressPrefixT

// UnmarshalYAML ...
func (i CompressPrefix) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return util.UnmarshalYAML(i, unmarshal)
}

// MarshalYAML ...
func (i CompressPrefix) MarshalYAML() (interface{}, error) {
	return util.MarshalYAML(i)
}

// FromMap ...
func (i CompressPrefix) FromMap(m map[string]interface{}) error {

	enabledV := util.ExtractFromMap(m, "enabled")
	if enabledV != nil {
		i.Enabled = util.ToBool(enabledV)
	}

	separatorsV := util.ExtractFromMap(m, "separators")
	if separatorsV != nil {
		i.Separators.Parse(separatorsV)
	}

	whiteListV := util.ExtractFromMap(m, "white-list")
	if whiteListV != nil {
		i.WhiteList.Parse(whiteListV)
	}

	actionV := util.ExtractFromMap(m, "action")
	if actionV != nil {
		i.Action = ParseCompressPrefixAction(strutil.MustString(actionV))
	}

	return nil
}

// ToMap ...
func (i CompressPrefix) ToMap() map[string]interface{} {
	r := make(map[string]interface{})
	r["enabled"] = i.Enabled
	r["separators"] = i.Separators.String()
	r["white-list"] = i.WhiteList.String()
	r["action"] = i.Action.String()
	return r
}

// Reset ...
func (i CompressPrefix) Reset() {
	i.Enabled = false
	i.Separators = &StringSetT{}
	i.WhiteList = &StringSetT{}
	i.Action = CompressPrefixActionDefault
}

var _compressCache4RemoveNonFirstLetter = make(map[string]string)
var _compressCache4Remove = make(map[string]string)

func (i CompressPrefix) detectSeparator(text string) (string, []string) {
	var sepMap map[string]bool
	if i.Separators.CaseSensitive {
		sepMap = i.Separators.ValueMap
	} else {
		sepMap = i.Separators.LowercasedValueMap
	}

	for separator := range sepMap {
		separated := strings.Split(text, separator)
		if len(separated) > 1 {
			return separator, separated
		}
	}

	return "", []string{text}
}

// Compress ...
func (i CompressPrefix) Compress(text string) string {
	if i.Action == CompressPrefixActionRemoveNonFirstLetter {
		if existingOne, ok := _compressCache4RemoveNonFirstLetter[text]; ok {
			return existingOne
		}

		separator, separated := i.detectSeparator(text)

		var r string
		if len(separated) > 1 {
			indexOfLast := len(separated) - 1
			for index, item := range separated[:indexOfLast] {
				separated[index] = string([]byte(item)[0])
			}

			r = strings.Join(separated, separator)
		} else {
			r = text
		}

		_compressCache4RemoveNonFirstLetter[text] = r
		return r
	}

	if i.Action == CompressPrefixActionRemove {
		if existingOne, ok := _compressCache4Remove[text]; ok {
			return existingOne
		}

		_, separated := i.detectSeparator(text)

		var r string
		if len(separated) > 1 {
			r = separated[len(separated)-1]
		} else {
			r = text
		}
		_compressCache4Remove[text] = r
		return r
	}

	return text
}