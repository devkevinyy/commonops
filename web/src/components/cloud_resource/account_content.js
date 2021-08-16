import React, { Component, Fragment } from "react";
import {
    Button,
    Layout,
    Table,
    Modal,
    Form,
    Input,
    Select,
    Popconfirm,
    Divider,
    Typography,
} from "antd";
import { message } from "antd";
import {
    deleteCloudAccouts,
    getCloudAccouts,
    postCloudAccouts,
    putCloudAccouts,
} from "../../api/cloud";
import OpsBreadcrumbPath from "../breadcrumb_path";

const { Text, Paragraph } = Typography;
const { Content } = Layout;
const { Option } = Select;

class CloudContent extends Component {
    constructor(props) {
        super(props);
        this.onShowSizeChange = this.onShowSizeChange.bind(this);
        this.formRef = React.createRef();
        this.state = {
            columns: [
                {
                    title: "云类型",
                    dataIndex: "accountType",
                    key: "accountType",
                    className: "small_font",
                    render: (value) => {
                        return <Text ellipsis={true}>{value}</Text>;
                    },
                },
                {
                    title: "账号名称",
                    dataIndex: "accountName",
                    key: "accountName",
                    className: "small_font",
                    render: (value) => {
                        return <Text ellipsis={true}>{value}</Text>;
                    },
                },
                {
                    title: "账号密码",
                    dataIndex: "accountPwd",
                    key: "accountPwd",
                    className: "small_font",
                    render: (value) => {
                        return <Text ellipsis={true}>{value}</Text>;
                    },
                },
                {
                    title: "访问key",
                    dataIndex: "accountKey",
                    key: "accountKey",
                    className: "small_font",
                    render: (value) => {
                        return (
                            <Paragraph style={{ marginBottom: 0 }} copyable>
                                {value}
                            </Paragraph>
                        );
                    },
                },
                {
                    title: "访问secret",
                    dataIndex: "accountSecret",
                    key: "accountSecret",
                    className: "small_font",
                    render: (value) => {
                        return (
                            <Paragraph style={{ marginBottom: 0 }} copyable>
                                {value}
                            </Paragraph>
                        );
                    },
                },
                {
                    title: "所属区域",
                    dataIndex: "accountRegion",
                    key: "accountRegion",
                    className: "small_font",
                    render: (value) => {
                        return <Text ellipsis={true}>{value}</Text>;
                    },
                },
                {
                    title: "操作",
                    key: "operation",
                    fixed: "right",
                    width: 160,
                    align: "center",
                    className: "small_font",
                    render: (text, record) => {
                        let superOper;
                        if (this.props.isSuperAdmin) {
                            superOper = (
                                <Fragment>
                                    <Button
                                        type="primary"
                                        size="small"
                                        onClick={this.handleEditAccountInfo.bind(
                                            this,
                                            record,
                                        )}
                                        disabled={!this.props.aclAuthMap["PUT:/cloud/accounts"]}
                                    >
                                        编辑
                                    </Button>
                                    <Divider type="vertical" />
                                    <Popconfirm
                                        title="确定删除该项内容?"
                                        onConfirm={this.confirm.bind(
                                            this,
                                            record,
                                        )}
                                        okText="确认"
                                        cancelText="取消"
                                        disabled={!this.props.aclAuthMap["DELETE:/cloud/accounts"]}
                                    >
                                        <Button type="danger" size="small" disabled={!this.props.aclAuthMap["DELETE:/cloud/accounts"]}>
                                            删除
                                        </Button>
                                    </Popconfirm>
                                </Fragment>
                            );
                        } else {
                            superOper = <span>无操作权限</span>;
                        }
                        return <div>{superOper}</div>;
                    },
                },
            ],
            tableData: [],
            tableLoading: false,
            add_new_account_visible: false,
            current_data_id: null,
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

    confirm = (e) => {
        this.handleDeleteAccountInfo(e);
    };

    showAddAccountModal = () => {
        this.setState({
            current_data_id: null,
            add_new_account_visible: true,
        });
    };

    handleAddAccountSubmit = (e) => {
        this.formRef.current.validateFields().then((values) => {
            if (this.state.current_data_id != null) {
                values["id"] = this.state.current_data_id;
                putCloudAccouts(values)
                    .then((res) => {
                        if (res.code === 0) {
                            this.setState({
                                add_new_account_visible: false,
                            });
                            message.success("更新成功！");
                        }
                        this.refreshTableData();
                    })
                    .catch((err) => {
                        console.log(err);
                    });
            } else {
                postCloudAccouts(values)
                    .then((res) => {
                        if (res.code === 0) {
                            this.setState({
                                add_new_account_visible: false,
                            });
                            message.success("新增成功！");
                        }
                        this.refreshTableData();
                    })
                    .catch((err) => {
                        console.log(err);
                    });
            }
        });
    };

    handleAddAccountCancel = (e) => {
        this.setState({
            add_new_account_visible: false,
        });
    };

    handleEditAccountInfo = (e) => {
        setTimeout(() => {
            this.formRef.current.setFieldsValue({
                id: e.id,
                accountType: e.accountType,
                accountName: e.accountName,
                accountPwd: e.accountPwd,
                accountKey: e.accountKey,
                accountSecret: e.accountSecret,
                accountRegion: e.accountRegion,
                bankAccount: e.bankAccountId,
            });
        }, 300);

        this.setState({
            current_data_id: e.id,
            add_new_account_visible: true,
        });
    };

    handleDeleteAccountInfo = (e) => {
        deleteCloudAccouts({ id: e.id })
            .then((res) => {
                if (res.code === 0) {
                    message.success("删除成功！");
                }
                this.refreshTableData();
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

    refreshTableData = () => {
        this.setState({ tableLoading: true });
        getCloudAccouts(
            this.state.pagination.page,
            this.state.pagination.pageSize,
        )
            .then((res) => {
                const pagination = this.state.pagination;
                pagination.total = parseInt(res.data.total);
                pagination.page = parseInt(res.data.page);
                pagination.showTotal(parseInt(res.data.total));
                this.setState({
                    pagination,
                });
                let data = res["data"]["accounts"];
                let tableData = [];
                for (let i = 0; i < data.length; i++) {
                    tableData.push({
                        id: data[i]["id"],
                        key: data[i]["id"],
                        accountType: data[i]["accountType"],
                        accountName: data[i]["accountName"],
                        accountPwd: data[i]["accountPwd"],
                        accountKey: data[i]["accountKey"],
                        accountSecret: data[i]["accountSecret"],
                        accountRegion: data[i]["accountRegion"],
                        bankAccountId: data[i]["bankAccountId"],
                    });
                }
                this.setState({ tableData: tableData, tableLoading: false });
            })
            .catch((err) => {
                console.log(err);
            });
    };

    componentDidMount() {
        this.refreshTableData();
    }

    render() {
        const formItemLayout = {
            labelCol: { span: 7 },
            wrapperCol: { span: 14 },
        };
        return (
            <Content
                style={{
                    background: "#fff",
                    padding: 20,
                    margin: 0,
                    height: "100%",
                }}
            >
                <OpsBreadcrumbPath
                    pathData={["云资源", "云账号", "账号列表"]}
                />
                {this.props.isSuperAdmin ? (
                    <div style={{ marginBottom: 10 }}>
                        <Button
                            type="primary"
                            onClick={this.showAddAccountModal}
                            disabled={!this.props.aclAuthMap["POST:/cloud/accounts"]}
                        >
                            新增账号
                        </Button>
                    </div>
                ) : (
                    ""
                )}
                <Modal
                    title="云账号信息详情"
                    visible={this.state.add_new_account_visible}
                    onOk={this.handleAddAccountSubmit}
                    onCancel={this.handleAddAccountCancel}
                    okText="确认"
                    cancelText="取消"
                >
                    <Form
                        ref={this.formRef}
                        initialValues={{ accountType: "私有云" }}
                    >
                        <Form.Item
                            label="云账号类型"
                            {...formItemLayout}
                            name="accountType"
                            rules={[
                                {
                                    required: true,
                                    message: "账号类型不能为空！",
                                },
                            ]}
                        >
                            <Select size="default">
                                <Option value="私有云">私有云</Option>
                                <Option value="阿里云">阿里云</Option>
                            </Select>
                        </Form.Item>
                        <Form.Item
                            label="云账号名称"
                            {...formItemLayout}
                            name="accountName"
                            rules={[
                                {
                                    required: true,
                                    message: "账号名称不能为空！",
                                },
                            ]}
                        >
                            <Input />
                        </Form.Item>
                        <Form.Item
                            label="云账号密码"
                            {...formItemLayout}
                            name="accountPwd"
                            rules={[
                                // {
                                //     required: true,
                                //     message: "账号密码不能为空！",
                                // },
                            ]}
                        >
                            <Input />
                        </Form.Item>
                        <Form.Item
                            label="云账号key"
                            {...formItemLayout}
                            name="accountKey"
                            rules={[
                                // {
                                //     required: true,
                                //     message: "账号key不能为空！",
                                // },
                            ]}
                        >
                            <Input />
                        </Form.Item>
                        <Form.Item
                            label="云账号secret"
                            {...formItemLayout}
                            name="accountSecret"
                            rules={[
                                // {
                                //     required: true,
                                //     message: "账号secret不能为空！",
                                // },
                            ]}
                        >
                            <Input />
                        </Form.Item>
                        <Form.Item
                            label="所属区域"
                            {...formItemLayout}
                            name="accountRegion"
                            rules={[
                                {
                                    required: true,
                                    message: "所属区域不能为空！",
                                },
                            ]}
                        >
                            <Input />
                        </Form.Item>
                    </Form>
                </Modal>
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

export default CloudContent;
