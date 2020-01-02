import req from '../utils/axios';

/*
    用户登录
 */
const postUserLogin = (data) => {
    return req.post('user/login', data)
};
export {postUserLogin}

/*
    用户创建
 */
const postUserCreate = (data) => {
    return req.post('user/create', data)
};
export {postUserCreate}

/*
    更新用户状态
 */
const putUserCreate = (data) => {
    return req.put('user/active', data)
};
export {putUserCreate}

/*
    用户 jwt token 刷新
 */
const getUserTokenRefresh = () => {
    return req.get('user/tokenRefresh')
};
export {getUserTokenRefresh}

/*
    获取注册用户列表
 */
const getUsersList = (page, size) =>{
    return req.get('user/list', {page, size})
};
export {getUsersList}

/*
    用户修改密码
 */
const postUpdatePassword = (data) => {
    return req.post('user/updatePassword', data)
};
export {postUpdatePassword}

/*
    用户提交需求&bug反馈
 */
const postUserFeedback = (data) => {
    return req.post('user/feedback', data)
};
export {postUserFeedback}