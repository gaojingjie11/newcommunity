package repository

import (
	"regexp"
	"strings"
)

var standardBluePlatePattern = regexp.MustCompile(`^[京津沪渝冀豫云辽黑湘皖鲁新苏浙赣鄂桂甘晋蒙陕吉闽贵粤青藏川宁琼][A-Z][A-Z0-9]{5}$`)

func NormalizeCarPlate(carPlate string) string {
	return strings.ToUpper(strings.TrimSpace(carPlate))
}

func IsValidStandardBluePlate(carPlate string) bool {
	return standardBluePlatePattern.MatchString(NormalizeCarPlate(carPlate))
}
