import req from "../utils/axios";

const postCdProcessTemplate = (params) => {
    return req.post("cd/processTemplate", params);
};
export { postCdProcessTemplate };

const getCdProcessTemplateData = (params) => {
    return req.get("cd/processTemplateList", params);
};
export { getCdProcessTemplateData };

const postCdProcessLog = (params) => {
    return req.post("cd/processLog", params);
};
export { postCdProcessLog };

const getCdProcessLog = (params) => {
    return req.get("cd/processLog", params);
};
export { getCdProcessLog };