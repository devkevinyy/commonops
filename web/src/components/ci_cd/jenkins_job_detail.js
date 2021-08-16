import React, { Component } from "react";
import {
    Layout,
    message,
    Descriptions,
    List,
    Tag,
    Typography,
    Button,
    Drawer,
    Icon,
    Row,
    Col,
} from "antd";
import OpsBreadcrumbPath from "../breadcrumb_path";
import {
    getJenkinsJobBuildList,
    getJenkinsJobBuildLog,
    postJenkinsDisableJob,
    postJenkinsEnableJob,
    postJenkinsStartJob,
} from "../../api/jenkins";
import moment from "moment";
import Spin from "antd/es/spin";

const { Content } = Layout;
const { Text } = Typography;

class Jenkins_job_detail extends Component {
    constructor(props) {
        super(props);
        this.onDrawerClose = this.onDrawerClose.bind(this);
        this.startBuildJob = this.startBuildJob.bind(this);
        this.refreshJobInfo = this.refreshJobInfo.bind(this);
        this.disableJob = this.disableJob.bind(this);
        this.enableJob = this.enableJob.bind(this);
        this.loadJenkinsJobBuildList = this.loadJenkinsJobBuildList.bind(this);
        this.state = {
            buildInfo: [],
            buildConsoleLog: "",
            logLoading: false,
            jobBuilding: false,
            refreshInfoLoading: false,
            buildLog: [],
            start: 0,
            timer: null,
            progressLoading: "none",
            nextBuildNumber: 0,
        };
    }

    componentDidMount() {
        this.loadJenkinsJobBuildList({
            name: this.props.location.state.jobName,
        });
        this.setState({
            buildName: this.props.location.state.jobName,
        });
    }

    loadJenkinsJobBuildList(params) {
        this.setState({ refreshInfoLoading: true });
        getJenkinsJobBuildList(params)
            .then((res) => {
                if (res.code === 0) {
                    let buildLog = [];
                    if (
                        res.data["buildLog"] !== undefined &&
                        res.data["buildLog"] !== null
                    ) {
                        buildLog = res.data["buildLog"];
                    }
                    this.setState({
                        buildInfo: res.data,
                        buildLog: buildLog,
                        nextBuildNumber: parseInt(res.data["nextBuildNumber"]),
                    });
                } else {
                    message.error(res.msg);
                }
                this.setState({ refreshInfoLoading: false });
            })
            .catch((err) => {
                message.error(err.toLocaleString());
            });
    }

    refreshJobInfo() {
        this.loadJenkinsJobBuildList({
            name: this.props.location.state.jobName,
        });
    }

    onDrawerClose() {
        this.setState({
            drawerVisible: false,
            progressLoading: "none",
            buildConsoleLog: "",
            start: 0,
        });
        clearTimeout(this.timer);
        this.timer = null;
        this.refreshJobInfo();
    }

    showBuildConsoleLog(buildId) {
        this.setState({
            drawerVisible: true,
            logLoading: true,
            jobBuilding: true,
        });
        getJenkinsJobBuildLog({
            name: this.state.buildName,
            buildId: buildId,
            start: parseInt(this.state.start),
        })
            .then((res) => {
                if (res.code === 0) {
                    this.setState({
                        buildConsoleLog:
                            this.state.buildConsoleLog + res.data["content"],
                        logLoading: false,
                        jobBuilding: false,
                        start: res.data["start"],
                    });
                    if (res.data["hasMore"] === "true") {
                        this.setState({ progressLoading: "inline-block" });
                        this.timer = setTimeout(() => {
                            this.showBuildConsoleLog(buildId);
                        }, 2000);
                    } else {
                        this.setState({ progressLoading: "none" });
                        clearTimeout(this.timer);
                        this.timer = null;
                        this.refreshJobInfo();
                    }
                } else {
                    message.error(res.msg);
                }
            })
            .catch((err) => {
                message.error(err.toLocaleString());
            });
    }

    startBuildJob() {
        this.setState({ jobBuilding: true });
        postJenkinsStartJob({
            name: this.state.buildName,
        })
            .then((res) => {
                if (res.code === 0) {
                    setTimeout(() => {
                        this.setState({ jobBuilding: false });
                        this.refreshJobInfo();
                    }, 4000);
                } else {
                    message.error(res.msg);
                }
            })
            .catch((err) => {
                message.error(err.toLocaleString());
            });
    }

    disableJob() {
        postJenkinsDisableJob({ name: this.state.buildName })
            .then((res) => {
                if (res.code === 0) {
                    message.success("禁用成功");
                    this.refreshJobInfo();
                } else {
                    message.error(res.msg);
                }
            })
            .catch((err) => {
                message.error(err.toLocaleString());
            });
    }

    enableJob() {
        postJenkinsEnableJob({ name: this.state.buildName })
            .then((res) => {
                if (res.code === 0) {
                    message.success("启用成功");
                    this.refreshJobInfo();
                } else {
                    message.error(res.msg);
                }
            })
            .catch((err) => {
                message.error(err.toLocaleString());
            });
    }

    render() {
        let titleContent = (
            <div>
                构建日志
                <Icon
                    type="loading"
                    style={{
                        marginLeft: "20px",
                        display: this.state.progressLoading,
                    }}
                />
            </div>
        );
        let ableBuild = true;
        let ableContent = (
            <Button
                size="small"
                onClick={this.disableJob}
                style={{ marginLeft: "20px" }}
            >
                禁 用
            </Button>
        );
        if (this.state.buildInfo["buildable"] === false) {
            ableBuild = false;
            ableContent = (
                <Button
                    size="small"
                    onClick={this.enableJob}
                    style={{ marginLeft: "20px" }}
                >
                    启 用
                </Button>
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
                <OpsBreadcrumbPath pathData={["构建部署", "任务详情"]} />
                <Row style={{ padding: "0px 0px 10px 0px" }}>
                    <Col span={2}>
                        <Button type="primary" onClick={this.refreshJobInfo}>
                            刷 新
                        </Button>
                    </Col>
                </Row>

                <Spin spinning={this.state.refreshInfoLoading}>
                    <Descriptions bordered size="small" column={2}>
                        <Descriptions.Item label="任务名">
                            {this.state.buildInfo.displayName}
                        </Descriptions.Item>
                        <Descriptions.Item label="任务描述">
                            {this.state.buildInfo.description}
                        </Descriptions.Item>
                        <Descriptions.Item label="任务操作">
                            <Button
                                size="small"
                                disabled={!ableBuild}
                                loading={this.state.jobBuilding}
                                onClick={this.startBuildJob}
                            >
                                构建
                            </Button>
                            {ableContent}
                        </Descriptions.Item>
                    </Descriptions>
                </Spin>

                <Drawer
                    title={titleContent}
                    placement="left"
                    closable={false}
                    onClose={this.onDrawerClose}
                    visible={this.state.drawerVisible}
                    width={700}
                >
                    <Spin spinning={this.state.logLoading}>
                        <Text style={{ whiteSpace: "pre" }}>
                            {this.state.buildConsoleLog}
                        </Text>
                    </Spin>
                </Drawer>

                <List
                    header={<div>近5次构建状态</div>}
                    footer={null}
                    bordered
                    dataSource={this.state.buildLog}
                    style={{ marginTop: "10px", width: "50%" }}
                    renderItem={(item) => {
                        let data = JSON.parse(item);
                        let buildTag = <Tag color="#2db7f5">构建中</Tag>;
                        if (data.result !== undefined && data.result !== null) {
                            if (data.result === "SUCCESS") {
                                buildTag = <Tag color="#108ee9">成 功</Tag>;
                            } else if (data.result === "FAILURE") {
                                buildTag = <Tag color="#f50">失 败</Tag>;
                            } else if (data.result === "ABORTED") {
                                buildTag = <Tag color="#f50">中 止</Tag>;
                            } else {
                                buildTag = (
                                    <Tag color="magenta">{data.result}</Tag>
                                );
                            }
                        }

                        let timeTag = (
                            <Text type="secondary">
                                {moment(data.timestamp).format(
                                    "MM-DD HH:mm:ss",
                                )}
                            </Text>
                        );
                        let actionContent = (
                            <div
                                style={{
                                    display: "inline-block",
                                    float: "right",
                                }}
                            >
                                <Button
                                    type="link"
                                    size="small"
                                    onClick={this.showBuildConsoleLog.bind(
                                        this,
                                        data.id,
                                    )}
                                >
                                    构建日志
                                </Button>
                            </div>
                        );
                        return (
                            <List.Item>
                                {buildTag} - {data.fullDisplayName} - {timeTag}
                                {actionContent}
                            </List.Item>
                        );
                    }}
                />
            </Content>
        );
    }
}

export default Jenkins_job_detail;
