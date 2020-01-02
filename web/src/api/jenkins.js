import req from '../utils/axios';

const getJenkinsAllJobs = (params) =>{
    return req.get('jenkins/jobs', params)
};
export {getJenkinsAllJobs}

const getJenkinsJobBuildList = (params) =>{
    return req.get('jenkins/jobBuildList', params)
};
export {getJenkinsJobBuildList}

const getJenkinsJobLastBuildLog = (params) =>{
    return req.get('jenkins/jobLastBuildLog', params)
};
export {getJenkinsJobLastBuildLog}

const getJenkinsJobBuildLog = (params) =>{
    return req.get('jenkins/jobBuildLog', params)
};
export {getJenkinsJobBuildLog}

const postJenkinsStartJob = (params) =>{
    return req.post('jenkins/jobBuild', params)
};
export {postJenkinsStartJob}

const postJenkinsEnableJob = (params) =>{
    return req.post('jenkins/enableJob', params)
};
export {postJenkinsEnableJob}

const postJenkinsDisableJob = (params) =>{
    return req.post('jenkins/disableJob', params)
};
export {postJenkinsDisableJob}