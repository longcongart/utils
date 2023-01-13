package utils

// 创建验证器
import (
	"strconv"

	"github.com/go-playground/validator/v10"
)

// NewValidator func 从model fields创建新验证器.
func NewValidator() *validator.Validate {
	// 为Product等模型创建新的验证器
	validate := validator.New()

	// 自定义验证器 采用雪花算法字段.
	_ = validate.RegisterValidation("id", func(fl validator.FieldLevel) bool {
		field := fl.Field().String()
		if uid, err := strconv.ParseInt(field, 10, 64); uid <= 0 || err != nil {
			return true
		}
		return false
	})

	return validate
}

// ValidatorErrors 用于显示每个无效字段的验证错误.
func ValidatorErrors(err error) map[string]string {
	// 定义字段字典.
	fields := map[string]string{}

	// 为每个无效字段生成错误消息.
	for _, err := range err.(validator.ValidationErrors) {
		fields[err.Field()] = err.Error()
	}

	return fields
}
