import React, { Component } from "react";
import {
    Button,
    Col,
    DatePicker,
    Descriptions,
    Divider,
    Drawer,
    Form,
    Typography,
    Input,
    InputNumber,
    Layout,
    message,
    Modal,
    Radio,
    Row,
    Select,
    Spin,
    Table,
    Tabs,
    Popconfirm,
    Badge,
    Tooltip,
    Tag,
    Tree,
} from "antd";
import {
    SearchOutlined,
    PlusCircleOutlined,
    EyeInvisibleOutlined,
    EyeTwoTone,
    CodeOutlined,
} from "@ant-design/icons";
import OpsBreadcrumbPath from "../breadcrumb_path";
import moment from "moment";
import "../../assets/css/main.css";
import {
    deleteCloudServer,
    getCloudAccouts,
    getCloudMonitorEcs,
    getCloudServerDetail,
    getCloudServers,
    postCloudServer,
    putCloudServer,
    getCloudServersAllData,
    postCloudServerBatchSSH,
} from "../../api/cloud";
import { Controlled as CodeMirror } from "react-codemirror2";
import "codemirror/lib/codemirror.css";
import "moment/locale/zh-cn";
import ExtraInfoModal from "./common/extra_info_modal";
import LineChart from "./common/line_chart";
import { LinuxSvg, WindowsSvg } from "../../assets/Icons";
import { Terminal } from "xterm";
import { AttachAddon } from "xterm-addon-attach";
import { FitAddon } from "xterm-addon-fit";
import "../../../node_modules/xterm/css/xterm.css";
import ReconnectingWebSocket from "reconnecting-websocket";
import { WSBase } from "../../config";
moment.locale("zh-cn");

const TabPane = Tabs.TabPane;
const { Text, Paragraph, Title } = Typography;
const Option = Select.Option;
const { Content } = Layout;

class ServerInfoModal extends Component {
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
                title="服务器信息"
                destroyOnClose={true}
                visible={this.props.server_info_modal_visible}
                onOk={this.props.handlePostServerInfoSubmit}
                onCancel={this.props.handlePostServerInfoCancel}
                okText="确认"
                cancelText="取消"
                centered={true}
                width={600}
            >
                <Form
                    ref={this.props.formRef}
                    initialValues={{
                        osType: "linux",
                        createTime: moment(),
                        expiredTime: moment(
                            "2099-12-30 00:00:00",
                            "YYYY-MM-DD HH:mm:ss",
                        ),
                    }}
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
                        label="主机名"
                        {...formItemLayout}
                        name="hostName"
                        rules={[
                            { required: true, message: "hostName不能为空！" },
                        ]}
                    >
                        <Input />
                    </Form.Item>
                    <Form.Item
                        label="实例Id"
                        {...formItemLayout}
                        name="instanceId"
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
                        label="实例描述"
                        {...formItemLayout}
                        name="description"
                        rules={[
                            { required: true, message: "实例描述不能为空！" },
                        ]}
                    >
                        <Input />
                    </Form.Item>
                    <Form.Item
                        label="内网IP"
                        {...formItemLayout}
                        name="innerIpAddress"
                        rules={[
                            { required: true, message: "内网IP不能为空！" },
                        ]}
                    >
                        <Input />
                    </Form.Item>
                    <Form.Item
                        label="外网IP"
                        {...formItemLayout}
                        name="publicIpAddress"
                        rules={[
                            { required: true, message: "外网IP不能为空！" },
                        ]}
                    >
                        <Input />
                    </Form.Item>
                    <Form.Item
                        label="系统类型"
                        {...formItemLayout}
                        name="osType"
                        rules={[
                            { required: true, message: "系统类型不能为空！" },
                        ]}
                    >
                        <Select>
                            <Option value="linux">Linux</Option>
                            <Option value="windows">Windows</Option>
                        </Select>
                    </Form.Item>
                    <Form.Item
                        label="CPU核数(个)"
                        {...formItemLayout}
                        name="cpu"
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
                        label="内存(G)"
                        {...formItemLayout}
                        name="memory"
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
                        label="磁盘(G)"
                        {...formItemLayout}
                        name="disk"
                        rules={[{ type: "integer" }]}
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
                        name="expiredTime"
                    >
                        <DatePicker format="YYYY-MM-DD" />
                    </Form.Item>
                    <Form.Item
                        label="SSH登录端口"
                        {...formItemLayout}
                        name="sshPort"
                        rules={[{ type: "integer" }]}
                    >
                        <InputNumber />
                    </Form.Item>
                    <Form.Item
                        label="SSH登录用户"
                        {...formItemLayout}
                        name="sshUser"
                    >
                        <Input />
                    </Form.Item>
                    <Form.Item
                        label="SSH登录密码"
                        {...formItemLayout}
                        name="sshPwd"
                    >
                        <Input.Password
                            iconRender={(visible) =>
                                visible ? (
                                    <EyeTwoTone />
                                ) : (
                                    <EyeInvisibleOutlined />
                                )
                            }
                        />
                    </Form.Item>
                    <Form.Item label="机器标签" {...formItemLayout} name="tags">
                        <Select mode="tags" tokenSeparators={[","]}></Select>
                    </Form.Item>
                </Form>
            </Modal>
        );
    }
}

class ServerContent extends Component {
    constructor(props) {
        super(props);
        this.handlePostServerInfoSubmit = this.handlePostServerInfoSubmit.bind(
            this,
        );
        this.handlePostServerInfoCancel = this.handlePostServerInfoCancel.bind(
            this,
        );
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
                },
                {
                    title: "实例名称",
                    dataIndex: "InstanceName",
                    key: "InstanceName",
                    width: 200,
                    textWrap: "word-break",
                    render: (value) => {
                        return (
                            <Tooltip title={value}>
                                <Text
                                    ellipsis={true}
                                    style={{ width: "200px" }}
                                >
                                    {value}
                                </Text>
                            </Tooltip>
                        );
                    },
                },
                {
                    title: "ip",
                    dataIndex: "ip",
                    key: "ip",
                    width: 200,
                    render: (value, record) => {
                        let innerContent = "";
                        let privateContent = "";
                        let publicContent = "";
                        if (record.InnerIpAddress) {
                            innerContent = (
                                <div>
                                    <Paragraph
                                        style={{
                                            marginBottom: 0,
                                            display: "inline",
                                        }}
                                        copyable={record.InnerIpAddress !== ""}
                                    >
                                        {record.InnerIpAddress}
                                    </Paragraph>
                                    (内网)
                                </div>
                            );
                        }
                        if (record.PrivateIpAddress) {
                            privateContent = (
                                <div>
                                    <Paragraph
                                        style={{
                                            marginBottom: 0,
                                            display: "inline",
                                        }}
                                        copyable={
                                            record.PrivateIpAddress !== ""
                                        }
                                    >
                                        {record.PrivateIpAddress}
                                    </Paragraph>
                                    (私有)
                                </div>
                            );
                        }
                        if (record.PublicIpAddress) {
                            publicContent = (
                                <div>
                                    <Paragraph
                                        style={{
                                            marginBottom: 0,
                                            display: "inline",
                                        }}
                                        copyable={record.PublicIpAddress !== ""}
                                    >
                                        {record.PublicIpAddress}
                                    </Paragraph>
                                    (外网)
                                </div>
                            );
                        }
                        return (
                            <div className="ip_column">
                                {innerContent}
                                {privateContent}
                                {publicContent}
                            </div>
                        );
                    },
                },
                {
                    title: "实例标签",
                    dataIndex: "Tags",
                    key: "Tags",
                    render: (value) => {
                        let tagList = value.split(",");
                        return (
                            <div>
                                {tagList.map((item) => {
                                    return <Tag color="geekblue">{item}</Tag>;
                                })}
                            </div>
                        );
                    },
                },
                {
                    title: "配置",
                    dataIndex: "配置",
                    key: "配置",
                    width: 100,
                    render: (value, record) => {
                        let cpuContent = (
                            <Paragraph
                                style={{ marginBottom: 0, display: "inline" }}
                            >
                                {record.Cpu}核
                            </Paragraph>
                        );
                        let memoryContent = (
                            <Paragraph
                                style={{ marginBottom: 0, display: "inline" }}
                            >
                                {record.Memory}G
                            </Paragraph>
                        );
                        let trafficType = "";
                        if (record.InternetChargeType === "PayByTraffic") {
                            trafficType = "流量";
                        }
                        if (record.InternetChargeType === "PayByBandwidth") {
                            trafficType = "带宽";
                        }
                        let trafficOutContent = (
                            <div>
                                <Paragraph
                                    style={{
                                        marginBottom: 0,
                                        display: "inline",
                                    }}
                                >
                                    {record.InternetMaxBandwidthOut}Mbps(
                                    {trafficType})
                                </Paragraph>
                            </div>
                        );
                        return (
                            <div className="ip_column">
                                {cpuContent} &nbsp;
                                {memoryContent}
                                {trafficOutContent}
                            </div>
                        );
                    },
                },
                {
                    title: "系统类型",
                    dataIndex: "OSType",
                    key: "OSType",
                    align: "center",
                    width: 80,
                    render: (value, record) => {
                        let status = "error";
                        if (record.Status === "Running") {
                            status = "processing";
                        }
                        if (value === "windows") {
                            return (
                                <div>
                                    <WindowsSvg />
                                    <Badge
                                        status={status}
                                        style={{
                                            marginLeft: "5px",
                                            position: "relative",
                                            top: "-10px",
                                        }}
                                    />
                                </div>
                            );
                        } else if (value === "linux") {
                            return (
                                <div>
                                    <LinuxSvg />
                                    <Badge
                                        status={status}
                                        style={{
                                            marginLeft: "5px",
                                            position: "relative",
                                            top: "-10px",
                                        }}
                                    />
                                </div>
                            );
                        } else {
                            return <Text ellipsis={true}>{value}</Text>;
                        }
                    },
                },
                {
                    title: "区域",
                    dataIndex: "ZoneId",
                    key: "ZoneId",
                    width: 120,
                    render: (value) => {
                        return <Text ellipsis={true}>{value}</Text>;
                    },
                },
                {
                    title: "过期时间",
                    dataIndex: "ExpiredTime",
                    key: "ExpiredTime",
                    width: 120,
                    render: (value) => {
                        return <Text ellipsis={true}>{value}</Text>;
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
                                    disabled={
                                        !this.props.aclAuthMap[
                                            "PUT:/cloud/servers"
                                        ]
                                    }
                                    onClick={this.serverEdit.bind(this, record)}
                                >
                                    编辑
                                </Button>
                                <Divider type="vertical" />
                                <Popconfirm
                                    title="确定删除该项内容?"
                                    onConfirm={this.serverDelete.bind(
                                        this,
                                        record,
                                    )}
                                    okText="确认"
                                    cancelText="取消"
                                    disabled={
                                        !this.props.aclAuthMap[
                                            "DELETE:/cloud/servers"
                                        ]
                                    }
                                >
                                    <Button
                                        type="danger"
                                        size="small"
                                        disabled={
                                            !this.props.aclAuthMap[
                                                "DELETE:/cloud/servers"
                                            ]
                                        }
                                    >
                                        删除
                                    </Button>
                                </Popconfirm>
                            </div>
                        );
                    },
                },
            ],
            tableLoading: false,
            webSocketReady: false,
            chartData: [],
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
            drawerVisible: false,
            drawerPlacement: "left",
            instanceId: "",
            timeTagValue: "1h",
            metricTagValue: "CPUUtilization",
            chartFormat: "%",
            currentServerDetail: {},
            msgContent: "",
            server_info_modal_visible: false,
            queryExpiredTime: null,
            queryKeyword: "",
            queryCloudAccount: "0",
            queryManageUser: "0",
            queryDefineGroup: "",
            cloudAccountList: [],
            selectedRowKeys: [],
            idsList: [],
            updateMode: "single",
            extraInfoData: {},
            cmdInput: "",
            cmdResult: [],
            batch_ssh_modal_visible: false,
            cmdRunning: false,
            selectedIds: [],
            serverNodeTreeData: [
                {
                    title: "服务器资源",
                    key: "服务器资源",
                    disabled: true,
                    children: [],
                },
            ],
        };
        this.timer = null;
        this.terminal = null;
        this.rws = null;
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

    componentWillUnmount() {
        if (this.rws !== null) {
            this.rws.close();
        }
        if (this.terminal !== null) {
            this.terminal.dispose();
        }
    }

    initWsConnection() {
        this.rws.addEventListener("open", () => {
            console.log("connect success");
        });

        this.rws.addEventListener("close", () => {
            console.log("close");
        });

        this.rws.addEventListener("message", (e) => {
            console.log("message: ", e);
        });

        this.rws.addEventListener("error", () => {
            console.log("error");
        });
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

    serverEdit(data) {
        getCloudServerDetail(data.ID)
            .then((res) => {
                if (res["code"] !== 0) {
                    message.error(res["msg"]);
                } else {
                    let extraInfoData = {
                        instanceId: res.data["InstanceId"],
                        innerIpAddress: res.data["InnerIpAddress"],
                        publicIpAddress: res.data["PublicIpAddress"],
                        privateIpAddress: res.data["PrivateIpAddress"],
                        instanceName: res.data["InstanceName"],
                        cpu: res.data["Cpu"],
                        memory: (res.data["Memory"] / 1024).toString(),
                        expiredTime:
                            res.data["ExpiredTime"] !== ""
                                ? moment(res.data["ExpiredTime"])
                                : "",
                        resForm:
                            res.data["DataStatus"] === 1
                                ? "阿里云"
                                : "手动添加",
                        sshPort: res.data["SshPort"] + "",
                        sshUser: res.data["SshUser"],
                        sshPwd: res.data["SshPwd"],
                        tags: res.data["Tags"].split(","),
                    };
                    this.setState({
                        extraInfoModalVisible: true,
                        ecsId: data.ID,
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
            targetId = String(this.state.ecsId);
        } else {
            targetId = this.state.idsList.join(",");
        }
        putCloudServer({
            ...data,
            id: targetId,
        })
            .then((res) => {
                if (res.code === 0) {
                    this.setState({
                        extraInfoModalVisible: false,
                        selectedRowKeys: [],
                    });
                    message.success("修改成功");
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

    serverDelete(data) {
        deleteCloudServer(data.ID)
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
                    : this.state.queryExpiredTime.format("YYYY-MM-DD HH:mm:ss"),
            queryKeyword: this.state.queryKeyword,
            queryCloudAccount: this.state.queryCloudAccount,
        };
        getCloudServers(queryParams)
            .then((res) => {
                const pagination = this.state.pagination;
                pagination.total = parseInt(res.data.total);
                pagination.page = parseInt(res.data.page);
                pagination.showTotal(parseInt(res.data.total));
                this.setState({
                    pagination: { ...pagination },
                });
                let data = res["data"]["servers"];
                let tableData = [];
                for (let i = 0; i < data.length; i++) {
                    tableData.push({
                        key: data[i]["ID"],
                        ID: data[i]["ID"],
                        Memory: data[i]["Memory"] / 1024,
                        Cpu: data[i]["Cpu"],
                        HostName: data[i]["HostName"],
                        InstanceId: data[i]["InstanceId"],
                        InnerIpAddress: data[i]["InnerIpAddress"],
                        PublicIpAddress: data[i]["PublicIpAddress"],
                        PrivateIpAddress: data[i]["PrivateIpAddress"],
                        InternetMaxBandwidthIn:
                            data[i]["InternetMaxBandwidthIn"],
                        InternetMaxBandwidthOut:
                            data[i]["InternetMaxBandwidthOut"],
                        InternetChargeType: data[i]["InternetChargeType"],
                        InstanceName: data[i]["InstanceName"],
                        OSType: data[i]["OSType"],
                        ZoneId: data[i]["ZoneId"],
                        OSName: data[i]["OSName"],
                        ExpiredTime: moment(data[i]["ExpiredTime"]).format(
                            "YYYY-MM-DD",
                        ),
                        Status: data[i]["Status"],
                        CloudAccountName: data[i]["CloudAccountName"],
                        DataStatus: data[i]["DataStatus"],
                        SshPort: data[i]["SshPort"] + "",
                        SshUser: data[i]["SshUser"],
                        SshPwd: data[i]["SshPwd"],
                        Tags: data[i]["Tags"],
                    });
                }
                this.setState({ tableData: tableData, tableLoading: false });
            })
            .catch((err) => {
                message.error(err);
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
        getCloudMonitorEcs(instanceId, timeTagValue, metricTagValue)
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
        getCloudServerDetail(this.state.currentServerDetail.ID).then((res) => {
            if (res["code"] !== 0) {
                message.error(res["msg"]);
            } else {
                this.setState(
                    {
                        serverDetailLoading: true,
                        currentServerDetail: res["data"],
                    },
                    () => {
                        this.setState({ serverDetailLoading: false });
                    },
                );
            }
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
            case "CPUUtilization":
                this.setState({ chartFormat: "%" });
                break;
            case "memory_usedutilization":
                this.setState({ chartFormat: "%" });
                break;
            case "diskusage_utilization":
                this.setState({ chartFormat: "%" });
                break;
            case "cpu_total":
                this.setState({ chartFormat: "%" });
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

    // 新增自定义机器信息
    handleAdd = () => {
        this.setState({ server_info_modal_visible: true, ecsId: 0 });
    };

    handlePostServerInfoSubmit() {
        this.formRef.current.validateFields().then((values) => {
            postCloudServer({
                ...values,
                createTime: values.createTime.format("YYYY-MM-DD HH:mm:ss"),
                expiredTime:
                    values.expiredTime === undefined
                        ? undefined
                        : values.expiredTime.format("YYYY-MM-DD HH:mm:ss"),
            })
                .then((res) => {
                    if (res.code === 0) {
                        message.success(
                            "添加成功，请到权限管理中增加访问权限！",
                        );
                        this.setState({ server_info_modal_visible: false });
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

    handlePostServerInfoCancel() {
        this.setState({ server_info_modal_visible: false });
    }

    handleBatchSSH = () => {
        getCloudServersAllData().then((res) => {
            if (res.code === 0) {
                let childrenData = [];
                for (let i = 0; i < res.data.length; i++) {
                    childrenData.push({
                        key: res.data[i]["ID"],
                        title:
                            res.data[i]["InstanceName"] +
                            " - " +
                            res.data[i]["PublicIpAddress"] +
                            " - " +
                            res.data[i]["Tags"],
                    });
                }
                this.setState({
                    serverNodeTreeData: [
                        {
                            title: "服务器资源",
                            key: "服务器资源",
                            disabled: true,
                            children: childrenData,
                        },
                    ],
                });
            } else {
                message.error(res.msg);
            }
        });
        this.setState({ batch_ssh_modal_visible: true });
    };

    onCloseBatchSshDrawer = () => {
        this.setState({ batch_ssh_modal_visible: false });
    };

    onSelectServerNode = (selectedKeys, info) => {
        this.setState({ selectedIds: selectedKeys });
    };

    submitBatchSsh = () => {
        this.setState({ cmdRunning: true });
        postCloudServerBatchSSH({
            ids: this.state.selectedIds,
            command: this.state.cmdInput,
        }).then((res) => {
            if (res.code === 0) {
                this.setState({ cmdResult: res.data });
            } else {
                message.error(res.msg);
            }
            this.setState({ cmdRunning: false });
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

    initTerminal = () => {
        let dom = document.getElementById("server_terminal");
        if (dom) {
            this.rws = new ReconnectingWebSocket(
                WSBase +
                    "ws/cloud/ssh?serverId=" +
                    this.state.currentServerDetail.ID +
                    "&token=" +
                    localStorage.getItem("ops_token"),
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
            const attachAddon = new AttachAddon(this.rws);
            this.terminal.loadAddon(attachAddon);
            const fitAddon = new FitAddon();
            this.terminal.loadAddon(fitAddon);
            this.terminal.open(document.getElementById("server_terminal"));
            this.terminal.writeln("Welcome to use Web Terminal.");
            this.terminal.prompt();
            fitAddon.fit();
            this.initWsConnection();
            this.terminal.focus();
            if (!this.timer) {
                clearTimeout(this.timer);
            }
        } else {
            this.timer = setTimeout(this.initTerminal, 0);
        }
    };

    detailTabChange(key) {
        if (key === "3") {
            this.initTerminal();
        }
    }

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
                <OpsBreadcrumbPath
                    pathData={["云资源", "云服务器", "服务器列表"]}
                />

                <Row style={{ padding: "0px 0px 10px 0px" }}>
                    <Col span={3} className="col-span">
                        <DatePicker
                            style={{ width: "100%" }}
                            defaultValue={this.state.queryExpiredTime}
                            placeholder="到期截止时间"
                            onChange={this.onExpiredDateChange}
                        />
                    </Col>
                    <Col span={5} className="col-span">
                        <Input
                            placeholder="输入实例id\ip\实例名称\标签查询"
                            value={this.state.queryKeyword}
                            onChange={this.keywordOnChange}
                        />
                    </Col>
                    <Col span={4} className="col-span">
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
                            type="primary"
                            icon={<SearchOutlined />}
                            onClick={this.handleQuery}
                        >
                            查 询
                        </Button>
                    </Col>
                    <Col span={2} className="col-span">
                        <Button
                            icon={<PlusCircleOutlined />}
                            onClick={this.handleAdd}
                            disabled={
                                !this.props.aclAuthMap["POST:/cloud/servers"]
                            }
                        >
                            新 增
                        </Button>
                    </Col>
                    <Col span={3} className="col-span">
                        <Button
                            icon={<CodeOutlined />}
                            onClick={this.handleBatchSSH}
                            disabled={
                                !this.props.aclAuthMap[
                                    "POST:/cloud/servers/batch/ssh"
                                ]
                            }
                        >
                            批量SSH
                        </Button>
                    </Col>
                </Row>

                <ServerInfoModal
                    formRef={this.formRef}
                    server_info_modal_visible={
                        this.state.server_info_modal_visible
                    }
                    handlePostServerInfoSubmit={this.handlePostServerInfoSubmit}
                    handlePostServerInfoCancel={this.handlePostServerInfoCancel}
                />

                {/*完善信息组件*/}
                <ExtraInfoModal
                    editData={this.state.extraInfoData}
                    resType="ecs"
                    updateMode={this.state.updateMode}
                    resFrom={this.state.resFrom}
                    visible={this.state.extraInfoModalVisible}
                    handleOk={this.handleExtraInfoOk}
                    handleCancel={this.handleExtraInfoCancel}
                />

                {/*云服务器列表*/}
                <Drawer
                    title="服务器详情及监控数据"
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
                        onChange={this.detailTabChange.bind(this)}
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
                                            onChange={this.handleTimeTagChange}
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
                                            value={this.state.metricTagValue}
                                            onChange={
                                                this.handleMetricTagChange
                                            }
                                        >
                                            <Radio.Button value="CPUUtilization">
                                                cpu使用率
                                            </Radio.Button>
                                            <Radio.Button value="memory_usedutilization">
                                                内存使用率
                                            </Radio.Button>
                                            <Radio.Button value="diskusage_utilization">
                                                磁盘使用率
                                            </Radio.Button>
                                            <Radio.Button value="cpu_total">
                                                平均负载
                                            </Radio.Button>
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
                                    <Descriptions.Item label="主机名">
                                        {
                                            this.state.currentServerDetail
                                                .HostName
                                        }
                                    </Descriptions.Item>
                                    <Descriptions.Item label="机器描述">
                                        {
                                            this.state.currentServerDetail
                                                .Description
                                        }
                                    </Descriptions.Item>
                                    <Descriptions.Item label="实例ID">
                                        {
                                            this.state.currentServerDetail
                                                .InstanceId
                                        }
                                    </Descriptions.Item>
                                    <Descriptions.Item label="内网IP">
                                        {
                                            this.state.currentServerDetail
                                                .InnerIpAddress
                                        }
                                    </Descriptions.Item>
                                    <Descriptions.Item label="外网IP">
                                        {
                                            this.state.currentServerDetail
                                                .PublicIpAddress
                                        }
                                    </Descriptions.Item>
                                    <Descriptions.Item label="私有IP">
                                        {
                                            this.state.currentServerDetail
                                                .PrivateIpAddress
                                        }
                                    </Descriptions.Item>
                                    <Descriptions.Item label="Cpu">
                                        {this.state.currentServerDetail.Cpu}核
                                    </Descriptions.Item>
                                    <Descriptions.Item label="内存">
                                        {this.state.currentServerDetail.Memory /
                                            1024}
                                        G
                                    </Descriptions.Item>
                                    <Descriptions.Item label="公网入带宽">
                                        {
                                            this.state.currentServerDetail
                                                .InternetMaxBandwidthIn
                                        }
                                        Mbps
                                    </Descriptions.Item>
                                    <Descriptions.Item label="公网出带宽">
                                        {
                                            this.state.currentServerDetail
                                                .InternetMaxBandwidthOut
                                        }
                                        Mbps
                                    </Descriptions.Item>
                                    <Descriptions.Item label="网络计费">
                                        {
                                            this.state.currentServerDetail
                                                .InternetChargeType
                                        }
                                    </Descriptions.Item>
                                    <Descriptions.Item label="创建时间">
                                        {
                                            this.state.currentServerDetail
                                                .CreationTime
                                        }
                                    </Descriptions.Item>
                                    <Descriptions.Item label="过期时间">
                                        {
                                            this.state.currentServerDetail
                                                .ExpiredTime
                                        }
                                    </Descriptions.Item>
                                    <Descriptions.Item label="镜像ID">
                                        {this.state.currentServerDetail.ImageId}
                                    </Descriptions.Item>
                                    <Descriptions.Item label="付费类型">
                                        {
                                            this.state.currentServerDetail
                                                .InstanceChargeType
                                        }
                                    </Descriptions.Item>
                                    <Descriptions.Item label="网络类型">
                                        {
                                            this.state.currentServerDetail
                                                .InstanceNetworkType
                                        }
                                    </Descriptions.Item>
                                    <Descriptions.Item label="实例类型">
                                        {
                                            this.state.currentServerDetail
                                                .InstanceType
                                        }
                                    </Descriptions.Item>
                                    <Descriptions.Item label="系统名称">
                                        {this.state.currentServerDetail.OSName}
                                    </Descriptions.Item>
                                </Descriptions>
                            </Spin>
                        </TabPane>
                        <TabPane tab="终端" key="3">
                            <div id="server_terminal" />
                        </TabPane>
                    </Tabs>
                </Drawer>

                <Drawer
                    title="批量远程操作"
                    placement="left"
                    closable={false}
                    width={800}
                    onClose={this.onCloseBatchSshDrawer}
                    visible={this.state.batch_ssh_modal_visible}
                >
                    <Title level={5}>选择待操作的服务器</Title>
                    <Tree
                        checkable
                        defaultExpandedKeys={["服务器资源"]}
                        onCheck={this.onSelectServerNode}
                        treeData={this.state.serverNodeTreeData}
                    />
                    <Title level={5}>输入命令: </Title>
                    <div style={{ marginBottom: 10, marginTop: 10 }}>
                        <CodeMirror
                            className="sqlEditor"
                            options={{
                                showCursorWhenSelecting: true,
                                option: {
                                    autofocus: true,
                                },
                                lineWrapping: true,
                            }}
                            value={this.state.cmdInput}
                            onBeforeChange={(editor, data, value) => {
                                this.setState({ cmdInput: value });
                            }}
                        />
                    </div>
                    <div>
                        <Button type="primary" onClick={this.submitBatchSsh}>
                            提交执行
                        </Button>
                    </div>
                    <Title level={5}>执行结果: </Title>
                    <Spin
                        tip="远程命令执行中..."
                        spinning={this.state.cmdRunning}
                    >
                        <pre
                            style={{
                                marginTop: 20,
                                minHeight: 200,
                                paddingLeft: 10,
                            }}
                            className="preJenkinsLog"
                        >
                            {this.state.cmdResult.map((item) => {
                                let serverName = "";
                                let result = "";
                                console.log(item);
                                for (var server in item) {
                                    serverName = server;
                                    result = item[server];
                                    break;
                                }
                                return (
                                    <div style={{ marginBottom: 10 }}>
                                        <div>目标机器: {serverName}</div>
                                        <div>
                                            执行结果: <br />
                                            {result}
                                        </div>
                                        <div>---------------------</div>
                                    </div>
                                );
                            })}
                        </pre>
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

export default ServerContent;
