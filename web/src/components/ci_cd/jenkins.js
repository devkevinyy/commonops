import React, {Component} from 'react';
import {Layout} from "antd";
import {Route, Switch} from "react-router-dom";
import JenkinsJobsContent from "./jenkins_jobs";
import Jenkins_job_detail from "./jenkins_job_detail";


const {Content} = Layout;

class JenkinsContent extends Component {

    render() {
        return (
            <Content style={{
                background: '#fff', padding: 0, margin: 0, height: "100%",
            }}>
                <Switch>
                    <Route path="/admin/cicd/jobs" component={JenkinsJobsContent} />
                    <Route path="/admin/cicd/job_detail" component={Jenkins_job_detail} />
                </Switch>
            </Content>
        )
    }
}

export default JenkinsContent;