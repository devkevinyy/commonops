import React, { Component } from "react";
import {Layout, message} from 'antd';
import { Route, Switch } from 'react-router-dom';
import WelcomeContent from "./WelcomeContent";
import ServerContent from "../components/cloud_resource/server_content";
import RdsContent from "../components/cloud_resource/rds_content";
import KvContent from "../components/cloud_resource/kv_content";
import SlbContent from "../components/cloud_resource/slb_content";
import CloudContent from "../components/cloud_resource/account_content";
import Task_content from "../components/task/task_content";
import UserManager from "../components/permissions/user_manager";
import RolesManager from "../components/permissions/role_manager";
import PermissionsManager from "../components/permissions/permission_manager";
import PasswordManager from "../components/permissions/password_manager";
import {getUserTokenRefresh} from "../api/user";
import jwt_decode from "jwt-decode";
import SyncAliyunContent from "../components/data/sync_aliyun";
import OtherContent from "../components/cloud_resource/other_content";
import SettingsContent from "../components/system/settings_content";
import LogoutContent from "../components/logout";
import UserFeedbackManager from "../components/system/feedback_manager";
import JenkinsContent from "../components/ci_cd/jenkins";
import KubernetesContent from "../components/kubernetes/kubernetes";

const { Content } = Layout;

class ContentLayout extends Component {

    constructor(props) {
        super(props);
        this.state = {
            isSuperAdmin: false,
        }
    }

    componentWillMount() {
        getUserTokenRefresh().then((res)=>{
            if(res.code === 0){
                localStorage.setItem('ops_token', res.data.token);
                let info = jwt_decode(res.data.token);
                this.setState({
                   isSuperAdmin: info['userInfo']['isSuperAdmin']
                });
            }
        }).catch((err)=>{
            message.error(err.toLocaleString());
        });
    }

    render() {
        return (
            <Content style={{ height: "100%", minWidth: '900px', overflow: 'scroll' }}>
                <Switch>
                    <Route exact path="/admin" component={WelcomeContent}/>
                    <Route path="/admin/cloud_resource/cloud_server" render={() => (<ServerContent isSuperAdmin={this.state.isSuperAdmin}/>)} />
                    <Route path="/admin/cloud_resource/cloud_rds" render={() => (<RdsContent isSuperAdmin={this.state.isSuperAdmin}/>)} />
                    <Route path="/admin/cloud_resource/cloud_kv" render={() => (<KvContent isSuperAdmin={this.state.isSuperAdmin}/>)} />
                    <Route path="/admin/cloud_resource/cloud_slb" render={() => (<SlbContent isSuperAdmin={this.state.isSuperAdmin}/>)} />
                    <Route path="/admin/cloud_resource/cloud_account" render={() => (<CloudContent isSuperAdmin={this.state.isSuperAdmin}/>)} />
                    <Route path="/admin/cloud_resource/cloud_other" render={() => (<OtherContent isSuperAdmin={this.state.isSuperAdmin}/>)} />
                    <Route path="/admin/task" render={() => (<Task_content isSuperAdmin={this.state.isSuperAdmin}/>)} />
                    <Route path="/admin/cicd" render={() => (<JenkinsContent/>)} />
                    <Route path="/admin/k8s_cluster" render={() => (<KubernetesContent/>)} />
                    <Route path="/admin/data/syncAliyun" component={SyncAliyunContent}/>
                    <Route path="/admin/permission/users" component={UserManager}/>
                    <Route path="/admin/permission/roles" component={RolesManager}/>
                    <Route path="/admin/permission/permissions" component={PermissionsManager}/>
                    <Route path="/admin/permission/password" component={PasswordManager}/>
                    <Route path="/admin/system/setting" component={SettingsContent}/>
                    <Route path="/admin/system/user_feedback" component={UserFeedbackManager}/>
                    <Route path="/admin/logout" component={LogoutContent}/>
                </Switch>
            </Content>
        )
    }

}

export default ContentLayout;
