import React, { Component, Fragment } from "react";
import OpsBreadcrumbPath from "../breadcrumb_path";
import {
    Steps,
    Button,
    message,
    Input,
    Carousel,
    Form,
    Cascader,
    Radio,
    Tooltip,
    Card,
    Alert,
    Row,
    Select,
    Layout,
} from "antd";
import "../../assets/css/job.css";
import Col from "antd/es/grid/col";
import {
    SolutionOutlined,
    HighlightOutlined,
    PlusCircleOutlined,
    MinusCircleOutlined,
} from "@ant-design/icons";
import { postAddDailyJob } from "../../api/daily_task";

const Step = Steps.Step;
const { Content } = Layout;
const { Option } = Select;
const { TextArea } = Input;

const options = [
    {
        value: "日常工作",
        label: "日常工作",
        children: [
            {
                value: "协助",
                label: "协助",
            },
        ],
    },
];

const ConfigTemplate = {
    "": [
        {
            field: "task_content",
            text: "问题描述",
            type: "text",
            errorMessage: "请输入具体问题描述",
            required: true,
        },
    ],
    "日常工作-协助": [
        {
            field: "task_content",
            text: "工作内容",
            type: "text",
            errorMessage: "请输入具体工作内容",
        },
        {
            field: "open_deploy_auto_config",
            text: "自定义项",
            type: "kvinput",
        },
    ],
};

class KvIput extends Component {
    constructor(props) {
        super(props);
        this.addNewConfig = this.addNewConfig.bind(this);
        this.state = {
            configCount: 0,
            configData: [],
        };
    }

    componentDidMount() {
        this.onChange = this.props.onChange;
    }

    deleteConfigItem(data) {
        let index = data.index;
        let configData = this.state.configData;
        configData.splice(index, 1);
        this.setState({ configData });
        this.onChange(this.state.configData);
    }

    addNewConfig() {
        let configData = this.state.configData;
        configData.push({
            key: "",
            value: "",
        });
        this.setState(configData);
    }

    handleKeyChange(data, e) {
        let index = data["index"];
        let configData = this.state.configData;
        configData[index]["key"] = e.target.value;
        this.setState(configData);
        this.onChange(this.state.configData);
    }

    handleValueChange(data, e) {
        let index = data["index"];
        let configData = this.state.configData;
        configData[index]["value"] = e.target.value;
        this.setState(configData);
        this.onChange(this.state.configData);
    }

    render() {
        let configItem = this.state.configData.map((item, index) => {
            return (
                <Row>
                    <Col span={10} style={{ marginRight: 10 }}>
                        <Input
                            key="k{index}"
                            addonBefore="键"
                            value={item.key}
                            onChange={this.handleKeyChange.bind(this, {
                                index,
                            })}
                        />
                    </Col>
                    <Col span={10} style={{ marginRight: 10 }}>
                        <Input
                            key="v{index}"
                            addonBefore="值"
                            value={item.value}
                            onChange={this.handleValueChange.bind(this, {
                                index,
                            })}
                        />
                    </Col>
                    <Col span={2}>
                        <Button
                            type="sdanger"
                            icon={<MinusCircleOutlined />}
                            shape="circle"
                            onClick={this.deleteConfigItem.bind(this, {
                                index,
                            })}
                        />
                    </Col>
                </Row>
            );
        });
        return (
            <Fragment>
                {configItem}
                <Row>
                    <Col span={5}>
                        <Button
                            type="primary"
                            icon={<PlusCircleOutlined />}
                            onClick={this.addNewConfig}
                        >
                            添加新项
                        </Button>
                    </Col>
                </Row>
            </Fragment>
        );
    }
}

class JobBaseInfo extends Component {
    constructor(props) {
        super(props);
        this.state = {
            deptData: [],
        };
    }

    render() {
        const formItemLayout = {
            labelCol: { span: 8 },
            wrapperCol: { span: 10 },
        };
        return (
            <Content
                style={{
                    height: "100%",
                }}
            >
                <Form
                    ref={this.props.formRef}
                    initialValues={{ importantDegree: "普通" }}
                >
                    <Form.Item
                        {...formItemLayout}
                        label="工单名称"
                        name="jobName"
                        rules={[
                            { required: true, message: "请输入工单任务名称" },
                        ]}
                    >
                        <Input placeholder="请输入工单任务名称" />
                    </Form.Item>
                    <Form.Item
                        {...formItemLayout}
                        label="工单类型"
                        name="jobType"
                        rules={[{ required: true, message: "请选择工单类型" }]}
                    >
                        <Cascader
                            placeholder="请选择工单类型"
                            options={options}
                        />
                    </Form.Item>

                    <Form.Item
                        {...formItemLayout}
                        label="该任务紧急程度"
                        style={{ textAlign: "left" }}
                        name="importantDegree"
                        rules={[
                            {
                                required: true,
                                message: "请选择任务紧急程度",
                            },
                        ]}
                    >
                        <Radio.Group buttonStyle="solid">
                            <Tooltip
                                placement="topLeft"
                                title="等待相关人员根据自身任务项安排时间处理"
                            >
                                <Radio.Button value="普通">普通</Radio.Button>
                            </Tooltip>
                            <Tooltip
                                placement="topLeft"
                                title="需要在半小时内处理完成"
                            >
                                <Radio.Button value="紧急">紧急</Radio.Button>
                            </Tooltip>
                            <Tooltip
                                placement="topLeft"
                                title="需要立即安排人员处理"
                            >
                                <Radio.Button value="非常紧急">
                                    非常紧急
                                </Radio.Button>
                            </Tooltip>
                        </Radio.Group>
                    </Form.Item>
                </Form>
            </Content>
        );
    }
}

class JobConfigInfo extends Component {
    constructor(props) {
        super(props);
        this.getFormInputContent = this.getFormInputContent.bind(this);
        this.state = {};
    }

    getFormInputContent(input) {
        let res;
        switch (input.type) {
            case "char":
                res = <Input />;
                break;
            case "text":
                res = <TextArea rows={4} />;
                break;
            case "select":
                res = (
                    <Select style={{ width: 120 }}>
                        {input.options.map((item, index) => {
                            return (
                                <Option key={index} value={item}>
                                    {item}
                                </Option>
                            );
                        })}
                    </Select>
                );
                break;
            case "kvinput":
                res = (
                    // <KvIput setFieldsValue={this.props.form.setFieldsValue} />
                    <KvIput />
                );
                break;
            default:
                res = <Input />;
                break;
        }
        return res;
    }

    render() {
        const formItemLayout = {
            labelCol: { span: 4 },
            wrapperCol: { span: 18 },
        };
        let configData = ConfigTemplate[this.props.jobType];
        if (configData === undefined) {
            configData = ConfigTemplate[""];
        }
        let configDataForm;
        configDataForm = configData.map((item, index) => {
            return (
                <Form.Item
                    key={index}
                    {...formItemLayout}
                    label={item.text}
                    style={{ textAlign: "left" }}
                    name={item.field}
                    rules={[
                        {
                            required: item.required,
                            message: item.errorMessage,
                        },
                    ]}
                >
                    {this.getFormInputContent(item)}
                </Form.Item>
            );
        });
        return (
            <Content
                style={{
                    height: "100%",
                }}
            >
                <Form ref={this.props.formRef}>
                    <Form.Item
                        {...formItemLayout}
                        label="关健信息"
                        required={true}
                    >
                        <Card
                            span={20}
                            size="small"
                            title="填写完善和准确可以提高工单处理效率哦"
                            headStyle={{
                                backgroundColor: "#e6f7ff",
                                fontWeight: "350",
                            }}
                        >
                            {configDataForm}
                        </Card>
                    </Form.Item>
                    <Row style={{ marginBottom: 10 }}>
                        <Col span={4} />
                        <Col span={18}>
                            <Alert
                                message="你可以在备注信息中填写注意事项，可以避免执行过程中出错哦"
                                type="info"
                                closable
                            />
                        </Col>
                    </Row>
                    <Form.Item
                        {...formItemLayout}
                        label="备注信息"
                        name="remark"
                    >
                        <TextArea rows={4} />
                    </Form.Item>
                </Form>
            </Content>
        );
    }
}

class Deploy_project_content extends Component {
    constructor(props) {
        super(props);
        this.submitJobForm = this.submitJobForm.bind(this);
        this.baseInfoFormRef = React.createRef();
        this.configInfoFormRef = React.createRef();
        this.state = {
            current: 0,
            stepStatus: "wait",
            jobTypeStr: "",
            form1Data: null,
            form2Data: null,
            submitLoading: false,
        };
    }

    componentDidMount() {}

    next() {
        if (this.state.current === 0) {
            // 工单基本信息
            this.baseInfoFormRef.current
                .validateFields()
                .then((values) => {
                    const current = this.state.current + 1;
                    let form1Data = Object.assign(
                        {},
                        this.state.form1Data,
                        values,
                    );
                    this.setState({
                        current,
                        stepStatus: "process",
                        jobTypeStr: values.jobType.join("-"),
                        form1Data: form1Data,
                    });
                    this.carousel.next();
                })
                .catch((err) => {
                    this.setState({ stepStatus: "error" });
                });
        }
    }

    prev() {
        switch (this.state.current) {
            case 2:
                this.setState({ form3Data: null });
                break;
            case 1:
                this.setState({ form2Data: null });
                break;
            default:
                break;
        }
        const current = this.state.current - 1;
        this.setState({ current });
        this.carousel.prev();
    }

    submitJobForm() {
        this.configInfoFormRef.current
            .validateFields()
            .then((values) => {
                let form2Data = Object.assign({}, this.state.form2Data, values);
                this.setState(
                    {
                        stepStatus: "finish",
                        form2Data: form2Data,
                    },
                    () => {
                        this.setState({ submitLoading: true });
                        let reqData = {
                            ...this.state.form1Data,
                            ...this.state.form2Data,
                            jobType: this.state.form1Data["jobType"].join("-"),
                            open_deploy_auto_config: JSON.stringify(
                                this.state.form2Data.open_deploy_auto_config,
                            ),
                        };
                        postAddDailyJob(reqData)
                            .then((res) => {
                                if (res.code === 0) {
                                    this.setState({ submitLoading: false });
                                    message.success(res.msg);
                                    this.props.history.push("/admin/task/jobs");
                                } else {
                                    this.setState({ submitLoading: false });
                                    message.error(res.msg);
                                }
                            })
                            .catch((err) => {
                                this.setState({ submitLoading: false });
                                message.error(err.toLocaleString());
                            });
                    },
                );
            })
            .catch((err) => {
                this.setState({ stepStatus: "error" });
            });
    }

    render() {
        const { current } = this.state;
        return (
            <Content
                style={{
                    background: "#fff",
                    margin: 0,
                    height: "100%",
                }}
            >
                <OpsBreadcrumbPath pathData={["工作协助", "日常工单"]} />
                <Steps
                    current={current}
                    size="small"
                    style={{ marginTop: 20 }}
                    status={this.state.stepStatus}
                >
                    <Step title="任务基本信息" icon={<SolutionOutlined />} />
                    <Step title="任务详情" icon={<HighlightOutlined />} />
                </Steps>
                <div
                    className="steps-content"
                    style={{ textAlign: "center", padding: 20 }}
                >
                    <Carousel
                        ref={(carousel) => {
                            this.carousel = carousel;
                        }}
                        dots={false}
                        effect="fade"
                        adaptiveHeight={true}
                        style={{ textAlign: "center" }}
                    >
                        <Fragment key={1}>
                            <JobBaseInfo formRef={this.baseInfoFormRef} />
                        </Fragment>
                        <Fragment key={2}>
                            <JobConfigInfo
                                formRef={this.configInfoFormRef}
                                jobType={this.state.jobTypeStr}
                                loadConfigDataTemplateSpin={
                                    this.state.loadConfigDataTemplateSpin
                                }
                            />
                        </Fragment>
                    </Carousel>
                </div>
                <div
                    className="steps-action"
                    style={{ textAlign: "center", marginTop: 20 }}
                >
                    {current < 1 && (
                        <Button type="primary" onClick={() => this.next()}>
                            下一步
                        </Button>
                    )}
                    {current === 1 && (
                        <Button
                            type="primary"
                            loading={this.state.submitLoading}
                            onClick={this.submitJobForm}
                        >
                            提 交
                        </Button>
                    )}
                    {current > 0 && (
                        <Button
                            style={{ marginLeft: 8 }}
                            onClick={() => this.prev()}
                        >
                            上一步
                        </Button>
                    )}
                </div>
            </Content>
        );
    }
}

export default Deploy_project_content;
