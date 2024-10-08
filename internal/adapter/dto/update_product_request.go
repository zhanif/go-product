package dto

type UpdateProductRequest struct {
	Name  string `json:"name" validate:"required,min=3,max=255"`
	Price uint32 `json:"price" validate:"required,min=0"`
	Stock uint32 `json:"stock" validate:"required,min=0"`
}

func (dto *UpdateProductRequest) Validate() error {
	return validate.Struct(dto)
}
