const standardBluePlatePattern = /^[京津沪渝冀豫云辽黑湘皖鲁新苏浙赣鄂桂甘晋蒙陕吉闽贵粤青藏川宁琼][A-Z][A-Z0-9]{5}$/

export function normalizeCarPlate(value = '') {
  return String(value).trim().toUpperCase()
}

export function isValidStandardBluePlate(value = '') {
  return standardBluePlatePattern.test(normalizeCarPlate(value))
}

export function getCarPlateValidationMessage() {
  return '请输入标准蓝牌车牌号，如：辽A12345'
}
