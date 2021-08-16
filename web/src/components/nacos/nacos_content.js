import React, { Component } from "react";
import {
    Layout,
    Row,
    message,
    Button,
    Col,
    Typography,
    Select,
    Form,
    Drawer,
    Input,
    Radio,
    Table,
    Modal,
    Divider,
    Descriptions,
    Popconfirm,
    Tree,
    Tag,
} from "antd";
import OpsBreadcrumbPath from "../breadcrumb_path";
import { withRouter } from "react-router-dom";
import {
    getNacosList,
    getNacosNamespaceList,
    getNacosConfigsList,
    putNacosConfig,
    postNacosConfig,
    postNacosServer,
    deleteNacosConfig,
    postNacosConfigCopy,
    getNacosConfigDetail,
    getNacosAllConfigs,
    postNacosConfigSync,
    getConfigTemplatesAll,
} from "../../api/nacos";

const { Paragraph, Text } = Typography;
const { Content } = Layout;
const { Option } = Select;
const { TextArea } = Input;

class NacosServerForm extends Component {
    render() {
        const formItemLayout = {
            labelCol: { span: 6 },
            wrapperCol: { span: 16 },
        };
        return (
            <Form ref={this.props.formRef} {...formItemLayout}>
                <Form.Item
                    label="集群名称"
                    name="alias"
                    rules={[{ required: true, message: "集群名称不能为空！" }]}
                >
                    <Input />
                </Form.Item>
                <Form.Item
                    label="访问地址"
                    name="endpoint"
                    rules={[{ required: true, message: "endpoint不能为空！" }]}
                >
                    <Input addonBefore="http://" />
                </Form.Item>
                <Form.Item
                    label="用户名"
                    name="username"
                    rules={[{ required: true, message: "username不能为空！" }]}
                >
                    <Input />
                </Form.Item>
                <Form.Item
                    label="密码"
                    name="password"
                    rules={[{ required: true, message: "password不能为空！" }]}
                >
                    <Input />
                </Form.Item>
            </Form>
        );
    }
}

class NacosConfigForm extends Component {
    constructor(props) {
        super(props);
        this.state = {
            isUserTemplate: false,
            fillFields: [],
        };
    }

    componentDidMount() {
        if (this.props.updateTemplateId !== 0) {
            this.setState({
                isUserTemplate: true,
                configContent: this.props.templatesMap[
                    this.props.updateTemplateId
                ]["ConfigContent"],
                fillFields: JSON.parse(
                    this.props.templatesMap[this.props.updateTemplateId][
                        "FillField"
                    ],
                ),
            });
        } else {
            this.setState({
                isUserTemplate: false,
                configContent: this.props.nacosConfigContent,
            });
        }
    }

    templateRadioOnChange = (e) => {
        if (e.target.value === "useTemplate") {
            this.setState({ isUserTemplate: true });
        } else {
            this.setState({ isUserTemplate: false });
        }
    };

    templateOnChange = (id) => {
        this.setState({
            configContent: this.props.templatesMap[id]["ConfigContent"],
            fillFields: JSON.parse(this.props.templatesMap[id]["FillField"]),
        });
    };

    render() {
        const formItemLayout = {
            labelCol: { span: 6 },
            wrapperCol: { span: 14 },
        };
        return (
            <Form ref={this.props.formRef} {...formItemLayout}>
                <Form.Item
                    label="是否使用模板"
                    name="isUseTemplate"
                    // rules={[
                    //     { required: true, message: "是否使用模板不能为空！" },
                    // ]}
                >
                    <Radio.Group
                        defaultValue="notUserTemplate"
                        onChange={this.templateRadioOnChange}
                    >
                        <Radio value="notUserTemplate">不使用模板</Radio>
                        <Radio value="useTemplate">使用模板</Radio>
                    </Radio.Group>
                </Form.Item>
                <Form.Item
                    label="模板列表"
                    name="templateId"
                    rules={[
                        {
                            required: this.state.isUserTemplate,
                            message: "请选择模板！",
                        },
                    ]}
                >
                    <Select
                        showSearch={true}
                        disabled={!this.state.isUserTemplate}
                        onChange={this.templateOnChange}
                    >
                        <Option key={0} value={0}>
                            不使用模板
                        </Option>
                        {this.props.templatesList.map((item, index) => {
                            return (
                                <Option key={index} value={item.Id}>
                                    {item.Name}
                                </Option>
                            );
                        })}
                    </Select>
                </Form.Item>
                <Form.Item
                    label="ConfigId"
                    name="configId"
                    hidden={true}
                ></Form.Item>
                <Form.Item label="Id" name="id" hidden={true}></Form.Item>
                <Form.Item
                    label="DataId"
                    name="dataId"
                    rules={[{ required: true, message: "DataId不能为空！" }]}
                >
                    <Input />
                </Form.Item>
                <Form.Item
                    label="Group"
                    name="group"
                    rules={[{ required: true, message: "Group不能为空！" }]}
                >
                    <Input />
                </Form.Item>
                <Form.Item
                    label="格式"
                    name="configType"
                    // rules={[{ required: true, message: "格式不能为空！" }]}
                >
                    <Radio.Group defaultValue="yaml">
                        <Radio value="yaml">YAML</Radio>
                        {/* <Radio value="properties">Properties</Radio>
                        <Radio value="text">TEXT</Radio>
                        <Radio value="json">JSON</Radio>
                        <Radio value="xml">XML</Radio>
                        <Radio value="html">HTML</Radio> */}
                    </Radio.Group>
                </Form.Item>
                {this.state.isUserTemplate ? (
                    <Form.Item label="配置内容" name="content">
                        <pre className="preJenkinsLog">
                            {this.state.configContent}
                        </pre>
                    </Form.Item>
                ) : (
                    <Form.Item label="配置内容" name="content">
                        <TextArea rows={4} />
                    </Form.Item>
                )}
                {this.state.isUserTemplate &&
                    this.state.fillFields.map((item, index) => {
                        return (
                            <Form.Item
                                label={item}
                                name={item}
                                rules={[
                                    {
                                        required: true,
                                        message: item + "不能为空！",
                                    },
                                ]}
                            >
                                <Input />
                            </Form.Item>
                        );
                    })}
            </Form>
        );
    }
}

class NacosConfigCopyForm extends Component {
    render() {
        const formItemLayout = {
            labelCol: { span: 6 },
            wrapperCol: { span: 15 },
        };
        return (
            <Form ref={this.props.formRef} {...formItemLayout}>
                <Form.Item
                    label="srcNamespace"
                    name="srcNamespace"
                    hidden={true}
                />
                <Form.Item label="srcDataId" name="srcDataId" hidden={true} />
                <Form.Item label="srcGroup" name="srcGroup" hidden={true} />
                <Form.Item
                    label="Namespace"
                    name="dstNamespace"
                    rules={[{ required: true, message: "Namespace不能为空！" }]}
                >
                    <Select>
                        {this.props.nsList.map((item, index) => {
                            return (
                                <Option
                                    key={index}
                                    value={item.NamespaceShowName}
                                >
                                    {item.NamespaceShowName}
                                </Option>
                            );
                        })}
                    </Select>
                </Form.Item>
                <Form.Item
                    label="DataId"
                    name="dstDataId"
                    rules={[{ required: true, message: "DataId不能为空！" }]}
                >
                    <Input />
                </Form.Item>
                <Form.Item
                    label="Group"
                    name="dstGroup"
                    rules={[{ required: true, message: "Group不能为空！" }]}
                >
                    <Input />
                </Form.Item>
            </Form>
        );
    }
}

class NacosContent extends Component {
    constructor(props) {
        super(props);
        this.configFormRef = React.createRef();
        this.configCopyFormRef = React.createRef();
        this.serverFormRef = React.createRef();
        this.state = {
            addNacosModalVisible: false,
            addNacosConfigModalVisible: false,
            copyNacosConfigModalVisible: false,
            syncNacosConfigModalVisible: false,
            currentNacosServerId: "",
            currentNs: "",
            nacosList: [],
            nsList: [],
            columns: [
                {
                    title: "Id",
                    dataIndex: "Id",
                    key: "Id",
                    render: (value) => {
                        return <Text ellipsis={true}>{value}</Text>;
                    },
                },
                {
                    title: "Namespace",
                    dataIndex: "Namespace",
                    key: "Namespace",
                    render: (value) => {
                        return <Text ellipsis={true}>{value}</Text>;
                    },
                },
                {
                    title: "是否使用模板",
                    dataIndex: "TemplateId",
                    key: "TemplateId",
                    render: (value) => {
                        if (value === 0) {
                            return (
                                <Text ellipsis={true}>
                                    <Tag>未使用模板</Tag>
                                </Text>
                            );
                        }
                        return (
                            <Text ellipsis={true}>
                                <Tag color="#87d068">使用模板</Tag>
                            </Text>
                        );
                    },
                },
                {
                    title: "DataId",
                    dataIndex: "DataId",
                    key: "dataId",
                    render: (value) => {
                        return <Text ellipsis={true}>{value}</Text>;
                    },
                },
                {
                    title: "Group",
                    dataIndex: "ConfigGroup",
                    key: "ConfigGroup",
                    render: (value) => {
                        return <Text ellipsis={true}>{value}</Text>;
                    },
                },
                {
                    title: "模板名称",
                    dataIndex: "Name",
                    key: "Name",
                    render: (value) => {
                        return <Text ellipsis={true}>{value}</Text>;
                    },
                },
                {
                    title: "配置类型",
                    dataIndex: "type",
                    key: "type",
                    render: (value) => {
                        return <Text ellipsis={true}>yaml</Text>;
                    },
                },
                {
                    title: "操作",
                    dataIndex: "操作",
                    key: "操作",
                    align: "center",
                    width: 300,
                    render: (value, record) => {
                        return (
                            <div style={{ textAlign: "center" }}>
                                <Button
                                    type="link"
                                    size="small"
                                    onClick={this.showConfigContent.bind(
                                        this,
                                        record,
                                    )}
                                >
                                    详情
                                </Button>
                                <Divider type="vertical" />
                                <Popconfirm
                                    title="该操作会更新Nacos对应配置，是否继续?"
                                    onConfirm={this.updateConfigContent.bind(
                                        this,
                                        record,
                                    )}
                                    okText="确认"
                                    cancelText="取消"
                                >
                                    <Button
                                        type="link"
                                        size="small"
                                        disabled={
                                            !this.props.aclAuthMap[
                                                "PUT:/configCenter/nacos/config"
                                            ]
                                        }
                                    >
                                        修改
                                    </Button>
                                </Popconfirm>
                                {/* <Divider type="vertical" />
                                <Button
                                    type="link"
                                    size="small"
                                    onClick={this.copyConfigContent.bind(
                                        this,
                                        record,
                                    )}
                                    disabled={
                                        !this.props.aclAuthMap[
                                            "POST:/configCenter/nacos/config/copy"
                                        ]
                                    }
                                >
                                    复制
                                </Button> */}
                                <Divider type="vertical" />
                                <Popconfirm
                                    title="该操作会删除Nacos对应配置，是否继续?"
                                    onConfirm={this.deleteConfigContent.bind(
                                        this,
                                        record,
                                    )}
                                    okText="确认"
                                    cancelText="取消"
                                >
                                    <Button
                                        type="danger"
                                        size="small"
                                        disabled={
                                            !this.props.aclAuthMap[
                                                "DELETE:/configCenter/nacos/config"
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
            tableData: [],
            pagination: {
                showSizeChanger: true,
                pageSizeOptions: ["10", "20"],
                onShowSizeChange: (current, size) =>
                    this.onShowSizeChange(current, size),
                showQuickJumper: false,
                showTotal: (total) => `共 ${total} 条`,
                pageSize: 10,
                page: 1,
                total: 0,
                onChange: (page, pageSize) => this.changePage(page, pageSize),
            },
            tableLoading: false,
            configContentDrawerVisible: false,
            configDetail: {},
            treeData: [],
            checkedKeys: [],
            allConfigTemplatesMap: {},
            allConfigTemplatesList: [],
            updateTemplateId: 0,
        };
    }

    componentDidMount() {
        this.loadNacosServerList();
        this.loadConfigAllTemplates();
    }

    loadNacosServerList() {
        getNacosList().then((res) => {
            if (res.code === 0) {
                this.setState({ nacosList: res.data });
            } else {
                message.error(res.msg);
            }
        });
    }

    loadNacosNamespaceList() {
        if (this.state.currentNacosServerId !== "") {
            getNacosNamespaceList({
                clusterId: this.state.currentNacosServerId,
            }).then((res) => {
                if (res.code === 0) {
                    this.setState({ nsList: res.data.data });
                } else {
                    message.error(res.msg);
                }
            });
        }
    }

    loadNamespaceAllConfigs() {
        getNacosAllConfigs({
            clusterId: "" + this.state.currentNacosServerId,
        }).then((res) => {
            if (res.code === 0) {
                let treeData = [];
                let data = res.data;
                for (let i = 0; i < data.length; i++) {
                    let configs = data[i]["configs"];
                    if (configs !== null) {
                        let children = [];
                        for (let j = 0; j < configs.length; j++) {
                            children.push({
                                title:
                                    configs[j]["dataId"] +
                                    "  |  " +
                                    configs[j]["group"],
                                key:
                                    configs[j]["dataId"] +
                                    "  |  " +
                                    configs[j]["group"],
                                namespace: data[i]["namespace"],
                                dataId: configs[j]["dataId"],
                                group: configs[j]["group"],
                            });
                        }
                        treeData.push({
                            title: data[i]["namespace"],
                            key: data[i]["namespace"],
                            children: children,
                        });
                    }
                }
                this.setState({ treeData: treeData });
            } else {
                message.error(res.msg);
            }
        });
    }

    loadConfigAllTemplates() {
        getConfigTemplatesAll().then((res) => {
            if (res.code === 0) {
                let allConfigTemplatesMap = {};
                for (let i = 0; i < res.data.length; i++) {
                    allConfigTemplatesMap[res.data[i]["Id"]] = res.data[i];
                }
                this.setState({
                    allConfigTemplatesMap: allConfigTemplatesMap,
                    allConfigTemplatesList: res.data,
                });
            }
        });
    }

    loadNamespaceConfigs() {
        if (this.state.currentNacosServerId === "") {
            message.info("请选择集群!");
            return;
        }
        if (this.state.currentNs === "") {
            message.info("请选择命名空间!");
            return;
        }
        this.setState({ tableLoading: true });
        const queryParams = {
            page: parseInt(this.state.pagination.page),
            size: parseInt(this.state.pagination.pageSize),
            clusterId: this.state.currentNacosServerId,
            namespace: this.state.currentNs,
        };
        getNacosConfigsList(queryParams).then((res) => {
            if (res.code === 0) {
                const pagination = this.state.pagination;
                pagination.total = parseInt(res.data.count);
                pagination.page = parseInt(res.data.currentPage);
                pagination.showTotal(parseInt(res.data.count));
                this.setState({
                    pagination,
                });
                this.setState({ tableData: res.data.data });
            } else {
                message.error(res.msg);
            }
            this.setState({ tableLoading: false });
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
                this.loadNamespaceConfigs();
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
                this.loadNamespaceConfigs();
            },
        );
    };

    showConfigContent(record) {
        let params = {
            clusterId: "" + this.state.currentNacosServerId,
            namespace: record.Namespace,
            dataId: record.DataId,
            group: record.ConfigGroup,
        };
        getNacosConfigDetail(params).then((res) => {
            if (res.code === 0) {
                this.setState({
                    configContentDrawerVisible: true,
                    configDetail: res.data["nacosData"],
                });
            } else {
                message.error(res.msg);
            }
        });
    }

    onCloseConfigContentDrawer = () => {
        this.setState({
            configContentDrawerVisible: false,
        });
    };

    // copyConfigContent(record) {
    //     let that = this;
    //     this.setState({ copyNacosConfigModalVisible: true }, () => {
    //         setTimeout(() => {
    //             that.configCopyFormRef.current.setFieldsValue({
    //                 srcNamespace: record.tenant,
    //                 srcDataId: record.dataId,
    //                 srcGroup: record.group,
    //             });
    //         }, 300);
    //     });
    // }

    updateConfigContent(record) {
        let that = this;
        let params = {
            clusterId: "" + this.state.currentNacosServerId,
            namespace: record.Namespace,
            dataId: record.DataId,
            group: record.ConfigGroup,
        };
        getNacosConfigDetail(params).then((res) => {
            if (res.code === 0) {
                let originData = res.data["originData"];
                let nacosData = res.data["nacosData"];
                this.setState(
                    {
                        addNacosConfigModalVisible: true,
                        updateTemplateId: originData.TemplateId,
                        nacosConfigContent: nacosData["Content"],
                    },
                    () => {
                        setTimeout(() => {
                            let fillData =
                                originData.FillData !== ""
                                    ? JSON.parse(originData.FillData)
                                    : {};
                            that.configFormRef.current.setFieldsValue({
                                configId: nacosData["Id"],
                                id: originData["Id"],
                                dataId: nacosData.DataId,
                                group: nacosData.Group,
                                configType: nacosData.Type,
                                content:
                                    originData.TemplateId !== 0
                                        ? originData.ConfigContent
                                        : nacosData["Content"],
                                templateId: originData.TemplateId,
                                isUseTemplate:
                                    originData.TemplateId !== 0
                                        ? "useTemplate"
                                        : "notUserTemplate",
                                templateId:
                                    originData.TemplateId === 0
                                        ? 0
                                        : originData.TemplateId,
                                ...fillData,
                                // configTags:
                                //     res.data.ConfigTags === "static"
                                //         ? "static"
                                //         : "dynamic",
                            });
                        }, 200);
                    },
                );
            } else {
                message.error(res.msg);
            }
        });
    }

    deleteConfigContent(record) {
        deleteNacosConfig({
            clusterId: "" + this.state.currentNacosServerId,
            id: parseInt(record.Id),
            namespace: this.state.currentNs,
            dataId: record.DataId,
            group: record.ConfigGroup,
        }).then((res) => {
            if (res.code === 0) {
                message.success("删除成功");
                this.loadNamespaceConfigs();
            } else {
                message.error(res.msg);
            }
        });
    }

    addNacosServer() {
        this.setState({ addNacosModalVisible: true });
    }

    addNacosConfig() {
        if (this.state.currentNacosServerId === "") {
            message.info("请选择集群!");
            return;
        }
        if (this.state.currentNs === "") {
            message.info("请选择命名空间!");
            return;
        }
        this.setState({ addNacosConfigModalVisible: true });
    }

    handleNacosChange(e) {
        this.setState({ currentNacosServerId: e, currentNs: "" }, () => {
            this.loadNacosNamespaceList();
        });
    }

    handleNsChange(e) {
        this.setState({ currentNs: e }, () => {
            this.loadNamespaceConfigs();
        });
    }

    nacosConfigQuery() {
        this.loadNamespaceConfigs();
    }

    submitAddNacosServer = (e) => {
        e.preventDefault();
        this.serverFormRef.current.validateFields().then((values) => {
            let params = {
                ...values,
            };
            postNacosServer(params).then((res) => {
                if (res.code === 0) {
                    this.setState({ addNacosModalVisible: false });
                    this.loadNacosServerList();
                } else {
                    message.error(res.msg);
                }
            });
        });
    };

    cancelAddNacosServer() {
        this.setState({ addNacosModalVisible: false });
    }

    submitCreateConfig = (e) => {
        e.preventDefault();
        this.configFormRef.current.validateFields().then((values) => {
            let params = {
                ...values,
                clusterId: "" + this.state.currentNacosServerId,
                namespace: this.state.currentNs,
            };
            if (values.configId === undefined) {
                postNacosConfig(params).then((res) => {
                    if (res.code === 0) {
                        message.success("创建成功");
                        this.setState({ addNacosConfigModalVisible: false });
                        this.loadNamespaceConfigs();
                    } else {
                        message.error(res.msg);
                    }
                });
            } else {
                putNacosConfig(params).then((res) => {
                    if (res.code === 0) {
                        message.success("修改成功");
                        this.setState({ addNacosConfigModalVisible: false });
                        this.loadNamespaceConfigs();
                    } else {
                        message.error(res.msg);
                    }
                });
            }
        });
    };

    submitCreateConfigCopy = (e) => {
        e.preventDefault();
        this.configCopyFormRef.current.validateFields().then((values) => {
            let params = {
                ...values,
                clusterId: "" + this.state.currentNacosServerId,
            };
            postNacosConfigCopy(params).then((res) => {
                if (res.code === 0) {
                    message.success("复制成功!");
                    this.setState({ copyNacosConfigModalVisible: false });
                } else {
                    message.error(res.msg);
                }
            });
        });
    };

    cancelCreateConfig() {
        this.setState({ addNacosConfigModalVisible: false });
    }

    cancelCreateConfigCopy() {
        this.setState({ copyNacosConfigModalVisible: false });
    }

    submitSyncConfigContent = (e) => {
        e.preventDefault();
        let params = {
            clusterId: "" + this.state.currentNacosServerId,
            srcNamespace: this.state.configDetail.Tenant,
            srcDataId: this.state.configDetail.DataId,
            srcGroup: this.state.configDetail.Group,
            dstConfigs: this.state.checkedKeys,
        };
        postNacosConfigSync(params).then((res) => {
            if (res.code === 0) {
                message.success("操作成功");
                this.setState({ syncNacosConfigModalVisible: false });
            } else {
                message.error(res.msg);
            }
        });
    };

    cancelSyncConfigContent() {
        this.setState({ syncNacosConfigModalVisible: false });
    }

    onCheck = (checkedKeys, info) => {
        let leafNodes = [];
        for (var i = 0; i < info.checkedNodes.length; i++) {
            if (!("children" in info.checkedNodes[i])) {
                leafNodes.push(info.checkedNodes[i]);
            }
        }
        this.setState({ checkedKeys: leafNodes });
    };

    syncStaticConfigToOthers() {
        this.loadNamespaceAllConfigs();
        this.setState({ syncNacosConfigModalVisible: true });
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
                <OpsBreadcrumbPath pathData={["配置中心", "Nacos管理"]} />
                <Row style={{ marginBottom: 20 }}>
                    <Col span={7}>
                        选择集群:&nbsp;&nbsp;
                        <Select
                            style={{ width: "200px" }}
                            onChange={this.handleNacosChange.bind(this)}
                        >
                            {this.state.nacosList.map((item, index) => {
                                return (
                                    <Option key={index} value={item.Id}>
                                        {item.Alias}
                                    </Option>
                                );
                            })}
                        </Select>
                    </Col>
                    <Col span={7}>
                        选择命名空间:&nbsp;&nbsp;
                        <Select
                            style={{ width: "200px" }}
                            onChange={this.handleNsChange.bind(this)}
                        >
                            {this.state.nsList.map((item, index) => {
                                return (
                                    <Option
                                        key={index}
                                        value={item.NamespaceShowName}
                                    >
                                        {item.NamespaceShowName}
                                    </Option>
                                );
                            })}
                        </Select>
                    </Col>
                    <Col span={10}>
                        <Button
                            onClick={this.nacosConfigQuery.bind(this)}
                            disabled={
                                !this.props.aclAuthMap[
                                    "GET:/configCenter/nacos/configs"
                                ]
                            }
                        >
                            查询
                        </Button>
                        &nbsp;&nbsp;
                        <Button
                            type="primary"
                            onClick={this.addNacosConfig.bind(this)}
                            disabled={
                                !this.props.aclAuthMap[
                                    "POST:/configCenter/nacos/config"
                                ]
                            }
                        >
                            新增配置
                        </Button>
                        &nbsp;&nbsp;
                        <Button
                            type="primary"
                            onClick={this.addNacosServer.bind(this)}
                            disabled={
                                !this.props.aclAuthMap[
                                    "POST:/configCenter/nacos"
                                ]
                            }
                        >
                            添加Nacos集群
                        </Button>
                    </Col>
                </Row>

                <Modal
                    title="新增Nacos集群"
                    visible={this.state.addNacosModalVisible}
                    onOk={this.submitAddNacosServer}
                    onCancel={this.cancelAddNacosServer.bind(this)}
                    destroyOnClose={true}
                >
                    <NacosServerForm formRef={this.serverFormRef} />
                </Modal>

                <Modal
                    title="配置信息"
                    visible={this.state.addNacosConfigModalVisible}
                    onOk={this.submitCreateConfig}
                    onCancel={this.cancelCreateConfig.bind(this)}
                    width={800}
                    destroyOnClose={true}
                >
                    <NacosConfigForm
                        formRef={this.configFormRef}
                        templatesList={this.state.allConfigTemplatesList}
                        templatesMap={this.state.allConfigTemplatesMap}
                        fillData={this.state.fillData}
                        updateTemplateId={this.state.updateTemplateId}
                    />
                </Modal>

                {/* <Modal
                    title="配置复制"
                    visible={this.state.copyNacosConfigModalVisible}
                    onOk={this.submitCreateConfigCopy}
                    onCancel={this.cancelCreateConfigCopy.bind(this)}
                    width={500}
                    destroyOnClose={true}
                >
                    <NacosConfigCopyForm
                        nsList={this.state.nsList}
                        formRef={this.configCopyFormRef}
                    />
                </Modal> */}

                {/* <Modal
                    title="配置同步"
                    visible={this.state.syncNacosConfigModalVisible}
                    onOk={this.submitSyncConfigContent}
                    onCancel={this.cancelSyncConfigContent.bind(this)}
                    width={600}
                    destroyOnClose={true}
                >
                    <Tree
                        checkable
                        defaultExpandAll={true}
                        onCheck={this.onCheck}
                        treeData={this.state.treeData}
                    />
                </Modal> */}

                <Drawer
                    title="配置信息"
                    placement="left"
                    width={800}
                    closable={false}
                    onClose={this.onCloseConfigContentDrawer}
                    visible={this.state.configContentDrawerVisible}
                >
                    <Descriptions title="配置详情" size="small" column={2}>
                        <Descriptions.Item label="配置类型">
                            {this.state.configDetail.Type}
                        </Descriptions.Item>
                        {/* <Descriptions.Item label="创建人">
                            {this.state.configDetail.CreateUser}
                        </Descriptions.Item> */}
                        <Descriptions.Item label="创建IP">
                            {this.state.configDetail.CreateIp}
                        </Descriptions.Item>
                        {/* <Descriptions.Item label="配置标签">
                            {this.state.configDetail.ConfigTags}
                        </Descriptions.Item> */}
                    </Descriptions>
                    <Divider />
                    {/* <Alert
                        message={
                            <div>
                                <span>
                                    目前支持text、yaml、properties类型的静态配置的一键同步
                                </span>
                                <br />
                                <span>
                                    同步时时请确认其它配置中未包含当前配置内容，否则会重复
                                </span>
                            </div>
                        }
                        type="warning"
                    /> */}
                    <Text strong>配置内容: </Text>
                    <Paragraph
                        style={{ display: "inline-block", width: "40px" }}
                        copyable={{
                            icon: <span>复制</span>,
                            text: this.state.configDetail.Content,
                        }}
                    />
                    {/* <Button
                        type="link"
                        disabled={
                            !(
                                this.state.configDetail.ConfigTags ===
                                    "static" &&
                                (this.state.configDetail.Type === "text" ||
                                    this.state.configDetail.Type === "yaml" ||
                                    this.state.configDetail.Type ===
                                        "properties")
                            )
                        }
                        onClick={this.syncStaticConfigToOthers.bind(this)}
                    >
                        一键同步
                    </Button> */}
                    <pre
                        style={{
                            backgroundColor: "rgb(36, 35, 35)",
                            color: "#eee",
                            padding: "10px 10px",
                        }}
                    >
                        {this.state.configDetail.Content}
                    </pre>
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

export default withRouter(NacosContent);
