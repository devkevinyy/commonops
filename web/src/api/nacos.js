import req from "../utils/axios";

const postNacosServer = (params) => {
    return req.post("configCenter/nacos", params);
};
export { postNacosServer };

const getNacosList = (params) => {
    return req.get("configCenter/nacos/list", params);
};
export { getNacosList };

const getNacosNamespaceList = (params) => {
    return req.get("configCenter/nacos/namespaces", params);
};
export { getNacosNamespaceList };

const getNacosConfigsList = (params) => {
    return req.get("configCenter/nacos/configs", params);
};
export { getNacosConfigsList };

const getNacosConfigDetail = (params) => {
    return req.get("configCenter/nacos/config", params);
};
export { getNacosConfigDetail };

const getNacosAllConfigs = (params) => {
    return req.get("configCenter/nacos/configs/all", params);
};
export { getNacosAllConfigs };

const postNacosConfig = (params) => {
    return req.post("configCenter/nacos/config", params);
};
export { postNacosConfig };

const postNacosConfigSync = (params) => {
    return req.post("configCenter/nacos/config/sync", params);
};
export { postNacosConfigSync };

const putNacosConfig = (params) => {
    return req.put("configCenter/nacos/config", params);
};
export { putNacosConfig };

const deleteNacosConfig = (params) => {
    return req.delete("configCenter/nacos/config", params);
};
export { deleteNacosConfig };

const postNacosConfigCopy = (params) => {
    return req.post("configCenter/nacos/config/copy", params);
};
export { postNacosConfigCopy };

const getConfigTemplates = (params) => {
    return req.get("configCenter/configTemplates", params);
};
export { getConfigTemplates };

const postConfigTemplate = (params) => {
    return req.post("configCenter/configTemplate", params);
};
export { postConfigTemplate };

const putConfigTemplate = (params) => {
    return req.put("configCenter/configTemplate", params);
};
export { putConfigTemplate };

const deleteConfigTemplate = (params) => {
    return req.delete("configCenter/configTemplate", params);
};
export { deleteConfigTemplate };

const getConfigTemplatesAll = (params) => {
    return req.get("configCenter/configTemplates/all", params);
};
export { getConfigTemplatesAll };
