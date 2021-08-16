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
} from "antd";
import {
    deleteAuthLink,
    getPermissionsList,
    postAddAuthLink,
} from "../../api/role";
import OpsBreadcrumbPath from "../breadcrumb_path";

const { Content } = Layout;

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
                            name="path"
                            rules={[
                                {
                                    required: true,
                                    message: "请输入权限链接",
                                },
                            ]}
                        >
                            <Input placeholder="请输入权限链接" />
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
                    width: 360,
                },
                {},
                {
                    title: "操作",
                    key: "operation",
                    fixed: "right",
                    width: 60,
                    align: "center",
                    render: (text, record) => {
                        return (
                            <div>
                                <Popconfirm
                                    disabled={record['canDelete']===0}
                                    title="确定删除该项内容吗?"
                                    onConfirm={this.confirmDeleteAuthLink.bind(
                                        this,
                                        record,
                                    )}
                                    okText="确认"
                                    cancelText="取消"
                                >
                                    <Button type="danger" size="small" disabled={record['canDelete']===0}>
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
                {/*<div style={{ padding: "0px 0px 10px 0px" }}>*/}
                {/*    <Row>*/}
                {/*        <Col>*/}
                {/*            <Button*/}
                {/*                type="primary"*/}
                {/*                onClick={this.createAuthLink}*/}
                {/*            >*/}
                {/*                新建权限链接*/}
                {/*            </Button>*/}
                {/*        </Col>*/}
                {/*    </Row>*/}
                {/*</div>*/}
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
            </Content>
        );
    }
}

export default PermissionsManager;
