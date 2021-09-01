import React, { Component } from "react";
import {
  Button,
  Col,
  Divider,
  Layout,
  message,
  Row,
  Spin,
  Tree,
  Typography,
} from "antd";
import { LeftOutlined, RightOutlined } from '@ant-design/icons';
import OpsBreadcrumbPath from "../breadcrumb_path";

import { Controlled as CodeMirror } from "react-codemirror2";
import "codemirror/lib/codemirror.css";
import "codemirror/mode/shell/shell";

import { getCurrentUserId } from "../../services/common";
import { getAllNodeTree } from "../../api/batch";
import { postCloudServerBatchSSH } from "../../api/cloud";

const { Content } = Layout;
const { Title } = Typography;

class CmdsContent extends Component {
  constructor(props) {
    super(props);
    this.formRef = React.createRef();
    this.state = {
      collapsed: false,
      nodeTreeData: [],
      span_left: 5,
      span_right: 17,
      cmdInput: "",
      cmdResult: [],
      cmdRunning: false,
      selectedIds: [],
    };
  }

  // 刷新节点树
  refreshNodeTree = () => {
    const userId = getCurrentUserId();
    getAllNodeTree(userId)
    .then((res) => {
      if (res.code === 0) {
        this.setState({
          nodeTreeData: res.data
        });
      } else {
        message.error(res.msg);
      }
    });
  }

  // 折叠节点数
  toggleCollapsed = () => {
    if (this.state.collapsed) {
      this.setState({
        collapsed: !this.state.collapsed,
        span_left: 5,
        span_right: 17,
      });
    } else {
      this.setState({
        collapsed: !this.state.collapsed,
        span_left: 0,
        span_right: 22,
      });
    }
  };

  // 选择主机
  onSelectServerNode = (selectedKeys, info) => {
    this.setState({
      selectedIds: selectedKeys
    });
  }


  // 执行命令
  submitBatchSsh = () => {
    this.setState({
      cmdRunning: true
    });
    // 由于选择父节点全选主机，key会包含父节点，并且子节点key这里是字符串了，需要处理下。
    // ['p1', '1', '2'] -> [1, 2]
    const ids = this.state.selectedIds.map((id) => parseInt(id)).filter((id) => !isNaN(id));
    postCloudServerBatchSSH({
      ids: ids,
      command: this.state.cmdInput,
    }).then((res) => {
      if (res.code === 0) {
        this.setState({
          cmdResult: res.data
        });
      } else {
        message.error(res.msg);
      }
      this.setState({
        cmdRunning: false
      });
    });
  };

  componentDidMount() {
    this.refreshNodeTree();
  }

  render() {
    return (
      <Content
        style={{
          background: "#fff",
          padding: 20,
          margin: 0,
          height: "100%",
        }}
      >
        <OpsBreadcrumbPath pathData={["批量执行", "远程命令"]} />
        <Divider />
        <Row gutter={1}>
          <Col span={this.state.span_left}>
            <div style={{ backgroundColor: "#F7F9F9" }}>
            <Title level={5}>选择云主机</Title>
            <Tree
              style={{ backgroundColor: "#F7F9F9" }}
              checkable
              onCheck={this.onSelectServerNode}
              treeData={this.state.nodeTreeData}
            />
            </div>
          </Col>
          <Col>
            <Button type="primary" size="small" onClick={this.toggleCollapsed}>
              {React.createElement(this.state.collapsed ? RightOutlined : LeftOutlined)}
            </Button>
          </Col>
          <Col span={this.state.span_right}>
              <Title level={5}>执行结果</Title>
              <Spin
                tip="远程命令执行中..."
                spinning={this.state.cmdRunning}
              >
                <pre
                  style={{
                    marginTop: 10,
                    minHeight: 500,
                    paddingLeft: 10,
                  }}
                  className="preJenkinsLog"
                >
                  {this.state.cmdResult.map((item) => {
                    let serverName = "";
                    let result = "";
                    for (var server in item) {
                      serverName = server;
                      result = item[server];
                      break;
                    }
                    return (
                      <div style={{ marginBottom: 10 }}>
                        <div>目标机器: {serverName}</div>
                        <div>
                          执行结果：<br/>
                          {result}
                        </div>
                        <div><br/>---------------------<br/></div>
                      </div>
                    );
                  })}
                </pre>
              </Spin>
              <Divider />
              <Row align="top">
                <Col flex="8">
                  <div style={{ marginBottom: 10, marginTop: 10 }}>
                    <CodeMirror
                      className="sqlEditor"
                      options={{
                        showCursorWhenSelecting: true,
                        option: {
                          autofocus: true,
                        },
                        lineWrapping: true,
                        mode: "shell",
                      }}
                      value={this.state.cmdInput}
                      onBeforeChange={(editor, data, value) => {
                        this.setState({ cmdInput: value });
                      }}
                    />
                  </div>
                </Col>
                <Col flex="auto">
                  <div style={{ marginBottom: 10, marginTop: 10, textAlign: "right" }}>
                    <Button
                      type="primary"
                      disabled={
                        !this.props.aclAuthMap["POST:/cloud/servers/batch/ssh"]
                      }
                      onClick={this.submitBatchSsh}
                    >
                      执行命令
                    </Button>
                  </div>
                </Col>
              </Row>
          </Col>
        </Row>
      </Content>
    );
  }
}

export default CmdsContent;
