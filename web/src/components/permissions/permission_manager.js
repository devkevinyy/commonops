import React, { Component, Fragment } from "react";
import {
    Layout,
    Table,
    Form,
    Popconfirm,
    Row,
    Col,
    Button,
    Modal,
    Input,
    message,
    Tag,
    Select,
    Divider,
} from "antd";
import {
    getPermissionsList,
    getAuthLink,
    putAuthLink,
    postAddAuthLink,
    deleteAuthLink,
} from "../../api/permission"
import OpsBreadcrumbPath from "../breadcrumb_path";
import ExtraInfoModal from "./common/extra_info_modal";

const { Content } = Layout;
const { Option } = Select;

let columnStyle = {
    overFlow: "hidden",
    textOverflow: "ellipsis",
    display: "-webkit-box",
};

class AuthLinkModal extends Component {
    render() {
        const formItemLayout = {
            labelCol: { span: 6 },
            wrapperCol: { span: 14 },
        };
        return (
            <Fragment>
                <Modal
                    title="新建权限链接"
                    destroyOnClose="true"
                    visible={this.props.authLinkModalVisible}
                    onOk={this.props.handleAddAuthLink}
                    onCancel={this.props.handleCancelAddAuthLink}
                    okText="确认"
                    cancelText="取消"
                >
                    <Form ref={this.props.formRef}>
                        <Form.Item
                            {...formItemLayout}
                            label="权限名称"
                            name="name"
                            rules={[
                                {
                                    required: true,
                                    message: "请输入路径权限名称",
                                },
                            ]}
                        >
                            <Input placeholder="请输入路径权限名称" />
                        </Form.Item>
                        <Form.Item
                            {...formItemLayout}
                            label="权限描述"
                            name="description"
                            rules={[
                                {
                                    required: true,
                                    message: "请输入权限描述",
                                },
                            ]}
                        >
                            <Input placeholder="请输入权限描述" />
                        </Form.Item>
                        <Form.Item
                            {...formItemLayout}
                            label="权限链接"
                            name="urlPath"
                            rules={[
                                {
                                    required: true,
                                    message: "请输入权限链接",
                                },
                            ]}
                        >
                            <Input placeholder="请输入权限链接" />
                        </Form.Item>
                        <Form.Item
                            {...formItemLayout}
                            label="权限类型"
                            name="authType"
                            rules={[
                                {
                                    required: true,
                                    message: "请输入权限类型",
                                },
                            ]}
                        >
                            <Select
                                placeholder="选择权限类型"
                            >
                              <Option value="菜单">菜单</Option>
                              <Option value="操作">操作</Option>
                            </Select>
                        </Form.Item>
                    </Form>
                </Modal>
            </Fragment>
        );
    }
}

class PermissionsManager extends Component {
    constructor(props) {
        super(props);
        this.refreshTableData = this.refreshTableData.bind(this);
        this.createAuthLink = this.createAuthLink.bind(this);
        this.handleAddAuthLink = this.handleAddAuthLink.bind(this);
        this.handleCancelAddAuthLink = this.handleCancelAddAuthLink.bind(this);
        this.onShowSizeChange = this.onShowSizeChange.bind(this);
        this.formRef = React.createRef();
        this.state = {
            columns: [
                {
                    title: "ID",
                    dataIndex: "Id",
                    key: "Id",
                    width: 60,
                },
                {
                    title: "权限类型",
                    dataIndex: "authType",
                    key: "authType",
                    width: 100,
                    render: (text => {
                        if(text==="菜单") {
                            return <Tag color="volcano">{text}</Tag>;
                        }
                        return <Tag color="blue">{text}</Tag>;
                    })
                },
                {
                    title: "权限名称",
                    dataIndex: "name",
                    key: "name",
                    width: 200,
                },
                {
                    title: "权限描述",
                    dataIndex: "description",
                    key: "description",
                    className: { columnStyle },
                    width: 260,
                },
                {
                    title: "权限路径",
                    dataIndex: "urlPath",
                    key: "urlPath",
                    className: { columnStyle },
                    width: 300,
                },
                {},
                {
                    title: "操作",
                    key: "operation",
                    fixed: "right",
                    width: 160,
                    align: "center",
                    render: (text, record) => {
                        return (
                            <div>
                                <Button
                                    type="info"
                                    size="small"
                                    disabled={
                                        !this.props.aclAuthMap["PUT:/permissions/authLink"]
                                    }
                                    onClick={this.authLinkEdit.bind(this, record)}
                                >
                                    编辑
                                </Button>
                                <Divider type="vertical" />
                                <Popconfirm
                                    title="确定删除该项内容吗?"
                                    onConfirm={this.confirmDeleteAuthLink.bind(
                                        this,
                                        record,
                                    )}
                                    okText="确认"
                                    cancelText="取消"
                                    disabled={
                                        !this.props.aclAuthMap["DELETE:/permissions/authLink"]
                                    }
                                >
                                    <Button
                                        type="danger"
                                        size="small"
                                        disabled={record['canDelete']===0}
                                    >
                                        删除
                                    </Button>
                                </Popconfirm>
                            </div>
                        );
                    },
                },
            ],
            tableData: [],
            tableLoading: false,
            authLinkModalVisible: false,
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
            extraInfoData: {},
            updateMode: "single",
            idsList: [],
            selectedRowKeys: [],
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

    confirmDeleteAuthLink = (e) => {
        deleteAuthLink({
            id: e.Id,
        })
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
        getPermissionsList(
            this.state.pagination.page,
            this.state.pagination.pageSize,
        )
            .then((res) => {
                if(res.code===0) {
                    const pagination = this.state.pagination;
                    pagination.total = parseInt(res.data.total);
                    pagination.page = parseInt(res.data.page);
                    pagination.showTotal(parseInt(res.data.total));
                    this.setState({
                        pagination,
                    });
                    this.setState({ tableData: res["data"]["permissions"]});
                } else {
                    message.error(res.msg);
                }
                this.setState({tableLoading: false });
            })
    };

    componentDidMount() {
        this.refreshTableData();
    }

    createAuthLink() {
        this.setState({
            authLinkModalVisible: true,
        });
    }

    handleCancelAddAuthLink() {
        this.setState({
            authLinkModalVisible: false,
        });
    }

    handleAddAuthLink() {
        this.formRef.current.validateFields().then((values) => {
            postAddAuthLink(values)
                .then((res) => {
                    if (res.code === 0) {
                        message.success("创建成功");
                        this.setState({
                            authLinkModalVisible: false,
                        });
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

    authLinkEdit = (data) => {
        getAuthLink(data.Id)
            .then((res) => {
                if (res["code"] !== 0) {
                    message.error(res["msg"]);
                } else {
                    let extraInfoData = {
                      Id: res.data["Id"],
                      name: res.data["name"],
                      description: res.data["description"],
                      urlPath: res.data["urlPath"],
                      canDelete: res.data["canDelete"],
                      authType: res.data["authType"],
                    };
                    this.setState({
                      extraInfoModalVisible: true,
                      authLinkId: data.Id,
                      updateMode: "single",
                      extraInfoData: extraInfoData,
                    });
                }
            })
            .catch((err) => {
                message.error(err.toLocaleString());
            });
    }

    handleExtraInfoOk = (data) => {
        let targetId = "";
        // makge sure it is an int
        data["canDelete"] = parseInt(data["canDelete"]);
        // convert id to string
        data["Id"] = data["Id"].toString();
        if (this.state.updateMode === "single") {
            targetId = String(this.state.authLinkId);
        } else {
            targetId = this.state.idsList.join(",");
        }
        putAuthLink({
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
                message.error(err.toLocaleString())
            });
    }

    handleExtraInfoCancel = (data) => {
        this.setState({ extraInfoModalVisible: false });
    }

    render() {
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
                    pathData={["权限管理", "链接权限", "权限链接列表"]}
                />
                <div style={{ padding: "0px 0px 10px 0px" }}>
                    <Row>
                        <Col>
                            <Button
                                type="primary"
                                onClick={this.createAuthLink}
                                disabled={
                                    !this.props.aclAuthMap["POST:/permissions/authLink"]
                                }
                            >
                                新建权限链接
                            </Button>
                        </Col>
                    </Row>
                </div>
                <AuthLinkModal
                    formRef={this.formRef}
                    authLinkModalVisible={this.state.authLinkModalVisible}
                    handleAddAuthLink={this.handleAddAuthLink}
                    handleCancelAddAuthLink={this.handleCancelAddAuthLink}
                />
                <Table
                    columns={this.state.columns}
                    dataSource={this.state.tableData}
                    scroll={{ x: "max-content" }}
                    pagination={this.state.pagination}
                    loading={this.state.tableLoading}
                    size="small"
                />

                {/*完善信息组件*/}
                <ExtraInfoModal
                    editData={this.state.extraInfoData}
                    resType="authLink"
                    updateMode={this.state.updateMode}
                    resFrom={this.state.resFrom}
                    visible={this.state.extraInfoModalVisible}
                    handleOk={this.handleExtraInfoOk}
                    handleCancel={this.handleExtraInfoCancel}
                />
            </Content>
        );
    }
}

export default PermissionsManager;
