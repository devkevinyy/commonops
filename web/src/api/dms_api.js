import req from "../utils/axios";

const getDmsInstanceData = (params) => {
    return req.get("dms/instances", params);
};
export { getDmsInstanceData };

const postDmsInstanceData = (params) => {
    return req.post("dms/instance", params);
};
export { postDmsInstanceData };

const deleteDmsInstanceData = (params) => {
    return req.delete("dms/instance", params);
};
export { deleteDmsInstanceData };

const getAllDmsInstanceData = (params) => {
    return req.get("dms/instances/all", params);
};
export { getAllDmsInstanceData };

const deleteDmsInstanceDbData = (params) => {
    return req.delete("dms/databaseData", params);
};
export { deleteDmsInstanceDbData };

const postDmsInstanceDbData = (params) => {
    return req.post("dms/databaseData", params);
};
export { postDmsInstanceDbData };

const getDmsAuthData = (params) => {
    return req.get("dms/authData", params);
};
export { getDmsAuthData };

const getDmsDatabaseData = (params) => {
    return req.get("dms/databaseData", params);
};
export { getDmsDatabaseData };

const postDmsUserAuth = (params) => {
    return req.post("dms/auth", params);
};
export { postDmsUserAuth };

const getUserDmsInstanceData = (params) => {
    return req.get("dms/userInstanceData", params);
};
export { getUserDmsInstanceData };

const getUserDmsDatabaseData = (params) => {
    return req.get("dms/userDatabaseData", params);
};
export { getUserDmsDatabaseData };

const getUserDmsLog = (params) => {
    return req.get("dms/userLog", params);
};
export { getUserDmsLog };

const postDmsUserExecSQL = (params) => {
    return req.post("dms/userExecSQL", params);
};
export { postDmsUserExecSQL };

const deleteDmsUserAuth = (params) => {
    return req.delete("dms/auth", params);
};
export { deleteDmsUserAuth };
