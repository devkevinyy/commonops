import React, { Component } from "react";
import {
    Layout,
    Table,
    Button,
    Drawer,
    Tabs,
    Row,
    Col,
    message,
    Radio,
    Divider,
    Spin,
    DatePicker,
    Input,
    Descriptions,
    Modal,
    Form,
    Select,
    InputNumber,
    Typography,
    Popconfirm,
    Badge,
    Tooltip,
} from "antd";
import OpsBreadcrumbPath from "../breadcrumb_path";
import {
    deleteCloudKv,
    getCloudAccouts,
    getCloudKv,
    getCloudKvDetail,
    getCloudMonitorKv,
    postCloudKv,
    putCloudKv,
} from "../../api/cloud";
import { SearchOutlined, PlusCircleOutlined } from "@ant-design/icons";
import moment from "moment";
import ExtraInfoModal from "./common/extra_info_modal";
import LineChart from "./common/line_chart";

const { Text, Paragraph } = Typography;
const { Option } = Select;
const { Content } = Layout;
const TabPane = Tabs.TabPane;

class KvInfoModal extends Component {
    constructor(props) {
        super(props);
        this.state = {
            cloudAccountList: [],
        };
    }

    componentDidMount() {
        let that = this;
        getCloudAccouts(1, 100)
            .then((res) => {
                if (res.code === 0) {
                    that.setState({
                        cloudAccountList: res.data.accounts,
                    });
                } else {
                    message.error(res.msg);
                }
            })
            .catch((err) => {
                message.error(err.toLocaleString());
            });
    }

    render() {
        const formItemLayout = {
            labelCol: { span: 7 },
            wrapperCol: { span: 14 },
        };
        let accountOptions;
        accountOptions = this.state.cloudAccountList.map((item) => {
            return (
                <Option key={item.id} value={item.id}>
                    {item.accountName}
                </Option>
            );
        });
        return (
            <Modal
                title="Redis信息"
                destroyOnClose={true}
                visible={this.props.kv_info_modal_visible}
                onOk={this.props.handlePostKvInfoSubmit}
                onCancel={this.props.handlePostKvInfoCancel}
                okText="确认"
                cancelText="取消"
                width={600}
            >
                <Form
                    ref={this.props.formRef}
                    initialValues={{ instanceType: "Redis" }}
                >
                    <Form.Item
                        label="云账号"
                        {...formItemLayout}
                        name="cloudAccountId"
                        rules={[
                            { required: true, message: "云账号不能为空！" },
                        ]}
                    >
                        <Select>{accountOptions}</Select>
                    </Form.Item>
                    <Form.Item
                        label="实例Id"
                        {...formItemLayout}
                        name="instanceId"
                        rules={[
                            { required: true, message: "实例Id不能为空！" },
                        ]}
                    >
                        <Input />
                    </Form.Item>
                    <Form.Item
                        label="实例名称"
                        {...formItemLayout}
                        name="instanceName"
                        rules={[
                            { required: true, message: "实例名称不能为空！" },
                        ]}
                    >
                        <Input />
                    </Form.Item>
                    <Form.Item
                        label="实例类型"
                        {...formItemLayout}
                        name="instanceType"
                        rules={[
                            { required: true, message: "数据库类型不能为空！" },
                        ]}
                    >
                        <Select>
                            <Option value="Redis">Redis</Option>
                        </Select>
                    </Form.Item>
                    <Form.Item
                        label="连接串"
                        {...formItemLayout}
                        name="connectionString"
                        rules={[{ required: true, message: "请输入连接串！" }]}
                    >
                        <Input />
                    </Form.Item>
                    <Form.Item
                        label="端口"
                        {...formItemLayout}
                        name="port"
                        rules={[{ required: true, message: "请输入端口！" }]}
                    >
                        <Input />
                    </Form.Item>
                    <Form.Item
                        label="容量(G)"
                        {...formItemLayout}
                        name="capacity"
                        rules={[
                            {
                                type: "integer",
                                required: true,
                                message: "请输入数值型数据！",
                            },
                        ]}
                    >
                        <InputNumber />
                    </Form.Item>
                    <Form.Item
                        label="带宽(M)"
                        {...formItemLayout}
                        name="bandwidth"
                        rules={[
                            {
                                type: "integer",
                                required: true,
                                message: "请输入数值型数据！",
                            },
                        ]}
                    >
                        <InputNumber />
                    </Form.Item>
                    <Form.Item
                        label="创建时间"
                        {...formItemLayout}
                        name="createTime"
                        rules={[
                            { required: true, message: "创建时间不能为空！" },
                        ]}
                    >
                        <DatePicker format="YYYY-MM-DD" />
                    </Form.Item>
                    <Form.Item
                        label="过期时间"
                        {...formItemLayout}
                        name="endTime"
                    >
                        <DatePicker format="YYYY-MM-DD" />
                    </Form.Item>
                </Form>
            </Modal>
        );
    }
}

class KvContent extends Component {
    constructor(props) {
        super(props);
        this.handlePostKvInfoSubmit = this.handlePostKvInfoSubmit.bind(this);
        this.handlePostKvInfoCancel = this.handlePostKvInfoCancel.bind(this);
        this.handleExtraInfoOk = this.handleExtraInfoOk.bind(this);
        this.handleExtraInfoCancel = this.handleExtraInfoCancel.bind(this);
        let operWidth = this.props.isSuperAdmin ? 220 : 100;
        this.formRef = React.createRef();
        this.state = {
            columns: [
                {
                    title: "Id",
                    dataIndex: "ID",
                    key: "ID",
                    width: 50,
                    render: (value) => {
                        return <Text ellipsis={true}>{value}</Text>;
                    },
                },
                {
                    title: "实例id",
                    dataIndex: "InstanceId",
                    key: "InstanceId",
                    width: 200,
                    render: (value) => {
                        return <Text ellipsis={true}>{value}</Text>;
                    },
                },
                {
                    title: "实例名称",
                    dataIndex: "InstanceName",
                    key: "InstanceName",
                    width: 250,
                    render: (value, record) => {
                        return (
                            <Tooltip title={value}>
                                <Text ellipsis={true} style={{ width: "100%" }}>
                                    {value}
                                </Text>
                            </Tooltip>
                        );
                    },
                },
                {
                    title: "云账号",
                    dataIndex: "CloudAccountName",
                    key: "CloudAccountName",
                    width: 100,
                    render: (value) => {
                        return <Text ellipsis={true}>{value}</Text>;
                    },
                },
                {
                    title: "实例类型",
                    dataIndex: "InstanceType",
                    key: "InstanceType",
                    width: 100,
                    render: (value) => {
                        return <Text ellipsis={true}>{value}</Text>;
                    },
                },
                {
                    title: "连接串",
                    dataIndex: "ConnectionDomain",
                    key: "ConnectionDomain",
                    width: 400,
                    render: (value) => {
                        return (
                            <Paragraph
                                style={{ marginBottom: 0, fontSize: 13 }}
                                copyable
                            >
                                {value}
                            </Paragraph>
                        );
                    },
                },
                {
                    title: "带宽",
                    dataIndex: "Bandwidth",
                    key: "Bandwidth",
                    width: 100,
                    render: (value) => {
                        return <Text ellipsis={true}>{value}MB/s</Text>;
                    },
                },
                {
                    title: "连接数",
                    dataIndex: "Connections",
                    key: "Connections",
                    width: 100,
                    render: (value) => {
                        return <Text ellipsis={true}>{value}</Text>;
                    },
                },
                {
                    title: "容量",
                    dataIndex: "Capacity",
                    key: "Capacity",
                    width: 50,
                    render: (value) => {
                        return <Text ellipsis={true}>{value / 1024}G</Text>;
                    },
                },
                {
                    title: "QPS",
                    dataIndex: "QPS",
                    key: "QPS",
                    width: 100,
                    render: (value) => {
                        return <Text ellipsis={true}>{value}</Text>;
                    },
                },
                {
                    title: "过期时间",
                    dataIndex: "EndTime",
                    key: "EndTime",
                    width: 100,
                    render: (value) => {
                        return <Text ellipsis={true}>{value}</Text>;
                    },
                },
                {
                    title: "实例状态",
                    dataIndex: "InstanceStatus",
                    key: "InstanceStatus",
                    align: "center",
                    width: 100,
                    render: (value) => {
                        if (value === "Normal") {
                            return <Badge status="processing" />;
                        } else {
                            return <Badge status="error" />;
                        }
                    },
                },
                {
                    title: "操作",
                    key: "operation",
                    fixed: "right",
                    width: { operWidth },
                    align: "center",
                    render: (text, record) => {
                        return (
                            <div>
                                <Button
                                    type="primary"
                                    size="small"
                                    onClick={this.openMonitorDrawer.bind(
                                        this,
                                        record,
                                    )}
                                >
                                    监控
                                </Button>
                                <Divider type="vertical" />
                                <Button
                                    type="info"
                                    size="small"
                                    onClick={this.kvEdit.bind(this, record)}
                                    disabled={!this.props.aclAuthMap["PUT:/cloud/kv"]}
                                >
                                    编辑
                                </Button>
                                <Divider type="vertical" />
                                <Popconfirm
                                    title="确定删除该项内容?"
                                    onConfirm={this.kvDelete.bind(this, record)}
                                    okText="确认"
                                    cancelText="取消"
                                    disabled={!this.props.aclAuthMap["DELETE:/cloud/kv"]}
                                >
                                    <Button type="danger" size="small" disabled={!this.props.aclAuthMap["DELETE:/cloud/kv"]}>
                                        删除
                                    </Button>
                                </Popconfirm>
                            </div>
                        );
                    },
                },
            ],
            extraInfoModalVisible: false,
            tableLoading: false,
            tableData: [],
            chartData: [],
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
            drawerVisible: false,
            drawerPlacement: "left",
            instanceId: "",
            timeTagValue: "1h",
            metricTagValue: "StandardConnectionUsage",
            chartFormat: "%",
            currentServerDetail: {},
            queryExpiredTime: null,
            queryKeyword: "",
            queryCloudAccount: "0",
            queryManageUser: "0",
            queryDefineGroup: "",
            cloudAccountList: [],
            kv_info_modal_visible: false,
            selectedRowKeys: [],
            idsList: [],
            updateMode: "single",
            extraInfoData: {},
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

    componentDidMount() {
        this.refreshTableData();
        this.loadCloudAccountsData();
    }

    loadCloudAccountsData() {
        let that = this;
        getCloudAccouts(1, 100)
            .then((res) => {
                if (res.code === 0) {
                    that.setState({
                        cloudAccountList: res.data.accounts,
                    });
                } else {
                    message.error(res.msg);
                }
            })
            .catch((err) => {
                message.error(err.toLocaleString());
            });
    }

    kvEdit(data) {
        getCloudKvDetail(data.ID)
            .then((res) => {
                if (res["code"] !== 0) {
                    message.error(res["msg"]);
                } else {
                    let renewStatus = false;
                    if (res.data["RenewStatus"] === 1) {
                        renewStatus = true;
                    }
                    let extraInfoData = {
                        renewStatus: renewStatus,
                        renewSiteId: parseInt(res.data["RenewSiteId"])
                            ? parseInt(res.data["RenewSiteId"])
                            : null,
                        bankAccount: parseInt(res.data["BankAccount"])
                            ? parseInt(res.data["BankAccount"])
                            : null,
                        manageUser: parseInt(res.data["ManageUser"])
                            ? parseInt(res.data["ManageUser"])
                            : null,
                        cost: res.data["Cost"] ? res.data["Cost"] : "0",
                        costType: res.data["CostType"]
                            ? res.data["CostType"]
                            : "元/月",

                        kvInstanceName: res.data["InstanceName"],
                        kvBandwidth: res.data["Bandwidth"].toString(),
                        kvCapacity: (res.data["Capacity"] / 1024).toString(),
                        kvExpiredTime:
                            res.data["EndTime"] !== ""
                                ? moment(res.data["EndTime"])
                                : "",
                        resForm:
                            res.data["DataStatus"] === 1
                                ? "阿里云"
                                : "手动添加",
                    };
                    this.setState({
                        extraInfoModalVisible: true,
                        kvId: data.ID,
                        updateMode: "single",
                        resFrom:
                            data["DataStatus"] === 1 ? "阿里云" : "手动添加",
                        extraInfoData: extraInfoData,
                    });
                }
            })
            .catch((err) => {
                message.error(err.toLocaleString());
            });
    }

    handleExtraInfoOk(data) {
        let targetId = "";
        if (this.state.updateMode === "single") {
            targetId = String(this.state.kvId);
        } else {
            targetId = this.state.idsList.join(",");
        }
        putCloudKv({
            ...data,
            id: targetId,
        })
            .then((res) => {
                if (res.code === 0) {
                    message.success("修改成功");
                    this.setState({
                        extraInfoModalVisible: false,
                        selectedRowKeys: [],
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

    handleExtraInfoCancel(data) {
        this.setState({ extraInfoModalVisible: false });
    }

    kvDelete(data) {
        deleteCloudKv(data.ID)
            .then((res) => {
                if (res.code === 0) {
                    message.success("删除成功");
                    this.refreshTableData();
                } else {
                    message.error(res.msg);
                }
            })
            .catch((err) => {
                message.error(err.toLocaleString());
            });
    }

    refreshTableData = () => {
        this.setState({ tableLoading: true });
        const queryParams = {
            page: this.state.pagination.page,
            size: this.state.pagination.pageSize,
            queryExpiredTime:
                this.state.queryExpiredTime === null
                    ? null
                    : this.state.queryExpiredTime.format("YYYY-YY-DD HH:mm:ss"),
            queryKeyword: this.state.queryKeyword,
            queryCloudAccount: this.state.queryCloudAccount,
        };
        getCloudKv(queryParams)
            .then((res) => {
                const pagination = this.state.pagination;
                pagination.total = parseInt(res.data.total);
                pagination.page = parseInt(res.data.page);
                pagination.showTotal(parseInt(res.data.total));
                this.setState({
                    pagination,
                });
                let data = res["data"]["kv"];
                let tableData = [];
                for (let i = 0; i < data.length; i++) {
                    tableData.push({
                        key: data[i]["ID"],
                        ID: data[i]["ID"],
                        InstanceId: data[i]["InstanceId"],
                        InstanceName: data[i]["InstanceName"],
                        InstanceType: data[i]["InstanceType"],
                        ConnectionDomain: data[i]["ConnectionDomain"],
                        EndTime: moment(data[i]["EndTime"]).format(
                            "YYYY-MM-DD",
                        ),
                        InstanceStatus: data[i]["InstanceStatus"],
                        Capacity: data[i]["Capacity"],
                        QPS: data[i]["QPS"],
                        Bandwidth: data[i]["Bandwidth"],
                        Connections: data[i]["Connections"],
                        CloudAccountName: data[i]["CloudAccountName"],
                        DataStatus: data[i]["DataStatus"],
                    });
                }
                this.setState({ tableData: tableData, tableLoading: false });
            })
            .catch((err) => {
                console.log(err);
            });
    };

    openMonitorDrawer = (data) => {
        this.setState(
            {
                drawerVisible: true,
                instanceId: data.InstanceId,
                currentServerDetail: data,
            },
            () => {
                this.refreshMonitorData(
                    data.InstanceId,
                    this.state.timeTagValue,
                    this.state.metricTagValue,
                );
                this.refreshSeverDetail();
            },
        );
    };

    refreshMonitorData = (instanceId, timeTagValue, metricTagValue) => {
        this.setState({ chartLoading: true });
        getCloudMonitorKv(instanceId, timeTagValue, metricTagValue)
            .then((res) => {
                if (res["code"] !== 0) {
                    message.error(res["msg"]);
                    this.setState({ chartLoading: false });
                    return;
                }
                if (res["data"]["Datapoints"] === "") {
                    message.warn(
                        "未获取到监控数据，可能是非阿里云机器或其它原因！",
                    );
                    this.setState({ chartLoading: false });
                    return;
                }
                let dataPoints = JSON.parse(res["data"]["Datapoints"]);
                let chartData = [];
                for (let i = 0; i < dataPoints.length; i++) {
                    chartData.push({
                        date: moment(dataPoints[i]["timestamp"]).format(
                            "DD日HH:mm",
                        ),
                        value: dataPoints[i]["Average"],
                    });
                }
                this.setState({ chartLoading: false, chartData: chartData });
            })
            .catch((err) => {
                message.error(err.toLocaleString());
            });
    };

    // 获取服务器的详细信息
    refreshSeverDetail = (e) => {
        this.setState({ serverDetailLoading: true });
        getCloudKvDetail(this.state.currentServerDetail.ID)
            .then((res) => {
                if (res["code"] !== 0) {
                    message.error(res["msg"]);
                }
                this.setState({ currentServerDetail: res["data"] }, () => {
                    this.setState({ serverDetailLoading: false });
                });
            })
            .catch((err) => {
                console.log(err);
            });
    };

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

    onExpiredDateChange = (moment) => {
        if (moment == null) {
            this.setState({ queryExpiredTime: null });
        } else {
            this.setState({ queryExpiredTime: moment });
        }
    };

    keywordOnChange = (e) => {
        this.setState({ queryKeyword: e.target.value });
    };

    handleCloudAccountChange = (queryCloudAccount) => {
        this.setState({ queryCloudAccount });
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
                console.log(this.state.pagination.page);
                this.refreshTableData();
            },
        );
    };

    onCloseDrawer = () => {
        this.setState({ drawerVisible: false });
    };

    handleTimeTagChange = (e) => {
        this.setState({ timeTagValue: e.target.value });
        this.refreshMonitorData(
            this.state.instanceId,
            e.target.value,
            this.state.metricTagValue,
        );
    };

    handleMetricTagChange = (e) => {
        this.setState({ metricTagValue: e.target.value });
        switch (e.target.value) {
            case "StandardConnectionUsage":
                this.setState({ chartFormat: "%" });
                break;
            case "StandardCpuUsage":
                this.setState({ chartFormat: "%" });
                break;
            case "StandardAvgRt":
                this.setState({ chartFormat: "us" });
                break;
            case "Standardappend":
                this.setState({ chartFormat: "Count/Second" });
                break;
            default:
                this.setState({ chartFormat: "" });
                break;
        }
        this.refreshMonitorData(
            this.state.instanceId,
            this.state.timeTagValue,
            e.target.value,
        );
    };

    // 新增自定义redis信息
    handleAdd = () => {
        this.setState({ kv_info_modal_visible: true, kvId: 0 });
    };

    handlePostKvInfoSubmit() {
        this.formRef.current.validateFields().then((values) => {
            postCloudKv({
                ...values,
                createTime: values.createTime.format("YYYY-MM-DD HH:mm:ss"),
                endTime:
                    values.endTime === undefined
                        ? undefined
                        : values.endTime.format("YYYY-MM-DD HH:mm:ss"),
            })
                .then((res) => {
                    if (res.code === 0) {
                        message.success(
                            "添加成功，请到权限管理中增加访问权限！",
                        );
                        this.setState({ kv_info_modal_visible: false });
                        this.refreshTableData();
                    } else {
                        message.error(res.msg);
                    }
                })
                .catch((err) => {
                    message.error(err.toLocaleString());
                });
        });
    }

    handlePostKvInfoCancel() {
        this.setState({ kv_info_modal_visible: false });
    }

    handleBatchUpdate = () => {
        const idList = this.state.selectedRowKeys;
        if (idList.length === 0) {
            message.warn("请先选择数据");
            return;
        }
        this.setState({
            extraInfoModalVisible: true,
            updateMode: "batch",
            idsList: idList,
        });
    };

    onSelectChange = (selectedRowKeys) => {
        this.setState({ selectedRowKeys });
    };

    render() {
        let accountOptions;
        accountOptions = this.state.cloudAccountList.map((item) => {
            return (
                <Option key={item.id} value={item.id}>
                    {item.accountName}
                </Option>
            );
        });
        return (
            <Content
                style={{
                    background: "#fff",
                    padding: "5px 20px",
                    margin: 0,
                    height: "100%",
                }}
            >
                <OpsBreadcrumbPath pathData={["云资源", "KV存储", "KV列表"]} />
                <Row style={{ padding: "0px 0px 10px 0px" }}>
                    <Col span={3} className="col-span">
                        <DatePicker
                            style={{ width: "100%" }}
                            placeholder="到期截止时间"
                            defaultValue={this.state.queryExpiredTime}
                            onChange={this.onExpiredDateChange}
                        />
                    </Col>
                    <Col span={5} className="col-span">
                        <Input
                            placeholder="输入实例id\名称\连接串查询"
                            value={this.state.queryKeyword}
                            onChange={this.keywordOnChange}
                        />
                    </Col>
                    <Col span={3} className="col-span">
                        <Select
                            defaultValue={this.state.queryCloudAccount}
                            style={{ width: "100%" }}
                            onChange={this.handleCloudAccountChange}
                        >
                            <Option value="0">所有云账号</Option>
                            {accountOptions}
                        </Select>
                    </Col>
                    <Col span={2} className="col-span">
                        <Button
                            style={{ width: "100%" }}
                            type="primary"
                            icon={<SearchOutlined />}
                            onClick={this.handleQuery}
                        >
                            查 询
                        </Button>
                    </Col>
                    <Col span={2} className="col-span">
                        <Button
                            style={{ width: "100%" }}
                            icon={<PlusCircleOutlined />}
                            onClick={this.handleAdd}
                            disabled={!this.props.aclAuthMap["POST:/cloud/kv"]}
                        >
                            新 增
                        </Button>
                    </Col>
                </Row>

                <KvInfoModal
                    formRef={this.formRef}
                    kv_info_modal_visible={this.state.kv_info_modal_visible}
                    handlePostKvInfoSubmit={this.handlePostKvInfoSubmit}
                    handlePostKvInfoCancel={this.handlePostKvInfoCancel}
                />

                {/*完善信息组件*/}
                <ExtraInfoModal
                    editData={this.state.extraInfoData}
                    resType="kv"
                    updateMode={this.state.updateMode}
                    resFrom={this.state.resFrom}
                    visible={this.state.extraInfoModalVisible}
                    handleOk={this.handleExtraInfoOk}
                    handleCancel={this.handleExtraInfoCancel}
                />

                <div>
                    <Drawer
                        title="实例详情及监控数据"
                        placement={this.state.drawerPlacement}
                        closable={true}
                        destroyOnClose={true}
                        onClose={this.onCloseDrawer}
                        visible={this.state.drawerVisible}
                        width={950}
                    >
                        <Tabs
                            defaultActiveKey="1"
                            tabPosition="left"
                            style={{ marginLeft: -30 }}
                        >
                            <TabPane tab="监控详情" key="1">
                                <Spin
                                    tip="图表生成中..."
                                    spinning={this.state.chartLoading}
                                >
                                    <Row style={{ marginBottom: "10px" }}>
                                        <Col
                                            span={3}
                                            style={{ lineHeight: "30px" }}
                                        >
                                            时间维度：
                                        </Col>
                                        <Col span={15}>
                                            <Radio.Group
                                                value={this.state.timeTagValue}
                                                onChange={
                                                    this.handleTimeTagChange
                                                }
                                            >
                                                <Radio.Button value="1h">
                                                    1小时
                                                </Radio.Button>
                                                <Radio.Button value="6h">
                                                    6小时
                                                </Radio.Button>
                                                <Radio.Button value="12h">
                                                    12小时
                                                </Radio.Button>
                                                <Radio.Button value="1d">
                                                    1 天
                                                </Radio.Button>
                                                <Radio.Button value="3d">
                                                    3 天
                                                </Radio.Button>
                                                <Radio.Button value="7d">
                                                    7 天
                                                </Radio.Button>
                                                <Radio.Button value="14d">
                                                    14 天
                                                </Radio.Button>
                                            </Radio.Group>
                                        </Col>
                                    </Row>
                                    <Row>
                                        <Col
                                            span={3}
                                            style={{ lineHeight: "30px" }}
                                        >
                                            监控维度：
                                        </Col>
                                        <Col span={16}>
                                            <Radio.Group
                                                value={
                                                    this.state.metricTagValue
                                                }
                                                onChange={
                                                    this.handleMetricTagChange
                                                }
                                            >
                                                <Radio.Button value="StandardConnectionUsage">
                                                    连接数使用率
                                                </Radio.Button>
                                                <Radio.Button value="StandardCpuUsage">
                                                    CPU 使用率
                                                </Radio.Button>
                                                <Radio.Button value="StandardAvgRt">
                                                    平均响应时间
                                                </Radio.Button>
                                                {/*<Radio.Button value="Standardappend">*/}
                                                {/*    append 命令的执行频率*/}
                                                {/*</Radio.Button>*/}
                                            </Radio.Group>
                                        </Col>
                                    </Row>
                                    <Row style={{ marginTop: 20 }}>
                                        <Col>
                                            <LineChart
                                                width={800}
                                                height={400}
                                                data={this.state.chartData}
                                            />
                                        </Col>
                                    </Row>
                                </Spin>
                            </TabPane>
                            <TabPane tab="信息详情" key="2">
                                <Spin
                                    tip="数据获取中..."
                                    spinning={this.state.serverDetailLoading}
                                >
                                    <Descriptions
                                        title="基本信息"
                                        bordered
                                        size="small"
                                        column={2}
                                    >
                                        <Descriptions.Item label="实例ID">
                                            {
                                                this.state.currentServerDetail
                                                    .InstanceId
                                            }
                                        </Descriptions.Item>
                                        <Descriptions.Item label="实例名称">
                                            {
                                                this.state.currentServerDetail
                                                    .InstanceName
                                            }
                                        </Descriptions.Item>
                                        <Descriptions.Item label="实例类型">
                                            {
                                                this.state.currentServerDetail
                                                    .InstanceType
                                            }
                                        </Descriptions.Item>
                                        <Descriptions.Item label="架构类型">
                                            {
                                                this.state.currentServerDetail
                                                    .ArchitectureType
                                            }
                                        </Descriptions.Item>
                                        <Descriptions.Item label="实例带宽">
                                            {
                                                this.state.currentServerDetail
                                                    .Bandwidth
                                            }
                                        </Descriptions.Item>
                                        <Descriptions.Item label="实例容量">
                                            {
                                                this.state.currentServerDetail
                                                    .Capacity
                                            }
                                        </Descriptions.Item>
                                        <Descriptions.Item label="带宽">
                                            {
                                                this.state.currentServerDetail
                                                    .Bandwidth
                                            }
                                        </Descriptions.Item>
                                        <Descriptions.Item label="连接地址">
                                            {
                                                this.state.currentServerDetail
                                                    .ConnectionDomain
                                            }
                                        </Descriptions.Item>
                                        <Descriptions.Item label="端口">
                                            {
                                                this.state.currentServerDetail
                                                    .Port
                                            }
                                        </Descriptions.Item>
                                        <Descriptions.Item label="私有IP">
                                            {
                                                this.state.currentServerDetail
                                                    .PrivateIp
                                            }
                                        </Descriptions.Item>
                                        <Descriptions.Item label="连接数限制">
                                            {
                                                this.state.currentServerDetail
                                                    .Connections
                                            }
                                        </Descriptions.Item>
                                        <Descriptions.Item label="创建时间">
                                            {
                                                this.state.currentServerDetail
                                                    .CreateTime
                                            }
                                        </Descriptions.Item>
                                        <Descriptions.Item label="过期时间">
                                            {
                                                this.state.currentServerDetail
                                                    .EndTime
                                            }
                                        </Descriptions.Item>
                                        <Descriptions.Item label="QPS">
                                            {this.state.currentServerDetail.QPS}
                                        </Descriptions.Item>
                                    </Descriptions>
                                </Spin>
                            </TabPane>
                        </Tabs>
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
                </div>
            </Content>
        );
    }
}

export default KvContent;
