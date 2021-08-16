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
                    <Route path="/admin/k8s_cluster/info">
                        <ClusterManageContent {...this.props}/>
                    </Route>
                    <Route path="/admin/k8s_cluster/manage">
                        <K8sNamespacesContent {...this.props}/>
                    </Route>
                </Switch>
            </Content>
        )
    }
}

export default KubernetesContent;
