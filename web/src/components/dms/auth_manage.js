import React, { Component } from "react";
import {
    Button,
    Col,
    Layout,
    message,
    Modal,
    Row,
    Select,
    Table,
    Typography,
    Tree,
    Input,
    Form,
    Radio,
    InputNumber,
    DatePicker,
    Popconfirm,
} from "antd";
import OpsBreadcrumbPath from "../breadcrumb_path";
import {
    deleteDmsUserAuth,
    getDmsAuthData,
    getDmsDatabaseData,
    getAllDmsInstanceData,
    postDmsUserAuth,
} from "../../api/dms_api";
import { SearchOutlined, PlusCircleOutlined } from "@ant-design/icons";
import { getUsersList } from "../../api/user";
import { OpsIcon } from "../../assets/Icons";
const { Content } = Layout;

const { Text } = Typography;
const { Option } = Select;
const { TreeNode } = Tree;
const { TextArea } = Input;

class AuthManageContent extends Component {
    constructor(props) {
        super(props);
        this.onShowSizeChange = this.onShowSizeChange.bind(this);
        this.formRef = React.createRef();
        this.state = {
            columns: [
                {
                    title: "Id",
                    dataIndex: "ID",
                    key: "ID",
                    width: 60,
                    render: (value) => {
                        return <Text ellipsis={true}>{value}</Text>;
                    },
                },
                {
                    title: "权限类型",
                    dataIndex: "OperType",
                    key: "OperType",
                    align: "center",
                    width: 100,
                    render: (value, record) => {
                        let operType = "可见";
                        if (value === 2) {
                            operType = "查询";
                        }
                        if (value === 3) {
                            operType = "修改";
                        }
                        return <Text ellipsis={true}>{operType}</Text>;
                    },
                },
                {
                    title: "实例名称",
                    dataIndex: "InstanceName",
                    key: "InstanceName",
                    align: "center",
                    width: 200,
                    render: (value) => {
                        return (
                            <Text ellipsis={true} style={{ width: "100%" }}>
                                {value}
                            </Text>
                        );
                    },
                },
                {
                    title: "库名称",
                    dataIndex: "DatabaseName",
                    key: "DatabaseName",
                    align: "center",
                    width: 150,
                    render: (value) => {
                        return (
                            <Text ellipsis={true} style={{ width: "100%" }}>
                                {value}
                            </Text>
                        );
                    },
                },
                {
                    title: "授权用户",
                    dataIndex: "Username",
                    key: "Username",
                    align: "center",
                    width: 150,
                    render: (value) => {
                        return (
                            <Text ellipsis={true} style={{ width: "100%" }}>
                                {value}
                            </Text>
                        );
                    },
                },
                {
                    title: "剩余次数",
                    dataIndex: "OperCount",
                    key: "OperCount",
                    align: "center",
                    width: 100,
                    render: (value) => {
                        return (
                            <Text ellipsis={true} style={{ width: "100%" }}>
                                {value}
                            </Text>
                        );
                    },
                },
                {
                    title: "过期时间",
                    dataIndex: "ValidTime",
                    key: "ValidTime",
                    align: "center",
                    width: 100,
                    render: (value) => {
                        return (
                            <Text ellipsis={true} style={{ width: "100%" }}>
                                {value}
                            </Text>
                        );
                    },
                },
                {
                    title: "白名单表",
                    dataIndex: "AllowTables",
                    key: "AllowTables",
                    align: "center",
                    width: 100,
                    render: (value) => {
                        return (
                            <Text ellipsis={true} style={{ width: "100%" }}>
                                {value}
                            </Text>
                        );
                    },
                },
                {
                    title: "操作",
                    key: "operation",
                    fixed: "right",
                    align: "center",
                    width: 100,
                    render: (text, record) => {
                        return (
                            <div>
                                {/*<Button type="primary" size="small" >禁用</Button>*/}
                                {/*<Divider type="vertical" />*/}
                                <Popconfirm
                                    title="确定删除该项内容?"
                                    onConfirm={this.deleteUserAuth.bind(
                                        this,
                                        record,
                                    )}
                                    okText="确认"
                                    cancelText="取消"
                                    disabled={!this.props.aclAuthMap["DELETE:/dms/auth"]}
                                >
                                    <Button type="danger" size="small" disabled={!this.props.aclAuthMap["DELETE:/dms/auth"]}>
                                        删除
                                    </Button>
                                </Popconfirm>
                            </div>
                        );
                    },
                },
            ],
            tableLoading: false,
            addAuthModalVisible: false,
            tableData: [],
            usersData: [],
            instanceOptions: [],
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
            expandedKeys: [],
            searchValue: "",
            autoExpandParent: true,
            queryInstanceId: "",
            queryEmpId: "",
            queryOperType: "",
            treeData: [],
            selectedNodeId: "",
            selectedNodeType: "",
        };
    }

    componentDidMount() {
        this.getUserAuthData();
        this.loadAllUsersData();
        this.loadAllInstanceData();
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
                this.getUserAuthData();
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
                this.getUserAuthData();
            },
        );
    };

    getUserAuthData() {
        this.setState({ tableLoading: true });
        const queryParams = {
            page: this.state.pagination.page,
            size: this.state.pagination.pageSize,
            instanceId: this.state.queryInstanceId.trim(),
            operType: this.state.queryOperType.trim(),
            empId: this.state.queryEmpId,
        };
        getDmsAuthData(queryParams)
            .then((res) => {
                if (res.code === 0) {
                    const pagination = this.state.pagination;
                    pagination.total = parseInt(res.data.total);
                    pagination.page = parseInt(res.data.page);
                    pagination.showTotal(parseInt(res.data.total));
                    this.setState({ pagination });
                    this.setState({
                        tableData: res.data.authData,
                        tableLoading: false,
                    });
                } else {
                    message.error(res.msg);
                }
            })
            .catch((err) => {
                message.error(err.toLocaleString());
            });
    }

    loadAllUsersData() {
        getUsersList(1, 1000)
            .then((res) => {
                if (res.code === 0) {
                    let data = res.data.users;
                    let optionsList = [];
                    for (let i = 0; i < data.length; i++) {
                        optionsList.push(
                            <Option key={data[i].empId} value={data[i].empId}>
                                {data[i].username}
                            </Option>,
                        );
                    }
                    this.setState({
                        usersData: optionsList,
                    });
                } else {
                    message.error(res.msg);
                }
            })
            .catch((err) => {
                message.error(err.toLocaleString());
            });
    }

    loadAllInstanceData() {
        getAllDmsInstanceData()
            .then((res) => {
                if (res.code === 0) {
                    let optionsList = [];
                    let instanceTreeNode = [];
                    for (let i = 0; i < res.data.length; i++) {
                        optionsList.push(
                            <Option
                                key={res.data[i].InstanceId}
                                value={res.data[i].InstanceId}
                            >
                                {res.data[i].InstanceAlias}
                            </Option>,
                        );
                        instanceTreeNode.push({
                            title: res.data[i].InstanceAlias,
                            key: res.data[i].InstanceId,
                            type: "instance",
                            instance_type:
                                res.data[i].InstanceType === "2"
                                    ? "mysql"
                                    : "sqlserver",
                        });
                    }
                    this.setState({
                        instanceOptions: optionsList,
                        treeData: instanceTreeNode,
                    });
                } else {
                    message.error(res.msg);
                }
            })
    }

    handleInstanceChange = (queryInstanceId) => {
        this.setState({ queryInstanceId });
    };

    handleEmpIdChange = (queryEmpId) => {
        this.setState({ queryEmpId });
    };

    handleOperTypeChange = (queryOperType) => {
        this.setState({ queryOperType });
    };

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
                this.getUserAuthData();
            },
        );
    };

    handleAddAuth() {
        this.setState({ addAuthModalVisible: true });
    }

    handleUserAuthCancel = () => {
        this.setState({ addAuthModalVisible: false });
    };

    handleSubmitUserAuth = () => {
        this.formRef.current.validateFields().then((values) => {
            if (
                this.state.selectedNodeId === "" ||
                this.state.selectedNodeType === ""
            ) {
                message.warn("请选择具体的数据表！");
                return;
            }
            let reqParams = {
                ...values,
                validTime: values["validTime"].format("YYYY-MM-DD HH:mm:ss"),
                selectedNodeId: this.state.selectedNodeId,
                selectedNodeType: this.state.selectedNodeType,
            };
            postDmsUserAuth(reqParams)
                .then((res) => {
                    if (res.code === 0) {
                        message.success("添加成功!");
                        this.getUserAuthData();
                    } else {
                        message.error(res.msg);
                    }
                    this.setState({ addAuthModalVisible: false });
                })
        });
    };

    deleteUserAuth(data) {
        deleteDmsUserAuth({ id: "" + data.ID }).then((res) => {
            if (res.code === 0) {
                message.success("删除成功");
                this.getUserAuthData();
            } else {
                message.error(res.msg);
            }
        });
    }

    onLoadData = (treeNode) => {
        return new Promise((resolve) => {
            if (treeNode.props.children) {
                resolve();
                return;
            }
            if (treeNode.props.dataRef.type === "instance") {
                getDmsDatabaseData({ instanceId: treeNode.props.dataRef.key })
                    .then((res) => {
                        if (res.code === 0) {
                            let instanceChildren = [];
                            for (let i = 0; i < res.data.length; i++) {
                                instanceChildren.push({
                                    title: res.data[i].SchemaName,
                                    key: res.data[i].DatabaseId,
                                    type: "database",
                                    instance_type: res.data[i].InstanceType,
                                    isLeaf: true,
                                });
                            }
                            treeNode.props.dataRef.children = instanceChildren;
                            this.setState({
                                treeData: [...this.state.treeData],
                            });
                            resolve();
                        } else {
                            message.error(res.msg);
                        }
                    })
            }
        });
    };

    renderTreeNodes = (data) =>
        data.map((item) => {
            let iconType = "icondatabase";
            if (item.instance_type === "mysql") {
                iconType = "iconmysql";
            }
            if (item.instance_type === "sqlserver") {
                iconType = "iconsqlserver";
            }
            if (item.children) {
                return (
                    <TreeNode
                        icon={
                            <OpsIcon
                                style={{ fontSize: "20px", color: "#08c" }}
                                type={iconType}
                            />
                        }
                        title={item.title}
                        key={item.key}
                        dataRef={item}
                    >
                        {this.renderTreeNodes(item.children)}
                    </TreeNode>
                );
            }
            return (
                <TreeNode
                    icon={
                        <OpsIcon
                            style={{ fontSize: "20px", color: "#08c" }}
                            type={iconType}
                        />
                    }
                    key={item.key}
                    {...item}
                    dataRef={item}
                />
            );
        });

    onTreeNodeSelect = (selectedKeys, e) => {
        this.setState({
            selectedNodeId: e.selectedNodes[0].dataRef.key,
            selectedNodeType: e.selectedNodes[0].dataRef.type,
        });
    };

    resetSelectedNodeValue = () => {
        this.setState({
            selectedNodeId: "",
            selectedNodeType: "",
        });
    };

    render() {
        const formItemLayout = {
            labelCol: { span: 6 },
            wrapperCol: { span: 17 },
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
                <OpsBreadcrumbPath pathData={["DMS", "权限管理"]} />
                <Row style={{ padding: "0px 0px 10px 0px" }}>
                    <Col span={5} className="col-span">
                        <Select
                            placeholder="请选择所属实例"
                            style={{ width: "100%" }}
                            onChange={this.handleInstanceChange}
                        >
                            <Option value="">所有实例</Option>
                            {this.state.instanceOptions}
                        </Select>
                    </Col>
                    <Col span={4} className="col-span">
                        <Select
                            placeholder="请选择权限类型"
                            style={{ width: "100%" }}
                            onChange={this.handleOperTypeChange}
                        >
                            <Option value={2}>查询</Option>
                            <Option value={3}>修改</Option>
                        </Select>
                    </Col>
                    <Col span={4} className="col-span">
                        <Select
                            showSearch
                            placeholder="请选择用户"
                            style={{ width: "100%" }}
                            optionFilterProp="children"
                            filterOption={(input, option) =>
                                option.props.children
                                    .toLowerCase()
                                    .indexOf(input.toLowerCase()) >= 0
                            }
                            onChange={this.handleEmpIdChange}
                        >
                            <Option value="">所有用户</Option>
                            {this.state.usersData}
                        </Select>
                    </Col>
                    <Col span={3} className="col-span">
                        <Button
                            style={{ width: "100%" }}
                            type="primary"
                            icon={<SearchOutlined />}
                            onClick={this.handleQuery}
                        >
                            查 询
                        </Button>
                    </Col>
                    <Col span={3} className="col-span">
                        <Button
                            style={{ width: "100%" }}
                            icon={<PlusCircleOutlined />}
                            onClick={this.handleAddAuth.bind(this)}
                            disabled={!this.props.aclAuthMap["POST:/dms/auth"]}
                        >
                            添加权限
                        </Button>
                    </Col>
                </Row>

                <Modal
                    title="用户权限设置"
                    visible={this.state.addAuthModalVisible}
                    destroyOnClose={true}
                    width={700}
                    onOk={this.handleSubmitUserAuth}
                    onCancel={this.handleUserAuthCancel}
                    afterClose={this.resetSelectedNodeValue}
                >
                    <Row>
                        <Col span={10}>
                            <div
                                style={{
                                    minHeight: 450,
                                    maxHeight: "100%",
                                    overflowY: "scroll",
                                }}
                            >
                                <Tree
                                    showIcon={true}
                                    loadData={this.onLoadData}
                                    showLine={true}
                                    style={{
                                        maxHeight: "450px",
                                        overflowY: "scroll",
                                    }}
                                    onSelect={this.onTreeNodeSelect}
                                >
                                    {this.renderTreeNodes(this.state.treeData)}
                                </Tree>
                            </div>
                        </Col>
                        <Col span={2}></Col>
                        <Col span={12}>
                            <Form
                                ref={this.formRef}
                                {...formItemLayout}
                                initialValues={{ operType: 2, operCount: 5 }}
                            >
                                <Form.Item
                                    label="指定用户："
                                    name="empId"
                                    rules={[
                                        {
                                            required: true,
                                            message: "请选择用户!",
                                        },
                                    ]}
                                >
                                    <Select
                                        showSearch
                                        placeholder="请选择用户"
                                        style={{ width: "100%" }}
                                        optionFilterProp="children"
                                        filterOption={(input, option) =>
                                            option.props.children
                                                .toLowerCase()
                                                .indexOf(input.toLowerCase()) >=
                                            0
                                        }
                                    >
                                        {this.state.usersData}
                                    </Select>
                                </Form.Item>
                                {/*<Form.Item*/}
                                {/*    label="审批人："*/}
                                {/*    name="approveEmpId"*/}
                                {/*    rules={[*/}
                                {/*        {*/}
                                {/*            required: true,*/}
                                {/*            message: "请选择审批人!",*/}
                                {/*        },*/}
                                {/*    ]}*/}
                                {/*>*/}
                                {/*    <Select*/}
                                {/*        showSearch*/}
                                {/*        placeholder="请选择审批人"*/}
                                {/*        style={{ width: "100%" }}*/}
                                {/*        optionFilterProp="children"*/}
                                {/*        filterOption={(input, option) =>*/}
                                {/*            option.props.children*/}
                                {/*                .toLowerCase()*/}
                                {/*                .indexOf(input.toLowerCase()) >=*/}
                                {/*            0*/}
                                {/*        }*/}
                                {/*    >*/}
                                {/*        {this.state.usersData}*/}
                                {/*    </Select>*/}
                                {/*</Form.Item>*/}
                                <Form.Item
                                    label="权限类型："
                                    name="operType"
                                    rules={[
                                        {
                                            required: true,
                                            message: "请设置权限!",
                                        },
                                    ]}
                                >
                                    <Radio.Group>
                                        <Radio value={2}>查询</Radio>
                                        <Radio value={3}>修改</Radio>
                                    </Radio.Group>
                                </Form.Item>
                                <Form.Item
                                    label="申请原因："
                                    name="authReason"
                                    rules={[
                                        {
                                            required: true,
                                            message: "请输入申请原因!",
                                        },
                                    ]}
                                >
                                    <TextArea
                                        placeholder="请务必记录用户申请权限的原因"
                                        rows={3}
                                    />
                                </Form.Item>
                                <Form.Item
                                    label="操作次数："
                                    name="operCount"
                                    rules={[
                                        {
                                            required: true,
                                            message: "请输入操作次数!",
                                        },
                                    ]}
                                >
                                    <InputNumber min={1} />
                                    &nbsp;次
                                </Form.Item>
                                <Form.Item
                                    label="有效时间："
                                    name="validTime"
                                    rules={[
                                        {
                                            required: true,
                                            message: "请选择有效时间!",
                                        },
                                    ]}
                                >
                                    <DatePicker showTime />
                                </Form.Item>
                                <Form.Item
                                    label="白名单表："
                                    name="allowTables"
                                    rules={[
                                        {
                                            required: false,
                                        },
                                    ]}
                                >
                                    <TextArea
                                        placeholder="输入对用户开放的数据表，多表之间用英文分号分割"
                                        rows={2}
                                    />
                                </Form.Item>
                            </Form>
                        </Col>
                    </Row>
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

export default AuthManageContent;
