import React, { Component } from "react";
import { Layout } from "antd";
import { Route, Switch } from "react-router-dom";
import OpsBreadcrumbPath from "../breadcrumb_path";
import { JobsManageContent, BuildsManageContent } from "./jenkins_jobs";

const { Content } = Layout;

class CiContent extends Component {

    render() {
        return (
            <Content
                style={{
                    background: "#fff",
                    padding: "5px 20px",
                    margin: 0,
                    height: "100%",
                }}
            >
                <OpsBreadcrumbPath pathData={["集成部署", "持续构建"]} />
                <Switch>
                    <Route
                        exact={true}
                        path="/admin/cicd/ci/"
                    >
                        <JobsManageContent {...this.props}/>
                    </Route>
                    <Route
                        path="/admin/cicd/ci/jobs"
                    >
                        <BuildsManageContent {...this.props}/>
                    </Route>
                </Switch>
            </Content>
        );
    }
}

export default CiContent;
