package util

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

// IsValid kiểm tra xem số điện thoại có đúng định dạng hay không
func IsValidPhoneNumber(phoneNumber string) bool {
	// Sử dụng biểu thức chính quy để kiểm tra định dạng số điện thoại
	// Đây chỉ là một biểu thức chính quy đơn giản, bạn có thể điều chỉnh để phù hợp với định dạng số điện thoại của bạn
	regex := `^\+[1-9]\d{1,14}$`
	match, _ := regexp.MatchString(regex, string(phoneNumber))
	return match
}

func IsValidEmail(email string) bool {
	// Biểu thức chính quy để kiểm tra định dạng email
	// Đây là một biểu thức chính quy đơn giản, không phải là một biểu thức chính quy hoàn chỉnh cho địa chỉ email
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// Sử dụng hàm MatchString để kiểm tra định dạng
	match, _ := regexp.MatchString(regex, email)

	return match
}
