import req from '../utils/axios';

const getNodes = (params) =>{
    return req.get('/kubernetes/nodes', params)
};
export {getNodes}

const getNamespaces = (params) =>{
    return req.get('/kubernetes/namespaces', params)
};
export {getNamespaces}

const getDeployments = (params) =>{
    return req.get('/kubernetes/deployments', params)
};
export {getDeployments}

const getReplicationControllers = (params) =>{
    return req.get('/kubernetes/replication_controllers', params)
};
export {getReplicationControllers}

const getReplicaSets = (params) =>{
    return req.get('/kubernetes/replica_sets', params)
};
export {getReplicaSets}

const getServices = (params) =>{
    return req.get('/kubernetes/services', params)
};
export {getServices}

const getPods = (params) =>{
    return req.get('/kubernetes/pods', params)
};
export {getPods}

const getPodLogs = (params) =>{
    return req.get('/kubernetes/pod/log', params)
};
export {getPodLogs}

const postApplyYaml = (data) =>{
    return req.post('/kubernetes/yaml_resource', data)
};
export {postApplyYaml}

const getResourceYaml = (params) =>{
    return req.get('/kubernetes/yaml', params)
};
export {getResourceYaml}

const putResourceScale = (data) =>{
    return req.put('/kubernetes/scale', data)
};
export {putResourceScale}

const deleteResource = (data) =>{
    return req.delete('/kubernetes/resource', data)
};
export {deleteResource}

const getConfigDict = (params) =>{
    return req.get('/kubernetes/config_dict', params)
};
export {getConfigDict}

const getSecretDict = (params) =>{
    return req.get('/kubernetes/secret_dict', params)
};
export {getSecretDict}

const deleteConfigMap = (data) =>{
    return req.delete('/kubernetes/config_map', data)
};
export {deleteConfigMap}

const deleteSecret = (data) =>{
    return req.delete('/kubernetes/secret', data)
};
export {deleteSecret}

const getClusterData = (params) =>{
    return req.get('/kubernetes/cluster', params)
};
export {getClusterData}

const postCluster = (data) =>{
    return req.post('/kubernetes/cluster', data)
};
export {postCluster}

const deleteCluster = (data) =>{
    return req.delete('/kubernetes/cluster', data)
};
export {deleteCluster}

const getNodeMetrics = (params) =>{
    return req.get('/kubernetes/metrics/nodes', params)
};
export {getNodeMetrics}

const postPrometheus = (data) =>{
    return req.post('/kubernetes/metrics/prometheus', data)
};
export {postPrometheus}