import React, { Component } from 'react';
import {Layout} from 'antd';
import { BrowserRouter as Router, Route, withRouter } from 'react-router-dom';
import OpsHeader from '../components/header';
import OpsSider from '../components/sidebar';
import OpsFooter from '../components/footer';
import ContentLayout from "./ContentLayout";

class AdminContent extends Component {

    constructor(props) {
        super(props);
        this.handleSiderMenu = this.handleSiderMenu.bind(this);
        this.state = {
            collapsed: false,
        };
    }

    handleSiderMenu() {
        this.setState({collapsed: !this.state.collapsed});
    };

    onExit = () => {
        this.setState(() => ({ stepsEnabled: false }));
    };

    render() {
        return (
            <Router>
                <Layout style={{ minHeight: '100vh' }}>
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

