import req from "../utils/axios";

const getLogAuth = (params) =>{
    return req.get('log/log_auth', params)
};
export {getLogAuth}

const getLogCommonQuery = (params) =>{
    return req.get('log/common_query', params)
};
export {getLogCommonQuery}