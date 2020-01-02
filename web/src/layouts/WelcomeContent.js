import React, { Component } from 'react';
import IntroduceContent from "../components/upgrade_introduce/introduce_content";
import {Layout} from "antd";

const { Content } = Layout;

class WelcomeContent extends Component {
  render() {
    return (
      <Content style={{ height: '100%' }}>
        <IntroduceContent />
        <div className="base_png" />
      </Content>
    );
  }
}

export default WelcomeContent;
