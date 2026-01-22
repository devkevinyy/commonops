import req from "../utils/axios"

// 根据用户权限获取目录树列表
const getAllNodeTree = (userId) => {
  return req.get("batch/allNodeTree", {userId: userId});
};
export { getAllNodeTree }
