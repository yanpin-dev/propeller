package middleware

import (
	"github.com/yanpin-dev/propeller/pkg/app"
	"github.com/yanpin-dev/propeller/pkg/logger"
	"github.com/yanpin-dev/propeller/pkg/xerrors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	zh2 "github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/translations/zh"
	"reflect"
	"strings"
	"sync"

	"net/http"
)

// ErrorMiddleware 统一异常处理
type ErrorMiddleware struct {
	log logger.LogInfoFormat

	uni      *ut.UniversalTranslator
	validate *validator.Validate

	once sync.Once
}

func NewErrorMiddleware(log logger.LogInfoFormat) app.Middleware {
	return &ErrorMiddleware{
		log: log,
	}
}

func (h *ErrorMiddleware) init() {
	h.once.Do(func() {
		h.uni = ut.New(zh2.New(), zh2.New())

		if validate, ok := binding.Validator.Engine().(*validator.Validate); ok {
			validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
				name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
				if name == "-" {
					return fld.Name
				}
				return name
			})
			h.validate = validate
		} else {
			h.log.Warn("failed to init validate")
		}
	})
}

func (h *ErrorMiddleware) NewHandler() gin.HandlerFunc {
	h.init()

	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		if h.handleValidateErrors(c, err) {
			return
		}

		if h.handleBizErrors(c, err) {
			return
		}

		if h.handleUnknownErrors(c, err) {
			return
		}

	}
}

func (h *ErrorMiddleware) handleValidateErrors(c *gin.Context, err error) bool {
	if _, ok := err.(*validator.InvalidValidationError); ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code": xerrors.ValidateFailed,
			"msg":  fmt.Sprintf("%v", err),
		})
		return true
	}

	if errs, ok := err.(validator.ValidationErrors); ok {
		//local := c.DefaultQuery("local", "zh")
		tans, _ := h.uni.GetTranslator("zh")
		zh.RegisterDefaultTranslations(h.validate, tans)
		fieldErrors := make(map[string]string, len(errs))
		for _, v := range errs {
			fieldErrors[v.Field()] = v.Translate(tans)
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":        xerrors.InvalidInput,
			"msg":         "输入数据校验异常",
			"errorFields": fieldErrors,
		})
		return true
	}

	return false
}

func (h *ErrorMiddleware) handleBizErrors(c *gin.Context, err error) bool {
	var bizErr *xerrors.BizError
	var ok bool
	if bizErr, ok = err.(*xerrors.BizError); !ok {
		return false
	}
	if bizErr.Code == xerrors.RecordNotFound {
		c.AbortWithStatusJSON(http.StatusNotFound, bizErr)
		return true
	}
	if bizErr.Code == xerrors.InvalidToken {
		c.AbortWithStatusJSON(http.StatusUnauthorized, bizErr)
		return true
	}
	c.AbortWithStatusJSON(http.StatusBadRequest, bizErr)
	return true
}

func (h *ErrorMiddleware) handleUnknownErrors(c *gin.Context, err error) bool {
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		"code":  xerrors.Unknown,
		"msg":   "未知异常",
		"error": fmt.Sprintf("%s", err.Error()),
	})

	return true
}
