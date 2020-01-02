import req from "../utils/axios";

const SyncTimeOut = 90000;

const getSyncAliyunEcs = (params) =>{
    return req.get('/data/syncAliyunEcs', params, SyncTimeOut)
};
export {getSyncAliyunEcs}

const getSyncAliyunRds = (params) =>{
    return req.get('/data/syncAliyunRds', params, SyncTimeOut)
};
export {getSyncAliyunRds}

const getSyncAliyunKv = (params) =>{
    return req.get('/data/syncAliyunKv', {}, SyncTimeOut)
};
export {getSyncAliyunKv}

const getSyncAliyunSlb = (params) =>{
    return req.get('/data/syncAliyunSlb', {}, SyncTimeOut)
};
export {getSyncAliyunSlb}

const getAliyunStatisData = (params) =>{
    return req.get('/data/syncAliyunStatisData')
};
export {getAliyunStatisData}