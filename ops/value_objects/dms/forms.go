package dms

/**
 * @Author: yangchujie
 * @Author: chujieyang@gmail.com
 * @Date: 6/15/21 2:59 PM
 * @Desc:
 */

type DmsAuthDataResp struct {
	Total    uint              `json:"total"`
	Page     uint              `json:"page"`
	AuthData interface{} `json:"authData"`
}

type DmsInstanceDataResp struct {
	Total        uint                  `json:"total"`
	Page         uint                  `json:"page"`
	InstanceData  interface{} `json:"instanceData"`
}

type DmsInstanceForm struct {
	InstanceId    string   `form:"instanceId" json:"instanceId"`
	InstanceType  string   `form:"instanceType" json:"instanceType"`
	InstanceAlias string   `form:"instanceAlias" json:"instanceAlias"`
	Host          string   `form:"host" json:"host"`
	Port          int      `form:"port" json:"port"`
	OperUser      string   `form:"operUser" json:"operUser"`
	OperPwd       string   `form:"operPwd" json:"operPwd"`
	Database      []string `form:"database" json:"database"`
}

type DmsInstanceQueryForm struct {
	Query    string `form:"query" json:"query"`
	Page     uint    `form:"page" json:"page" binding:"required"`
	PageSize uint    `form:"pageSize" json:"pageSize" binding:"required"`
}

type DmsInstanceDbDeleteForm struct {
	InstanceId string `form:"instanceId" json:"instanceId"`
	DatabaseId string `form:"databaseId" json:"databaseId"`
}

type DmsInstanceDbAddForm struct {
	InstanceId string `form:"instanceId" json:"instanceId"`
	DbName     string `form:"dbName" json:"dbName"`
}

type DmsUserAuthQueryForm struct {
	Page       uint    `form:"page" json:"page" binding:"required"`
	Size       uint    `form:"size" json:"size" binding:"required"`
	InstanceId string `form:"instanceId" json:"instanceId"`
	EmpId      string `form:"empId" json:"empId"`
	OperType   string `form:"operType" json:"operType"`
}

type DmsDeleteUserAuthForm struct {
	Id string `form:"id" json:"id" binding:"required"`
}

type DmsUserExecSqlForm struct {
	SelectedNodeId   string `form:"selectedNodeId" json:"selectedNodeId" binding:"required"`
	SelectedNodeType string `form:"selectedNodeType" json:"selectedNodeType" binding:"required"`
	SqlInput         string `form:"sqlInput" json:"sqlInput" binding:"required"`
	SqlDescription   string `form:"sqlDescription" json:"sqlDescription"`
}

type DmsUserQueryLog struct {
	Page     uint `form:"page" json:"page" binding:"required"`
	PageSize uint `form:"pageSize" json:"pageSize" binding:"required"`
}

type DmsAddUserAuthForm struct {
	EmpId            string `json:"empId" binding:"required"`
	SelectedNodeId   string `json:"selectedNodeId" binding:"required"`
	SelectedNodeType string `json:"selectedNodeType" binding:"required"`
	OperType         int    `json:"operType" binding:"required"`
	OperCount        int    `json:"operCount" binding:"required"`
	ValidTime        string `json:"validTime" binding:"required"`
	AuthReason       string `json:"authReason" binding:"required"`
	AllowTables      string `json:"allowTables"`
	ApproveEmpId     string `json:"approveEmpId"`
}

