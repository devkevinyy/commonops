import React, { Component, Fragment } from "react";
import {
    Card,
    Layout,
    Row,
    message,
    Button,
    Col,
    Tooltip,
    Table,
    Divider,
    Typography,
    Drawer,
    Form,
    Input,
    Select,
    Modal,
    Menu,
    Dropdown,
    Collapse,
    Popconfirm,
    Empty,
} from "antd";
import moment from "moment";
import { Controlled as CodeMirror } from "react-codemirror2";
import "codemirror/lib/codemirror.css";
import {
    LeftCircleOutlined,
    CloseCircleOutlined,
    ClockCircleOutlined,
    CheckCircleOutlined,
    PlusOutlined,
    InfoCircleOutlined,
    SettingOutlined,
} from "@ant-design/icons";
import { withRouter } from 'react-router-dom';
import {
    getJobBuildList,
    getJenkinsAllJobs,
    postJenkinsStartJob,
    getJenkinsJobBuildLog,
    getJenkinsCredentialsList,
    postJenkinsAddCredential,
    postJenkinsJob,
    getJenkinsJob,
    putJenkinsJob,
    deleteJenkinsJob,
    deleteJenkinsJobBuildLog,
    getJenkinsJobBuildStages,
    getJenkinsJobBuildStageLog,
    getJenkinsJobBuildArchiveArtifactsInfo,
} from "../../api/jenkins";

const { Content } = Layout;
const { Option } = Select;
const { Text, Paragraph } = Typography;
const { Panel } = Collapse;
const { confirm } = Modal;

class JobsManageContent extends Component {
    constructor(props) {
        super(props);
        this.formRef = React.createRef();
        this.addCredentialFormRef = React.createRef();
        this.state = {
            jobs: [],
            createJobDrawerVisible: false,
            addCredentialModalVisible: false,
            updateJobConfigDrawerVisible: false,
            jobLoading: false,
            credentialsList: [],
            env: {
                GIT_REPO: "",
                GIT_BRANCH: "",
                GIT_CREDENTIALSId: "",
                DOCKER_IMAGE_REPO: "",
                DOCKER_IAMGE_CREDENTIALSId: "",
                DOCKER_IMAGE_NAME: "",
            },
            currentJobName: "",
            currentJobConfigXml: "",
        };
    }

    componentDidMount() {
        this.loadJenkinsAllJobs();
        this.loadAllJenkinsCredentialsList();
    }

    loadJenkinsAllJobs() {
        this.setState({ jobLoading: true });
        getJenkinsAllJobs()
            .then((res) => {
                if (res.code === 0) {
                    this.setState({ jobs: res.data["jobs"] });
                } else {
                    message.error(res.msg);
                }
                this.setState({ jobLoading: false });
            })
    }

    loadAllJenkinsCredentialsList() {
        getJenkinsCredentialsList().then((res) => {
            if (res.code === 0) {
                this.setState({ credentialsList: res.data["credentials"] });
            } else {
                message.error(res.msg);
            }
        });
    }

    jumpJobDetail(jobName) {
        this.props.history.push({
            pathname: "/admin/cicd/ci/jobs",
            state: { jobName: jobName },
        });
    }

    editJobConfig(jobName) {
        this.setState({ updateJobConfigDrawerVisible: true });
        getJenkinsJob({ jobName: jobName }).then((res) => {
            if (res.code === 0) {
                let script = res.data.match(
                    /<script\b[^>]*>([\s\S]*)<\/script>/,
                );
                let temp = document.createElement("div");
                temp.innerHTML = script[1];
                let scriptContent = temp.innerText || temp.textContent;
                this.setState({
                    currentJobName: jobName,
                    currentJobConfigXml: scriptContent,
                });
            } else {
                message.error(res.msg);
            }
        });
    }

    deleteJobItem(jobName) {
        deleteJenkinsJob({ jobName: jobName }).then((res) => {
            if (res.code === 0) {
                message.success("删除成功!");
                this.loadJenkinsAllJobs();
            } else {
                message.error(res.msg);
            }
        });
    }

    onCloseUpdateJobDrawer() {
        this.setState({
            currentJobName: "",
            currentJobConfigXml: "",
            updateJobConfigDrawerVisible: false,
        });
    }

    updateJenkinsFile() {
        putJenkinsJob({
            jobName: this.state.currentJobName,
            pipelineScript: this.state.currentJobConfigXml,
        }).then((res) => {
            if (res.code === 0) {
                message.success("更新成功!");
                this.onCloseUpdateJobDrawer();
            } else {
                message.error(res.msg);
            }
        });
    }

    onCloseCreateJobDrawer() {
        this.setState({ createJobDrawerVisible: false });
    }

    createNewBuildJob() {
        this.setState({ createJobDrawerVisible: true });
    }

    addNewCredential() {
        this.setState({ addCredentialModalVisible: true });
    }

    handleAddCredentialOk() {
        this.addCredentialFormRef.current.validateFields().then((values) => {
            postJenkinsAddCredential(values).then((res) => {
                if (res.code === 0) {
                    message.success("创建成功");
                    this.loadAllJenkinsCredentialsList();
                } else {
                    message.error(res.msg);
                }
                this.setState({ addCredentialModalVisible: false });
            });
        });
    }

    handleAddCredentialCancel() {
        this.setState({ addCredentialModalVisible: false });
    }

    handelAddNewJob() {
        this.formRef.current.validateFields().then((values) => {
            postJenkinsJob(values).then((res) => {
                if (res.code === 0) {
                    message.success("创建成功");
                    this.loadJenkinsAllJobs();
                } else {
                    message.error(res.msg);
                }
                this.setState({ createJobDrawerVisible: false });
            });
        });
    }

    onChangeGitRepo(e) {
        this.setState({
            env: {
                ...this.state.env,
                GIT_REPO: e.target.value.trim(),
            },
        });
    }

    onChangeGitBranch(e) {
        this.setState({
            env: {
                ...this.state.env,
                GIT_BRANCH: e.target.value.trim(),
            },
        });
    }

    onChangeGitCredentials(e) {
        this.setState({
            env: {
                ...this.state.env,
                GIT_CREDENTIALSId: e.trim(),
            },
        });
    }

    onChangeDockerImageRepo(e) {
        this.setState({
            env: {
                ...this.state.env,
                DOCKER_IMAGE_REPO: e.target.value.trim(),
            },
        });
    }

    onChangeDockerImageCredentials(e) {
        this.setState({
            env: {
                ...this.state.env,
                DOCKER_IAMGE_CREDENTIALSId: e.trim(),
            },
        });
    }

    onChangeDockerImageName(e) {
        this.setState({
            env: {
                ...this.state.env,
                DOCKER_IMAGE_NAME: e.target.value.trim(),
            },
        });
    }

    render() {
        let env = this.state.env;
        let configTemplate = String.raw`
pipeline {
  agent any
  tools {
    maven &apos;maven_3.3.3&apos;
  }
  environment {
    GIT_REPO = <span class="flashing">"${env.GIT_REPO}"</span>
    BRANCH = <span class="flashing">"${env.GIT_BRANCH}"</span>
    GIT_CREDENTIALSId = <span class="flashing">"${env.GIT_CREDENTIALSId}"</span>
    DOCKER_IMAGE_REPO = <span class="flashing">"${env.DOCKER_IMAGE_REPO}"</span>
    DOCKER_IAMGE_CREDENTIALSId = <span class="flashing">"${env.DOCKER_IAMGE_CREDENTIALSId}"</span>
    IMAGE_NAME = <span class="flashing">"${env.DOCKER_IMAGE_NAME}"</span>
  }
  stages {
    stage(&quot;检出&quot;) {
      steps {
        checkout(
          [$class: &apos;GitSCM&apos;,
          userRemoteConfigs: [[
            url: &quot;${env.GIT_REPO}&quot;,
            credentialsId: &quot;${env.GIT_CREDENTIALSId}&quot;
          ]],
          branches: [[name: &quot;${env.GIT_BRANCH}&quot;]],
          ]
        )
      }
    }

    stage(&apos;Maven Build&apos;) {
        steps {
            sh &apos;cd test; mvn -B -DskipTests clean package&apos;
        }
    }

    stage(&apos;构建镜像并推送到镜像库&apos;) {
      steps {
        script {
          sh &apos;echo &quot;build and push&quot;&apos;
          docker.withRegistry(
            &quot;${env.DOCKER_IMAGE_REPO}&quot;,
            &quot;${env.DOCKER_IAMGE_CREDENTIALSId}&quot;
          ) {
            def dockerImage = docker.build(&quot;${env.DOCKER_IMAGE_NAME}:<BUILD_NUMBER>&quot;, &quot;-f Dockerfile .&quot;)
            dockerImage.push()
          }
        }
      }
    }
  }
}`;
        const formItemLayout = {
            labelCol: { span: 24 },
            wrapperCol: { span: 24 },
        };
        const formHorizontalItemLayout = {
            labelCol: { span: 6 },
            wrapperCol: { span: 16 },
        };
        return (
            <Content
                style={{
                    background: "#fff",
                    margin: 0,
                }}
            >
                <Row>
                    <Col span={24}>
                        <Card
                            title="当前任务"
                            loading={this.state.jobLoading}
                            size="small"
                            extra={
                                <Fragment>
                                    <Button
                                        size="small"
                                        onClick={this.loadJenkinsAllJobs.bind(
                                            this,
                                        )}
                                    >
                                        刷新
                                    </Button>
                                    <Divider type="vertical" />
                                    <Button
                                        type="primary"
                                        size="small"
                                        onClick={this.createNewBuildJob.bind(
                                            this,
                                        )}
                                        disabled={!this.props.aclAuthMap["POST:/ci/job"]}
                                    >
                                        新建计划
                                    </Button>
                                </Fragment>
                            }
                        >
                            {this.state.jobs.map((item, index) => {
                                return (
                                    <Card.Grid
                                        key={index}
                                        style={{
                                            width: "25%",
                                            textAlign: "center",
                                            padding: 10,
                                        }}
                                    >
                                        <Row>
                                            <Col span={8} offset={8}>
                                                <Tooltip
                                                    placement="top"
                                                    title={"点击查看详情"}
                                                >
                                                    <Button
                                                        type="link"
                                                        onClick={this.jumpJobDetail.bind(
                                                            this,
                                                            item.name,
                                                        )}
                                                    >
                                                        {item.name}
                                                    </Button>
                                                </Tooltip>
                                            </Col>
                                            <Col span={4} offset={4}>
                                                <div
                                                    style={{
                                                        float: "right",
                                                        paddingTop: 5,
                                                    }}
                                                >
                                                    <Dropdown
                                                        overlay={
                                                            <Menu>
                                                                <Menu.Item>
                                                                    <Button
                                                                        type="text"
                                                                        size="small"
                                                                        onClick={this.editJobConfig.bind(
                                                                            this,
                                                                            item.name,
                                                                        )}
                                                                    >
                                                                        编辑任务
                                                                    </Button>
                                                                </Menu.Item>
                                                                <Menu.Item
                                                                    danger
                                                                >
                                                                    <Popconfirm
                                                                        title="删除后不可恢复! 确认删除吗?"
                                                                        onConfirm={this.deleteJobItem.bind(
                                                                            this,
                                                                            item.name,
                                                                        )}
                                                                        okText="确定"
                                                                        cancelText="取消"
                                                                    >
                                                                        <Button
                                                                            type="text"
                                                                            size="small"
                                                                        >
                                                                            删除任务
                                                                        </Button>
                                                                    </Popconfirm>
                                                                </Menu.Item>
                                                            </Menu>
                                                        }
                                                    >
                                                        <a
                                                            className="ant-dropdown-link"
                                                            onClick={(e) =>
                                                                e.preventDefault()
                                                            }
                                                        >
                                                            <SettingOutlined />
                                                        </a>
                                                    </Dropdown>
                                                </div>
                                            </Col>
                                        </Row>
                                    </Card.Grid>
                                );
                            })}
                        </Card>
                    </Col>
                </Row>

                <Drawer
                    title="新建构建计划"
                    placement="left"
                    width={800}
                    destroyOnClose={true}
                    closable={true}
                    onClose={this.onCloseCreateJobDrawer.bind(this)}
                    visible={this.state.createJobDrawerVisible}
                >
                    <Row>
                        <Col span={8}>
                            <Form
                                {...formItemLayout}
                                layout="vertical"
                                ref={this.formRef}
                                initialValues={{}}
                            >
                                <Form.Item
                                    label="构建计划名称"
                                    name="jobName"
                                    rules={[
                                        {
                                            required: true,
                                            message: "输入构建计划名称!",
                                        },
                                    ]}
                                >
                                    <Input placeholder="输入构建计划名称" />
                                </Form.Item>
                                <Form.Item
                                    label="Git项目地址"
                                    name="gitRepo"
                                    rules={[
                                        {
                                            required: true,
                                            message: "输入代码仓库地址!",
                                        },
                                    ]}
                                >
                                    <Input
                                        placeholder="输入代码仓库地址"
                                        onChange={this.onChangeGitRepo.bind(
                                            this,
                                        )}
                                    />
                                </Form.Item>
                                <Form.Item
                                    label="构建分支"
                                    name="gitBranch"
                                    rules={[
                                        {
                                            required: true,
                                            message: "输入代码构建分支!",
                                        },
                                    ]}
                                >
                                    <Input
                                        placeholder="输入代码构建分支"
                                        onChange={this.onChangeGitBranch.bind(
                                            this,
                                        )}
                                    />
                                </Form.Item>
                                <Form.Item
                                    label="Git访问凭证"
                                    name="gitCredentials"
                                    rules={[
                                        {
                                            required: true,
                                            message: "选择Git访问凭证!",
                                        },
                                    ]}
                                >
                                    <Select
                                        onChange={this.onChangeGitCredentials.bind(
                                            this,
                                        )}
                                        placeholder="选择Git访问凭证"
                                        dropdownRender={(menu) => (
                                            <div>
                                                {menu}
                                                <Divider
                                                    style={{ margin: "4px 0" }}
                                                />
                                                <div
                                                    style={{
                                                        display: "flex",
                                                        flexWrap: "nowrap",
                                                    }}
                                                >
                                                    <a
                                                        style={{
                                                            flex: "none",
                                                            display: "block",
                                                            padding: 5,
                                                            cursor: "pointer",
                                                        }}
                                                        onClick={this.addNewCredential.bind(
                                                            this,
                                                        )}
                                                    >
                                                        <PlusOutlined />{" "}
                                                        新建凭证
                                                    </a>
                                                </div>
                                            </div>
                                        )}
                                    >
                                        {this.state.credentialsList.map(
                                            (item, index) => {
                                                return (
                                                    <Option
                                                        key={index}
                                                        value={item.id}
                                                    >
                                                        <Text
                                                            style={{
                                                                fontSize: 13,
                                                            }}
                                                        >
                                                            {item.displayName}
                                                        </Text>
                                                    </Option>
                                                );
                                            },
                                        )}
                                    </Select>
                                </Form.Item>
                                <Form.Item
                                    label="Docker镜像仓库地址"
                                    name="dockerImageRepo"
                                    rules={[
                                        {
                                            required: true,
                                            message: "输入Docker镜像仓库地址!",
                                        },
                                    ]}
                                >
                                    <Input
                                        placeholder="输入Docker镜像仓库地址"
                                        onChange={this.onChangeDockerImageRepo.bind(
                                            this,
                                        )}
                                    />
                                </Form.Item>
                                <Form.Item
                                    label="Docker镜像名称"
                                    name="dockerImageName"
                                    rules={[
                                        {
                                            required: true,
                                            message: "输入Docker镜像名称!",
                                        },
                                    ]}
                                >
                                    <Input
                                        placeholder="输入Docker镜像名称"
                                        onChange={this.onChangeDockerImageName.bind(
                                            this,
                                        )}
                                    />
                                </Form.Item>
                                <Form.Item
                                    label="Docker镜像仓库访问凭证"
                                    name="dockerImageCredentials"
                                    rules={[
                                        {
                                            required: true,
                                            message:
                                                "选择Docker镜像仓库访问凭证!",
                                        },
                                    ]}
                                >
                                    <Select
                                        onChange={this.onChangeDockerImageCredentials.bind(
                                            this,
                                        )}
                                        placeholder="选择Docker镜像仓库访问凭证"
                                        dropdownRender={(menu) => (
                                            <div>
                                                {menu}
                                                <Divider
                                                    style={{ margin: "4px 0" }}
                                                />
                                                <div
                                                    style={{
                                                        display: "flex",
                                                        flexWrap: "nowrap",
                                                    }}
                                                >
                                                    <a
                                                        style={{
                                                            flex: "none",
                                                            display: "block",
                                                            padding: 5,
                                                            cursor: "pointer",
                                                        }}
                                                        onClick={this.addNewCredential.bind(
                                                            this,
                                                        )}
                                                    >
                                                        <PlusOutlined />{" "}
                                                        新建凭证
                                                    </a>
                                                </div>
                                            </div>
                                        )}
                                    >
                                        {this.state.credentialsList.map(
                                            (item, index) => {
                                                return (
                                                    <Option
                                                        key={index}
                                                        value={item.id}
                                                    >
                                                        <Text
                                                            style={{
                                                                fontSize: 13,
                                                            }}
                                                        >
                                                            {item.displayName}
                                                        </Text>
                                                    </Option>
                                                );
                                            },
                                        )}
                                    </Select>
                                </Form.Item>
                                <Form.Item>
                                    <Button
                                        type="primary"
                                        onClick={this.handelAddNewJob.bind(
                                            this,
                                        )}
                                    >
                                        提交
                                    </Button>
                                </Form.Item>
                            </Form>
                        </Col>
                        <Col span={15}>
                            <pre
                                style={{
                                    color: "#687168",
                                    fontSize: "11px",
                                    paddingLeft: 10,
                                    marginLeft: 20,
                                }}
                            >
                                <code
                                    dangerouslySetInnerHTML={{
                                        __html: configTemplate,
                                    }}
                                ></code>
                            </pre>
                        </Col>
                    </Row>
                </Drawer>

                <Drawer
                    title="配置详情"
                    placement="left"
                    width={800}
                    destroyOnClose={true}
                    closable={true}
                    onClose={this.onCloseUpdateJobDrawer.bind(this)}
                    visible={this.state.updateJobConfigDrawerVisible}
                    bodyStyle={{ paddingBottom: 30 }}
                    footer={
                        <div
                            style={{
                                textAlign: "right",
                            }}
                        >
                            <Button
                                type="primary"
                                onClick={this.updateJenkinsFile.bind(this)}
                            >
                                更新
                            </Button>
                        </div>
                    }
                >
                    <Row>
                        <Col span={24}>
                            <CodeMirror
                                className="jenkinsEditor"
                                value={this.state.currentJobConfigXml}
                                onBeforeChange={(editor, data, value) => {
                                    this.setState({
                                        currentJobConfigXml: value,
                                    });
                                }}
                                onChange={(editor, data, value) => {
                                    this.setState({
                                        currentJobConfigXml: value,
                                    });
                                }}
                            />
                        </Col>
                    </Row>
                </Drawer>

                <Modal
                    title="添加凭证"
                    destroyOnClose={true}
                    visible={this.state.addCredentialModalVisible}
                    onOk={this.handleAddCredentialOk.bind(this)}
                    onCancel={this.handleAddCredentialCancel.bind(this)}
                >
                    <Form
                        {...formHorizontalItemLayout}
                        initialValues={{}}
                        ref={this.addCredentialFormRef}
                    >
                        <Form.Item
                            label="凭证ID"
                            name="credentialId"
                            tooltip={{
                                title: "输入自定义的唯一凭证标识",
                                icon: <InfoCircleOutlined />,
                            }}
                            rules={[
                                {
                                    required: true,
                                    message: "输入凭证ID!",
                                },
                            ]}
                        >
                            <Input placeholder="输入凭证ID" />
                        </Form.Item>
                        <Form.Item
                            label="用户名"
                            name="username"
                            rules={[
                                {
                                    required: true,
                                    message: "输入用户名!",
                                },
                            ]}
                        >
                            <Input placeholder="输入用户名" />
                        </Form.Item>
                        <Form.Item
                            label="密码"
                            name="password"
                            rules={[
                                {
                                    required: true,
                                    message: "输入密码!",
                                },
                            ]}
                        >
                            <Input placeholder="输入密码" />
                        </Form.Item>
                        <Form.Item
                            label="凭证描述"
                            name="description"
                            rules={[
                                {
                                    required: true,
                                    message: "输入凭证描述!",
                                },
                            ]}
                        >
                            <Input placeholder="输入凭证描述" />
                        </Form.Item>
                    </Form>
                </Modal>
            </Content>
        );
    }
}

class BuildsManageContent extends Component {
    constructor(props) {
        super(props);
        this.state = {
            recordColumns: [
                {
                    title: "构建记录",
                    dataIndex: "fullDisplayName",
                    key: "fullDisplayName",
                    render: (value) => {
                        return <Text ellipsis={true}>{value}</Text>;
                    },
                },
                {
                    title: "构建结果",
                    dataIndex: "result",
                    key: "result",
                    align: "center",
                    render: (value, record) => {
                        let resultIcon = (
                            <ClockCircleOutlined
                                style={{ color: "blue", fontSize: 15 }}
                            />
                        );
                        if (record.building === false && value === "FAILURE") {
                            resultIcon = (
                                <CloseCircleOutlined
                                    style={{ color: "red", fontSize: 15 }}
                                />
                            );
                        }
                        if (record.building === false && value === "SUCCESS") {
                            resultIcon = (
                                <CheckCircleOutlined
                                    style={{ color: "green", fontSize: 15 }}
                                />
                            );
                        }
                        return <Text>{resultIcon}</Text>;
                    },
                },
                // {
                //     title: "任务描述",
                //     dataIndex: "description",
                //     key: "description",
                //     render: (value) => {
                //         return <Text ellipsis={true}>{value}</Text>;
                //     },
                // },
                {
                    title: "代码信息",
                    dataIndex: "git",
                    key: "git",
                    render: (value, record) => {
                        let branch = "";
                        let commitId = "";
                        let gitUrl = "";
                        for (let i = 0; i < record.actions.length; i++) {
                            if (
                                record.actions[i]["_class"] ===
                                "hudson.plugins.git.util.BuildData"
                            ) {
                                gitUrl = record.actions[i]["remoteUrls"][0];
                                commitId =
                                    record.actions[i]["lastBuiltRevision"][
                                        "branch"
                                    ][0]["SHA1"];
                                branch =
                                    record.actions[i]["lastBuiltRevision"][
                                        "branch"
                                    ][0]["name"];
                            }
                        }
                        return (
                            <Text
                                ellipsis={true}
                                style={{ fontSize: 0, color: "#3385ff" }}
                            >
                                <span style={{ fontSize: 11 }}>
                                    {branch}-{commitId}
                                </span>
                                <br />
                                <span style={{ fontSize: 11 }}>{gitUrl}</span>
                            </Text>
                        );
                    },
                },
                {
                    title: "创建时间",
                    dataIndex: "timestamp",
                    key: "timestamp",
                    render: (value) => {
                        return (
                            <Text ellipsis={true}>
                                {moment.unix(value / 1000).format("MM/DD/YYYY")}
                            </Text>
                        );
                    },
                },
                {
                    title: "耗时",
                    dataIndex: "duration",
                    key: "duration",
                    render: (value) => {
                        let du = moment.duration(value, "ms");
                        let mins = du.get("minutes");
                        let ss = du.get("seconds");
                        return (
                            <Text ellipsis={true}>
                                {mins + "分" + ss + "秒"}
                            </Text>
                        );
                    },
                },
                {
                    title: "操作",
                    key: "operation",
                    fixed: "right",
                    align: "center",
                    render: (text, record) => {
                        return (
                            <div>
                                <Button
                                    type="link"
                                    size="small"
                                    onClick={this.showBuildProgress.bind(
                                        this,
                                        record.id,
                                    )}
                                >
                                    构建日志
                                </Button>
                                <Divider type="vertical" />
                                <Button
                                    type="link"
                                    size="small"
                                    onClick={this.showBuildStageInfo.bind(
                                        this,
                                        record.id,
                                    )}
                                >
                                    阶段日志
                                </Button>
                                <Divider type="vertical" />
                                <Button
                                    type="link"
                                    size="small"
                                    onClick={this.getBuildArchiveArtifactsInfo.bind(
                                        this,
                                        record.id,
                                    )}
                                >
                                    制品发布
                                </Button>

                                <Divider type="vertical" />
                                <Button
                                    type="link"
                                    size="small"
                                    onClick={this.deleteBuildItem.bind(
                                        this,
                                        record.id,
                                    )}
                                >
                                    删除记录
                                </Button>
                            </div>
                        );
                    },
                },
            ],
            recordPagination: {
                showSizeChanger: false,
                showQuickJumper: false,
                showTotal: (total) => `共 ${total} 条`,
                pageSize: 10,
                total: 0,
            },
            jobName: this.props.location.state.jobName,
            buildNum: 0,
            buildInfo: undefined,
            buildList: [],
            recordTableLoading: false,
            buildStageInfo: {},
            buildStageDrawerVisible: false,
            currentFlowNodesLog: [],
        };
    }

    componentDidMount() {
        this.loadCurrentJobBuildList();
    }

    loadCurrentJobBuildList() {
        this.setState({ recordTableLoading: true });
        getJobBuildList({ jobName: this.state.jobName })
            .then((res) => {
                if (res.code === 0) {
                    if (res.data.buildsDetails) {
                        this.setState({
                            buildList: res.data.buildsDetails,
                            recordPagination: {
                                ...this.state.recordPagination,
                                total: res.data.buildsDetails.length,
                            },
                        });
                    } else {
                        this.setState({
                            buildList: [],
                            recordPagination: {
                                ...this.state.recordPagination,
                                total: 0,
                            },
                        });
                    }
                    this.setState({ buildInfo: res.data });
                } else {
                    message.error(res.msg);
                }
                this.setState({ recordTableLoading: false });
            })
            .catch((err) => {
                this.setState({ recordTableLoading: false });
                message.error(err.toLocaleString());
            });
    }

    backToJobsContent() {
        this.props.history.push({
            pathname: "/admin/cicd/ci/",
        });
    }

    startBuild() {
        postJenkinsStartJob({ jobName: this.state.jobName }).then((res) => {
            if (res.code === 0) {
                message.success("构建请求已提交，可手动刷新查看构建记录");
                setTimeout(() => {
                    this.loadCurrentJobBuildList();
                }, 1000);
            } else {
                message.error(res.msg);
            }
        });
    }

    showBuildProgress(number) {
        this.setState({ buildDetailDrawerVisible: true });
        getJenkinsJobBuildLog({
            jobName: this.state.jobName,
            number: number,
        }).then((res) => {
            if (res.code === 0) {
                this.setState({
                    buildLog: res.data,
                });
            } else {
                message.error(res.msg);
            }
        });
    }

    showBuildStageInfo(number) {
        this.setState({
            buildStageDrawerVisible: true,
            buildNum: Number(number),
        });
        getJenkinsJobBuildStages({
            jobName: this.state.jobName,
            number: number,
        }).then((res) => {
            if (res.code === 0) {
                this.setState({
                    buildStageInfo: JSON.parse(res.data),
                });
            } else {
                message.error(res.msg);
            }
        });
    }

    getBuildArchiveArtifactsInfo(number) {
        let that = this;
        getJenkinsJobBuildArchiveArtifactsInfo({
            jobName: this.state.jobName,
            number: number,
        }).then((res) => {
            if (res.code === 0) {
                confirm({
                    title: "容器镜像制品信息",
                    width: 500,
                    content: (
                        <div style={{ fontSize: 12 }}>
                            <Paragraph copyable={{ tooltips: false }}>
                                {res.data}
                            </Paragraph>
                        </div>
                    ),
                    okText: "去发布",
                    cancelText: "关闭",
                    onOk() {
                        that.props.history.push({
                            pathname: "/admin/cicd/cd/",
                            state: {
                                jobName: that.state.jobName,
                                buildNumber: number,
                                imageName: res.data,
                            },
                        });
                    },
                });
            } else {
                message.error(res.msg);
            }
        });
    }

    closeBuildDetailDrawer() {
        this.setState({ buildDetailDrawerVisible: false });
    }

    closeBuildStageDrawer() {
        this.setState({ buildStageDrawerVisible: false });
    }

    deleteBuildItem(number) {
        deleteJenkinsJobBuildLog({
            jobName: this.state.jobName,
            buildNum: Number(number),
        }).then((res) => {
            if (res.code === 0) {
                message.success("删除成功!");
                this.loadCurrentJobBuildList();
            } else {
                message.error(res.msg);
            }
        });
    }

    onChangeStageCollapse(nodeId) {
        this.setState({ currentFlowNodesLog: [] });
        if (!nodeId) {
            return;
        }
        getJenkinsJobBuildStageLog({
            jobName: this.state.jobName,
            buildNum: Number(this.state.buildNum),
            nodeId: Number(nodeId),
        }).then((res) => {
            if (res.code === 0) {
                this.setState({
                    currentFlowNodesLog: res.data["stageFlowNodes"],
                });
            } else {
                message.error(res.msg);
            }
        });
    }

    render() {
        return (
            <Content
                style={{
                    background: "#fff",
                    margin: 0,
                }}
            >
                <Row>
                    <Col span={24}>
                        <Card
                            title={
                                <span>
                                    <LeftCircleOutlined
                                        onClick={this.backToJobsContent.bind(
                                            this,
                                        )}
                                    />
                                    &nbsp;&nbsp; 构建记录
                                </span>
                            }
                            extra={
                                <Fragment>
                                    <Button
                                        size="small"
                                        onClick={this.loadCurrentJobBuildList.bind(
                                            this,
                                        )}
                                    >
                                        刷新
                                    </Button>
                                    <Divider type="vertical" />
                                    <Button
                                        type="primary"
                                        size="small"
                                        onClick={this.startBuild.bind(this)}
                                    >
                                        立即构建
                                    </Button>
                                </Fragment>
                            }
                            size="small"
                        >
                            <Table
                                columns={this.state.recordColumns}
                                dataSource={this.state.buildList}
                                scroll={{ x: "max-content" }}
                                pagination={this.state.recordPagination}
                                loading={this.state.recordTableLoading}
                                bordered
                                size="small"
                            />
                        </Card>
                    </Col>
                </Row>
                <Drawer
                    title="构建日志"
                    placement="left"
                    width={800}
                    bodyStyle={{ paddingTop: 0 }}
                    destroyOnClose={true}
                    closable={true}
                    onClose={this.closeBuildDetailDrawer.bind(this)}
                    visible={this.state.buildDetailDrawerVisible}
                >
                    <Text>
                        <pre class="preJenkinsLog">{this.state.buildLog}</pre>
                    </Text>
                </Drawer>

                <Drawer
                    title="阶段日志"
                    placement="left"
                    width={800}
                    destroyOnClose={true}
                    closable={true}
                    onClose={this.closeBuildStageDrawer.bind(this)}
                    visible={this.state.buildStageDrawerVisible}
                >
                    <Text style={{ fontSize: 11 }}>
                        {this.state.buildStageInfo &&
                        this.state.buildStageInfo["stages"] &&
                        this.state.buildStageInfo["stages"].length > 0 ? (
                            <Collapse
                                accordion={true}
                                bordered={false}
                                onChange={this.onChangeStageCollapse.bind(this)}
                            >
                                {this.state.buildStageInfo["stages"].map(
                                    (item, index) => {
                                        return (
                                            <Panel
                                                header={item.name}
                                                key={item.id}
                                            >
                                                {this.state.currentFlowNodesLog
                                                    ? this.state.currentFlowNodesLog.map(
                                                          (item, index) => {
                                                              return (
                                                                  <div>
                                                                      {item.log !==
                                                                      "" ? (
                                                                          <pre className="preJenkinsLog">
                                                                              {
                                                                                  item.log
                                                                              }
                                                                          </pre>
                                                                      ) : (
                                                                          "当前阶段无日志"
                                                                      )}
                                                                  </div>
                                                              );
                                                          },
                                                      )
                                                    : "暂无日志"}
                                            </Panel>
                                        );
                                    },
                                )}
                            </Collapse>
                        ) : (
                            <Empty />
                        )}
                    </Text>
                </Drawer>
            </Content>
        );
    }
}

JobsManageContent = withRouter(JobsManageContent);
BuildsManageContent = withRouter(BuildsManageContent);

export { JobsManageContent, BuildsManageContent};
