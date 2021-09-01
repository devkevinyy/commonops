package permission

// AuthLinkInfoForm is frontend form of authLink
type AuthLinkInfoForm struct {
	Id          string `json:"id"`
	CanDelete   int    `json:"canDelete"`
	AuthType    string `json:"authType"`
	Name        string `json:"name"`
	Description string `json:"description"`
	UrlPath     string `json:"urlPath"`
}
