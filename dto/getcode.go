package dto

type GetCodeRequestDto struct {
	latitude         float32
	longitude        float32
	localityLanguage *string
}

type GetCodeResponseDto struct {
	City string
}
