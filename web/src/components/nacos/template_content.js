import React, { Component } from "react";
import {
    Layout,
    message,
    Button,
    Typography,
    Select,
    Input,
    Table,
    Divider,
    Row,
    Col,
    Modal,
    Form,
    Tag,
    Drawer,
    Popconfirm,
} from "antd";
import OpsBreadcrumbPath from "../breadcrumb_path";
import moment from "moment";
import { withRouter } from "react-router-dom";
import { MinusCircleOutlined, PlusOutlined } from "@ant-design/icons";
import { Controlled as CodeMirror } from "react-codemirror2";
import "codemirror/lib/codemirror.css";
import {
    getConfigTemplates,
    postConfigTemplate,
    deleteConfigTemplate,
    putConfigTemplate,
} from "../../api/nacos";

const { Paragraph, Text } = Typography;
const { Content } = Layout;
const { Option } = Select;
const { TextArea } = Input;

class NacosTemplateContent extends Component {
    constructor(props) {
        super(props);
        this.templateFormRef = React.createRef();
        this.state = {
            columns: [
                {
                    title: "模板名称",
                    dataIndex: "Name",
                    key: "Name",
                    align: "center",
                    render: (value) => {
                        return <Text ellipsis={true}>{value}</Text>;
                    },
                },
                {
                    title: "更新时间",
                    dataIndex: "UpdateTime",
                    key: "UpdateTime",
                    align: "center",
                    render: (value) => {
                        return (
                            <Text ellipsis={true}>
                                {moment(value).format("YYYY-MM-DD HH:mm:ss")}
                            </Text>
                        );
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
                                    onClick={this.showTemplateDetail.bind(
                                        this,
                                        record,
                                    )}
                                >
                                    详情
                                </Button>
                                <Divider type="vertical" />
                                <Button
                                    type="link"
                                    size="small"
                                    onClick={this.updateConfigTemplate.bind(
                                        this,
                                        record,
                                    )}
                                    disabled={
                                        !this.props.aclAuthMap[
                                            "PUT:/configCenter/configTemplate"
                                        ]
                                    }
                                >
                                    修改
                                </Button>
                                <Divider type="vertical" />
                                <Popconfirm
                                    title="删除会取消配置与模板的关联，确定删除吗?"
                                    onConfirm={this.deleteConfigTemplate.bind(
                                        this,
                                        record,
                                    )}
                                    okText="确定删除"
                                    cancelText="取消"
                                >
                                    <Button
                                        type="danger"
                                        size="small"
                                        disabled={
                                            !this.props.aclAuthMap[
                                                "DELETE:/configCenter/configTemplate"
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
            queryName: "",
            addNewTemplateModalVisible: false,
            templateInput: "",
            templateDetail: { ConfigContent: "", FillField: [] },
        };
    }

    componentDidMount() {
        this.loadConfigTemplates();
    }

    loadConfigTemplates() {
        this.setState({ tableLoading: true });
        const queryParams = {
            page: parseInt(this.state.pagination.page),
            size: parseInt(this.state.pagination.pageSize),
            name: this.state.queryName,
        };
        getConfigTemplates(queryParams).then((res) => {
            if (res.code === 0) {
                var data = res.data;
                const pagination = this.state.pagination;
                pagination.total = parseInt(data.count);
                pagination.page = parseInt(this.state.pagination.page);
                pagination.showTotal(parseInt(data.count));
                this.setState({
                    pagination,
                });
                this.setState({ tableData: data.data });
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
                this.loadConfigTemplates();
            },
        );
    };

    queryNameOnChange = (e) => {
        this.setState({ queryName: e.target.value });
    };

    configTemplatesQuery = () => {
        this.loadConfigTemplates();
    };

    addNewTemplates = () => {
        this.setState({ addNewTemplateModalVisible: true });
    };

    handleNewTemplateSubmit = () => {
        this.templateFormRef.current.validateFields().then((values) => {
            let params = {
                ...values,
                configContent: this.state.templateInput,
            };
            if (params.id === undefined) {
                postConfigTemplate(params).then((res) => {
                    if (res.code === 0) {
                        message.success("创建成功");
                        this.setState({
                            addNewTemplateModalVisible: false,
                            templateInput: "",
                        });
                        this.loadConfigTemplates();
                    } else {
                        message.error(res.msg);
                    }
                });
            } else {
                putConfigTemplate(params).then((res) => {
                    if (res.code === 0) {
                        message.success(res.msg);
                        this.setState({
                            addNewTemplateModalVisible: false,
                            templateInput: "",
                        });
                        this.loadConfigTemplates();
                    } else {
                        message.error(res.msg);
                    }
                });
            }
        });
    };

    handleNewTemplateCancel = () => {
        this.setState({ addNewTemplateModalVisible: false });
    };

    showTemplateDetail(record) {
        this.setState({
            templateDetailDrawerVisible: true,
            templateDetail: {
                ConfigContent: record.ConfigContent,
                FillField: JSON.parse(record.FillField),
            },
        });
    }

    updateConfigTemplate(record) {
        this.setState(
            {
                templateInput: record.ConfigContent,
                addNewTemplateModalVisible: true,
            },
            () => {
                setTimeout(() => {
                    this.templateFormRef.current.setFieldsValue({
                        id: record.Id,
                        name: record.Name,
                        fillField: JSON.parse(record.FillField),
                    });
                }, 300);
            },
        );
    }

    deleteConfigTemplate(record) {
        deleteConfigTemplate({ id: record.Id }).then((res) => {
            if (res.code === 0) {
                message.success("删除成功");
                this.loadConfigTemplates();
            } else {
                message.error(res.msg);
            }
        });
    }

    templateDetailDrawerOnClose = () => {
        this.setState({ templateDetailDrawerVisible: false });
    };

    render() {
        const formItemLayout = {
            labelCol: {
                xs: { span: 24 },
                sm: { span: 4 },
            },
            wrapperCol: {
                xs: { span: 24 },
                sm: { span: 20 },
            },
        };
        const formItemLayoutWithOutLabel = {
            wrapperCol: {
                xs: { span: 24, offset: 0 },
                sm: { span: 20, offset: 4 },
            },
        };
        return (
            <Content
                style={{
                    background: "#fff",
                    padding: "5px 20px",
                    margin: 0,
                    height: "100%",
                }}
            >
                <OpsBreadcrumbPath pathData={["配置中心", "配置模板管理"]} />
                <Row style={{ marginBottom: 20 }}>
                    <Col>
                        <Text>模板名称: </Text>&nbsp;&nbsp;
                        <Input
                            style={{ width: 200 }}
                            placeholder="输入模板名称"
                            onChange={this.queryNameOnChange}
                        />
                    </Col>
                    &nbsp;&nbsp;
                    <Button type="primary" onClick={this.configTemplatesQuery}>
                        查询
                    </Button>
                    &nbsp;&nbsp;
                    <Button onClick={this.addNewTemplates}>新增模板</Button>
                </Row>

                <Modal
                    title="模板信息"
                    visible={this.state.addNewTemplateModalVisible}
                    onOk={this.handleNewTemplateSubmit}
                    onCancel={this.handleNewTemplateCancel}
                    destroyOnClose={true}
                >
                    <Form ref={this.templateFormRef} {...formItemLayout}>
                        <Form.Item hidden={true} name="id"></Form.Item>
                        <Form.Item
                            label="模板名称"
                            name="name"
                            rules={[
                                {
                                    required: true,
                                    message: "模板名称不能为空！",
                                },
                            ]}
                        >
                            <Input />
                        </Form.Item>
                        <Text>模板内容:</Text>
                        <CodeMirror
                            className="jenkinsEditor"
                            options={{
                                showCursorWhenSelecting: true,
                                option: {
                                    autofocus: true,
                                },
                                lineWrapping: true,
                            }}
                            value={this.state.templateInput}
                            onBeforeChange={(editor, data, value) => {
                                this.setState({ templateInput: value });
                            }}
                        />
                        <Form.List name="fillField">
                            {(fields, { add, remove }) => (
                                <>
                                    {fields.map((field, index) => (
                                        <Form.Item
                                            {...(index === 0
                                                ? formItemLayout
                                                : formItemLayoutWithOutLabel)}
                                            label={
                                                index === 0 ? "动态配置项" : ""
                                            }
                                            style={{
                                                marginTop: 10,
                                                marginBottom: 5,
                                            }}
                                            required={false}
                                            key={field.key}
                                        >
                                            <Form.Item
                                                {...field}
                                                validateTrigger={[
                                                    "onChange",
                                                    "onBlur",
                                                ]}
                                                style={{ marginBottom: 5 }}
                                                rules={[
                                                    {
                                                        required: true,
                                                        whitespace: true,
                                                        message:
                                                            "输入动态配置项名称或删除该项",
                                                    },
                                                ]}
                                                noStyle
                                            >
                                                <Input
                                                    placeholder="输入动态配置项名称或删除该项"
                                                    style={{ width: "90%" }}
                                                />
                                            </Form.Item>
                                            <MinusCircleOutlined
                                                className="dynamic-delete-button"
                                                onClick={() =>
                                                    remove(field.name)
                                                }
                                            />
                                        </Form.Item>
                                    ))}
                                    <Form.Item>
                                        <Button
                                            type="dashed"
                                            onClick={() => add()}
                                            icon={<PlusOutlined />}
                                        >
                                            新增动态配置项
                                        </Button>
                                    </Form.Item>
                                </>
                            )}
                        </Form.List>
                    </Form>
                </Modal>

                <Drawer
                    title="模板详情"
                    placement="left"
                    closable={true}
                    width={700}
                    onClose={this.templateDetailDrawerOnClose}
                    visible={this.state.templateDetailDrawerVisible}
                >
                    <Text strong>配置内容: </Text>
                    <br />
                    <pre class="preJenkinsLog">
                        {this.state.templateDetail.ConfigContent}
                    </pre>
                    <Divider />
                    <Text strong>动态填充项: </Text>
                    <br />
                    {this.state.templateDetail.FillField.map((item) => {
                        return (
                            <div>
                                <Tag color="#108ee9">{item}</Tag>
                            </div>
                        );
                    })}
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

export default withRouter(NacosTemplateContent);
