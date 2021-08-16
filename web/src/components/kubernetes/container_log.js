import React, { Component } from "react";
import { Button, Icon, Layout, message, Row, Spin } from "antd";
import { ReloadOutlined } from "@ant-design/icons";
import { getPodLogs } from "../../api/kubernetes";

const { Content } = Layout;

class ContainerLogContent extends Component {
    constructor(props) {
        super(props);
        this.goBack = this.goBack.bind(this);
        this.syncLogs = this.syncLogs.bind(this);
        this.state = {
            logLoading: false,
            containerInfo: this.props.location.state,
        };
    }

    componentDidMount() {
        this.syncLogs();
    }

    goBack() {
        this.setState({ logContent: "" }, () => {
            this.props.history.goBack();
        });
    }

    syncLogs() {
        this.setState({ logLoading: true });
        getPodLogs({
            namespace: this.state.containerInfo.namespace,
            podName: this.state.containerInfo.podName,
            containerName: this.state.containerInfo.containerName,
        })
            .then((res) => {
                if (res.code === 0) {
                    this.setState(
                        { logContent: res.data, logLoading: false },
                        () => {
                            this.scrollToBottom();
                        },
                    );
                } else {
                    message.error(res.msg);
                }
            })
            .catch((err) => {
                message.error(err.toLocaleString());
            });
    }

    scrollToBottom() {
        let panel = document.getElementById("logPanel");
        panel.scrollTop = panel.scrollHeight;
    }

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
                <Row style={{ padding: "0px 0px 10px 0px" }}>
                    <Button type="link" onClick={this.goBack}>
                        <Icon type="left" />
                        返回上级
                    </Button>
                    <Button
                        type="link"
                        icon={<ReloadOutlined />}
                        onClick={this.syncLogs}
                        style={{ float: "right" }}
                    >
                        刷 新
                    </Button>
                </Row>
                <Spin spinning={this.state.logLoading}>
                    <div style={{ height: "85vh" }}>
                        <div
                            id="logPanel"
                            style={{
                                whiteSpace: "pre-line",
                                backgroundColor: "rgb(19, 19, 19)",
                                color: "#fff",
                                fontSize: "14px",
                                lineHeight: "20px",
                                height: "100%",
                                padding: "10px",
                                overflow: "scroll",
                            }}
                        >
                            {this.state.logContent}
                        </div>
                    </div>
                </Spin>
            </Content>
        );
    }
}

export default ContainerLogContent;
