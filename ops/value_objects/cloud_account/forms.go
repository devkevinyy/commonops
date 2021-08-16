package cloud_account

/**
 * @Author: yangchujie
 * @Author: chujieyang@gmail.com
 * @Date: 6/11/21 4:33 PM
 * @Desc:
 */

type CloudAccountResp struct {
	Total    uint                   `json:"total"`
	Page     uint                   `json:"page"`
	Accounts interface{}         `json:"accounts"`
}

type CloudAccountForm struct {
	Id            uint    `form:"id" json:"id"`
	AccountType   string `form:"accountType" json:"accountType"`
	AccountName   string `form:"accountName" json:"accountName"`
	AccountPwd    string `form:"accountPwd" json:"accountPwd"`
	AccountKey    string `form:"accountKey" json:"accountKey"`
	AccountSecret string `form:"accountSecret" json:"accountSecret"`
	AccountRegion string `form:"accountRegion" json:"accountRegion"`
	BankAccount   int    `form:"bankAccount" json:"bankAccount"`
}