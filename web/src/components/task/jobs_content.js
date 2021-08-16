import React, { Component, Fragment } from "react";
import OpsBreadcrumbPath from "../breadcrumb_path";
import {
    Button,
    Col,
    Divider,
    Input,
    Layout,
    Row,
    Table,
    Tag,
    message,
    Drawer,
    Spin,
    Descriptions,
    Alert,
    Card,
    Typography,
    Modal,
    DatePicker,
} from "antd";
import {
    getDailyJobDetail,
    getDailyJobs,
    putDailyJob,
    putDailyJobExecutorUser,
} from "../../api/daily_task";
import moment from "moment";
import { SearchOutlined } from "@ant-design/icons";
import "moment/locale/zh-cn";
import { getCurrentUserId } from "../../services/common";
import { ServerBase } from "../../config";
moment.locale("zh-cn");

const { TextArea } = Input;
const { Text, Paragraph } = Typography;
const { Content } = Layout;

const NULL_TIMESTAMP = "0001-01-01 00:00:00";

class Jobs_content extends Component {
    constructor(props) {
        super(props);
        this.downloadFile = this.downloadFile.bind(this);
        this.ChangeExecutorSelect = this.ChangeExecutorSelect.bind(this);
        this.changeRefuseReason = this.changeRefuseReason.bind(this);
        this.autoCreateJob = this.autoCreateJob.bind(this);
        this.state = {
            columns: [
                {
                    title: "Id",
                    dataIndex: "id",
                    key: "id",
                    width: 50,
                    render: (value) => {
                        return <Text ellipsis={true}>{value}</Text>;
                    },
                },
                {
                    title: "紧急度",
                    dataIndex: "important_degree",
                    key: "important_degree",
                    width: 70,
                    render: (value) => {
                        let content;
                        switch (value) {
                            case "非常紧急":
                                content = <Tag color="#f50">{value}</Tag>;
                                break;
                            case "紧急":
                                content = (
                                    <Tag color="rgb(255, 147, 137)">
                                        {value}
                                    </Tag>
                                );
                                break;
                            default:
                                content = <Tag color="#108ee9">{value}</Tag>;
                                break;
                        }
                        return content;
                    },
                },
                {
                    title: "状态",
                    dataIndex: "status",
                    key: "status",
                    width: 50,
                    render: (value) => {
                        const content = this.getStatusText(value);
                        return <Text ellipsis={true}>{content}</Text>;
                    },
                },
                {
                    title: "任务名",
                    dataIndex: "job_name",
                    key: "job_name",
                    width: 200,
                    render: (value) => {
                        return <Text ellipsis={true}>{value}</Text>;
                    },
                },
                {
                    title: "任务类型",
                    dataIndex: "job_type",
                    key: "job_type",
                    width: 100,
                    render: (value) => {
                        return <Text ellipsis={true}>{value}</Text>;
                    },
                },
                {
                    title: "创建人",
                    dataIndex: "creator_user_name",
                    key: "creator_user_name",
                    width: 70,
                    render: (value) => {
                        return <Text ellipsis={true}>{value}</Text>;
                    },
                },
                {
                    title: "创建时间",
                    dataIndex: "create_time",
                    key: "create_time",
                    width: 150,
                    render: (value) => {
                        return (
                            <Text ellipsis={true}>
                                {moment(value).format("MM月DD日HH:mm")}
                            </Text>
                        );
                    },
                },
                {
                    title: "执行人",
                    dataIndex: "executor_user_name",
                    key: "executor_user_name",
                    width: 70,
                    render: (value) => {
                        return <Text ellipsis={true}>{value}</Text>;
                    },
                },
                {
                    title: "响应时间",
                    dataIndex: "accept_time",
                    key: "accept_time",
                    width: 150,
                    render: (value) => {
                        if (value === NULL_TIMESTAMP) {
                            return <Text ellipsis={true}>-</Text>;
                        }
                        return (
                            <Text ellipsis={true}>
                                {moment(value).format("MM月DD日HH:mm")}
                            </Text>
                        );
                    },
                },
                {
                    title: "结束时间",
                    dataIndex: "end_time",
                    key: "end_time",
                    width: 150,
                    render: (value) => {
                        if (value === NULL_TIMESTAMP) {
                            return <Text ellipsis={true}>-</Text>;
                        }
                        return (
                            <Text ellipsis={true}>
                                {moment(value).format("MM月DD日HH:mm")}
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
                            <Fragment>
                                <Button
                                    type="primary"
                                    size="small"
                                    onClick={this.showJobDetail.bind(
                                        this,
                                        record,
                                    )}
                                >
                                    详情
                                </Button>
                                <Divider type="vertical" />
                                <Button
                                    type="info"
                                    size="small"
                                    disabled={!record.isGetEnable}
                                    onClick={this.startJob.bind(
                                        this,
                                        record.id,
                                    )}
                                >
                                    领取
                                </Button>
                                <Divider type="vertical" />
                                <Button
                                    type="info"
                                    size="small"
                                    disabled={!record.isFinishEnable}
                                    onClick={this.finishJob.bind(
                                        this,
                                        record.id,
                                    )}
                                >
                                    完成
                                </Button>
                                <Divider type="vertical" />
                                <Button
                                    type="danger"
                                    size="small"
                                    disabled={!record.isRefuseEnable}
                                    onClick={this.refuseJob.bind(
                                        this,
                                        record.id,
                                    )}
                                >
                                    拒绝
                                </Button>
                                <Divider type="vertical" />
                                <Button
                                    type="danger"
                                    size="small"
                                    disabled={!record.isDeleteEnable}
                                    onClick={this.deleteJob.bind(
                                        this,
                                        record.id,
                                    )}
                                >
                                    删除
                                </Button>
                            </Fragment>
                        );
                    },
                },
            ],
            tableLoading: false,
            drawerSpinning: false,
            drawerVisible: false,
            opsUsersList: [],
            changeExecutorModalVisible: false,
            refuseJobModalVisible: false,
            refuseId: 0,
            refuseReason: "",
            changeToUserId: 0,
            jobChangeId: 0,
            tableData: [],
            pagination: {
                showSizeChanger: true,
                pageSizeOptions: ["10", "20", "30", "100"],
                onShowSizeChange: (current, size) =>
                    this.onShowSizeChange(current, size),
                showQuickJumper: false,
                showTotal: (total) => `共 ${total} 条`,
                pageSize: 10,
                page: 1,
                total: 0,
                onChange: (page, pageSize) => this.changePage(page, pageSize),
            },
            job_file: "",
            queryKeyword: "",
            currentJob: "",
            currentJobDetail: "",
            isGetEnable: true,
            isFinishEnable: true,
            isDeleteEnable: true,
            queryCreateTime: null,
        };
    }

    onShowSizeChange(current, size) {
        let pagination = {
            ...this.state.pagination,
            page: 1,
            current: 1,
            pageSize: size,
        };
        this.setState(
            {
                pagination: pagination,
            },
            () => {
                this.refreshTableData();
            },
        );
    }

    changePage = (page, pageSize) => {
        this.setState(
            {
                pagination: {
                    ...this.state.pagination,
                    page: page,
                    current: page,
                    pageSize: pageSize,
                },
            },
            () => {
                this.refreshTableData();
            },
        );
    };

    componentDidMount() {
        this.refreshTableData();
    }

    showJobDetail(data) {
        this.setState({
            drawerVisible: true,
            drawerSpinning: true,
            currentJob: data,
        });
        getDailyJobDetail(data.id)
            .then((res) => {
                if (res.code === 0) {
                    this.setState(
                        {
                            ...res.data,
                            currentJobDetail: res.data,
                        },
                    );
                } else {
                    message.error(res.msg);
                }
                this.setState({ drawerSpinning: false });
            })
            .catch((err) => {
                message.error(err.toLocaleString());
            });
    }

    startJob(id) {
        putDailyJob({
            id: id,
            action: "getJob",
        })
            .then((res) => {
                if (res.code === 0) {
                    let currentJob = Object.assign({}, this.state.currentJob, {
                        isGetEnable: false,
                    });
                    message.success("领取成功");
                    this.setState({
                        currentJob: currentJob,
                    });
                    this.refreshTableData();
                } else {
                    message.error(res.msg);
                }
            })
            .catch((err) => {
                message.error(err.toLocaleString());
            });
    }

    ChangeExecutorSelect(data) {
        this.setState({ changeToUserId: data });
    }

    changeRefuseReason(e) {
        this.setState({ refuseReason: e.target.value });
    }

    finishJob(id) {
        putDailyJob({
            id: id,
            action: "finishJob",
        })
            .then((res) => {
                if (res.code === 0) {
                    message.success("成功完成，继续努力");
                    this.refreshTableData();
                } else {
                    message.error(res.msg);
                }
            })
    }

    refuseJob(id) {
        this.setState({
            refuseJobModalVisible: true,
            refuseId: id,
        });
    }

    deleteJob(id) {
        putDailyJob({
            id: id,
            action: "deleteJob",
        })
            .then((res) => {
                if (res.code === 0) {
                    message.success("删除成功");
                    this.refreshTableData();
                } else {
                    message.error(res.msg);
                }
            })
    }

    handlerRefuseJobCommit = () => {
        if (this.state.refuseId === 0) {
            message.error("未选择拒绝工单");
            return;
        }
        if (this.state.refuseReason === "") {
            message.error("必须填写拒绝原因");
            return;
        }
        putDailyJob({
            id: this.state.refuseId,
            refuseReason: this.state.refuseReason,
            action: "refuseJob",
        })
            .then((res) => {
                if (res.code === 0) {
                    this.setState({ refuseJobModalVisible: false });
                    message.success("操作成功");
                    this.refreshTableData();
                } else {
                    message.error(res.msg);
                }
            })
    };

    downloadFile() {
        window.open(
            ServerBase +
                "fileDownload?objectName=" +
                encodeURIComponent(this.state.job_file),
            "_blank",
        );
    }

    refreshTableData = () => {
        this.setState({ tableLoading: true });
        const queryParams = {
            page: this.state.pagination.page,
            size: this.state.pagination.pageSize,
            queryKeyword: this.state.queryKeyword,
            queryCreateTime: this.state.queryCreateTime,
        };
        getDailyJobs(queryParams)
            .then((res) => {
                const pagination = this.state.pagination;
                pagination.total = parseInt(res.data.total);
                pagination.page = parseInt(res.data.page);
                pagination.showTotal(parseInt(res.data.total));
                this.setState({
                    pagination,
                });
                let data = res["data"]["jobs"];
                let tableData = [];
                for (let i = 0; i < data.length; i++) {
                    tableData.push({
                        key: data[i]["id"],
                        id: data[i]["id"],
                        job_name: data[i]["job_name"],
                        job_type: data[i]["job_type"],
                        important_degree: data[i]["important_degree"],
                        creator_user_name: data[i]["creator_user_name"],
                        executor_user_name: data[i]["executor_user_name"],
                        status: data[i]["status"],
                        create_time: data[i]["create_time"],
                        accept_time: data[i]["accept_time"],
                        end_time: data[i]["end_time"],
                        approve_id: data[i]["approve_id"],
                        approve_content: data[i]["approve_content"],
                        refuse_reason: data[i]["refuse_reason"],
                        isExecutorChangable:
                            data[i]["status"] === 2 &&
                            data[i]["executor_user_id"] === getCurrentUserId(),
                        isGetEnable: data[i]["status"] === 1,
                        isFinishEnable:
                            data[i]["status"] === 2 &&
                            data[i]["executor_user_id"] === getCurrentUserId(),
                        isDeleteEnable:
                            data[i]["creator_user_id"] === getCurrentUserId() &&
                            data[i]["status"] < 2,
                        isRefuseEnable: data[i]["status"] === 1,
                    });
                }
                this.setState({ tableData: tableData, tableLoading: false });
            })
            .catch((err) => {
                console.log(err);
            });
    };

    keywordOnChange = (e) => {
        this.setState({
            queryKeyword: e.target.value,
        });
    };

    // 用户自定义查询
    handleQuery = () => {
        this.setState(
            {
                pagination: {
                    ...this.state.pagination,
                    page: 1,
                    current: 1,
                },
            },
            () => {
                this.refreshTableData();
            },
        );
    };

    onDrawerClose = () => {
        this.setState({
            drawerVisible: false,
        });
    };

    autoCreateJob() {
        localStorage.setItem(
            "job_template",
            JSON.stringify(this.state.currentJobDetail),
        );
        this.props.history.push({ pathname: "/admin/task/deploy_project" });
    }

    onQueryDateChange = (moment) => {
        if (moment == null) {
            this.setState({ queryCreateTime: null });
        } else {
            this.setState({ queryCreateTime: moment.format("YYYY-MM-DD") });
        }
    };

    getStatusText(status) {
        let content;
        switch (status) {
            case 1:
                content = <Tag color="volcano">待领取</Tag>;
                break;
            case 2:
                content = <Tag color="geekblue">处理中</Tag>;
                break;
            case 3:
                content = <Tag color="green">已完成</Tag>;
                break;
            case 0:
                content = <Tag color="magenta">已拒绝</Tag>;
                break;
            default:
                content = <Tag color="lime">已删除</Tag>;
                break;
        }
        return content;
    }

    render() {
        let waitContent;
        if (this.state.status > 0 && this.state.status < 3) {
            const m1 = moment(this.state.create_time);
            const m2 = moment(moment().format("YYYY-MM-DD HH:mm:ss"));
            let du = moment.duration(m2 - m1, "ms");
            waitContent =
                this.state.creator_user_name + " 已等待 " + du.humanize();
        } else {
            const m1 = moment(this.state.create_time);
            const m2 = moment(this.state.end_time);
            let du = moment.duration(m2 - m1, "ms");
            waitContent = "工单用时: " + du.humanize();
        }
        let k8sConfigContent = (
            <Descriptions.Item label="无配置">用户未设置</Descriptions.Item>
        );
        if (
            this.state.open_deploy_auto_config !== undefined &&
            this.state.open_deploy_auto_config !== ""
        ) {
            let ConfigList = JSON.parse(this.state.open_deploy_auto_config);
            k8sConfigContent = ConfigList.map((item, index) => {
                return (
                    <Descriptions.Item key={index} label={item.key}>
                        <Paragraph
                            style={{ marginBottom: 0, width: "300px" }}
                            ellipsis={true}
                            copyable={item.value !== ""}
                        >
                            {item.value}
                        </Paragraph>
                    </Descriptions.Item>
                );
            });
        }
        let refuseContent = "";
        if (this.state.status === 0) {
            refuseContent = (
                <Alert
                    style={{ marginTop: 10 }}
                    message={
                        <Fragment>
                            拒绝原因：{this.state.refuse_reason}
                        </Fragment>
                    }
                    type="error"
                />
            );
        }
        let titleContent = (
            <div>
                <span style={{ marginRight: "20px" }}>工单详情</span>
            </div>
        );
        return (
            <Content
                stmayle={{
                    background: "#fff",
                    margin: 0,
                    height: "100%",
                }}
            >
                <OpsBreadcrumbPath pathData={["工作协助", "工单列表"]} />
                <div style={{ marginBottom: 20 }}>
                    <Row>
                        <Col span={8} className="col-span">
                            <Input
                                placeholder="输入工单名称、创建人、执行人查询"
                                value={this.state.queryKeyword}
                                onChange={this.keywordOnChange}
                            />
                        </Col>
                        <Col span={3} className="col-span">
                            <DatePicker
                                style={{ width: "100%" }}
                                defaultValue={this.state.queryCreateTime}
                                placeholder="创建时间"
                                onChange={this.onQueryDateChange}
                            />
                        </Col>
                        <Col span={3}>
                            <Button
                                type="primary"
                                icon={<SearchOutlined />}
                                onClick={this.handleQuery}
                            >
                                查 询
                            </Button>
                        </Col>
                    </Row>
                </div>

                <Modal
                    title="你正在拒绝他人工单，请谨慎操作！"
                    visible={this.state.refuseJobModalVisible}
                    onOk={this.handlerRefuseJobCommit}
                    onCancel={() => {
                        this.setState({ refuseJobModalVisible: false });
                    }}
                    style={{ textAlign: "center" }}
                >
                    <Alert
                        type="warning"
                        message="请详细填写你的拒绝理由，将反馈给工单创建者并留存！"
                        style={{ marginBottom: 10 }}
                    />
                    <TextArea
                        rows={4}
                        placeholder="填写详细拒绝原因，以便对方后续完善"
                        onChange={this.changeRefuseReason}
                    />
                </Modal>

                <Drawer
                    title={titleContent}
                    placement="left"
                    closable={true}
                    onClose={this.onDrawerClose}
                    visible={this.state.drawerVisible}
                    width={800}
                >
                    <Spin
                        tip="数据加载中..."
                        spinning={this.state.drawerSpinning}
                    >
                        <Row style={{ marginBottom: 10 }}>
                            <Col span={24} style={{ fontSize: 15 }}>
                                <Alert
                                    message={
                                        <Fragment>
                                            {this.getStatusText(
                                                this.state.status,
                                            )}
                                            {this.state.job_name}
                                        </Fragment>
                                    }
                                    type="info"
                                />
                                {refuseContent}
                            </Col>
                        </Row>
                        <Row style={{ marginBottom: 10 }}>
                            <Col span={24}>
                                <Alert
                                    message={waitContent}
                                    type="warning"
                                    showIcon
                                />
                            </Col>
                        </Row>
                        <Row style={{ marginTop: "10px" }}>
                            <Card
                                size="small"
                                title="工单基本信息"
                                style={{ marginBottom: 10, width: "100%" }}
                            >
                                <Descriptions bordered size="small" column={2}>
                                    <Descriptions.Item label="任务类型">
                                        {this.state.job_type}
                                    </Descriptions.Item>
                                    <Descriptions.Item label="紧急度">
                                        {this.state.important_degree}
                                    </Descriptions.Item>
                                    <Descriptions.Item label="发起时间">
                                        {this.state.create_time}
                                    </Descriptions.Item>
                                </Descriptions>
                            </Card>
                        </Row>
                        <Row style={{ marginTop: "10px" }}>
                            <Card
                                size="small"
                                title="工作内容"
                                style={{ marginBottom: 10, width: "100%" }}
                            >
                                <Paragraph
                                    style={{ marginBottom: 0 }}
                                    ellipsis={true}
                                    copyable={this.state.task_content !== ""}
                                >
                                    {this.state.task_content}
                                </Paragraph>
                            </Card>
                            <Card
                                size="small"
                                title="自定义项"
                                style={{ width: "100%" }}
                            >
                                <Descriptions
                                    bordered
                                    size="small"
                                    column={1}
                                    style={{ marginTop: "10px" }}
                                >
                                    {k8sConfigContent}
                                </Descriptions>
                            </Card>
                            <Card
                                size="small"
                                title="用户备注"
                                style={{ width: "100%" }}
                            >
                                <Paragraph
                                    style={{ marginBottom: 0 }}
                                    ellipsis={true}
                                    copyable={this.state.remark !== ""}
                                >
                                    {this.state.remark}
                                </Paragraph>
                            </Card>
                        </Row>
                        <Row style={{ textAlign: "center", marginTop: 20 }}>
                            <Col span={7} />
                            <Col span={4}>
                                <Button
                                    type="primary"
                                    disabled={
                                        !this.state.currentJob.isGetEnable
                                    }
                                    onClick={this.startJob.bind(
                                        this,
                                        this.state.currentJob.id,
                                    )}
                                >
                                    一键领取
                                </Button>
                            </Col>
                            <Col offset={1} />
                            <Col span={4}>
                                <Button
                                    type="danger"
                                    disabled={
                                        !this.state.currentJob.isFinishEnable
                                    }
                                    onClick={this.finishJob.bind(
                                        this,
                                        this.state.currentJob.id,
                                    )}
                                >
                                    我已完成
                                </Button>
                            </Col>
                            <Col span={7} />
                        </Row>
                        <Row style={{ minHeight: 30 }} />
                    </Spin>
                </Drawer>
                <Table
                    columns={this.state.columns}
                    dataSource={this.state.tableData}
                    scroll={{ x: "max-content" }}
                    pagination={this.state.pagination}
                    loading={this.state.tableLoading}
                    bordered
                    size="small"
                />
            </Content>
        );
    }
}

export default Jobs_content;
