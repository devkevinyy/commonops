import React, { Component } from "react";
import { Button, Card, Icon, Layout, Alert } from "antd";
import { Terminal } from "xterm";
import { AttachAddon } from "xterm-addon-attach";
import "xterm/css/xterm.css";
import ReconnectingWebSocket from "reconnecting-websocket";
import { WSBase } from "../../config";

const { Content } = Layout;
let attachAddon;

class WebTerminalContent extends Component {
    constructor(props) {
        super(props);
        this.goBack = this.goBack.bind(this);
        this.rws = new ReconnectingWebSocket(
            WSBase + "ws/kubernetes/container_terminal",
        );
        this.terminal = new Terminal({
            rows: 36,
            fontSize: 14,
            cursorBlink: true,
            cursorStyle: "block",
            bellStyle: "sound",
            theme: "Console",
        });
        this.terminal.prompt = () => {
            this.terminal.write("\r\n$ ");
        };
        this.state = {
            containerDetail: this.props.location.state,
            connectionStatus: 0,
        };
    }

    componentDidMount() {
        this.terminal.open(document.getElementById("terminal"));
        this.terminal.writeln("Welcome to use Web Terminal.");
        this.terminal.prompt();
        this.initWsConnection();
        attachAddon = new AttachAddon(this.rws);
        this.terminal.loadAddon(attachAddon);
        this.terminal.focus();
    }

    componentWillUnmount() {
        this.rws.send(
            JSON.stringify({
                action: "connection_close",
                clusterId: localStorage.getItem("clusterId"),
                namespace: this.state.containerDetail.namespace,
                podName: this.state.containerDetail.podName,
                containerName: this.state.containerDetail.containerName,
            }),
        );
        this.rws.close();
        this.terminal.dispose();
    }

    initWsConnection() {
        let that = this;
        this.rws.addEventListener("open", () => {
            that.rws.send(
                JSON.stringify({
                    action: "init_connection",
                    clusterId: localStorage.getItem("clusterId"),
                    namespace: that.state.containerDetail.namespace,
                    podName: that.state.containerDetail.podName,
                    containerName: that.state.containerDetail.containerName,
                }),
            );
            that.setState({ connectionStatus: 1 });
        });

        this.rws.addEventListener("close", () => {
            that.setState({ connectionStatus: 0 });
        });

        this.rws.addEventListener("error", () => {
            that.setState({ connectionStatus: 0 });
        });
    }

    goBack() {
        this.rws.send(
            JSON.stringify({
                action: "connection_close",
                clusterId: localStorage.getItem("clusterId"),
                namespace: this.state.containerDetail.namespace,
                podName: this.state.containerDetail.podName,
                containerName: this.state.containerDetail.containerName,
            }),
        );
        this.rws.close();
        this.props.history.goBack();
    }

    render() {
        let connectContent = (
            <Alert message="远程连接未建立" type="error" showIcon />
        );
        if (this.state.connectionStatus === 1) {
            connectContent = (
                <Alert message="远程连接成功" type="success" showIcon />
            );
        }
        return (
            <Content
                style={{
                    background: "#fff",
                    padding: "5px 20px",
                    margin: 0,
                    height: "100%",
                }}
            >
                <Button type="link" onClick={this.goBack}>
                    <Icon type="left" />
                    返回上级
                </Button>
                {connectContent}
                <Card
                    size="small"
                    title="终端"
                    style={{ width: "100%" }}
                    bodyStyle={{ padding: 0 }}
                >
                    <div id="terminal" style={{ width: "100%" }} />
                </Card>
            </Content>
        );
    }
}

export default WebTerminalContent;
