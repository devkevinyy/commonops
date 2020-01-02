import React, { Component } from 'react';
import {Layout} from 'antd';
import { BrowserRouter as Router, Route, withRouter } from 'react-router-dom';
import OpsHeader from '../components/header';
import OpsSider from '../components/sidebar';
import OpsFooter from '../components/footer';
import ContentLayout from "./ContentLayout";
import {Steps} from "intro.js-react";

import 'intro.js/introjs.css';
import {isShowIntroPage} from "../services/common";

class AdminContent extends Component {

    constructor(props) {
        super(props);
        this.handleSiderMenu = this.handleSiderMenu.bind(this);
        this.state = {
            collapsed: false,
            stepsEnabled: false,
            initialStep: 0,
            options: {
                prevLabel: "上一步",
                nextLabel: "下一步",
                skipLabel: "跳过",
                doneLabel: "结束",
                hintButtonLabel: "开始使用"
            },
            steps: [
                {
                    element: '.introStep2',
                    intro: '用户可以点击这里进行bug&需求的反馈，系统管理员可在系统管理->用户反馈中进行查看！',
                },
                {
                    element: '.introStep4',
                    intro: '在这里进入各细分功能列表，完成相应的操作；管理员可以在权限管理中控制用户访问菜单的权限！',
                },
            ]
        };
    }

    componentDidMount() {
        let show = isShowIntroPage();
        if(show===1){
            this.setState({
                stepsEnabled: true
            })
        }
    }

    handleSiderMenu() {
        this.setState({collapsed: !this.state.collapsed});
    };

    onExit = () => {
        this.setState(() => ({ stepsEnabled: false }));
    };

    render() {
        const { stepsEnabled, steps, initialStep, options } = this.state;

        return (
            <Router>
                <Layout style={{ minHeight: '100vh' }}>
                    <Steps
                        enabled={stepsEnabled}
                        steps={steps}
                        initialStep={initialStep}
                        onExit={this.onExit}
                        options={options}
                    />
                    <OpsSider menuCollapsed={this.state.collapsed} />
                    <Layout>
                        <OpsHeader menuCollapsed={this.state.collapsed} handleSiderMenu={this.handleSiderMenu}/>
                        <Route path="/admin" component={ContentLayout}/>
                        <OpsFooter/>
                    </Layout>
                </Layout>
            </Router>
        );
    }
}


export default withRouter(AdminContent);

