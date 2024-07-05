package util

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func TryBind(obj any, c *gin.Context) (error, gin.H) {
	if err := c.ShouldBindJSON(&obj); err != nil {
		var allErrors validator.ValidationErrors
		hasErrors := errors.As(err, &allErrors)
		if hasErrors {
			errs := make([]string, 0)
			for _, e := range allErrors {
				errs = append(errs, e.Error())
			}
			return errors.New("invalid with errors"), gin.H{
				"errors": errs,
			}
		}
		return errors.New("invalid with error"), gin.H{
			"error": "invalid json",
		}
	}

	return nil, nil
}
