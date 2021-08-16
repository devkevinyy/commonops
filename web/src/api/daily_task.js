import req from "../utils/axios";

const postAddDailyJob = (data) =>{
    return req.post('dailyJob', data)
};
export {postAddDailyJob}

const getDailyJobs = (params) =>{
    return req.get('dailyJob/list',params)
};
export {getDailyJobs}

const getDailyJobDetail = (id) =>{
    return req.get('dailyJob/info', {"id": id})
};
export {getDailyJobDetail}

const putDailyJob = (data) =>{
    return req.put('dailyJob', data)
};
export {putDailyJob}

const putDailyJobExecutorUser = (data) =>{
    return req.put('dailyJob/executorUser', data)
};
export {putDailyJobExecutorUser}