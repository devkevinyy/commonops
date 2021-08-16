import React, { Component } from "react";
import {
    Button,
    Col,
    Input,
    Row,
    Table,
    Layout,
    DatePicker,
    Select,
    message,
    Modal,
    Form,
    Typography,
    Divider,
    Popconfirm,
    Tooltip,
} from "antd";
import OpsBreadcrumbPath from "../breadcrumb_path";
import {
    deleteCloudOtherRes,
    getCloudAccouts,
    getCloudOtherRes,
    postCloudOtherRes,
    putCloudOtherRes,
} from "../../api/cloud";
import { SearchOutlined, PlusCircleOutlined } from "@ant-design/icons";
import moment from "moment";
import ExtraInfoModal from "./common/extra_info_modal";

const { Option } = Select;
const { Text, Paragraph } = Typography;
const { Content } = Layout;

class OtherInfoModal extends Component {
    constructor(props) {
        super(props);
        this.switchChange = this.switchChange.bind(this);
        this.state = {
            sitesData: [],
            renewSwitch: true,
            cloudAccountList: [],
        };
    }

    componentDidMount() {
        this.loadCloudAccountData();
    }

    loadCloudAccountData() {
        getCloudAccouts(1, 100)
            .then((res) => {
                if (res.code === 0) {
                    this.setState({
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

    switchChange(value) {
        this.setState({ renewSwitch: value });
    }

    render() {
        const formItemLayout = {
            labelCol: { span: 6 },
            wrapperCol: { span: 17 },
        };
        let accountOptions;
        accountOptions = this.state.cloudAccountList.map((item) => {
            return (
                <Option key={item.id} value={item.id}>
                    {item.accountName}
                </Option>
            );
        });
        const { actionType } = this.props;
        let actionName = "新增";
        if (actionType === "update") {
            actionName = "编辑";
        }
        return (
            <Modal
                title={actionName + " - 资源信息"}
                destroyOnClose={true}
                visible={this.props.server_info_modal_visible}
                onOk={this.props.handlePostServerInfoSubmit}
                onCancel={this.props.handlePostServerInfoCancel}
                okText="确认"
                cancelText="取消"
                centered={true}
                width={1000}
            >
                <Form ref={this.props.formRef}>
                    <Row>
                        <Col span={11} offset={1}>
                            <Form.Item
                                label="资源类型"
                                {...formItemLayout}
                                name="resType"
                                rules={[
                                    {
                                        required: true,
                                        message: "请选择资源类型",
                                    },
                                ]}
                            >
                                <Select>
                                    <Option value="MongoDB">MongoDB</Option>
                                    <Option value="MaxCompute">
                                        MaxCompute
                                    </Option>
                                    <Option value="SSL">SSL</Option>
                                    <Option value="带宽">带宽</Option>
                                    <Option value="后付费">后付费</Option>
                                    <Option value="弹性IP">弹性IP</Option>
                                    <Option value="Memcached">Memcached</Option>
                                    <Option value="Kafka">Kafka</Option>
                                    <Option value="NAS">NAS</Option>
                                    <Option value="PolarDB">PolarDB</Option>
                                    <Option value="HBase">HBase</Option>
                                </Select>
                            </Form.Item>
                        </Col>
                        <Col span={11}>
                            <Form.Item
                                label="实例ID"
                                {...formItemLayout}
                                name="instanceId"
                            >
                                <Input />
                            </Form.Item>
                        </Col>
                    </Row>
                    <Row>
                        <Col span={11} offset={1}>
                            <Form.Item
                                label="实例名称"
                                {...formItemLayout}
                                name="instanceName"
                                rules={[
                                    {
                                        required: true,
                                        message: "请输入实例名称",
                                    },
                                ]}
                            >
                                <Input />
                            </Form.Item>
                        </Col>
                        <Col span={11}>
                            <Form.Item
                                label="连接地址"
                                {...formItemLayout}
                                name="connections"
                            >
                                <Input />
                            </Form.Item>
                        </Col>
                    </Row>

                    <Row>
                        <Col span={11} offset={1}>
                            <Form.Item
                                label="区域"
                                {...formItemLayout}
                                name="region"
                            >
                                <Input />
                            </Form.Item>
                        </Col>
                        <Col span={11}>
                            <Form.Item
                                label="Engine"
                                {...formItemLayout}
                                name="engine"
                            >
                                <Input />
                            </Form.Item>
                        </Col>
                    </Row>

                    <Row>
                        <Col span={11} offset={1}>
                            <Form.Item
                                label="CPU"
                                {...formItemLayout}
                                name="cpu"
                            >
                                <Input />
                            </Form.Item>
                        </Col>
                        <Col span={11}>
                            <Form.Item
                                label="磁盘"
                                {...formItemLayout}
                                name="disk"
                            >
                                <Input />
                            </Form.Item>
                        </Col>
                    </Row>

                    <Row>
                        <Col span={11} offset={1}>
                            <Form.Item
                                label="内存"
                                {...formItemLayout}
                                name="memory"
                            >
                                <Input />
                            </Form.Item>
                        </Col>
                        <Col span={11}>
                            <Form.Item
                                label="带宽"
                                {...formItemLayout}
                                name="bandwidth"
                            >
                                <Input />
                            </Form.Item>
                        </Col>
                    </Row>

                    <Row>
                        <Col span={11} offset={1}>
                            <Form.Item
                                label="创建时间"
                                {...formItemLayout}
                                name="createTime"
                                rules={[
                                    {
                                        required: true,
                                        message: "请输入创建时间",
                                    },
                                ]}
                            >
                                <DatePicker
                                    format="YYYY-MM-DD"
                                    style={{ width: "100%" }}
                                />
                            </Form.Item>
                        </Col>
                        <Col span={11}>
                            <Form.Item
                                label="过期时间"
                                {...formItemLayout}
                                name="expiredTime"
                                rules={[
                                    {
                                        required: true,
                                        message: "请输入过期时间",
                                    },
                                ]}
                            >
                                <DatePicker
                                    format="YYYY-MM-DD"
                                    style={{ width: "100%" }}
                                />
                            </Form.Item>
                        </Col>
                    </Row>
                    <Row>
                        <Col span={11} offset={1}>
                            <Form.Item
                                label="云账号"
                                {...formItemLayout}
                                name="cloudAccountId"
                            >
                                <Select>{accountOptions}</Select>
                            </Form.Item>
                        </Col>
                    </Row>
                </Form>
            </Modal>
        );
    }
}

class OtherContent extends Component {
    constructor(props) {
        super(props);
        this.keywordOnChange = this.keywordOnChange.bind(this);
        this.handleQuery = this.handleQuery.bind(this);
        this.handlePostOtherInfoSubmit = this.handlePostOtherInfoSubmit.bind(
            this,
        );
        this.handlePostOtherInfoCancel = this.handlePostOtherInfoCancel.bind(
            this,
        );
        this.changePage = this.changePage.bind(this);
        this.handleAdd = this.handleAdd.bind(this);
        this.handleResTypeChange = this.handleResTypeChange.bind(this);
        this.handleCloudAccountChange = this.handleCloudAccountChange.bind(
            this,
        );
        this.formRef = React.createRef();
        this.state = {
            columns: [
                {
                    title: "Id",
                    dataIndex: "ID",
                    key: "ID",
                    className: "small_font",
                    render: (value) => {
                        return <Text ellipsis={true}>{value}</Text>;
                    },
                },
                {
                    title: "资源类型",
                    dataIndex: "ResType",
                    key: "ResType",
                    className: "small_font",
                    render: (value) => {
                        return <Text ellipsis={true}>{value}</Text>;
                    },
                },
                {
                    title: "实例名称",
                    dataIndex: "InstanceName",
                    key: "InstanceName",
                    className: "small_font",
                    render: (value) => {
                        return (
                            <Tooltip title={value}>
                                <Text ellipsis={true}>{value}</Text>
                            </Tooltip>
                        );
                    },
                },
                {
                    title: "云账号",
                    dataIndex: "CloudAccountName",
                    key: "CloudAccountName",
                    className: "small_font",
                    render: (value) => {
                        return <Text ellipsis={true}>{value}</Text>;
                    },
                },
                {
                    title: "实例id",
                    dataIndex: "InstanceId",
                    key: "InstanceId",
                    className: "small_font",
                    render: (value) => {
                        return <Text ellipsis={true}>{value}</Text>;
                    },
                },
                {
                    title: "访问地址",
                    dataIndex: "Connections",
                    key: "Connections",
                    className: "small_font",
                    render: (value) => {
                        return (
                            <Paragraph
                                style={{ marginBottom: 0 }}
                                copyable={value !== ""}
                            >
                                {value}
                            </Paragraph>
                        );
                    },
                },
                {
                    title: "Region",
                    dataIndex: "Region",
                    key: "Region",
                    className: "small_font",
                    render: (value) => {
                        return <Text>{value}</Text>;
                    },
                },
                {
                    title: "Engine",
                    dataIndex: "Engine",
                    key: "Engine",
                    className: "small_font",
                    render: (value) => {
                        return <Text>{value}</Text>;
                    },
                },
                {
                    title: "Cpu",
                    dataIndex: "Cpu",
                    key: "Cpu",
                    className: "small_font",
                    render: (value) => {
                        return <Text>{value}</Text>;
                    },
                },
                {
                    title: "带宽",
                    dataIndex: "Bandwidth",
                    key: "Bandwidth",
                    className: "small_font",
                    render: (value) => {
                        return <Text ellipsis={true}>{value}</Text>;
                    },
                },
                {
                    title: "创建时间",
                    dataIndex: "CreateTime",
                    key: "CreateTime",
                    className: "small_font",
                    render: (value) => {
                        return <Text ellipsis={true}>{value}</Text>;
                    },
                },
                {
                    title: "过期时间",
                    dataIndex: "ExpiredTime",
                    key: "ExpiredTime",
                    className: "small_font",
                    render: (value) => {
                        return <Text ellipsis={true}>{value}</Text>;
                    },
                },
                {
                    title: "操作",
                    key: "operation",
                    fixed: "right",
                    align: "center",
                    className: "small_font",
                    render: (text, record) => {
                        return (
                            <div>
                                <Button
                                    type="info"
                                    size="small"
                                    onClick={this.serverEdit.bind(this, record)}
                                    disabled={!this.props.aclAuthMap["PUT:/cloud/other"]}
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
                                    disabled={!this.props.aclAuthMap["DELETE:/cloud/other"]}
                                >
                                    <Button type="danger" size="small" disabled={!this.props.aclAuthMap["DELETE:/cloud/other"]}>
                                        删除
                                    </Button>
                                </Popconfirm>
                            </div>
                        );
                    },
                },
            ],
            tableLoading: false,
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
            actionType: "add",
            editId: 0,
            extraInfoModalVisible: false,
            other_info_modal_visible: false,
            queryExpiredTime: null,
            queryKeyword: "",
            queryCloudAccount: "0",
            queryManageUser: "0",
            cloudAccountList: [],
            queryResType: "所有",
            selectedRowKeys: [],
            idsList: [],
            updateMode: "single",
        };
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

    refreshTableData = () => {
        this.setState({ tableLoading: true });
        const queryParams = {
            page: this.state.pagination.page,
            size: this.state.pagination.pageSize,
            queryKeyword: this.state.queryKeyword,
            queryResType: this.state.queryResType,
            queryExpiredTime:
                this.state.queryExpiredTime === null
                    ? null
                    : this.state.queryExpiredTime.format("YYYY-MM-DD HH:mm:ss"),
            queryCloudAccount: this.state.queryCloudAccount,
        };
        getCloudOtherRes(queryParams)
            .then((res) => {
                const pagination = this.state.pagination;
                pagination.total = parseInt(res.data.total);
                pagination.page = parseInt(res.data.page);
                pagination.showTotal(parseInt(res.data.total));
                this.setState({
                    pagination: { ...pagination },
                });
                let data = res["data"]["otherRes"];
                let tableData = [];
                for (let i = 0; i < data.length; i++) {
                    tableData.push({
                        key: data[i]["ID"],
                        ID: data[i]["ID"],
                        CloudAccountId: data[i]["CloudAccountId"],
                        CloudAccountName: data[i]["CloudAccountName"],
                        ResType: data[i]["ResType"],
                        InstanceId: data[i]["InstanceId"],
                        InstanceName: data[i]["InstanceName"],
                        Connections: data[i]["Connections"],
                        Region: data[i]["Region"],
                        Engine: data[i]["Engine"],
                        Cpu: data[i]["Cpu"],
                        Bandwidth: data[i]["Bandwidth"],
                        RenewStatus: data[i]["RenewStatus"],
                        RenewSiteId: data[i]["RenewSiteId"],
                        BankAccount: data[i]["BankAccount"],
                        CreateTime: moment(data[i]["CreateTime"]).format(
                            "YYYY-MM-DD",
                        ),
                        ExpiredTime: moment(data[i]["ExpiredTime"]).format(
                            "YYYY-MM-DD",
                        ),
                    });
                }
                this.setState({ tableData: tableData, tableLoading: false });
            })
            .catch((err) => {
                message.error(err);
            });
    };

    handleAdd() {
        this.setState({
            other_info_modal_visible: true,
            actionType: "add",
            editId: 0,
        });
    }

    handleQuery() {
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
    }

    onExpiredDateChange = (moment) => {
        if (moment == null) {
            this.setState({ queryExpiredTime: null });
        } else {
            this.setState({ queryExpiredTime: moment });
        }
    };

    keywordOnChange(e) {
        this.setState({ queryKeyword: e.target.value });
    }

    handleResTypeChange(value) {
        this.setState({ queryResType: value });
    }

    handleCloudAccountChange(value) {
        this.setState({ queryCloudAccount: value });
    }

    serverEdit(data) {
        const that = this;
        this.setState(
            {
                other_info_modal_visible: true,
                editId: data.ID,
                actionType: "update",
                updateMode: "single",
            },
            () => {
                setTimeout(() => {
                    that.formRef.current.setFieldsValue({
                        cloudAccountId:
                            data.CloudAccountId !== 0
                                ? data.CloudAccountId
                                : null,
                        resType: data.ResType,
                        instanceId: data.InstanceId,
                        instanceName: data.InstanceName,
                        connections: data.Connections,
                        region: data.Region,
                        engine: data.Engine,
                        cpu: data.Cpu,
                        bandwidth: data.Bandwidth,
                        createTime:
                            data.CreateTime !== ""
                                ? moment(data.CreateTime, "YYYY-MM-DD")
                                : undefined,
                        expiredTime:
                            data.ExpiredTime !== ""
                                ? moment(data.ExpiredTime, "YYYY-MM-DD")
                                : undefined,
                    });
                }, 500);
            },
        );
    }

    serverDelete(data) {
        deleteCloudOtherRes(data.ID)
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

    handlePostOtherInfoSubmit() {
        this.formRef.current.validateFields().then((values) => {
            postCloudOtherRes({
                ...values,
                id: this.state.editId,
                createTime: values.createTime.format("YYYY-MM-DD"),
                expiredTime:
                    values.expiredTime === undefined
                        ? undefined
                        : values.expiredTime.format("YYYY-MM-DD"),
            })
                .then((res) => {
                    if (res.code === 0) {
                        if (this.state.actionType === "add") {
                            message.success(
                                "添加成功，请到权限管理中增加访问权限！",
                            );
                        } else {
                            message.success("编辑成功");
                        }
                        this.setState({ other_info_modal_visible: false });
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

    handlePostOtherInfoCancel() {
        this.setState({ other_info_modal_visible: false });
    }

    handleExtraInfoOk(data) {
        let targetId = "";
        if (this.state.updateMode === "single") {
            targetId = String(this.state.editId);
        } else {
            targetId = this.state.idsList.join(",");
        }
        putCloudOtherRes({
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
                <OpsBreadcrumbPath
                    pathData={["云资源", "其它资源", "其它资源汇总列表"]}
                />
                <Row style={{ padding: "0px 0px 10px 0px" }}>
                    <Col span={4} className="col-span">
                        <DatePicker
                            style={{ width: "100%" }}
                            defaultValue={this.state.queryExpiredTime}
                            placeholder="到期截止时间"
                            onChange={this.onExpiredDateChange}
                        />
                    </Col>
                    <Col span={4} className="col-span">
                        <Input
                            placeholder="输入实例id/名称/地址查询"
                            value={this.state.queryKeyword}
                            onChange={this.keywordOnChange}
                        />
                    </Col>
                    <Col span={3} className="col-span">
                        <Select
                            defaultValue={this.state.queryResType}
                            style={{ width: "100%" }}
                            onChange={this.handleResTypeChange}
                        >
                            <Option value="所有">所有</Option>
                            <Option value="MongoDB">MongoDB</Option>
                            <Option value="MaxCompute">MaxCompute</Option>
                            <Option value="SSL">SSL</Option>
                            <Option value="带宽">带宽</Option>
                            <Option value="后付费">后付费</Option>
                            <Option value="弹性IP">弹性IP</Option>
                            <Option value="Memcached">Memcached</Option>
                            <Option value="Kafka">Kafka</Option>
                            <Option value="NAS">NAS</Option>
                            <Option value="PolarDB">PolarDB</Option>
                            <Option value="HBase">HBase</Option>
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
                            disabled={!this.props.aclAuthMap["POST:/cloud/other"]}
                        >
                            新 增
                        </Button>
                    </Col>
                </Row>

                {/*完善信息组件*/}
                <ExtraInfoModal
                    wrappedComponentRef={(form) =>
                        (this.extraInfoFormRef = form)
                    }
                    visible={this.state.extraInfoModalVisible}
                    handleOk={this.handleExtraInfoOk}
                    handleCancel={this.handleExtraInfoCancel}
                />

                <OtherInfoModal
                    formRef={this.formRef}
                    server_info_modal_visible={
                        this.state.other_info_modal_visible
                    }
                    actionType={this.state.actionType}
                    handlePostServerInfoSubmit={this.handlePostOtherInfoSubmit}
                    handlePostServerInfoCancel={this.handlePostOtherInfoCancel}
                />

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

export default OtherContent;
