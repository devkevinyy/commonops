import React, { Component } from "react";
import { Layout, message } from "antd";
import { Route, Switch } from "react-router-dom";
import WelcomeContent from "./WelcomeContent";
import ServerContent from "../components/cloud_resource/server_content";
import RdsContent from "../components/cloud_resource/rds_content";
import KvContent from "../components/cloud_resource/kv_content";
import SlbContent from "../components/cloud_resource/slb_content";
import CloudContent from "../components/cloud_resource/account_content";
import TaskContent from "../components/task/task_content";
import UserManager from "../components/permissions/user_manager";
import RolesManager from "../components/permissions/role_manager";
import PermissionsManager from "../components/permissions/permission_manager";
import PasswordManager from "../components/permissions/password_manager";
import { getUserTokenRefresh } from "../api/user";
import jwt_decode from "jwt-decode";
import SyncAliyunContent from "../components/data/sync_aliyun";
import OtherContent from "../components/cloud_resource/other_content";
import SettingsContent from "../components/system/settings_content";
import LogoutContent from "../components/logout";
import UserFeedbackManager from "../components/system/feedback_manager";
import CiCdContent from "../components/ci_cd/cicd";
import KubernetesContent from "../components/kubernetes/kubernetes";
import InstanceManageContent from "../components/dms/instance_manage";
import AuthManageContent from "../components/dms/auth_manage";
import DataManageContent from "../components/dms/data_manage";
import DomainManageContent from "../components/dns/domain_manage";
import NacosContent from "../components/nacos/nacos_content";
import NacosTemplateContent from "../components/nacos/template_content";
import { getUserPermissionsList } from "../api/permission";

const { Content } = Layout;

class ContentLayout extends Component {
    constructor(props) {
        super(props);
        this.state = {
            isSuperAdmin: false,
            actionPermissionMap: {},
        };
    }

    componentWillMount() {
        this.loadUserActionPermissions();
        this.loadUserTokenRefresh();
    }

    loadUserTokenRefresh() {
        getUserTokenRefresh().then((res) => {
            if (res.code === 0) {
                localStorage.setItem("ops_token", res.data.token);
                let info = jwt_decode(res.data.token);
                this.setState({
                    isSuperAdmin: info["userInfo"]["isSuperAdmin"],
                });
            } else {
                message.error(res.msg);
            }
        });
    }

    loadUserActionPermissions() {
        getUserPermissionsList({ authType: "操作" }).then((res) => {
            if (res.code === 0) {
                const dataList = res.data;
                let authMap = {};
                dataList.map((item) => {
                    if (item.authType === "操作") {
                        authMap[item.urlPath] = true;
                    }
                });
                this.setState({ actionPermissionMap: authMap });
            } else {
                message.error("获取用户菜单权限异常");
            }
        });
    }

    render() {
        return (
            <Content
                style={{
                    height: "100%",
                    minWidth: "900px",
                    overflow: "scroll",
                }}
            >
                <Switch>
                    <Route exact path="/admin" component={WelcomeContent} />
                    <Route
                        path="/admin/cloud_resource/cloud_server"
                        render={() => (
                            <ServerContent
                                aclAuthMap={this.state.actionPermissionMap}
                                isSuperAdmin={this.state.isSuperAdmin}
                            />
                        )}
                    />
                    <Route
                        path="/admin/cloud_resource/cloud_rds"
                        render={() => (
                            <RdsContent
                                aclAuthMap={this.state.actionPermissionMap}
                                isSuperAdmin={this.state.isSuperAdmin}
                            />
                        )}
                    />
                    <Route
                        path="/admin/cloud_resource/cloud_kv"
                        render={() => (
                            <KvContent
                                aclAuthMap={this.state.actionPermissionMap}
                                isSuperAdmin={this.state.isSuperAdmin}
                            />
                        )}
                    />
                    <Route
                        path="/admin/cloud_resource/cloud_slb"
                        render={() => (
                            <SlbContent
                                aclAuthMap={this.state.actionPermissionMap}
                                isSuperAdmin={this.state.isSuperAdmin}
                            />
                        )}
                    />
                    <Route
                        path="/admin/cloud_resource/cloud_account"
                        render={() => (
                            <CloudContent
                                aclAuthMap={this.state.actionPermissionMap}
                                isSuperAdmin={this.state.isSuperAdmin}
                            />
                        )}
                    />
                    <Route
                        path="/admin/cloud_resource/cloud_other"
                        render={() => (
                            <OtherContent
                                aclAuthMap={this.state.actionPermissionMap}
                                isSuperAdmin={this.state.isSuperAdmin}
                            />
                        )}
                    />
                    <Route
                        path="/admin/task"
                        render={() => (
                            <TaskContent
                                aclAuthMap={this.state.actionPermissionMap}
                                isSuperAdmin={this.state.isSuperAdmin}
                            />
                        )}
                    />
                    <Route
                        path="/admin/dns/domain_manage"
                        render={() => (
                            <DomainManageContent
                                aclAuthMap={this.state.actionPermissionMap}
                                isSuperAdmin={this.state.isSuperAdmin}
                            />
                        )}
                    />
                    <Route
                        path="/admin/dms/instance_manage"
                        render={() => (
                            <InstanceManageContent
                                aclAuthMap={this.state.actionPermissionMap}
                                isSuperAdmin={this.state.isSuperAdmin}
                            />
                        )}
                    />
                    <Route
                        path="/admin/dms/auth_manage"
                        render={() => (
                            <AuthManageContent
                                aclAuthMap={this.state.actionPermissionMap}
                                isSuperAdmin={this.state.isSuperAdmin}
                            />
                        )}
                    />
                    <Route
                        path="/admin/dms/data_manage"
                        render={() => (
                            <DataManageContent
                                aclAuthMap={this.state.actionPermissionMap}
                                isSuperAdmin={this.state.isSuperAdmin}
                            />
                        )}
                    />
                    <Route
                        path="/admin/cicd"
                        render={() => (
                            <CiCdContent
                                aclAuthMap={this.state.actionPermissionMap}
                            />
                        )}
                    />
                    <Route
                        path="/admin/k8s_cluster"
                        render={() => (
                            <KubernetesContent
                                aclAuthMap={this.state.actionPermissionMap}
                            />
                        )}
                    />
                    <Route
                        path="/admin/config_center/nacos"
                        render={() => (
                            <NacosContent
                                aclAuthMap={this.state.actionPermissionMap}
                            />
                        )}
                    />
                    <Route
                        path="/admin/config_center/config_template"
                        render={() => (
                            <NacosTemplateContent
                                aclAuthMap={this.state.actionPermissionMap}
                            />
                        )}
                    />

                    <Route
                        path="/admin/data/syncAliyun"
                        component={SyncAliyunContent}
                    />
                    <Route
                        path="/admin/permission/users"
                        component={UserManager}
                    />
                    <Route
                        path="/admin/permission/roles"
                        component={RolesManager}
                    />
                    <Route
                        path="/admin/permission/permissions"
                        component={PermissionsManager}
                    />
                    <Route
                        path="/admin/permission/password"
                        component={PasswordManager}
                    />
                    <Route
                        path="/admin/system/setting"
                        component={SettingsContent}
                    />
                    <Route
                        path="/admin/system/user_feedback"
                        component={UserFeedbackManager}
                    />
                    <Route path="/admin/logout" component={LogoutContent} />
                </Switch>
            </Content>
        );
    }
}

export default ContentLayout;
