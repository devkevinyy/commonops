import jwt_decode from "jwt-decode";

function getCurrentUserId() {
    let token = localStorage.getItem("ops_token");
    if(token===undefined || token===null || token===""){
        window.location.href = "/login";
        return
    }
    let info = jwt_decode(token);
    return info['userInfo']['userId'];
}

function isSuperAdmin() {
    let token = localStorage.getItem("ops_token");
    if(token===undefined || token===null || token===""){
        window.location.href = "/login";
        return
    }
    let info = jwt_decode(token);
    return info['userInfo']['isSuperAdmin'];
}

export { getCurrentUserId, isSuperAdmin }
