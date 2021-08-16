import React, { Component } from "react";
import {
    Button,
    Col,
    Layout,
    message,
    Modal,
    Row,
    Table,
    Typography,
    Input,
    Form,
    Radio,
    InputNumber,
    Popconfirm,
    Divider,
} from "antd";
import OpsBreadcrumbPath from "../breadcrumb_path";
import {
    postDmsInstanceData,
    getDmsInstanceData,
    deleteDmsInstanceData,
    getDmsDatabaseData,
    deleteDmsInstanceDbData,
    postDmsInstanceDbData,
} from "../../api/dms_api";
import { SearchOutlined, PlusCircleOutlined } from "@ant-design/icons";
import { EyeInvisibleOutlined, EyeTwoTone } from "@ant-design/icons";
const { Content } = Layout;

const { Text } = Typography;

class InstanceManageContent extends Component {
    constructor(props) {
        super(props);
        this.onShowSizeChange = this.onShowSizeChange.bind(this);
        this.formRef = React.createRef();
        this.state = {
            columns: [
                {
                    title: "实例类型",
                    dataIndex: "InstanceType",
                    key: "InstanceType",
                    align: "center",
                    width: 100,
                    render: (value, record) => {
                        let operType = "Mysql";
                        if (value === 2) {
                            operType = "Mysql";
                        }
                        if (value === 3) {
                            operType = "Sqlserver";
                        }
                        return <Text ellipsis={true}>{operType}</Text>;
                    },
                },
                {
                    title: "实例名称",
                    dataIndex: "InstanceAlias",
                    key: "InstanceAlias",
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
                    title: "实例地址",
                    dataIndex: "Host",
                    key: "Host",
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
                    title: "实例端口",
                    dataIndex: "Port",
                    key: "Port",
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
                    title: "数据库账号",
                    dataIndex: "OperUser",
                    key: "OperUser",
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
                    width: 200,
                    render: (text, record) => {
                        return (
                            <div>
                                <Button
                                    type="info"
                                    size="small"
                                    onClick={this.editInstance.bind(
                                        this,
                                        record,
                                    )}
                                    disabled={!this.props.aclAuthMap["PUT:/dms/instance"]}
                                >
                                    编辑
                                </Button>
                                <Divider type="vertical" />
                                <Popconfirm
                                    title="确定删除该项内容?"
                                    onConfirm={this.deleteInstance.bind(
                                        this,
                                        record,
                                    )}
                                    okText="确认"
                                    cancelText="取消"
                                    disabled={!this.props.aclAuthMap["DELETE:/dms/instance"]}
                                >
                                    <Button type="danger" size="small" disabled={!this.props.aclAuthMap["DELETE:/dms/instance"]}>
                                        删除
                                    </Button>
                                </Popconfirm>
                            </div>
                        );
                    },
                },
            ],
            expendColumns: [
                {
                    title: "库名称",
                    dataIndex: "SchemaName",
                    key: "SchemaName",
                    width: 100,
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
                                <Popconfirm
                                    title="确定删除该项内容?"
                                    onConfirm={this.deleteInstanceDb.bind(
                                        this,
                                        record,
                                    )}
                                    okText="确认"
                                    cancelText="取消"
                                    disabled={!this.props.aclAuthMap["DELETE:/dms/instance"]}
                                >
                                    <Button type="danger" size="small" disabled={!this.props.aclAuthMap["DELETE:/dms/instance"]}>
                                        删除
                                    </Button>
                                </Popconfirm>
                            </div>
                        );
                    },
                },
            ],
            tableLoading: false,
            addInstanceModalVisible: false,
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
            searchValue: "",
            initialValues: {},
            editInstanseId: "",
            addInstanceDbModalVisible: false,
            addInstanceId: "",
            inputDbName: "",
            expendTableData: {},
        };
    }

    componentDidMount() {
        this.loadInstanceData();
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
                this.loadInstanceData();
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
                this.loadInstanceData();
            },
        );
    };

    loadInstanceData() {
        this.setState({ tableLoading: true });
        const queryParams = {
            page: this.state.pagination.page,
            pageSize: this.state.pagination.pageSize,
            query: this.state.searchValue,
        };
        getDmsInstanceData(queryParams)
            .then((res) => {
                if (res.code === 0) {
                    const pagination = this.state.pagination;
                    pagination.total = parseInt(res.data.total);
                    pagination.page = parseInt(res.data.page);
                    pagination.showTotal(parseInt(res.data.total));
                    this.setState({ pagination });
                    let tableData = [];
                    for (let i = 0; i < res.data.instanceData.length; i++) {
                        tableData.push({
                            key: res.data.instanceData[i].ID,
                            ...res.data.instanceData[i],
                        });
                    }
                    this.setState({
                        tableData: tableData,
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

    handleQueryChange = (e) => {
        this.setState({ searchValue: e.target.value });
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
                this.loadInstanceData();
            },
        );
    };

    handleAddInstance() {
        this.setState({
            addInstanceModalVisible: true,
            editInstanseId: "",
            initialValues: {},
        });
    }

    handleAddInstanceCancel = () => {
        this.setState({
            addInstanceModalVisible: false,
            editInstanseId: "",
            initialValues: {},
        });
    };

    handleSubmitAddInstance = () => {
        this.formRef.current.validateFields().then((values) => {
            postDmsInstanceData({
                ...values,
                instanceId: this.state.editInstanseId,
            })
                .then((res) => {
                    if (res.code === 0) {
                        message.success("添加成功!");
                        this.loadInstanceData();
                    } else {
                        message.error(res.msg);
                    }
                    this.setState({
                        addInstanceModalVisible: false,
                    });
                })
                .catch((err) => {
                    message.error(err.toLocaleString());
                });
            this.setState({ editInstanseId: "", initialValues: {} });
        });
    };

    addInstanceDb(instanceId) {
        this.setState({
            addInstanceId: instanceId,
            addInstanceDbModalVisible: true,
        });
    }

    dbNameChange(e) {
        this.setState({ inputDbName: e.target.value });
    }

    handleSubmitAddInstanceDb() {
        postDmsInstanceDbData({
            instanceId: this.state.addInstanceId,
            dbName: this.state.inputDbName,
        }).then((res) => {
            if (res.code === 0) {
                message.success("添加成功");
                this.onExpand(true, { InstanceId: this.state.addInstanceId });
                this.setState({
                    addInstanceDbModalVisible: false,
                });
            } else {
                message.error(res.msg);
            }
        });
    }

    handleAddInstanceDbCancel() {
        this.setState({
            inputDbName: "",
            addInstanceId: "",
            addInstanceDbModalVisible: false,
        });
    }

    editInstance(data) {
        this.setState({
            addInstanceModalVisible: true,
            editInstanseId: data.InstanceId,
            initialValues: {
                instanceType: data.InstanceType,
                instanceAlias: data.InstanceAlias,
                port: data.Port,
                host: data.Host,
                operUser: data.OperUser,
                operPwd: data.OperPwd,
            },
        });
    }

    deleteInstance(data) {
        deleteDmsInstanceData({ instanceId: data.InstanceId }).then((res) => {
            if (res.code === 0) {
                message.success("删除成功");
                this.loadInstanceData();
            } else {
                message.error(res.msg);
            }
        });
    }

    deleteInstanceDb(data) {
        deleteDmsInstanceDbData({
            instanceId: data.InstanceId,
            databaseId: data.DatabaseId,
        }).then((res) => {
            if (res.code === 0) {
                message.success("删除成功");
                this.onExpand(true, data);
            } else {
                message.error(res.msg);
            }
        });
    }

    onExpand = (expanded, record) => {
        if (expanded) {
            getDmsDatabaseData({ instanceId: record.InstanceId }).then(
                (res) => {
                    if (res.code === 0) {
                        let expendTableData = this.state.expendTableData;
                        expendTableData[record.InstanceId] = res.data;
                        this.setState({
                            expendTableData: expendTableData,
                        });
                    } else {
                        message.error(res.msg);
                    }
                },
            );
        }
    };

    expandedRowRender = (record, index, indent, expanded) => {
        return (
            <Table
                style={{ width: 300 }}
                columns={this.state.expendColumns}
                dataSource={this.state.expendTableData[record.InstanceId]}
                pagination={false}
                size="small"
                footer={() => {
                    return (
                        <Button
                            size="small"
                            onClick={this.addInstanceDb.bind(
                                this,
                                record.InstanceId,
                            )}
                        >
                            添加
                        </Button>
                    );
                }}
            />
        );
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
                <OpsBreadcrumbPath pathData={["DMS", "实例管理"]} />
                <Row style={{ padding: "0px 0px 10px 0px" }}>
                    <Col span={5} className="col-span">
                        <Input
                            placeholder="请输入实例名称关键字"
                            style={{ width: "100%" }}
                            onChange={this.handleQueryChange}
                        />
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
                    <Col span={3} className="col-span">
                        <Button
                            style={{ width: "100%" }}
                            icon={<PlusCircleOutlined />}
                            onClick={this.handleAddInstance.bind(this)}
                            disabled={!this.props.aclAuthMap["POST:/dms/instance"]}
                        >
                            添加实例
                        </Button>
                    </Col>
                </Row>

                <Modal
                    title="添加实例信息"
                    visible={this.state.addInstanceModalVisible}
                    destroyOnClose={true}
                    width={700}
                    onOk={this.handleSubmitAddInstance}
                    onCancel={this.handleAddInstanceCancel}
                >
                    <Form
                        ref={this.formRef}
                        initialValues={this.state.initialValues}
                    >
                        <Form.Item
                            {...formItemLayout}
                            label="类型："
                            name="instanceType"
                            rules={[
                                {
                                    required: true,
                                    message: "请设置数据库实例类型!",
                                },
                            ]}
                        >
                            <Radio.Group>
                                <Radio value="2">Mysql</Radio>
                                <Radio value="3">Sqlserver</Radio>
                            </Radio.Group>
                        </Form.Item>
                        <Form.Item
                            {...formItemLayout}
                            label="实例名称："
                            name="instanceAlias"
                            rules={[
                                {
                                    required: true,
                                    message: "请设置数据库实例名称!",
                                },
                            ]}
                        >
                            <Input />
                        </Form.Item>
                        <Form.Item
                            {...formItemLayout}
                            label="实例地址："
                            name="host"
                            rules={[
                                {
                                    required: true,
                                    message: "请设置数据库实例地址!",
                                },
                            ]}
                        >
                            <Input />
                        </Form.Item>
                        <Form.Item
                            {...formItemLayout}
                            label="实例端口："
                            name="port"
                            rules={[
                                {
                                    required: true,
                                    message: "请设置数据库实例端口!",
                                },
                            ]}
                        >
                            <InputNumber min={1} />
                        </Form.Item>
                        <Form.Item
                            {...formItemLayout}
                            label="数据库账号："
                            name="operUser"
                            rules={[
                                {
                                    required: true,
                                    message: "请设置数据库实例账号!",
                                },
                            ]}
                        >
                            <Input />
                        </Form.Item>
                        <Form.Item
                            {...formItemLayout}
                            label="数据库密码："
                            name="operPwd"
                            rules={[
                                {
                                    required: true,
                                    message: "请设置数据库实例密码!",
                                },
                            ]}
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
                    </Form>
                </Modal>

                <Modal
                    title="添加库名"
                    visible={this.state.addInstanceDbModalVisible}
                    destroyOnClose={true}
                    width={300}
                    onOk={this.handleSubmitAddInstanceDb.bind(this)}
                    onCancel={this.handleAddInstanceDbCancel.bind(this)}
                >
                    <Input
                        placeholder="输入库名"
                        onChange={this.dbNameChange.bind(this)}
                    />
                </Modal>

                <Table
                    columns={this.state.columns}
                    dataSource={this.state.tableData}
                    scroll={{ x: "max-content" }}
                    pagination={this.state.pagination}
                    loading={this.state.tableLoading}
                    bordered
                    size="small"
                    expandable={{
                        onExpand: this.onExpand,
                        expandedRowRender: this.expandedRowRender,
                    }}
                />
            </Content>
        );
    }
}

export default InstanceManageContent;
