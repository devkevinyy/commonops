import React, { Component } from 'react';
import { Layout } from 'antd';
const { Footer} = Layout;


class OpsFooter extends Component {
  render() {
    return (
      <Footer style={{ textAlign: 'center', padding: '15px 50px' }}>
        ©2019 Created by KevinYang
      </Footer>
    );
  }
}

export default OpsFooter;
