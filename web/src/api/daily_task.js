import req from "../utils/axios";

const postAddDailyJob = (data) =>{
    return req.post('daily_job/', data)
};
export {postAddDailyJob}

const getDailyJobs = (params) =>{
    return req.get('daily_job/',params)
};
export {getDailyJobs}

const getDailyJobDetail = (id) =>{
    return req.get('daily_job/info/'+id)
};
export {getDailyJobDetail}

const putDailyJob = (data) =>{
    return req.put('daily_job/', data)
};
export {putDailyJob}

const putDailyJobExecutorUser = (data) =>{
    return req.put('daily_job/executorUser', data)
};
export {putDailyJobExecutorUser}