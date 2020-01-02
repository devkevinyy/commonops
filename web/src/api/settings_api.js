import req from "../utils/axios";

const getSettingsValue = (params) =>{
    return req.get('settings/value', params)
};
export {getSettingsValue}

const putSettingsValue = (data) => {
    return req.put('settings/value', data)
};
export {putSettingsValue}
