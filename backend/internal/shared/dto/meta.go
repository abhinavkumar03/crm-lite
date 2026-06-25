package dto

func BuildPaginationMeta(
	page int,
	pageSize int,
	totalItems int64,
) PaginationMeta {

	totalPages := int(totalItems) / pageSize

	if int(totalItems)%pageSize != 0 {
		totalPages++
	}

	return PaginationMeta{
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}
}
