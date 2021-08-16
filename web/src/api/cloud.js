import req from "../utils/axios";

const getCloudServers = (params) => {
    return req.get("cloud/servers", params);
};
export { getCloudServers };

const getCloudServerDetail = (id) => {
    return req.get("cloud/server", { id: id });
};
export { getCloudServerDetail };

const getCloudAccouts = (page, size) => {
    return req.get("cloud/accounts", { page, size });
};
export { getCloudAccouts };

const postCloudAccouts = (data) => {
    return req.post("cloud/accounts", data);
};
export { postCloudAccouts };

const putCloudAccouts = (data) => {
    return req.put("cloud/accounts", data);
};
export { putCloudAccouts };

const deleteCloudAccouts = (data) => {
    return req.delete("cloud/accounts", data);
};
export { deleteCloudAccouts };

const getCloudMonitorEcs = (instanceId, timeDimension, metricDimension) => {
    return req.get("cloud/monitor/ecs", {
        instanceId,
        timeDimension,
        metricDimension,
    });
};
export { getCloudMonitorEcs };

const getCloudMonitorRds = (instanceId, timeDimension, metricDimension) => {
    return req.get("cloud/monitor/rds", {
        instanceId,
        timeDimension,
        metricDimension,
    });
};
export { getCloudMonitorRds };

const getCloudMonitorKv = (instanceId, timeDimension, metricDimension) => {
    return req.get("cloud/monitor/kv", {
        instanceId,
        timeDimension,
        metricDimension,
    });
};
export { getCloudMonitorKv };

const getCloudMonitorSlb = (instanceId, timeDimension, metricDimension) => {
    return req.get("cloud/monitor/slb", {
        instanceId,
        timeDimension,
        metricDimension,
    });
};
export { getCloudMonitorSlb };

const getCloudRds = (params) => {
    return req.get("cloud/rds", params);
};
export { getCloudRds };

const getCloudKv = (params) => {
    return req.get("cloud/kv", params);
};
export { getCloudKv };

const getCloudSlb = (params) => {
    return req.get("cloud/slb", params);
};
export { getCloudSlb };

const getCloudRdsDetail = (id) => {
    return req.get("cloud/rds/detail", { id: id });
};
export { getCloudRdsDetail };

const getCloudKvDetail = (id) => {
    return req.get("cloud/kv/detail", { id: id });
};
export { getCloudKvDetail };

const getCloudSlbDetail = (id) => {
    return req.get("cloud/slb/detail", { id: id });
};
export { getCloudSlbDetail };

const postCloudServer = (data) => {
    return req.post("cloud/servers", data);
};
export { postCloudServer };

const putCloudServer = (data) => {
    return req.put("cloud/servers", data);
};
export { putCloudServer };

const deleteCloudServer = (id) => {
    return req.delete("cloud/servers", { id: id });
};
export { deleteCloudServer };

const deleteCloudKv = (id) => {
    return req.delete("cloud/kv", { id: id });
};
export { deleteCloudKv };

const postCloudKv = (data) => {
    return req.post("cloud/kv", data);
};
export { postCloudKv };

const putCloudKv = (data) => {
    return req.put("cloud/kv", data);
};
export { putCloudKv };

const putCloudRds = (data) => {
    return req.put("cloud/rds", data);
};
export { putCloudRds };

const postCloudRds = (data) => {
    return req.post("cloud/rds", data);
};
export { postCloudRds };

const deleteCloudRds = (id) => {
    return req.delete("cloud/rds", { id: id });
};
export { deleteCloudRds };

const deleteCloudSlb = (id) => {
    return req.delete("cloud/slb", { id: id });
};
export { deleteCloudSlb };

const getCloudOtherRes = (params) => {
    return req.get("cloud/other", params);
};
export { getCloudOtherRes };

const postCloudOtherRes = (data) => {
    return req.post("cloud/other", data);
};
export { postCloudOtherRes };

const putCloudOtherRes = (data) => {
    return req.put("cloud/other", data);
};
export { putCloudOtherRes };

const deleteCloudOtherRes = (id) => {
    return req.delete("cloud/other", { id: id });
};
export { deleteCloudOtherRes };

const getCloudServersAllData = (params) => {
    return req.get("cloud/servers/treedata", params);
};
export { getCloudServersAllData };

const postCloudServerBatchSSH = (data) => {
    return req.post("cloud/servers/batch/ssh", data);
};
export { postCloudServerBatchSSH };
