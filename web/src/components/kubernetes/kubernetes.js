import React, {Component} from 'react';
import {Layout} from "antd";
import {Route, Switch} from "react-router-dom";
import K8sNamespacesContent from "./namespaces";
import ClusterManageContent from "./k8s_cluster";

const {Content} = Layout;

class KubernetesContent extends Component {

    render() {
        return (
            <Content style={{
                background: '#fff', padding: 0, margin: 0, height: "100%",
            }}>
                <Switch>
                    <Route path="/admin/k8s_cluster/info" component={ClusterManageContent} />
                    <Route path="/admin/k8s_cluster/manage" component={K8sNamespacesContent} />
                </Switch>
            </Content>
        )
    }
}

export default KubernetesContent;
