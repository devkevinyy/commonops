import req from "../utils/axios";

const getJenkinsJob = (params) => {
    return req.get("ci/job", params);
};
export { getJenkinsJob };

const postJenkinsJob = (params) => {
    return req.post("ci/job", params);
};
export { postJenkinsJob };

const putJenkinsJob = (params) => {
    return req.put("ci/job", params);
};
export { putJenkinsJob };

const deleteJenkinsJob = (params) => {
    return req.delete("ci/job", params);
};
export { deleteJenkinsJob };

const getJenkinsAllJobs = (params) => {
    return req.get("ci/jobList", params);
};
export { getJenkinsAllJobs };

const getJobBuildList = (params) => {
    return req.get("ci/buildList", params);
};
export { getJobBuildList };

const postJenkinsStartJob = (params) => {
    return req.post("ci/build", params);
};
export { postJenkinsStartJob };

const deleteJenkinsJobBuildLog = (params) => {
    return req.delete("ci/build", params);
};
export { deleteJenkinsJobBuildLog };

const getJenkinsJobBuildLog = (params) => {
    return req.get("ci/buildLog", params);
};
export { getJenkinsJobBuildLog };

const getJenkinsJobBuildStages = (params) => {
    return req.get("ci/build/stages", params);
};
export { getJenkinsJobBuildStages };

const getJenkinsJobBuildStageLog = (params) => {
    return req.get("ci/build/stage/log", params);
};
export { getJenkinsJobBuildStageLog };

const getJenkinsJobBuildArchiveArtifactsInfo = (params) => {
    return req.get("ci/build/archiveArtifactsInfo", params);
};
export { getJenkinsJobBuildArchiveArtifactsInfo };

const getJenkinsCredentialsList = (params) => {
    return req.get("ci/credentials/list", params);
};
export { getJenkinsCredentialsList };

const postJenkinsAddCredential = (params) => {
    return req.post("ci/credential", params);
};
export { postJenkinsAddCredential };

const postJenkinsEnableJob = (params) => {
    return req.post("jenkins/enableJob", params);
};
export { postJenkinsEnableJob };

const postJenkinsDisableJob = (params) => {
    return req.post("jenkins/disableJob", params);
};
export { postJenkinsDisableJob };
