import React, { Component } from 'react';
import { Layout, Popover, Typography } from 'antd';
const { Footer} = Layout;

const { Text } = Typography;

class OpsFooter extends Component {

  constructor(props) {
    super(props);
    this.state = {
      currentYear: new Date().getFullYear(),
    };
  }

  render() {
    return (
      <Footer style={{ textAlign: 'center', padding: '15px 50px' }}>
        ©2019-{this.state.currentYear} Created by &nbsp;
        <Popover content={<div className="wechat"/>} >
            <Text underline>KevinYang</Text>
        </Popover>
      </Footer>
    );
  }
}

export default OpsFooter;
