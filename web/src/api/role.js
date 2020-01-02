import req from '../utils/axios';

/*
    权限管理 - 新建角色
 */
const postAddRole = (data) => {
    return req.post('roles/addRole', data)
};
export {postAddRole}

/*
    权限管理 - 更新角色
 */
const putUpdateRole = (data) => {
    return req.put('roles/updateRole', data)
};
export {putUpdateRole}

/*
    权限管理 - 删除角色
 */
const deleteRole = (data) => {
    return req.delete('roles/deleteRole', data)
};
export {deleteRole}

/*
    权限管理 - 获取角色的权限链接
 */
const getRoleAuthLinks = (roleId) => {
    return req.get('roles/authLink', {roleId})
};
export {getRoleAuthLinks}

/*
    权限管理 - 新增权限链接
 */
const postAddAuthLink = (data) => {
    return req.post('roles/authLink', data)
};
export {postAddAuthLink}

/*
    权限管理 - 新建角色的权限链接
 */
const postRoleAuthLinks = (data) => {
    return req.post('roles/authLinks', data)
};
export {postRoleAuthLinks}

/*
    权限管理 - 删除权限链接
 */
const deleteAuthLink = (data) => {
    return req.delete('roles/authLink', data)
};
export {deleteAuthLink}

/*
    权限管理 - 获取角色列表
 */
const getRolesList = (page, size) =>{
    return req.get('roles/list', {page, size})
};
export {getRolesList}

/*
    权限管理 - 获取角色对应的用户列表
 */
const getRoleUserList = (roleId) =>{
    return req.get('roles/users', {roleId})
};
export {getRoleUserList}

/*
    权限管理 - 新增角色对应的用户
 */
const postRoleUserList = (data) =>{
    return req.post('roles/users', data)
};
export {postRoleUserList}

/*
    权限管理 - 获取角色下的服务器资源列表
 */
const getRoleResourceList = (roleId, resourceType) =>{
    return req.get('roles/resources', {roleId, resourceType})
};
export {getRoleResourceList}

/*
    权限管理 - 新增角色下的服务器资源
 */
const postRoleResourcesList = (data) =>{
    return req.post('roles/resources', data)
};
export {postRoleResourcesList}

/*
    权限管理 - 获取所有的权限链接
 */
const getPermissionsList = (page, size) =>{
    return req.get('permissions/list', {page, size})
};
export {getPermissionsList}