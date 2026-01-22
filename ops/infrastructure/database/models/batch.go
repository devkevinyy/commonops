package models

import (
	"strconv"

	"github.com/chujieyang/commonops/ops/infrastructure/database"
)

type ecsId struct {
	Id uint
}

// Children is a child info
type Children struct {
	Title string `json:"title"`
	Key   string `json:"key"`
}

// NodeTree is a single parent tree.
type NodeTree struct {
	Title    string     `json:"title"`
	Key      string     `json:"key"`
	Children []Children `json:"children"`
}

// getEcsByUserRoleId gets ecs id by user role id for user
func getEcsByUserRoleId(userId uint) (ecsIdList []uint) {
	var ecsIds []ecsId
	querySql := "SELECT ecs.id FROM ecs " +
		"INNER JOIN role_resources AS rr ON ecs.id = rr.resource_id " +
		"INNER JOIN user_roles AS ur ON rr.role_id = ur.role_id " +
		"WHERE ecs.data_status > 0 AND rr.resource_type = 'ecs' AND ur.user_id = ?"

	database.Mysql().Raw(querySql, userId).Scan(&ecsIds)
	for _, v := range ecsIds {
		ecsIdList = append(ecsIdList, v.Id)
	}

	return ecsIdList
}

// xInSlice is equal to in of Python.
func xInSlice(a uint, list []uint) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// GetAllNodeTree get ecs by user role id for user
func GetAllNodeTree(userId uint) (nodeTreeData []NodeTree, err error) {
	isSuperAdminOrOps := IsUserSuperAdminOrOps(userId)
	accountList := GetAllCloudAccounts()
	for i := range accountList {
		p := NodeTree{
			Title: accountList[i].Name,
			Key:   "p" + strconv.Itoa(int(accountList[i].Id)),
		}
		ecsList, err := GetEcsByAccount(accountList[i].Id)
		if err != nil {
			return nodeTreeData, err
		}
		for j := range ecsList {
			c := Children{
				Title: ecsList[j].InstanceName + ":" + ecsList[j].PublicIpAddress,
				Key:   strconv.Itoa(int(ecsList[j].ID)),
			}
			// 管理员拥有所有主机的权限
			if isSuperAdminOrOps {
				p.Children = append(p.Children, c)
			} else {
				// 普通成员拥有的ecs id列表
				idList := getEcsByUserRoleId(userId)
				// ecs id
				id := ecsList[j].ID
				if xInSlice(id, idList) {
					p.Children = append(p.Children, c)
				}
			}
		}
		nodeTreeData = append(nodeTreeData, p)
	}

	return nodeTreeData, nil
}
