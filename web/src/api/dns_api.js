import req from "../utils/axios";

const getDnsDomainListData = (params) => {
    return req.get("dns/domainList", params);
};
export { getDnsDomainListData };

const getDnsDomainHistoryListData = (params) => {
    return req.get("dns/domainHistoryList", params);
};
export { getDnsDomainHistoryListData };

const getDnsDomainRecordListData = (params) => {
    return req.get("dns/domainRecordsList", params);
};
export { getDnsDomainRecordListData };

const postDnsDomain = (params) => {
    return req.post("dns/domain", params);
};
export { postDnsDomain };

const postDnsDomainRecord = (params) => {
    return req.post("dns/domainRecord", params);
};
export { postDnsDomainRecord };

const postDnsDomainRecordUpdate = (params) => {
    return req.post("dns/domainRecordUpdate", params);
};
export { postDnsDomainRecordUpdate };

const deleteDnsDomainRecord = (params) => {
    return req.delete("dns/domainRecord", params);
};
export { deleteDnsDomainRecord };

const postDnsDomainRecordStatus = (params) => {
    return req.post("dns/domainRecordStatus", params);
};
export { postDnsDomainRecordStatus };
