import React, { Component } from 'react';
import { Layout } from 'antd/lib/index';
import { Route, Switch } from 'react-router-dom';
import Deploy_project_content from "./deploy_project_content";
import Jobs_content from "./jobs_content";

const { Content } = Layout;

class Task_content extends Component {
    render() {
        return (
            <Content style={{
                background: '#fff', padding: 20, margin: 0, height: "100%",
            }}>
                <Switch>
                    <Route path="/admin/task/deploy_project" component={Deploy_project_content} />
                    <Route path="/admin/task/jobs" component={Jobs_content}/>
                </Switch>
            </Content>
        );
    }
}

export default Task_content;
