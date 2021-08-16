import req from "../utils/axios";

const getUserPermissionsList = (params) =>{
    return req.get('user/permissions', params)
};
export {getUserPermissionsList}
