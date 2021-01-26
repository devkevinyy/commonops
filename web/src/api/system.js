import req from "../utils/axios";

const getUserFeedbackList = (page, size) =>{
    return req.get('user/feedback', {page, size})
};
export {getUserFeedbackList}