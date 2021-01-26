import req from "../utils/axios";

const getUserPermissionsList = () =>{
    return req.get('user/permissions')
};
export {getUserPermissionsList}
