import req from "../utils/axios";

const getUserPermissionsList = (params) =>{
    return req.get('user/permissions', params)
};
export {getUserPermissionsList}

// 权限管理 - 获取所有的权限链接
const getPermissionsList = (page, size) =>{
    return req.get('permissions/list', {page, size})
};
export {getPermissionsList}

// 权限链接 - 获取权限链接的详情
const getAuthLink = (id) => {
    return req.get('permissions/authLink', { id: id });
};
export {getAuthLink}

// 权限链接 - 修改权限链接的详情
const putAuthLink = (data) => {
    return req.put("permissions/authLink", data);
};
export {putAuthLink}

// 权限管理 - 新增权限链接
const postAddAuthLink = (data) => {
    return req.post('permissions/authLink', data)
};
export {postAddAuthLink}

// 权限管理 - 删除权限链接
const deleteAuthLink = (data) => {
    return req.delete('permissions/authLink', data)
};
export {deleteAuthLink}
