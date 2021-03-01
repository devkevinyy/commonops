import axios from "axios";

axios.interceptors.request.use(function(config) {
    return config;
});

axios.interceptors.response.use(function(config) {
    return config;
});

const ServerBase = "http://commonops.com:9999/";
const WSBase = "ws://commonops.com:9999/";

export { ServerBase, WSBase };
