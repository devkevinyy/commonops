import React, { Component, Fragment } from "react";
import {
    Card,
    Layout,
    Row,
    message,
    Button,
    Divider,
    Col,
    Descriptions,
    Typography,
    Steps,
    Result,
    Select,
    Form,
    Carousel,
    Tabs,
    Input,
    InputNumber,
    Collapse,
} from "antd";
import { Controlled as CodeMirror } from "react-codemirror2";
import "codemirror/lib/codemirror.css";
import OpsBreadcrumbPath from "../breadcrumb_path";
import { withRouter } from 'react-router-dom';
import {
    getClusterData,
    getNamespaces,
    postApplyYaml,
} from "../../api/kubernetes";
import {
    postCdProcessTemplate,
    getCdProcessTemplateData,
    postCdProcessLog,
} from "../../api/cd";
import { getClusterId } from "../../utils/axios";

const { Paragraph, Text } = Typography;
const { Content } = Layout;
const { Step } = Steps;
const { Option } = Select;
const { TabPane } = Tabs;
const { Panel } = Collapse;

class ClusterConfigContent extends Component {
    constructor(props) {
        super(props);
        this.state = {
            clusterData: [],
            namespaceData: [],
        };
    }

    componentDidMount() {
        this.loadClusterData();
    }

    loadClusterData() {
        getClusterData()
            .then((res) => {
                if (res.code === 0) {
                    this.setState({ clusterData: res.data['k8sData'] });
                } else {
                    message.error(res.msg);
                }
            })
    }

    getK8sNamespaces() {
        getNamespaces()
            .then((res) => {
                if (res.code === 0) {
                    this.setState({ namespaceData: res.data.items });
                } else {
                    message.error(res.msg);
                }
            })
    }

    clusterOnChange(clusterId) {
        localStorage.setItem("clusterId", clusterId);
        this.getK8sNamespaces();
    }

    render() {
        const layout = {
            labelCol: { span: 8 },
            wrapperCol: { span: 10 },
        };
        return (
            <div>
                <Form {...layout} ref={this.props.formRef}>
                    <Form.Item
                        label="流程名称"
                        name="processName"
                        rules={[
                            {
                                required: true,
                                message: "请输入部署名称!",
                            },
                        ]}
                    >
                        <Input placeholder="输入部署名称" />
                    </Form.Item>
                    <Form.Item
                        label="发布集群"
                        name="clusterId"
                        rules={[
                            {
                                required: true,
                                message: "请选择要发布到的集群!",
                            },
                        ]}
                    >
                        <Select
                            style={{ width: "100%" }}
                            onChange={this.clusterOnChange.bind(this)}
                        >
                            {this.state.clusterData.map((item, index) => {
                                return (
                                    <Option
                                        key={item.clusterId}
                                        value={item.clusterId}
                                    >
                                        {item.name}
                                    </Option>
                                );
                            })}
                        </Select>
                    </Form.Item>
                    <Form.Item
                        label="命名空间"
                        name="namespace"
                        rules={[
                            {
                                required: true,
                                message: "请选择要发布到的命名空间!",
                            },
                        ]}
                    >
                        <Select style={{ width: "100%" }}>
                            {this.state.namespaceData.map((item, index) => {
                                return (
                                    <Option
                                        key={item["metadata"]["name"]}
                                        value={item["metadata"]["name"]}
                                    >
                                        {item["metadata"]["name"]}
                                    </Option>
                                );
                            })}
                        </Select>
                    </Form.Item>
                </Form>
            </div>
        );
    }
}

class K8sResConfigContent extends Component {
    constructor(props) {
        super(props);
        this.state = {
            PROJECT_NAME: this.props.cdInfo.jobName,
            REPLICAS: 1,
            IMAGE_NAME: this.props.cdInfo.imageName.trim(),
            CONTAINER_PORT: "",
            SERVICE_PORT: "",
            SERVICE_TARGET_PORT: "",
            INGRESS_HOST: "",
            deployYAML: "",
            serviceYAML: "",
            configMapYAML: "",
            ingressYAML: "",
        };
    }

    componentDidMount() {
        this.updateDeployYAMLTemplate();
        this.updateServiceYAMLTemplate();
        this.updateConfigMapYAMLTemplate();
        this.updateIngressYAMLTemplate();
    }

    updateYAMLTemplateCallback() {
        console.log("update: ", this.state.deployYAML);
        this.props.yamlCallback(
            this.state.deployYAML,
            this.state.serviceYAML,
            this.state.configMapYAML,
            this.state.ingressYAML,
        );
    }

    contianerPortOnChange(e) {
        this.setState(
            {
                CONTAINER_PORT: e.target.value.trim(),
            },
            () => {
                this.updateDeployYAMLTemplate();
            },
        );
    }

    replicasOnChange(value) {
        this.setState(
            {
                REPLICAS: value,
            },
            () => {
                this.updateDeployYAMLTemplate();
            },
        );
    }

    servicePortOnChange(e) {
        this.setState(
            {
                SERVICE_PORT: e.target.value.trim(),
            },
            () => {
                this.updateServiceYAMLTemplate();
            },
        );
    }

    serviceTargetPortOnChange(e) {
        this.setState(
            {
                SERVICE_TARGET_PORT: e.target.value.trim(),
            },
            () => {
                this.updateServiceYAMLTemplate();
            },
        );
    }

    ingressHostOnChange(e) {
        this.setState(
            {
                INGRESS_HOST: e.target.value.trim(),
            },
            () => {
                this.updateIngressYAMLTemplate();
            },
        );
    }

    updateDeployYAMLTemplate() {
        let yaml = String.raw`apiVersion: apps/v1
kind: Deployment
metadata:
  name: ${this.state.PROJECT_NAME}-deploy
spec:
  replicas: ${this.state.REPLICAS}
  selector:
    matchLabels:
      app: ${this.state.PROJECT_NAME}
  template:
    metadata:
      labels:
        app: ${this.state.PROJECT_NAME}
    spec:
      containers:
        - image: >-
            ${this.state.IMAGE_NAME}
          name: ${this.state.PROJECT_NAME}
          ports:
            - containerPort: ${this.state.CONTAINER_PORT}
      imagePullSecrets:
        - name: deploy-image-pull-cred-7993294`;
        this.setState({ deployYAML: yaml }, () => {
            this.updateYAMLTemplateCallback();
        });
    }

    updateServiceYAMLTemplate() {
        let yaml = String.raw`apiVersion: v1
kind: Service
metadata:
  name: ${this.state.PROJECT_NAME}-service
spec:
  ports:
  - name: http
    port: ${this.state.SERVICE_PORT}
    protocol: TCP
    targetPort: ${this.state.SERVICE_TARGET_PORT}
  selector:
    app: ${this.state.PROJECT_NAME}
  type: ClusterIP`;
        this.setState({ serviceYAML: yaml }, () => {
            this.updateYAMLTemplateCallback();
        });
    }

    updateConfigMapYAMLTemplate() {
        let yaml = String.raw`apiVersion: v1
kind: ConfigMap
metadata:
  name: ${this.state.PROJECT_NAME}
data:
  # 类属性键；每一个键都映射到一个值`;
        this.setState({ configMapYAML: yaml }, () => {
            this.updateYAMLTemplateCallback();
        });
    }

    updateIngressYAMLTemplate() {
        let yaml = String.raw`apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: ${this.state.PROJECT_NAME}-ingress
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /$1
spec:
  rules:
  - host: ${this.state.INGRESS_HOST}
    http:
      paths:
      - backend:
          serviceName: ${this.state.PROJECT_NAME}-service
          servicePort: ${this.state.SERVICE_PORT}
        path: /
  tls:
  - hosts:
    - ${this.state.INGRESS_HOST}
    secretName: ${this.state.INGRESS_HOST}-ssl`;
        this.setState({ ingressYAML: yaml }, () => {
            this.updateYAMLTemplateCallback();
        });
    }

    render() {
        const layout = {
            labelCol: { span: 8 },
            wrapperCol: { span: 12 },
        };
        return (
            <div style={{ marginTop: "-15px", marginLeft: "20px" }}>
                <Tabs defaultActiveKey="1" style={{ textAlign: "left" }}>
                    <TabPane tab="Deployment" key="1">
                        <Row>
                            <Col span={8}>
                                <Form {...layout}>
                                    <Form.Item
                                        label="容器端口"
                                        name="containerPort"
                                        rules={[
                                            {
                                                required: true,
                                                message: "请输入容器端口!",
                                            },
                                        ]}
                                    >
                                        <Input
                                            placeholder="请输入容器端口"
                                            onChange={this.contianerPortOnChange.bind(
                                                this,
                                            )}
                                        />
                                    </Form.Item>
                                    <Form.Item
                                        label="副本数"
                                        name="replicas"
                                        rules={[
                                            {
                                                required: true,
                                                message: "请输入副本数!",
                                            },
                                        ]}
                                    >
                                        <InputNumber
                                            min={1}
                                            max={10}
                                            defaultValue={1}
                                            onChange={this.replicasOnChange.bind(
                                                this,
                                            )}
                                        />
                                    </Form.Item>
                                </Form>
                            </Col>
                            <Col span={16} style={{ overflowX: "scroll" }}>
                                <Collapse bordered={false} ghost>
                                    <Panel
                                        header="查看或编辑YAML"
                                        key="1"
                                        style={{ marginTop: "-10px" }}
                                    >
                                        <CodeMirror
                                            className="jenkinsEditor"
                                            style={{ overflowX: "scroll" }}
                                            value={this.state.deployYAML}
                                            onBeforeChange={(
                                                editor,
                                                data,
                                                value,
                                            ) => {
                                                this.setState({
                                                    deployYAML: value,
                                                });
                                            }}
                                            onChange={(editor, data, value) => {
                                                this.setState(
                                                    {
                                                        deployYAML: value,
                                                    },
                                                    () => {
                                                        this.updateYAMLTemplateCallback();
                                                    },
                                                );
                                            }}
                                        />
                                    </Panel>
                                </Collapse>
                            </Col>
                        </Row>
                    </TabPane>
                    <TabPane tab="Service" key="2">
                        <Row>
                            <Col span={8}>
                                <Form {...layout}>
                                    <Form.Item
                                        label="服务端口"
                                        name="servicePort"
                                        rules={[
                                            {
                                                required: true,
                                                message: "请输入服务端口!",
                                            },
                                        ]}
                                    >
                                        <Input
                                            placeholder="请输入容器端口"
                                            onChange={this.servicePortOnChange.bind(
                                                this,
                                            )}
                                        />
                                    </Form.Item>
                                    <Form.Item
                                        label="目标端口"
                                        name="serviceTargetPort"
                                        rules={[
                                            {
                                                required: true,
                                                message: "请输入目标端口!",
                                            },
                                        ]}
                                    >
                                        <Input
                                            placeholder="请输入目标端口"
                                            onChange={this.serviceTargetPortOnChange.bind(
                                                this,
                                            )}
                                        />
                                    </Form.Item>
                                </Form>
                            </Col>
                            <Col span={16} style={{ overflowX: "scroll" }}>
                                <Collapse bordered={false} ghost>
                                    <Panel
                                        header="查看或编辑YAML"
                                        key="1"
                                        style={{ marginTop: "-10px" }}
                                    >
                                        <CodeMirror
                                            className="jenkinsEditor"
                                            style={{ overflowX: "scroll" }}
                                            value={this.state.serviceYAML}
                                            onBeforeChange={(
                                                editor,
                                                data,
                                                value,
                                            ) => {
                                                this.setState({
                                                    serviceYAML: value,
                                                });
                                            }}
                                            onChange={(editor, data, value) => {
                                                this.setState(
                                                    {
                                                        serviceYAML: value,
                                                    },
                                                    () => {
                                                        this.updateYAMLTemplateCallback();
                                                    },
                                                );
                                            }}
                                        />
                                    </Panel>
                                </Collapse>
                            </Col>
                        </Row>
                    </TabPane>
                    <TabPane tab="ConfigMap" key="3">
                        <Row>
                            <Col span={24} style={{ overflowX: "scroll" }}>
                                <Collapse bordered={false} ghost>
                                    <Panel
                                        header="查看或编辑YAML"
                                        key="1"
                                        style={{ marginTop: "-10px" }}
                                    >
                                        <CodeMirror
                                            className="jenkinsEditor"
                                            style={{ overflowX: "scroll" }}
                                            value={this.state.configMapYAML}
                                            onBeforeChange={(
                                                editor,
                                                data,
                                                value,
                                            ) => {
                                                this.setState({
                                                    configMapYAML: value,
                                                });
                                            }}
                                            onChange={(editor, data, value) => {
                                                this.setState(
                                                    {
                                                        configMapYAML: value,
                                                    },
                                                    () => {
                                                        this.updateYAMLTemplateCallback();
                                                    },
                                                );
                                            }}
                                        />
                                    </Panel>
                                </Collapse>
                            </Col>
                        </Row>
                    </TabPane>
                    <TabPane tab="Ingress" key="4">
                        <Row>
                            <Col span={8}>
                                <Form {...layout}>
                                    <Form.Item
                                        label="配置域名"
                                        name="ingressHost"
                                        rules={[
                                            {
                                                required: true,
                                                message: "请输入域名!",
                                            },
                                        ]}
                                    >
                                        <Input
                                            placeholder="请输入域名"
                                            onChange={this.ingressHostOnChange.bind(
                                                this,
                                            )}
                                        />
                                    </Form.Item>
                                </Form>
                            </Col>
                            <Col span={16} style={{ overflowX: "scroll" }}>
                                <Collapse bordered={false} ghost>
                                    <Panel
                                        header="查看或编辑YAML"
                                        key="1"
                                        style={{ marginTop: "-10px" }}
                                    >
                                        <CodeMirror
                                            className="jenkinsEditor"
                                            style={{ overflowX: "scroll" }}
                                            value={this.state.ingressYAML}
                                            onBeforeChange={(
                                                editor,
                                                data,
                                                value,
                                            ) => {
                                                this.setState({
                                                    ingressYAML: value,
                                                });
                                            }}
                                            onChange={(editor, data, value) => {
                                                this.setState(
                                                    {
                                                        ingressYAML: value,
                                                    },
                                                    () => {
                                                        this.updateYAMLTemplateCallback();
                                                    },
                                                );
                                            }}
                                        />
                                    </Panel>
                                </Collapse>
                            </Col>
                        </Row>
                    </TabPane>
                </Tabs>
            </div>
        );
    }
}

class CdContent extends Component {
    constructor(props) {
        super(props);
        this.clusterFormRef = React.createRef();
        this.deployYAML = "";
        this.serviceYAML = "";
        this.configMapYAML = "";
        this.ingressYAML = "";
        this.state = {
            jobs: [],
            ableToCd: false,
            cdInfo: {
                jobName: "",
                buildNumber: "",
                imageName: "",
                ...this.props.location.state,
            },
            currentStep: 0,
            clusterId: "",
            namespace: "",
            processName: "",
            deployProcessType: "1",
            deployProcessRecord: 0,
            processTemplateList: [],
        };
    }

    componentDidMount() {
        if (
            this.state.cdInfo.imageName === "" ||
            this.state.cdInfo.imageName === "暂无资源"
        ) {
            message.warn("没有构建的镜像制品，无法进行发布!");
        } else {
            this.setState({ ableToCd: true });
        }
        this.loadCdProcessTemplateList();
    }

    yamlCallback(deploy, service, cm, ingress) {
        this.deployYAML = deploy;
        this.serviceYAML = service;
        this.configMapYAML = cm;
        this.ingressYAML = ingress;
    }

    next() {
        let current = this.state.currentStep + 1;
        if (this.state.currentStep === 0) {
            this.clusterFormRef.current
                .validateFields()
                .then((values) => {
                    console.log("表单验证: ", values);
                    this.setState({
                        currentStep: current,
                        clusterId: values["clusterId"],
                        namespace: values["namespace"],
                        processName: values["processName"],
                    });
                    this.carousel.next();
                })
                .catch((errorInfo) => {
                    message.warn("请填写必填项中的内容!");
                });
        } else {
            this.setState({ currentStep: current });
            this.carousel.next();
        }
    }

    prev() {
        let current = this.state.currentStep - 1;
        this.setState({ currentStep: current });
        this.carousel.prev();
        // if (current === 0) {
        //     console.log("设置值");
        //     this.clusterFormRef.current.setFieldsValue({
        //         clusterId: this.state.clusterId,
        //         namespace: this.state.namespace,
        //         processName: this.state.processName,
        //     });
        // }
    }

    loadCdProcessTemplateList() {
        getCdProcessTemplateData().then((res) => {
            if (res.code === 0) {
                this.setState({ processTemplateList: res.data });
            } else {
                message.error(res.msg);
            }
        });
    }

    deployToK8sByYAML() {
        if (this.state.deployProcessType === "1") {
            let reqParams = {
                templateId: this.state.deployProcessRecord,
                imageName: this.state.cdInfo.imageName.trim(),
            };
            postCdProcessLog(reqParams).then((res) => {
                if (res.code === 0) {
                    message.success("项目发布成功");
                } else {
                    message.error(res.msg);
                }
            });
        } else {
            let reqParams = {
                jobName: this.state.cdInfo.jobName,
                templateName: this.state.processName,
                clusterId: getClusterId(),
                namespace: this.state.namespace,
                imageName: this.state.cdInfo.imageName,
                deployYaml: this.deployYAML,
                configmapYaml: this.configMapYAML,
                serviceYaml: this.serviceYAML,
                ingressYaml: this.ingressYAML,
            };
            postCdProcessTemplate(reqParams).then((res) => {
                if (res.code === 0) {
                    message.success("项目发布成功");
                } else {
                    message.error(res.msg);
                }
            });
        }
    }

    submitToDeploy() {
        this.deployToK8sByYAML();
        message.success("已提交部署");
    }

    onDeployProcessChange(e) {
        this.setState({ deployProcessType: e });
    }

    onDeployProcessRecordChange(e) {
        this.setState({ deployProcessRecord: e });
    }

    render() {
        const layout = {
            labelCol: { span: 8 },
            wrapperCol: { span: 10 },
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
                <OpsBreadcrumbPath pathData={["CI & CD", "项目部署"]} />
                <Row>
                    <Col span={24}>
                        <Card title="项目部署" size="small">
                            <Descriptions bordered size="small" column={1}>
                                <Descriptions.Item label="部署项目">
                                    {this.state.cdInfo.jobName}
                                </Descriptions.Item>
                                <Descriptions.Item label="构建编号">
                                    {this.state.cdInfo.buildNumber}
                                </Descriptions.Item>
                                <Descriptions.Item label="部署镜像">
                                    <Paragraph copyable>
                                        {this.state.cdInfo.imageName}
                                    </Paragraph>
                                </Descriptions.Item>
                            </Descriptions>
                            <Tabs
                                defaultActiveKey="1"
                                onChange={this.onDeployProcessChange.bind(this)}
                            >
                                <TabPane tab="已有部署" key="1">
                                    <Divider orientation="left" plain>
                                        选择已有部署流程
                                    </Divider>
                                    <Form {...layout} ref={this.props.formRef}>
                                        <Form.Item
                                            label="选择流程"
                                            name="processRecord"
                                            rules={[
                                                {
                                                    required: true,
                                                    message: "请选择部署流程!",
                                                },
                                            ]}
                                        >
                                            <Select
                                                style={{ width: "100%" }}
                                                onChange={this.onDeployProcessRecordChange.bind(
                                                    this,
                                                )}
                                            >
                                                {this.state.processTemplateList.map(
                                                    (item) => {
                                                        return (
                                                            <Option
                                                                key={
                                                                    item.TemplateName
                                                                }
                                                                value={item.ID}
                                                            >
                                                                {
                                                                    item.TemplateName
                                                                }
                                                            </Option>
                                                        );
                                                    },
                                                )}
                                            </Select>
                                        </Form.Item>
                                    </Form>
                                    {this.state.ableToCd === true ? (
                                        <Fragment>
                                            {this.state.deployProcessType ===
                                                "1" && (
                                                <Button
                                                    type="primary"
                                                    size="small"
                                                    onClick={this.submitToDeploy.bind(
                                                        this,
                                                    )}
                                                >
                                                    确认发布
                                                </Button>
                                            )}
                                        </Fragment>
                                    ) : (
                                        <Result
                                            status="warning"
                                            title={
                                                <Text
                                                    type="danger"
                                                    style={{ fontSize: "14px" }}
                                                >
                                                    该次构建未生成镜像制品，不可发布!
                                                </Text>
                                            }
                                        />
                                    )}
                                </TabPane>
                                <TabPane tab="新建部署" key="2">
                                    <Divider orientation="left" plain>
                                        新建部署配置
                                    </Divider>
                                    {this.state.ableToCd === true ? (
                                        <Fragment>
                                            <Steps
                                                current={this.state.currentStep}
                                                size="small"
                                            >
                                                <Step
                                                    key="K8S集群选择"
                                                    title="K8S集群选择"
                                                />
                                                <Step
                                                    key="K8S资源配置"
                                                    title="K8S资源配置"
                                                />
                                                <Step key="发布" title="发布" />
                                            </Steps>
                                            <div
                                                className="steps-content"
                                                style={{ paddingTop: "20px" }}
                                            >
                                                <Carousel
                                                    ref={(carousel) => {
                                                        this.carousel = carousel;
                                                    }}
                                                    dots={false}
                                                >
                                                    <div>
                                                        <ClusterConfigContent
                                                            formRef={
                                                                this
                                                                    .clusterFormRef
                                                            }
                                                        />
                                                    </div>
                                                    <div>
                                                        <K8sResConfigContent
                                                            yamlCallback={this.yamlCallback.bind(
                                                                this,
                                                            )}
                                                            cdInfo={
                                                                this.state
                                                                    .cdInfo
                                                            }
                                                        />
                                                    </div>
                                                    <div>
                                                        <Collapse
                                                            ghost
                                                            style={{
                                                                textAlign:
                                                                    "left",
                                                            }}
                                                        >
                                                            <Panel
                                                                header="Deployment YAML确认"
                                                                key="1"
                                                            >
                                                                <pre
                                                                    style={{
                                                                        textAlign:
                                                                            "left",
                                                                        fontSize: 10,
                                                                    }}
                                                                >
                                                                    {
                                                                        this
                                                                            .deployYAML
                                                                    }
                                                                </pre>
                                                            </Panel>
                                                            <Panel
                                                                header="Service YAML确认"
                                                                key="2"
                                                            >
                                                                <pre
                                                                    style={{
                                                                        textAlign:
                                                                            "left",
                                                                        fontSize: 10,
                                                                    }}
                                                                >
                                                                    {
                                                                        this
                                                                            .serviceYAML
                                                                    }
                                                                </pre>
                                                            </Panel>
                                                            <Panel
                                                                header="ConfigMap YAML确认"
                                                                key="3"
                                                            >
                                                                <pre
                                                                    style={{
                                                                        textAlign:
                                                                            "left",
                                                                        fontSize: 10,
                                                                    }}
                                                                >
                                                                    {
                                                                        this
                                                                            .configMapYAML
                                                                    }
                                                                </pre>
                                                            </Panel>
                                                            <Panel
                                                                header="Ingress YAML确认"
                                                                key="4"
                                                            >
                                                                <pre
                                                                    style={{
                                                                        textAlign:
                                                                            "left",
                                                                        fontSize: 10,
                                                                    }}
                                                                >
                                                                    {
                                                                        this
                                                                            .ingressYAML
                                                                    }
                                                                </pre>
                                                            </Panel>
                                                        </Collapse>
                                                    </div>
                                                </Carousel>
                                            </div>
                                            <div style={{ marginTop: 20 }}>
                                                {this.state.currentStep < 2 && (
                                                    <Button
                                                        size="small"
                                                        type="primary"
                                                        onClick={this.next.bind(
                                                            this,
                                                        )}
                                                    >
                                                        下一步
                                                    </Button>
                                                )}
                                                {this.state.currentStep ===
                                                    2 && (
                                                    <Button
                                                        type="primary"
                                                        size="small"
                                                        onClick={this.submitToDeploy.bind(
                                                            this,
                                                        )}
                                                    >
                                                        确认发布
                                                    </Button>
                                                )}
                                                {this.state.currentStep > 0 && (
                                                    <Button
                                                        size="small"
                                                        style={{
                                                            margin: "0 8px",
                                                        }}
                                                        onClick={this.prev.bind(
                                                            this,
                                                        )}
                                                    >
                                                        上一步
                                                    </Button>
                                                )}
                                            </div>
                                        </Fragment>
                                    ) : (
                                        <Result
                                            status="warning"
                                            title={
                                                <Text
                                                    type="danger"
                                                    style={{ fontSize: "14px" }}
                                                >
                                                    该次构建未生成镜像制品，不可发布!
                                                </Text>
                                            }
                                        />
                                    )}
                                </TabPane>
                            </Tabs>
                        </Card>
                    </Col>
                </Row>
            </Content>
        );
    }
}

export default withRouter(CdContent);
