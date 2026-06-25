package utils

import "github.com/abhinavkumar03/crm-lite/backend/internal/shared/constants"

func NormalizePagination(page, pageSize int) (int, int) {

	if page <= 0 {
		page = constants.DefaultPage
	}

	if pageSize <= 0 {
		pageSize = constants.DefaultPageSize
	}

	if pageSize > constants.MaxPageSize {
		pageSize = constants.MaxPageSize
	}

	return page, pageSize
}

func Offset(page, pageSize int) int {
	return (page - 1) * pageSize
}
