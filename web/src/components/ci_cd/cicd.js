import React, { Component } from "react";
import { Layout } from "antd";
import { Route, Switch } from "react-router-dom";
import CiContent from "./ci_content";
import CdContent from "./cd_content";
import CdRecordContent from "./cd_record_content";

const { Content } = Layout;

class CiCdContent extends Component {

    render() {
        return (
            <Content
                style={{
                    background: "#fff",
                    padding: 0,
                    margin: 0,
                    height: "100%",
                }}
            >
                <Switch>
                    <Route path="/admin/cicd/ci">
                        <CiContent aclAuthMap={this.props.aclAuthMap}/>
                    </Route>
                    <Route path="/admin/cicd/cd">
                        <CdContent aclAuthMap={this.props.aclAuthMap}/>
                    </Route>
                    <Route path="/admin/cicd/cd_record">
                        <CdRecordContent aclAuthMap={this.props.aclAuthMap}/>
                    </Route>
                </Switch>
            </Content>
        );
    }
}

export default CiCdContent;
