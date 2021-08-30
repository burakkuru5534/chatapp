package auth

type AuthResponse struct {
	UserID         int64
	UserName       string
	UserEmail      string
	UserFullName   string
	AccessToken    string
	Roles          []string
	Policy         interface{}
	DtFn           string
	LicenseType    string
	FiCompTitle    string
	ProductType    int
	Implementation string
}
