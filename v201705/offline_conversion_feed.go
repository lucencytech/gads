package v201705

type OfflineConversionService struct {
	Auth
}

func NewOfflineConversionService(auth *Auth) *OfflineConversionService {
	return &OfflineConversionService{Auth: *auth}
}
