import React, { Component } from "react";
import {
    Button,
    Col,
    Divider,
    Layout,
    message,
    Row,
    Select,
    Table,
    Typography,
    Form,
    Drawer,
    Modal,
    Input,
    Tag,
} from "antd";
import OpsBreadcrumbPath from "../breadcrumb_path";
import { getCloudAccouts } from "../../api/cloud";
import { SearchOutlined, PlusCircleOutlined } from "@ant-design/icons";
import {
    deleteDnsDomainRecord,
    getDnsDomainHistoryListData,
    getDnsDomainListData,
    getDnsDomainRecordListData,
    postDnsDomain,
    postDnsDomainRecord,
    postDnsDomainRecordStatus,
    postDnsDomainRecordUpdate,
} from "../../api/dns_api";
const { Content } = Layout;

const { Text } = Typography;
const { Option } = Select;

class DomainForm extends Component {
    render() {
        const formItemLayout = {
            labelCol: { span: 6 },
            wrapperCol: { span: 17 },
        };
        return (
            <Form ref={this.props.formRef} {...formItemLayout}>
                <Form.Item
                    label="域名"
                    name="domainName"
                    rules={[{ required: true, message: "域名不能为空！" }]}
                >
                    <Input />
                </Form.Item>
            </Form>
        );
    }
}

class DomainRecordForm extends Component {
    render() {
        const formItemLayout = {
            labelCol: { span: 6 },
            wrapperCol: { span: 17 },
        };
        return (
            <Form ref={this.props.formRef} {...formItemLayout}>
                <Form.Item
                    label="记录类型"
                    name="rType"
                    rules={[{ required: true, message: "记录类型不能为空！" }]}
                >
                    <Select style={{ width: "100%" }}>
                        <Option value="A">A-将域名指向一个IPV4地址</Option>
                        <Option value="CNAME">
                            CNAME-将域名指向另一个域名
                        </Option>
                        <Option value="AAAA">
                            AAAA-将域名指向一个IPV6地址
                        </Option>
                        <Option value="NS">
                            NS-将子域名指向其它DNS服务器解析
                        </Option>
                        <Option value="MX">MX-将域名指向邮件服务器地址</Option>
                        <Option value="SRV">
                            SRV-记录提供特定的服务的服务器
                        </Option>
                        <Option value="TXT">
                            TXT-文本长度限制512,通常做SPF记录
                        </Option>
                        <Option value="CAA">CAA-CA证书颁发机构授权校验</Option>
                        <Option value="REDIRECT_URL">
                            显性URL-将域名重定向到另外一个地址
                        </Option>
                        <Option value="FORWARD_URL">
                            隐性URL-与显性URL类似，但是会隐藏真实目标地址
                        </Option>
                    </Select>
                </Form.Item>
                <Form.Item
                    label="主机记录"
                    name="rr"
                    rules={[{ required: true, message: "主机记录不能为空！" }]}
                >
                    <Input addonAfter={"." + this.props.domainName} />
                </Form.Item>
                <Form.Item
                    label="记录值"
                    name="rValue"
                    rules={[{ required: true, message: "记录值不能为空！" }]}
                >
                    <Input />
                </Form.Item>
            </Form>
        );
    }
}

class DomainManageContent extends Component {
    constructor(props) {
        super(props);
        this.onShowSizeChange = this.onShowSizeChange.bind(this);
        this.queryDomainList = this.queryDomainList.bind(this);
        this.handleCloudAccountChange = this.handleCloudAccountChange.bind(
            this,
        );
        this.refreshTableData = this.refreshTableData.bind(this);
        this.domainFormRef = React.createRef();
        this.recordFormRef = React.createRef();
        this.state = {
            columns: [
                {
                    title: "域名",
                    dataIndex: "DomainName",
                    key: "DomainName",
                    className: "small_font",
                    align: "center",
                    render: (value) => {
                        return (
                            <Text
                                copyable={true}
                                ellipsis={true}
                                style={{ width: "100%" }}
                            >
                                {value}
                            </Text>
                        );
                    },
                },
                {
                    title: "记录数",
                    dataIndex: "RecordCount",
                    key: "RecordCount",
                    className: "small_font",
                    align: "center",
                    render: (value) => {
                        return (
                            <Text ellipsis={true} style={{ width: "100%" }}>
                                {value}
                            </Text>
                        );
                    },
                },
                {
                    title: "DNS Server",
                    dataIndex: "DnsServers",
                    key: "DnsServers",
                    className: "small_font",
                    align: "center",
                    render: (value, record) => {
                        return (
                            <Text ellipsis={true} style={{ width: "100%" }}>
                                {record["DnsServers"]["DnsServer"].join(",")}
                            </Text>
                        );
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
                                    type="link"
                                    size="small"
                                    onClick={this.showDomainRecord.bind(
                                        this,
                                        record.DomainName,
                                    )}
                                    disabled={!this.props.aclAuthMap["GET:/dns/domainRecordsList"]}
                                >
                                    解析设置
                                </Button>
                                <Divider type="vertical" />
                                <Button
                                    type="link"
                                    size="small"
                                    onClick={this.showDomainRecordHistory.bind(
                                        this,
                                        record.DomainName,
                                    )}
                                    disabled={!this.props.aclAuthMap["GET:/dns/domainHistoryList"]}
                                >
                                    解析历史
                                </Button>
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
                onShowSizeChange: this.onShowSizeChange,
                showQuickJumper: false,
                showTotal: (total) => `共 ${total} 条`,
                pageSize: 10,
                page: 1,
                total: 0,
                onChange: (current) => this.changePage(current),
            },
            cloudAccountList: [],
            queryCloudAccount: "",
            recordVisible: false,
            recordTableLoading: false,
            recordColumns: [
                {
                    title: "状态",
                    dataIndex: "Status",
                    key: "Status",
                    className: "small_font",
                    render: (value) => {
                        let color = "red";
                        let statusText = "禁用";
                        if (value === "ENABLE") {
                            color = "green";
                            statusText = "正常";
                        }
                        return (
                            <Text ellipsis={true} style={{ color: color }}>
                                {statusText}
                            </Text>
                        );
                    },
                },
                {
                    title: "域名",
                    dataIndex: "DomainName",
                    key: "DomainName",
                    className: "small_font",
                    render: (value) => {
                        return <Text ellipsis={true}>{value}</Text>;
                    },
                },
                {
                    title: "主机记录",
                    dataIndex: "RR",
                    key: "RR",
                    className: "small_font",
                    render: (value) => {
                        return <Text ellipsis={true}>{value}</Text>;
                    },
                },
                {
                    title: "记录类型",
                    dataIndex: "Type",
                    key: "Type",
                    className: "small_font",
                    render: (value) => {
                        return <Text ellipsis={true}>{value}</Text>;
                    },
                },
                {
                    title: "记录值",
                    dataIndex: "Value",
                    key: "Value",
                    className: "small_font",
                    render: (value) => {
                        return <Text ellipsis={true}>{value}</Text>;
                    },
                },
                {
                    title: "TTL",
                    dataIndex: "TTL",
                    key: "TTL",
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
                        let setContent = (
                            <Button
                                type="link"
                                size="small"
                                onClick={this.setRecordStatus.bind(
                                    this,
                                    record,
                                    "on",
                                )}
                                disabled={!this.props.aclAuthMap["POST:/dns/domainRecordStatus"]}
                            >
                                启用
                            </Button>
                        );
                        if (record.Status === "ENABLE") {
                            setContent = (
                                <Button
                                    type="link"
                                    size="small"
                                    onClick={this.setRecordStatus.bind(
                                        this,
                                        record,
                                        "off",
                                    )}
                                    disabled={!this.props.aclAuthMap["POST:/dns/domainRecordStatus"]}
                                >
                                    禁用
                                </Button>
                            );
                        }
                        return (
                            <div>
                                <Button
                                    type="link"
                                    size="small"
                                    onClick={this.modifyRecord.bind(
                                        this,
                                        "modify",
                                        record,
                                    )}
                                    disabled={!this.props.aclAuthMap["POST:/dns/domainRecordUpdate"]}
                                >
                                    修改
                                </Button>
                                <Divider type="vertical" />
                                {setContent}
                                <Divider type="vertical" />
                                <Button
                                    type="danger"
                                    size="small"
                                    onClick={this.deleteRecord.bind(
                                        this,
                                        record,
                                    )}
                                    disabled={!this.props.aclAuthMap["DELETE:/dns/domainRecord"]}
                                >
                                    删除
                                </Button>
                            </div>
                        );
                    },
                },
            ],
            recordPagination: {
                showSizeChanger: true,
                pageSizeOptions: ["10", "20", "30", "100"],
                onShowSizeChange: this.onShowRecordSizeChange,
                showQuickJumper: false,
                showTotal: (total) => `共 ${total} 条`,
                pageSize: 10,
                page: 1,
                total: 0,
                onChange: (current) => this.changeRecordPage(current),
            },
            domainName: "",
            recordList: [],
            addDomainModalVisible: false,
            recordModalVisible: false,
            recordAction: "",
            recordModalTitle: "",
            recordData: undefined,
            queryDomain: "",
            queryDomainRR: "",
        };
    }

    componentDidMount() {
        this.loadCloudAccountsData();
    }

    loadCloudAccountsData() {
        let that = this;
        getCloudAccouts(1, 100)
            .then((res) => {
                if (res.code === 0) {
                    that.setState({ cloudAccountList: res.data.accounts });
                } else {
                    message.error(res.msg);
                }
            })
    }

    onShowSizeChange(current, size) {
        this.setState(
            {
                pagination: {
                    ...this.state.pagination,
                    page: 1,
                    current: 1,
                    pageSize: size,
                },
            },
            () => {
                this.refreshTableData();
            },
        );
    }

    changePage = (e) => {
        this.setState(
            {
                pagination: {
                    ...this.state.pagination,
                    page: e,
                    current: e,
                },
            },
            () => {
                this.refreshTableData();
            },
        );
    };

    onShowRecordSizeChange(current, size) {
        this.setState(
            {
                recordPagination: {
                    ...this.state.recordPagination,
                    page: 1,
                    current: 1,
                    pageSize: size,
                },
            },
            () => {
                this.loadDomainRecordData();
            },
        );
    }

    changeRecordPage = (e) => {
        this.setState(
            {
                recordPagination: {
                    ...this.state.recordPagination,
                    page: e,
                    current: e,
                },
            },
            () => {
                this.loadDomainRecordData();
            },
        );
    };

    refreshTableData() {
        this.setState({ tableLoading: true });
        const queryParams = {
            domainName: this.state.queryDomain,
            pageNum: this.state.pagination.page,
            pageSize: this.state.pagination.pageSize,
            cloudAccountId: this.state.queryCloudAccount,
        };
        getDnsDomainListData(queryParams)
            .then((res) => {
                if(res.code===0) {
                    let resp = JSON.parse(res.data);
                    const pagination = this.state.pagination;
                    pagination.total = parseInt(resp.TotalCount);
                    pagination.page = parseInt(resp.PageNumber);
                    pagination.showTotal(parseInt(resp.TotalCount));
                    this.setState({ pagination });
                    let tableData = resp["Domains"]["Domain"];
                    this.setState({ tableData: tableData});
                } else {
                    message.error(res.msg);
                }
                this.setState({ tableLoading: false });
            })
    }

    handleCloudAccountChange(queryCloudAccount) {
        let that = this;
        this.setState(
            { queryCloudAccount: queryCloudAccount, queryDomain: "" },
            () => {
                that.refreshTableData();
            },
        );
    }

    domainQueryOnChange(e) {
        this.setState({ queryDomain: e.target.value });
    }

    domainRecordRROnChange(e) {
        this.setState({ queryDomainRR: e.target.value });
    }

    queryDomainList() {
        if (this.state.queryCloudAccount === "") {
            message.warn("请先选择阿里云子账号!");
            return;
        }
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

    onRecordDrawerClose = () => {
        this.setState({
            queryDomainRR: "",
            recordList: [],
            recordVisible: false,
        });
    };

    showDomainRecord(domainName) {
        let that = this;
        this.setState({ recordVisible: true, domainName: domainName }, () => {
            that.loadDomainRecordData();
        });
    }

    showDomainRecordHistory(domainName) {
        getDnsDomainHistoryListData({
            cloudAccountId: this.state.queryCloudAccount,
            domainName: domainName,
        })
            .then((res) => {
                if (res.code === 0) {
                    let resp = JSON.parse(res.data);
                    let historyLogs = resp["RecordLogs"]["RecordLog"];
                    let historyContent = historyLogs.map((item) => {
                        return (
                            <div
                                style={{
                                    paddingBottom: 5,
                                    marginBottom: 10,
                                    borderBottom: "1px solid grey",
                                }}
                            >
                                <span>
                                    <Tag>{item.Action}</Tag>{" "}
                                    <Tag>{item.ActionTime}</Tag>{" "}
                                    <Tag>{item.ClientIp}</Tag>
                                    <br />
                                    <Text strong>{item.Message}</Text>
                                </span>
                            </div>
                        );
                    });
                    if (historyLogs.length === 0) {
                        historyContent = "暂无解析历史";
                    }
                    Modal.info({
                        title: "近10次解析记录",
                        width: 700,
                        content: <div>{historyContent}</div>,
                        onOk() {},
                    });
                } else {
                    message.error(res.msg);
                }
            })
            .catch((err) => {
                message.error(err.toLocaleString());
            });
    }

    loadDomainRecordData = () => {
        const queryParams = {
            pageNum: this.state.recordPagination.page,
            pageSize: this.state.recordPagination.pageSize,
            cloudAccountId: this.state.queryCloudAccount,
            domainName: this.state.domainName,
            rr: this.state.queryDomainRR,
        };
        getDnsDomainRecordListData(queryParams)
            .then((res) => {
                if(res.code===0) {
                    let resp = JSON.parse(res.data);
                    const recordPagination = this.state.recordPagination;
                    recordPagination.total = parseInt(resp.TotalCount);
                    recordPagination.page = parseInt(resp.PageNumber);
                    recordPagination.showTotal(parseInt(resp.TotalCount));
                    let recordList = resp["DomainRecords"]["Record"];
                    this.setState({ recordPagination, recordList });
                } else {
                    message.error(res.msg);
                }
            })
    };

    showAddDomainModal = () => {
        if (this.state.queryCloudAccount === "") {
            message.warn("请先选择阿里云子账号!");
            return;
        }
        this.setState({ addDomainModalVisible: true });
    };

    handleAddDomainOk = (e) => {
        e.preventDefault();
        let that = this;
        this.domainFormRef.current.validateFields().then((values) => {
            postDnsDomain({
                cloudAccountId: this.state.queryCloudAccount,
                domainName: values.domainName,
            })
                .then((res) => {
                    if (res.code === 0) {
                        message.success("添加成功");
                    } else {
                        message.error(res.msg);
                    }
                    that.setState({ addDomainModalVisible: false });
                    that.refreshTableData();
                })
        });
    };

    handleAddDomainCancel = () => {
        this.setState({ addDomainModalVisible: false });
    };

    setRecordStatus(data, oper) {
        let status = "Disable";
        if (oper === "on") {
            status = "Enable";
        }
        postDnsDomainRecordStatus({
            cloudAccountId: this.state.queryCloudAccount,
            recordId: data.RecordId,
            status: status,
        })
            .then((res) => {
                if (res.code === 0) {
                    message.success("设置成功");
                } else {
                    message.error(res.msg);
                }
                this.loadDomainRecordData();
            })
    }

    deleteRecord(data) {
        deleteDnsDomainRecord({
            cloudAccountId: this.state.queryCloudAccount,
            recordId: data.RecordId,
        })
            .then((res) => {
                if (res.code === 0) {
                    message.success("删除成功");
                } else {
                    message.error(res.msg);
                }
                this.loadDomainRecordData();
            })
    }

    modifyRecord(action, data) {
        let recordModalTitle = "";
        let modifyRecordId = "";
        if (action === "add") {
            recordModalTitle = "新增解析记录";
        }
        if (action === "modify") {
            recordModalTitle = "修改解析记录";
            modifyRecordId = data.RecordId;
            let that = this;
            this.setState({ recordModalVisible: true }, () => {
                setTimeout(function() {
                    that.recordFormRef.current.setFieldsValue({
                        rr: data.RR,
                        rType: data.Type,
                        rValue: data.Value,
                    });
                }, 300);
            });
        }
        this.setState({
            recordAction: action,
            recordModalTitle: recordModalTitle,
            recordModalVisible: true,
            modifyRecordId: modifyRecordId,
        });
    }

    handleRecordOk = (e) => {
        e.preventDefault();
        let that = this;
        this.recordFormRef.current.validateFields().then((values) => {
            let reqParams = {
                cloudAccountId: this.state.queryCloudAccount,
                domainName: this.state.domainName,
                rr: values.rr,
                rType: values.rType,
                rValue: values.rValue,
                recordId: this.state.modifyRecordId,
            };
            if (this.state.recordAction === "add") {
                // 新增解析记录
                postDnsDomainRecord(reqParams)
                    .then((res) => {
                        if (res.code === 0) {
                            message.success("添加成功");
                        } else {
                            message.error(res.msg);
                        }
                        that.setState({ recordModalVisible: false });
                        that.loadDomainRecordData();
                    })
            }
            if (this.state.recordAction === "modify") {
                // 修改解析记录
                postDnsDomainRecordUpdate(reqParams)
                    .then((res) => {
                        if (res.code === 0) {
                            message.success("修改成功");
                        } else {
                            message.error(res.msg);
                        }
                        that.setState({ recordModalVisible: false });
                        that.loadDomainRecordData();
                    })
            }
        });
    };

    handleRecordCancel = () => {
        this.setState({ recordModalVisible: false });
    };

    render() {
        let accountOptions;
        accountOptions = this.state.cloudAccountList.map((item) => {
            if (item.accountType === "阿里云") {
                return (
                    <Option key={item.id} value={item.id}>
                        {item.accountName}
                    </Option>
                );
            } else {
                return null;
            }
        });
        return (
            <Content
                style={{
                    background: "#fff",
                    padding: 20,
                    margin: 0,
                    height: "100%",
                }}
            >
                <OpsBreadcrumbPath pathData={["DNS", "域名解析管理"]} />
                <Row style={{ padding: "0px 0px 10px 0px" }}>
                    <Col span={4} className="col-span">
                        选择所属阿里云账号:
                    </Col>
                    <Col span={5} className="col-span">
                        <Select
                            defaultValue={this.state.queryCloudAccount}
                            style={{ width: "100%" }}
                            onChange={this.handleCloudAccountChange}
                        >
                            {accountOptions}
                        </Select>
                    </Col>
                    <Col span={4} className="col-span">
                        <Input
                            placeholder="输入域名查询"
                            value={this.state.queryDomain}
                            onChange={this.domainQueryOnChange.bind(this)}
                        />
                    </Col>
                    <Col span={3} className="col-span">
                        <Button
                            style={{ width: "100%" }}
                            type="primary"
                            icon={<SearchOutlined />}
                            onClick={this.queryDomainList}
                        >
                            查 询
                        </Button>
                    </Col>
                    <Col span={3} className="col-span">
                        <Button
                            style={{ width: "100%" }}
                            icon={<PlusCircleOutlined />}
                            onClick={this.showAddDomainModal}
                            disabled={!this.props.aclAuthMap["POST:/dns/domain"]}
                        >
                            添加域名
                        </Button>
                    </Col>
                </Row>

                <Table
                    columns={this.state.columns}
                    dataSource={this.state.tableData}
                    scroll={{ x: "max-content" }}
                    pagination={this.state.pagination}
                    loading={this.state.tableLoading}
                    bordered
                    size="small"
                />

                <Drawer
                    title="解析设置"
                    placement="left"
                    closable={false}
                    onClose={this.onRecordDrawerClose}
                    visible={this.state.recordVisible}
                    width={800}
                >
                    <Row style={{ padding: "0px 0px 10px 0px" }}>
                        <Col span={5} className="col-span">
                            <Input
                                placeholder="输入主机记录关键字"
                                value={this.state.queryDomainRR}
                                onChange={this.domainRecordRROnChange.bind(
                                    this,
                                )}
                            />
                        </Col>
                        <Col span={3} className="col-span">
                            <Button
                                style={{ width: "100%" }}
                                type="primary"
                                icon={<SearchOutlined />}
                                onClick={this.loadDomainRecordData}
                            >
                                查 询
                            </Button>
                        </Col>
                        <Col span={4} className="col-span">
                            <Button
                                type="primary"
                                style={{ width: "100%" }}
                                icon={<PlusCircleOutlined />}
                                onClick={this.modifyRecord.bind(
                                    this,
                                    "add",
                                    undefined,
                                )}
                                disabled={!this.props.aclAuthMap["POST:/dns/domainRecord"]}
                            >
                                添加记录
                            </Button>
                        </Col>
                    </Row>
                    <Table
                        columns={this.state.recordColumns}
                        dataSource={this.state.recordList}
                        scroll={{ x: "max-content" }}
                        pagination={this.state.recordPagination}
                        loading={this.state.recordTableLoading}
                        bordered
                        size="small"
                    />
                    <Modal
                        title={this.state.recordModalTitle}
                        visible={this.state.recordModalVisible}
                        onOk={this.handleRecordOk}
                        onCancel={this.handleRecordCancel}
                        destroyOnClose={true}
                    >
                        <DomainRecordForm
                            domainName={this.state.domainName}
                            formRef={this.recordFormRef}
                        />
                    </Modal>
                </Drawer>

                <Modal
                    title="添加域名"
                    visible={this.state.addDomainModalVisible}
                    onOk={this.handleAddDomainOk}
                    onCancel={this.handleAddDomainCancel}
                    destroyOnClose={true}
                >
                    <DomainForm formRef={this.domainFormRef} />
                </Modal>
            </Content>
        );
    }
}

export default DomainManageContent;
