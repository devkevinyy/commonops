import axios from 'axios'

axios.interceptors.request.use(function(config) {
    return config
});

axios.interceptors.response.use(function(config) {
    return config
});

// const ServerBase = "https://you_host:443/api/";
// const WSBase = "wss://you_host:443/api/";
//
const ServerBase = "http://localhost:9999/";
const WSBase = "ws://localhost:9999/";

export { ServerBase, WSBase};
