import axios from "axios";
import { message } from "antd";
import { ServerBase } from "../config";

axios.defaults.headers.post["Content-Type"] = "application/json";

let service = axios.create({
    baseURL: ServerBase,
    timeout: 7000,
    withCredentials: false,
});

function getLocalToken() {
    return localStorage.getItem("ops_token");
}
export { getLocalToken };

function getClusterId() {
    return localStorage.getItem("clusterId");
}
export { getClusterId };

//添加请求拦截器
service.interceptors.request.use(
    function(config) {
        config.headers.Authorization = getLocalToken();
        config.headers["ClusterId"] = getClusterId();
        return config;
    },
    function(error) {
        return Promise.reject(error);
    },
);

//添加响应拦截器
service.interceptors.response.use(
    function(response) {
        return response["data"];
    },
    function(error) {
        if (!error.response) {
            message.error("Network Error");
            return Promise.reject(error);
        }
        let httpStatus = error.response.status;
        switch (httpStatus) {
            case 400:
                message.error("[400]客户端请求参数错误");
                break;
            case 401:
                window.location.href = "/login?code=401";
                break;
            case 403:
                window.location.href = "/login?code=403";
                break;
            case 500:
                message.error("[500]服务异常，请联系运维处理");
                break;
            default:
                message.error("服务异常，请联系运维处理");
                break;
        }
        return Promise.reject(error);
    },
);

export default {
    //get请求
    get(url, param, timeout) {
        return new Promise((resolve, reject) => {
            service({
                method: "get",
                url,
                params: param,
                timeout: timeout,
            })
                .then((res) => {
                    resolve(res);
                })
                .catch((err) => {
                    if (!err.response) {
                        message.error("Network Error");
                    } else {
                        message.error(err.response.data.message);
                    }
                    reject(err);
                });
        });
    },
    //post请求
    post(url, param, timeout) {
        return new Promise((resolve, reject) => {
            service({
                method: "post",
                url,
                data: param,
                timeout: timeout,
            })
                .then((res) => {
                    resolve(res);
                })
                .catch((err) => {
                    if (!err.response) {
                        message.error("Network Error");
                    } else {
                        message.error(err.response.data.message);
                    }
                    reject(err);
                });
        });
    },
    //put请求
    put(url, param, timeout) {
        return new Promise((resolve, reject) => {
            service({
                method: "put",
                url,
                data: param,
                timeout: timeout,
            })
                .then((res) => {
                    resolve(res);
                })
                .catch((err) => {
                    if (!err.response) {
                        message.error("Network Error");
                    } else {
                        message.error(err.response.data.message);
                    }
                    reject(err);
                });
        });
    },
    //delete请求
    delete(url, param, timeout) {
        return new Promise((resolve, reject) => {
            service({
                method: "delete",
                url,
                data: param,
                timeout: timeout,
            })
                .then((res) => {
                    resolve(res);
                })
                .catch((err) => {
                    if (!err.response) {
                        message.error("Network Error");
                    } else {
                        message.error(err.response.data.message);
                    }
                    reject(err);
                });
        });
    },
};
