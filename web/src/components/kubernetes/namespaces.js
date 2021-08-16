import React, { Component } from "react";
import {
    Layout,
    Row,
    message,
    Col,
    Select,
    Tabs,
    Card,
    Table,
    Tag,
    Button,
    Modal,
    Alert,
    Spin,
    Icon,
    Tooltip,
    Drawer,
    Menu,
    Dropdown,
    InputNumber,
    Collapse,
    Divider,
    Input,
} from "antd";
import OpsBreadcrumbPath from "../breadcrumb_path";
import { withRouter } from 'react-router-dom';
import {
    deleteConfigMap,
    deleteResource,
    deleteSecret,
    getConfigDict,
    getDeployments,
    getNamespaces,
    getNodes,
    getPods,
    getReplicaSets,
    getReplicationControllers,
    getResourceYaml,
    getSecretDict,
    getServices,
    postApplyYaml,
    putResourceScale,
    getIngresses,
    putApplyYaml,
    restartDeployments,
} from "../../api/kubernetes";
import moment from "moment";
import { Controlled as CodeMirror } from "react-codemirror2";
import "codemirror/lib/codemirror.css";
import K8sTemplate from "../../utils/k8s_template";
import { Route, Switch } from "react-router-dom";
import RcDetailContent from "./rc_detail";
import ServiceDetailContent from "./service_detail";
import PodDetailContent from "./pod_detail";
import WebTerminalContent from "./web_terminal";
import ContainerLogContent from "./container_log";
import NodeDetailContent from "./node_detail";
import RsDetailContent from "./rs_detail";
import DeploymentDetailContent from "./deployment_detail";
import IngressDetailContent from "./ingress_detail";
import { StatusDoneSvg, StatusProgressSvg } from "../../assets/Icons";
import Highlighter from "react-highlight-words";

const { TabPane } = Tabs;
const { Content } = Layout;
const { Option } = Select;
const { Panel } = Collapse;
const { confirm } = Modal;

const customPanelStyle = {
    background: "#f7f7f7",
    borderRadius: 4,
    marginBottom: 24,
    border: 0,
    overflow: "hidden",
};

class K8sNamespacesContent extends Component {
    constructor(props) {
        super(props);
        this.refreshNamespaceResources = this.refreshNamespaceResources.bind(
            this,
        );
        this.selectChange = this.selectChange.bind(this);
        this.yarmCreateK8sResource = this.yarmCreateK8sResource.bind(this);
        this.handleOk = this.handleOk.bind(this);
        this.handleCancel = this.handleCancel.bind(this);
        this.selectTemplateChange = this.selectTemplateChange.bind(this);
        this.onDrawerClose = this.onDrawerClose.bind(this);
        this.handleYamlDetailCancel = this.handleYamlDetailCancel.bind(this);
        this.handleScaleCancel = this.handleScaleCancel.bind(this);
        this.handleScaleCommit = this.handleScaleCommit.bind(this);
        this.onInputNumberChange = this.onInputNumberChange.bind(this);
        this.handleConfigCancel = this.handleConfigCancel.bind(this);
        this.updateYamlContent = this.updateYamlContent.bind(this);
        this.state = {
            nodeColumns: [
                {
                    title: "名称",
                    key: "name",
                    ...this.getColumnSearchProps("name"),
                    render: (text, record) => {
                        return (
                            <Button
                                type="link"
                                size="small"
                                onClick={this.showDetail.bind(
                                    this,
                                    "node",
                                    record,
                                )}
                            >
                                {record.metadata.name}
                            </Button>
                        );
                    },
                },
                {
                    title: "标签",
                    key: "labels",
                    render: (text, record) => {
                        let labels = [];
                        for (let key in record.metadata.labels) {
                            labels.push(
                                <div key={key}>
                                    <Tag color="geekblue">
                                        {key +
                                            ":" +
                                            record.metadata.labels[key]}
                                    </Tag>
                                </div>,
                            );
                        }
                        return labels;
                    },
                },
                {
                    title: "已就绪",
                    key: "ready",
                    render: (text, record) => {
                        let data = record.status.conditions;
                        let content = <Tag color="#f50">False</Tag>;
                        for (let i = 0; i < data.length; i++) {
                            if (data[i]["type"] === "Ready") {
                                if (data[i]["status"] !== "False") {
                                    content = <Tag color="#87d068">True</Tag>;
                                }
                                break;
                            }
                        }
                        return content;
                    },
                },
                {
                    title: "磁盘压力",
                    key: "disk",
                    render: (text, record) => {
                        let data = record.status.conditions;
                        let content = <Tag color="#87d068">无压力</Tag>;
                        for (let i = 0; i < data.length; i++) {
                            if (data[i]["type"] === "DiskPressure") {
                                if (data[i]["status"] !== "False") {
                                    content = <Tag color="#f50">有压力</Tag>;
                                }
                                break;
                            }
                        }
                        return content;
                    },
                },
                {
                    title: "内存压力",
                    key: "memory",
                    render: (text, record) => {
                        let data = record.status.conditions;
                        let content = <Tag color="#87d068">无压力</Tag>;
                        for (let i = 0; i < data.length; i++) {
                            if (data[i]["type"] === "MemoryPressure") {
                                if (data[i]["status"] !== "False") {
                                    content = <Tag color="#f50">有压力</Tag>;
                                }
                                break;
                            }
                        }
                        return content;
                    },
                },
                {
                    title: "PID压力",
                    key: "pid",
                    render: (text, record) => {
                        let data = record.status.conditions;
                        let content = <Tag color="#87d068">无压力</Tag>;
                        for (let i = 0; i < data.length; i++) {
                            if (data[i]["type"] === "PIDPressure") {
                                if (data[i]["status"] !== "False") {
                                    content = <Tag color="#f50">有压力</Tag>;
                                }
                                break;
                            }
                        }
                        return content;
                    },
                },
                {
                    title: "创建时间",
                    key: "create_time",
                    render: (text, record) =>
                        moment(record.metadata.creationTimestamp).format(
                            "YYYY-MM-DD",
                        ),
                },
            ],
            rcColumns: [
                {
                    title: "名称",
                    key: "name",
                    ...this.getColumnSearchProps("name"),
                    render: (text, record) => {
                        let icon = <StatusDoneSvg />;
                        if (
                            record.status.readyReplicas === undefined ||
                            record.status.readyReplicas < record.status.replicas
                        ) {
                            icon = <StatusProgressSvg />;
                        }
                        return (
                            <div>
                                <span
                                    style={{ position: "relative", top: "5px" }}
                                >
                                    {icon}
                                </span>
                                <Button
                                    type="link"
                                    size="small"
                                    onClick={this.showDetail.bind(
                                        this,
                                        "rc",
                                        record,
                                    )}
                                >
                                    {record.metadata.name}
                                </Button>
                            </div>
                        );
                    },
                },
                {
                    title: "命名空间",
                    key: "namespace",
                    render: (text, record) => {
                        return <span>{record.metadata.namespace}</span>;
                    },
                },
                {
                    title: "标签",
                    key: "labels",
                    render: (text, record) => {
                        let labels = [];
                        for (let key in record.metadata.labels) {
                            labels.push(
                                <div key={key}>
                                    <Tag color="geekblue">
                                        {key +
                                            ":" +
                                            record.metadata.labels[key]}
                                    </Tag>
                                </div>,
                            );
                        }
                        return labels;
                    },
                },
                {
                    title: "容器组",
                    key: "replicas",
                    render: (text, record) => {
                        let readyReplicas = " - ";
                        if (record.status.hasOwnProperty("readyReplicas")) {
                            readyReplicas = record.status.readyReplicas;
                        }
                        return readyReplicas + " / " + record.status.replicas;
                    },
                },
                {
                    title: "创建时间",
                    key: "create_time",
                    render: (text, record) =>
                        moment(record.metadata.creationTimestamp).format(
                            "YYYY-MM-DD",
                        ),
                },
                {
                    title: "镜像",
                    key: "image",
                    render: (text, record) => {
                        let data = record.spec.template.spec.containers;
                        let labels = [];
                        for (let i = 0; i < data.length; i++) {
                            labels.push(
                                <div key={i}>
                                    <Tag color="geekblue">{data[i].image}</Tag>
                                </div>,
                            );
                        }
                        return labels;
                    },
                },
                {
                    title: "操作",
                    key: "action",
                    fixed: "right",
                    render: (text, record) => {
                        const menu = (
                            <Menu>
                                <Menu.Item key="0">
                                    <Button
                                        type="link"
                                        size="small"
                                        onClick={this.autoScaleHandler.bind(
                                            this,
                                            "rc",
                                            record,
                                        )}
                                        disabled={!this.props.aclAuthMap["PUT:/kubernetes/scale"]}
                                    >
                                        {" "}
                                        伸缩
                                    </Button>
                                </Menu.Item>
                                <Menu.Item key="1">
                                    <Button
                                        type="link"
                                        onClick={this.deleteHandler.bind(
                                            this,
                                            "rc",
                                            record,
                                        )}
                                        disabled={!this.props.aclAuthMap["DELETE:/kubernetes/resource"]}
                                    >
                                        删除
                                    </Button>
                                </Menu.Item>
                                <Menu.Item key="2">
                                    <Button
                                        type="link"
                                        size="small"
                                        onClick={this.showYamlHandler.bind(
                                            this,
                                            "rc",
                                            record,
                                        )}
                                        disabled={!this.props.aclAuthMap["GET:/kubernetes/yaml"]}
                                    >
                                        查看yaml文件
                                    </Button>
                                </Menu.Item>
                            </Menu>
                        );
                        return (
                            <Dropdown overlay={menu} trigger={["click"]}>
                                <Button
                                    type="link"
                                    size="small"
                                    className="ant-dropdown-link"
                                >
                                    资源操作 <Icon type="down" />
                                </Button>
                            </Dropdown>
                        );
                    },
                },
            ],
            rsColumns: [
                {
                    title: "名称",
                    key: "name",
                    ...this.getColumnSearchProps("name"),
                    render: (text, record) => {
                        let icon = <StatusDoneSvg />;
                        if (
                            record.status.readyReplicas === undefined ||
                            record.status.readyReplicas < record.status.replicas
                        ) {
                            icon = <StatusProgressSvg />;
                        }
                        return (
                            <div>
                                <span
                                    style={{ position: "relative", top: "5px" }}
                                >
                                    {icon}
                                </span>
                                <Button
                                    type="link"
                                    size="small"
                                    onClick={this.showDetail.bind(
                                        this,
                                        "rs",
                                        record,
                                    )}
                                >
                                    {record.metadata.name}
                                </Button>
                            </div>
                        );
                    },
                },
                {
                    title: "命名空间",
                    key: "namespace",
                    render: (text, record) => {
                        return <span>{record.metadata.namespace}</span>;
                    },
                },
                {
                    title: "标签",
                    key: "labels",
                    render: (text, record) => {
                        let labels = [];
                        for (let key in record.metadata.labels) {
                            labels.push(
                                <div key={key}>
                                    <Tag color="geekblue">
                                        {key +
                                            ":" +
                                            record.metadata.labels[key]}
                                    </Tag>
                                </div>,
                            );
                        }
                        return labels;
                    },
                },
                {
                    title: "容器组",
                    key: "replicas",
                    render: (text, record) => {
                        let readyReplicas = " - ";
                        if (record.status.hasOwnProperty("readyReplicas")) {
                            readyReplicas = record.status.readyReplicas;
                        }
                        return readyReplicas + " / " + record.status.replicas;
                    },
                },
                {
                    title: "创建时间",
                    key: "create_time",
                    render: (text, record) =>
                        moment(record.metadata.creationTimestamp).format(
                            "YYYY-MM-DD",
                        ),
                },
                {
                    title: "镜像",
                    key: "image",
                    render: (text, record) => {
                        let data = record.spec.template.spec.containers;
                        let labels = [];
                        for (let i = 0; i < data.length; i++) {
                            labels.push(
                                <div key={i}>
                                    <Tag color="geekblue">{data[i].image}</Tag>
                                </div>,
                            );
                        }
                        return labels;
                    },
                },
                {
                    title: "操作",
                    key: "action",
                    fixed: "right",
                    render: (text, record) => {
                        const menu = (
                            <Menu>
                                <Menu.Item key="0">
                                    <Button
                                        type="link"
                                        size="small"
                                        onClick={this.autoScaleHandler.bind(
                                            this,
                                            "rs",
                                            record,
                                        )}
                                        disabled={!this.props.aclAuthMap["PUT:/kubernetes/scale"]}
                                    >
                                        伸缩
                                    </Button>
                                </Menu.Item>
                                <Menu.Item key="1">
                                    <Button
                                        type="link"
                                        size="small"
                                        onClick={this.deleteHandler.bind(
                                            this,
                                            "rs",
                                            record,
                                        )}
                                        disabled={!this.props.aclAuthMap["DELETE:/kubernetes/resource"]}
                                    >
                                        删除
                                    </Button>
                                </Menu.Item>
                                <Menu.Item key="2">
                                    <Button
                                        type="link"
                                        size="small"
                                        onClick={this.showYamlHandler.bind(
                                            this,
                                            "rs",
                                            record,
                                        )}
                                        disabled={!this.props.aclAuthMap["GET:/kubernetes/yaml"]}
                                    >
                                        查看yaml文件
                                    </Button>
                                </Menu.Item>
                            </Menu>
                        );
                        return (
                            <Dropdown overlay={menu} trigger={["click"]}>
                                <Button
                                    type="link"
                                    className="ant-dropdown-link"
                                >
                                    资源操作 <Icon type="down" />
                                </Button>
                            </Dropdown>
                        );
                    },
                },
            ],
            serviceColumns: [
                {
                    title: "名称",
                    key: "name",
                    ...this.getColumnSearchProps("name"),
                    render: (text, record) => {
                        return (
                            <Button
                                type="link"
                                size="small"
                                onClick={this.showDetail.bind(
                                    this,
                                    "service",
                                    record,
                                )}
                            >
                                {record.metadata.name}
                            </Button>
                        );
                    },
                },
                {
                    title: "命名空间",
                    key: "namespace",
                    render: (text, record) => {
                        return <span>{record.metadata.namespace}</span>;
                    },
                },
                {
                    title: "标签",
                    key: "labels",
                    render: (text, record) => {
                        let labels = [];
                        for (let key in record.metadata.labels) {
                            labels.push(
                                <div key={key}>
                                    <Tag color="geekblue">
                                        {key +
                                            ":" +
                                            record.metadata.labels[key]}
                                    </Tag>
                                </div>,
                            );
                        }
                        return labels;
                    },
                },
                {
                    title: "集群IP",
                    key: "replicas",
                    render: (text, record) => {
                        return record.spec.clusterIP;
                    },
                },
                {
                    title: "内部端点",
                    key: "image",
                    render: (text, record) => {
                        let data = record.spec.ports;
                        let labels = [];
                        for (let i = 0; i < data.length; i++) {
                            labels.push(
                                <div key={i}>
                                    <Tag color="geekblue">
                                        {record.metadata.name +
                                            "." +
                                            record.metadata.namespace +
                                            ":" +
                                            data[i]["port"] +
                                            " " +
                                            data[i]["protocol"]}
                                    </Tag>
                                </div>,
                            );
                        }
                        return labels;
                    },
                },
                {
                    title: "创建时间",
                    key: "create_time",
                    render: (text, record) =>
                        moment(record.metadata.creationTimestamp).format(
                            "YYYY-MM-DD",
                        ),
                },
                {
                    title: "操作",
                    key: "action",
                    fixed: "right",
                    render: (text, record) => {
                        const menu = (
                            <Menu>
                                <Menu.Item key="0">
                                    <Button
                                        type="link"
                                        size="small"
                                        onClick={this.deleteHandler.bind(
                                            this,
                                            "service",
                                            record,
                                        )}
                                        disabled={!this.props.aclAuthMap["DELETE:/kubernetes/resource"]}
                                    >
                                        删除
                                    </Button>
                                </Menu.Item>
                                <Menu.Item key="1">
                                    <Button
                                        type="link"
                                        size="small"
                                        onClick={this.showYamlHandler.bind(
                                            this,
                                            "service",
                                            record,
                                        )}
                                        disabled={!this.props.aclAuthMap["GET:/kubernetes/yaml"]}
                                    >
                                        查看yaml文件
                                    </Button>
                                </Menu.Item>
                            </Menu>
                        );
                        return (
                            <Dropdown overlay={menu} trigger={["click"]}>
                                <Button
                                    size="small"
                                    className="ant-dropdown-link"
                                    type="link"
                                >
                                    资源操作 <Icon type="down" />
                                </Button>
                            </Dropdown>
                        );
                    },
                },
            ],
            ingressColumns: [
                {
                    title: "名称",
                    key: "name",
                    ...this.getColumnSearchProps("name"),
                    render: (text, record) => {
                        return (
                            <Button
                                type="link"
                                size="small"
                                onClick={this.showDetail.bind(
                                    this,
                                    "ingress",
                                    record,
                                )}
                            >
                                {record.metadata.name}
                            </Button>
                        );
                    },
                },
                {
                    title: "命名空间",
                    key: "namespace",
                    render: (text, record) => {
                        return <span>{record.metadata.namespace}</span>;
                    },
                },
                {
                    title: "域名",
                    key: "domain",
                    render: (text, record) => {
                        let content = record.spec.rules.map((item, index) => {
                            return (
                                <span>
                                    {item.host}
                                    <br />
                                </span>
                            );
                        });
                        return <span>{content}</span>;
                    },
                },
                {
                    title: "证书",
                    key: "cert",
                    render: (text, record) => {
                        if (record.spec.tls === undefined) {
                            return <span>-</span>;
                        }
                        let content = record.spec.tls.map((item, index) => {
                            return (
                                <span>
                                    {item.secretName}
                                    <br />
                                </span>
                            );
                        });
                        return <span>{content}</span>;
                    },
                },
                {
                    title: "创建时间",
                    key: "create_time",
                    render: (text, record) =>
                        moment(record.metadata.creationTimestamp).format(
                            "YYYY-MM-DD",
                        ),
                },
                {
                    title: "操作",
                    key: "action",
                    fixed: "right",
                    render: (text, record) => {
                        const menu = (
                            <Menu>
                                <Menu.Item key="1">
                                    <Button
                                        type="link"
                                        size="small"
                                        onClick={this.deleteHandler.bind(
                                            this,
                                            "ingress",
                                            record,
                                        )}
                                        disabled={!this.props.aclAuthMap["DELETE:/kubernetes/resource"]}
                                    >
                                        删除
                                    </Button>
                                </Menu.Item>
                                <Menu.Item key="2">
                                    <Button
                                        type="link"
                                        size="small"
                                        onClick={this.showYamlHandler.bind(
                                            this,
                                            "ingress",
                                            record,
                                        )}
                                        disabled={!this.props.aclAuthMap["GET:/kubernetes/yaml"]}
                                    >
                                        查看yaml文件
                                    </Button>
                                </Menu.Item>
                            </Menu>
                        );
                        return (
                            <Dropdown overlay={menu} trigger={["click"]}>
                                <Button
                                    size="small"
                                    className="ant-dropdown-link"
                                    type="link"
                                >
                                    资源操作 <Icon type="down" />
                                </Button>
                            </Dropdown>
                        );
                    },
                },
            ],
            deploymentColumns: [
                {
                    title: "名称",
                    key: "name",
                    ...this.getColumnSearchProps("name"),
                    render: (text, record) => {
                        let icon = <StatusDoneSvg />;
                        if (
                            record.status.readyReplicas === undefined ||
                            record.status.readyReplicas < record.status.replicas
                        ) {
                            icon = <StatusProgressSvg />;
                        }
                        return (
                            <div>
                                <span
                                    style={{ position: "relative", top: "5px" }}
                                >
                                    {icon}
                                </span>
                                <Button
                                    type="link"
                                    size="small"
                                    onClick={this.showDetail.bind(
                                        this,
                                        "deployment",
                                        record,
                                    )}
                                >
                                    {record.metadata.name}
                                </Button>
                            </div>
                        );
                    },
                },
                {
                    title: "命名空间",
                    key: "namespace",
                    render: (text, record) => {
                        return <span>{record.metadata.namespace}</span>;
                    },
                },
                {
                    title: "标签",
                    key: "labels",
                    render: (text, record) => {
                        let labels = [];
                        for (let key in record.metadata.labels) {
                            labels.push(
                                <div key={key}>
                                    <Tag color="geekblue">
                                        {key +
                                            ":" +
                                            record.metadata.labels[key]}
                                    </Tag>
                                </div>,
                            );
                        }
                        return labels;
                    },
                },
                {
                    title: "容器组",
                    key: "replicas",
                    render: (text, record) => {
                        let readyReplicas = " - ";
                        if (record.status.hasOwnProperty("readyReplicas")) {
                            readyReplicas = record.status.readyReplicas;
                        }
                        return readyReplicas + " / " + record.status.replicas;
                    },
                },
                {
                    title: "创建时间",
                    key: "create_time",
                    render: (text, record) =>
                        moment(record.metadata.creationTimestamp).format(
                            "YYYY-MM-DD",
                        ),
                },
                {
                    title: "镜像",
                    key: "image",
                    render: (text, record) => {
                        let data = record.spec.template.spec.containers;
                        let labels = [];
                        for (let i = 0; i < data.length; i++) {
                            labels.push(
                                <div key={i}>
                                    <Tag color="geekblue">{data[i].image}</Tag>
                                </div>,
                            );
                        }
                        return labels;
                    },
                },
                {
                    title: "操作",
                    key: "action",
                    fixed: "right",
                    render: (text, record) => {
                        const menu = (
                            <Menu>
                                <Menu.Item key="0">
                                    <Button
                                        type="link"
                                        size="small"
                                        onClick={this.autoScaleHandler.bind(
                                            this,
                                            "deployment",
                                            record,
                                        )}
                                        disabled={!this.props.aclAuthMap["PUT:/kubernetes/scale"]}
                                    >
                                        伸缩
                                    </Button>
                                </Menu.Item>
                                <Menu.Item key="1">
                                    <Button
                                        type="link"
                                        size="small"
                                        onClick={this.restartDeployHandler.bind(
                                            this,
                                            record,
                                        )}
                                        disabled={!this.props.aclAuthMap["PUT:/kubernetes/deployment/restart"]}
                                    >
                                        重启
                                    </Button>
                                </Menu.Item>
                                <Menu.Item key="2">
                                    <Button
                                        type="link"
                                        size="small"
                                        onClick={this.deleteHandler.bind(
                                            this,
                                            "deployment",
                                            record,
                                        )}
                                        disabled={!this.props.aclAuthMap["DELETE:/kubernetes/resource"]}
                                    >
                                        删除
                                    </Button>
                                </Menu.Item>
                                <Menu.Item key="3">
                                    <Button
                                        type="link"
                                        size="small"
                                        onClick={this.showYamlHandler.bind(
                                            this,
                                            "deployment",
                                            record,
                                        )}
                                        disabled={!this.props.aclAuthMap["GET:/kubernetes/yaml"]}
                                    >
                                        查看yaml文件
                                    </Button>
                                </Menu.Item>
                            </Menu>
                        );
                        return (
                            <Dropdown overlay={menu} trigger={["click"]}>
                                <Button
                                    size="small"
                                    className="ant-dropdown-link"
                                    type="link"
                                >
                                    资源操作 <Icon type="down" />
                                </Button>
                            </Dropdown>
                        );
                    },
                },
            ],
            podColumns: [
                {
                    title: "名称",
                    key: "name",
                    ...this.getColumnSearchProps("name"),
                    render: (text, record) => {
                        return (
                            <Button
                                type="link"
                                size="small"
                                onClick={this.showDetail.bind(
                                    this,
                                    "pod",
                                    record,
                                )}
                            >
                                {record.metadata.name}
                            </Button>
                        );
                    },
                },
                {
                    title: "命名空间",
                    key: "namespace",
                    render: (text, record) => {
                        return <span>{record.metadata.namespace}</span>;
                    },
                },
                {
                    title: "节点",
                    key: "labels",
                    render: (text, record) => {
                        return <span>{record.spec.nodeName}</span>;
                    },
                },
                {
                    title: "状态",
                    key: "status",
                    dataIndex: "",
                    render: (text, record) => {
                        let value = record.status.phase;
                        let statusContent = <Tag color="#87d068">{value}</Tag>;
                        if (value === "Failed") {
                            statusContent = <Tag color="#f50">{value}</Tag>;
                        }
                        if (value === "Pending") {
                            statusContent = <Tag color="#108ee9">{value}</Tag>;
                        }
                        return statusContent;
                    },
                },
                {
                    title: "重启次数",
                    key: "count",
                    render: (text, record) => {
                        if (record.status.phase === "Failed") {
                            return "0";
                        }
                        let data = record.status.containerStatuses;
                        let labels = [];
                        for (let i = 0; i < data.length; i++) {
                            labels.push(
                                <div key={i}>
                                    <Tag color="geekblue">
                                        {data[i]["name"] +
                                            " - " +
                                            data[i]["restartCount"]}
                                    </Tag>
                                </div>,
                            );
                        }
                        return labels;
                    },
                },
                {
                    title: "创建时间",
                    key: "create_time",
                    render: (text, record) =>
                        moment(record.metadata.creationTimestamp).format(
                            "YYYY-MM-DD",
                        ),
                },
                {
                    title: "操作",
                    key: "action",
                    fixed: "right",
                    render: (text, record) => {
                        const menu = (
                            <Menu>
                                <Menu.Item key="1">
                                    <Button
                                        type="link"
                                        size="small"
                                        onClick={this.deleteHandler.bind(
                                            this,
                                            "pod",
                                            record,
                                        )}
                                        disabled={!this.props.aclAuthMap["DELETE:/kubernetes/resource"]}
                                    >
                                        删除
                                    </Button>
                                </Menu.Item>
                                {/* <Menu.Item key="2">
                                    <Button
                                        type="link"
                                        size="small"
                                        onClick={this.showYamlHandler.bind(
                                            this,
                                            "pod",
                                            record,
                                        )}
                                    >
                                        查看yaml文件
                                    </Button>
                                </Menu.Item> */}
                            </Menu>
                        );
                        return (
                            <Dropdown overlay={menu} trigger={["click"]}>
                                <Button
                                    className="ant-dropdown-link"
                                    type="link"
                                    size="small"
                                >
                                    资源操作 <Icon type="down" />
                                </Button>
                            </Dropdown>
                        );
                    },
                },
            ],
            configDictColumns: [
                {
                    title: "名称",
                    key: "name",
                    ...this.getColumnSearchProps("name"),
                    render: (text, record) => {
                        return <span>{record.metadata.name}</span>;
                    },
                },
                {
                    title: "命名空间",
                    key: "namespace",
                    render: (text, record) => {
                        return <span>{record.metadata.namespace}</span>;
                    },
                },
                {
                    title: "创建时间",
                    key: "creationTimestamp",
                    render: (text, record) => {
                        return (
                            <span>
                                {moment(
                                    record.metadata.creationTimestamp,
                                ).format("YYYY-MM-DD")}
                            </span>
                        );
                    },
                },
                {
                    title: "操作",
                    key: "data",
                    fixed: "right",
                    render: (text, record) => {
                        const menu = (
                            <Menu>
                                <Menu.Item key="0">
                                    <Button
                                        type="link"
                                        size="small"
                                        onClick={this.showDictValue.bind(
                                            this,
                                            record,
                                        )}
                                    >
                                        查看
                                    </Button>
                                </Menu.Item>
                                <Menu.Item key="1">
                                    <Button
                                        type="link"
                                        size="small"
                                        onClick={this.deleteDictValue.bind(
                                            this,
                                            "configmap",
                                            record,
                                        )}
                                        disabled={!this.props.aclAuthMap["DELETE:/kubernetes/resource"]}
                                    >
                                        删除
                                    </Button>
                                </Menu.Item>
                                <Menu.Item key="2">
                                    <Button
                                        type="link"
                                        size="small"
                                        onClick={this.showYamlHandler.bind(
                                            this,
                                            "configmap",
                                            record,
                                        )}
                                        disabled={!this.props.aclAuthMap["GET:/kubernetes/yaml"]}
                                    >
                                        查看yaml文件
                                    </Button>
                                </Menu.Item>
                            </Menu>
                        );
                        return (
                            <Dropdown overlay={menu} trigger={["click"]}>
                                <Button
                                    size="small"
                                    className="ant-dropdown-link"
                                    type="link"
                                >
                                    资源操作 <Icon type="down" />
                                </Button>
                            </Dropdown>
                        );
                    },
                },
            ],
            secretDictColumns: [
                {
                    title: "名称",
                    key: "name",
                    ...this.getColumnSearchProps("name"),
                    render: (text, record) => {
                        return <span>{record.metadata.name}</span>;
                    },
                },
                {
                    title: "命名空间",
                    key: "namespace",
                    render: (text, record) => {
                        return <span>{record.metadata.namespace}</span>;
                    },
                },
                {
                    title: "类型",
                    dataIndex: "type",
                    key: "type",
                },
                {
                    title: "创建时间",
                    key: "creationTimestamp",
                    render: (text, record) => {
                        return (
                            <span>
                                {moment(
                                    record.metadata.creationTimestamp,
                                ).format("YYYY-MM-DD")}
                            </span>
                        );
                    },
                },
                {
                    title: "操作",
                    key: "data",
                    fixed: "right",
                    render: (text, record) => {
                        return (
                            <div>
                                <Button
                                    type="link"
                                    size="small"
                                    onClick={this.showDictValue.bind(
                                        this,
                                        record,
                                    )}
                                >
                                    查看
                                </Button>
                                <Divider type="vertical" />
                                <Button
                                    type="link"
                                    size="small"
                                    onClick={this.deleteDictValue.bind(
                                        this,
                                        "secret",
                                        record,
                                    )}
                                    disabled={!this.props.aclAuthMap["DELETE:/kubernetes/resource"]}
                                >
                                    删除
                                </Button>
                            </div>
                        );
                    },
                },
            ],
            namespaceList: [],
            replicationControllerList: [],
            deploymentList: [],
            replicaSetList: [],
            serviceList: [],
            ingressList: [],
            podList: [],
            nodeList: [],
            configDictList: [],
            secretDictList: [],
            dictConfigValue: {},
            refreshDataLoading: false,
            detailDrawerVisible: false,
            yamlModalVisible: false,
            scaleModalVisible: false,
            configModalVisible: false,
            yamlDetail: "",
            currentNamespace: "default",
            yamlCode: "",
            yamlType: "",
            replicaCount: 0,
            scaleResType: "",
            scaleResData: {
                spec: {
                    replicas: 0,
                },
            },
        };
    }

    componentDidMount() {
        this.getK8sNamespaces();
        this.refreshNamespaceResources();
        this.getK8sNodes();
    }

    selectChange(value) {
        this.setState({ currentNamespace: value }, () => {
            this.refreshNamespaceResources();
        });
    }

    getK8sNamespaces() {
        getNamespaces()
            .then((res) => {
                if (res.code === 0) {
                    this.setState({ namespaceList: res.data.items });
                } else {
                    message.error(res.msg);
                }
            })
            .catch((err) => {
                message.error(err.toLocaleString());
            });
    }

    getK8sNodes() {
        getNodes()
            .then((res) => {
                if (res.code === 0) {
                    this.setState({
                        nodeList: res.data.items,
                    });
                } else {
                    message.error(res.msg);
                }
            })
            .catch((err) => {
                message.error(err.toLocaleString());
            });
    }

    refreshNamespaceResources() {
        this.getNamespaceDeployments();
        this.getNamespaceReplicationControllers();
        this.getNamespaceReplicaSets();
        this.getNamespaceServices();
        this.getNamespaceIngress();
        this.getNamespacePods();
        this.getNamespaceConfigDict();
        this.getNamespaceSecretDict();
        this.showLoading();
    }

    showLoading = () => {
        this.setState({ refreshDataLoading: true });
        setTimeout(() => {
            this.setState({ refreshDataLoading: false });
        }, 800);
    };

    getNamespaceDeployments() {
        getDeployments({
            namespace: this.state.currentNamespace,
        })
            .then((res) => {
                if (res.code === 0) {
                    this.setState({
                        deploymentList: res.data.items,
                    });
                } else {
                    message.error(res.msg);
                }
            })
            .catch((err) => {
                message.error(err.toLocaleString());
            });
    }

    getNamespaceReplicationControllers() {
        getReplicationControllers({
            namespace: this.state.currentNamespace,
        })
            .then((res) => {
                if (res.code === 0) {
                    this.setState({
                        replicationControllerList: res.data.items,
                    });
                } else {
                    message.error(res.msg);
                }
            })
            .catch((err) => {
                message.error(err.toLocaleString());
            });
    }

    getNamespaceReplicaSets() {
        getReplicaSets({
            namespace: this.state.currentNamespace,
        })
            .then((res) => {
                if (res.code === 0) {
                    this.setState({
                        replicaSetList: res.data.items,
                    });
                } else {
                    message.error(res.msg);
                }
            })
            .catch((err) => {
                message.error(err.toLocaleString());
            });
    }

    getNamespaceServices() {
        getServices({
            namespace: this.state.currentNamespace,
        })
            .then((res) => {
                if (res.code === 0) {
                    this.setState({
                        serviceList: res.data.items,
                    });
                } else {
                    message.error(res.msg);
                }
            })
            .catch((err) => {
                message.error(err.toLocaleString());
            });
    }

    getNamespacePods() {
        getPods({
            namespace: this.state.currentNamespace,
        })
            .then((res) => {
                if (res.code === 0) {
                    this.setState({
                        podList: res.data.items,
                    });
                } else {
                    message.error(res.msg);
                }
            })
            .catch((err) => {
                message.error(err.toLocaleString());
            });
    }

    getNamespaceIngress() {
        getIngresses({ namespace: this.state.currentNamespace })
            .then((res) => {
                if (res.code === 0) {
                    this.setState({
                        ingressList: res.data.items,
                    });
                } else {
                    message.error(res.msg);
                }
            })
            .catch((err) => {
                message.error(err.toLocaleString());
            });
    }

    getNamespaceConfigDict() {
        getConfigDict({
            namespace: this.state.currentNamespace,
        })
            .then((res) => {
                if (res.code === 0) {
                    this.setState({
                        configDictList: res.data.items,
                    });
                } else {
                    message.error(res.msg);
                }
            })
            .catch((err) => {
                message.error(err.toLocaleString());
            });
    }

    getNamespaceSecretDict() {
        getSecretDict({
            namespace: this.state.currentNamespace,
        })
            .then((res) => {
                if (res.code === 0) {
                    this.setState({
                        secretDictList: res.data.items,
                    });
                } else {
                    message.error(res.msg);
                }
            })
            .catch((err) => {
                message.error(err.toLocaleString());
            });
    }

    getColumnSearchProps = (dataIndex) => ({
        filterDropdown: ({
            setSelectedKeys,
            selectedKeys,
            confirm,
            clearFilters,
        }) => (
            <div style={{ padding: 8 }}>
                <Input
                    ref={(node) => {
                        this.searchInput = node;
                    }}
                    value={selectedKeys[0]}
                    onChange={(e) =>
                        setSelectedKeys(e.target.value ? [e.target.value] : [])
                    }
                    onPressEnter={() =>
                        this.handleSearch(selectedKeys, confirm)
                    }
                    style={{ width: 188, marginBottom: 8, display: "block" }}
                />
                <Button
                    type="primary"
                    onClick={() => this.handleSearch(selectedKeys, confirm)}
                    icon="search"
                    size="small"
                    style={{ width: 90, marginRight: 8 }}
                >
                    查找
                </Button>
                <Button
                    onClick={() => this.handleReset(clearFilters)}
                    size="small"
                    style={{ width: 90 }}
                >
                    重置
                </Button>
            </div>
        ),
        filterIcon: (filtered) => (
            <Icon
                type="search"
                style={{ color: filtered ? "#1890ff" : "#40a9ff" }}
            />
        ),
        onFilter: (value, record) => {
            return record["metadata"][dataIndex]
                .toString()
                .toLowerCase()
                .includes(value.toLowerCase());
        },
        onFilterDropdownVisibleChange: (visible) => {
            if (visible) {
                setTimeout(() => this.searchInput.select());
            }
        },
        render: (text) => (
            <Highlighter
                highlightStyle={{ backgroundColor: "#ffc069", padding: 0 }}
                searchWords={[this.state.searchText]}
                autoEscape
                textToHighlight={text.toString()}
            />
        ),
    });

    handleSearch = (selectedKeys, confirm) => {
        confirm();
        this.setState({ searchText: selectedKeys[0] });
    };

    handleReset = (clearFilters) => {
        clearFilters();
        this.setState({ searchText: "" });
    };

    showDictValue(data) {
        this.setState({
            configModalVisible: true,
            dictConfigValue: data.data,
        });
    }

    deleteDictValue(type, data) {
        switch (type) {
            case "configmap":
                deleteConfigMap({
                    namespace: this.state.currentNamespace,
                    resName: data.metadata.name,
                })
                    .then((res) => {
                        if (res.code === 0) {
                            message.success("删除成功");
                            this.refreshNamespaceResources();
                        } else {
                            message.error(res.msg);
                        }
                    })
                    .catch((err) => {
                        message.error(err.toLocaleString());
                    });
                break;
            case "secret":
                deleteSecret({
                    namespace: this.state.currentNamespace,
                    resName: data.metadata.name,
                })
                    .then((res) => {
                        if (res.code === 0) {
                            message.success("删除成功");
                            this.refreshNamespaceResources();
                        } else {
                            message.error(res.msg);
                        }
                    })
                    .catch((err) => {
                        message.error(err.toLocaleString());
                    });
                break;
            default:
                message.warn("不支持的类型");
                break;
        }
    }

    handleConfigCancel() {
        this.setState({
            configModalVisible: false,
        });
    }

    yarmCreateK8sResource() {
        this.setState({ createModalVisible: true });
    }

    handleOk() {
        if (this.state.yamlCode.trim() === "") {
            message.warn("需要填写文件内容！");
            return;
        }
        postApplyYaml({
            namespace: this.state.currentNamespace,
            yamlContent: this.state.yamlCode,
        })
            .then((res) => {
                if (res.code === 0) {
                    this.setState({ createModalVisible: false });
                    message.success("提交成功，请到对应资源栏查看详情");
                    this.refreshNamespaceResources();
                } else {
                    message.error(res.msg);
                }
            })
            .catch((err) => {
                message.error(err.toLocaleString());
            });
    }

    handleCancel() {
        this.setState({ createModalVisible: false });
    }

    selectTemplateChange(data) {
        this.setState({ yamlCode: K8sTemplate[data] });
    }

    showDetail(resourceType, data) {
        this.setState({ detailDrawerVisible: true });
        switch (resourceType) {
            case "node":
                this.props.history.push({
                    pathname: "/admin/k8s_cluster/manage/node_detail",
                    state: { data, aclAuthMap: this.props.aclAuthMap },
                });
                break;
            case "rc":
                this.props.history.push({
                    pathname: "/admin/k8s_cluster/manage/rc_detail",
                    state: { data },
                });
                break;
            case "rs":
                this.props.history.push({
                    pathname: "/admin/k8s_cluster/manage/rs_detail",
                    state: { data },
                });
                break;
            case "deployment":
                this.props.history.push({
                    pathname: "/admin/k8s_cluster/manage/deployment_detail",
                    state: { data },
                });
                break;
            case "service":
                this.props.history.push({
                    pathname: "/admin/k8s_cluster/manage/service_detail",
                    state: { data },
                });
                break;
            case "ingress":
                this.props.history.push({
                    pathname: "/admin/k8s_cluster/manage/ingress_detail",
                    state: { data },
                });
                break;
            case "pod":
                this.props.history.push({
                    pathname: "/admin/k8s_cluster/manage/pod_detail",
                    state: { data, aclAuthMap: this.props.aclAuthMap },
                });
                break;
            default:
                message.error("不支持的类型：" + resourceType);
                break;
        }
    }

    onDrawerClose() {
        this.setState({ detailDrawerVisible: false });
    }

    autoScaleHandler(resType, data) {
        this.setState({
            scaleModalVisible: true,
            scaleResType: resType,
            scaleResData: data,
        });
    }

    restartDeployHandler(record) {
        restartDeployments({
            namespace: record.metadata.namespace,
            resName: record.metadata.name,
        }).then((res) => {
            if (res.code === 0) {
                message.success("重启成功!");
                this.refreshNamespaceResources();
            } else {
                message.error(res.msg);
            }
        });
    }

    deleteHandler(resType, data) {
        const that = this;
        confirm({
            title: "操作警告?",
            content: "确定删除该 kubernetes 资源吗？",
            icon: <Icon type="warning" />,
            onOk() {
                deleteResource({
                    namespace: data.metadata.namespace,
                    resType: resType,
                    resName: data.metadata.name,
                })
                    .then((res) => {
                        if (res.code === 0) {
                            message.success("删除成功");
                            that.showLoading();
                        } else {
                            message.error(res.msg);
                        }
                    })
                    .catch((err) => {
                        message.error(err.toLocaleString());
                    });
            },
            onCancel() {
                console.log("Cancel");
            },
        });
    }

    showYamlHandler(resType, data) {
        this.setState({ yamlModalVisible: true, yamlType: resType });
        getResourceYaml({
            namespace: data.metadata.namespace,
            resType: resType,
            resName: data.metadata.name,
        })
            .then((res) => {
                if (res.code === 0) {
                    this.setState({ yamlDetail: res.data });
                } else {
                    message.error(res.msg);
                }
            })
    }

    updateYamlContent() {
        putApplyYaml({
            namespace: this.state.currentNamespace,
            resName: this.state.yamlType,
            yamlContent: this.state.yamlDetail,
        })
            .then((res) => {
                if (res.code === 0) {
                    this.setState({ yamlModalVisible: false });
                    message.success("更新成功");
                    this.refreshNamespaceResources();
                } else {
                    message.error(res.msg);
                }
            })
            .catch((err) => {
                message.error(err.toLocaleString());
            });
    }

    handleYamlDetailCancel() {
        this.setState({ yamlModalVisible: false });
    }

    handleScaleCancel() {
        this.setState({ scaleModalVisible: false });
    }

    handleScaleCommit() {
        putResourceScale({
            namespace: this.state.scaleResData.metadata.namespace,
            resType: this.state.scaleResType,
            resName: this.state.scaleResData.metadata.name,
            replicaCount: this.state.replicaCount,
        })
            .then((res) => {
                if (res.code === 0) {
                    this.setState({ scaleModalVisible: false });
                    message.success("提交成功");
                    this.showLoading();
                    switch (this.state.scaleResType) {
                        case "rc":
                            this.getNamespaceReplicationControllers();
                            break;
                        default:
                            break;
                    }
                } else {
                    message.error(res.msg);
                }
            })
            .catch((err) => {
                message.error(err.toLocaleString());
            });
    }

    onInputNumberChange(data) {
        this.setState({ replicaCount: data });
    }

    render() {
        let modelTitle = "新建资源 - " + this.state.currentNamespace;
        let alertTitle = (
            <div>
                填入 YAML 文件内容，将指定资源部署到当前所选命名空间。
                <Tooltip
                    title={
                        <div>
                            更多配置信息请查看:{" "}
                            <a
                                href="https://kubernetes.io/docs/home/"
                                target="_blank"
                                rel="noopener noreferrer"
                            >
                                官方文档
                            </a>
                        </div>
                    }
                >
                    <Icon type="question-circle" />
                </Tooltip>
            </div>
        );
        let configList = [];
        for (let key in this.state.dictConfigValue) {
            configList.push(
                <Panel header={key} key={key} style={customPanelStyle}>
                    <CodeMirror value={this.state.dictConfigValue[key]} />
                </Panel>,
            );
        }
        return (
            <Content
                style={{
                    background: "#fff",
                    padding: "5px 20px",
                    margin: 0,
                    height: "100%",
                }}
            >
                <OpsBreadcrumbPath pathData={["Kubernetes", "集群管理"]} />

                <Row>
                    <span
                        style={{
                            lineHeight: "30px",
                            fontWeight: "500",
                            marginRight: "20px",
                        }}
                    >
                        命名空间:{" "}
                    </span>
                    <Select
                        defaultValue="default"
                        style={{ width: "200px", marginRight: "20px" }}
                        onChange={this.selectChange}
                    >
                        {this.state.namespaceList.map((item, index) => {
                            let disabled = false;
                            if (item["status"]["phase"] !== "Active") {
                                disabled = true;
                            }
                            return (
                                <Option
                                    key={index}
                                    value={item.metadata.name}
                                    disabled={disabled}
                                >
                                    {item.metadata.name}
                                </Option>
                            );
                        })}
                    </Select>
                    <Button
                        onClick={this.refreshNamespaceResources}
                        style={{ marginRight: "20px" }}
                    >
                        刷 新
                    </Button>
                    <Button type="primary" onClick={this.yarmCreateK8sResource} disabled={!this.props.aclAuthMap["POST:/kubernetes/yaml_resource"]}>
                        新 建
                    </Button>
                </Row>

                <Modal
                    title={modelTitle}
                    destroyOnClose={true}
                    visible={this.state.createModalVisible}
                    onOk={this.handleOk}
                    onCancel={this.handleCancel}
                    width={700}
                    centered={true}
                >
                    <Alert
                        message={alertTitle}
                        banner
                        style={{ marginBottom: "10px", marginTop: "-20px" }}
                    />
                    <Row style={{ marginBottom: "10px" }}>
                        <Col span={4}>
                            <span
                                style={{
                                    lineHeight: "30px",
                                    fontWeight: "500",
                                }}
                            >
                                选择资源模板:{" "}
                            </span>
                        </Col>
                        <Col span={6} style={{ marginRight: "20px" }}>
                            <Select
                                defaultValue=""
                                style={{ width: "100%" }}
                                onChange={this.selectTemplateChange}
                            >
                                <Option key="0" value="">
                                    无
                                </Option>
                                <Option key="1" value="deployment">
                                    Deployment
                                </Option>
                                <Option key="2" value="rc">
                                    ReplicationController
                                </Option>
                                <Option key="3" value="rs">
                                    ReplicaSet
                                </Option>
                                <Option key="4" value="statefull_rs">
                                    StatefulSets
                                </Option>
                                <Option key="5" value="service">
                                    Service
                                </Option>
                                <Option key="6" value="ingress">
                                    Ingress
                                </Option>
                                <Option key="7" value="daemonset">
                                    DaemonSet
                                </Option>
                                <Option key="8" value="configmap">
                                    ConfigMap
                                </Option>
                                <Option key="9" value="secret">
                                    Secret
                                </Option>
                            </Select>
                        </Col>
                    </Row>
                    <CodeMirror
                        value={this.state.yamlCode}
                        options={{
                            mode: "xml",
                            theme: "material",
                        }}
                        onBeforeChange={(editor, data, value) => {
                            this.setState({ yamlCode: value });
                        }}
                    />
                </Modal>

                <Modal
                    title="YAML详情"
                    destroyOnClose={true}
                    visible={this.state.yamlModalVisible}
                    onCancel={this.handleYamlDetailCancel}
                    width={700}
                    centered={true}
                    footer={[
                        <Button type="primary" onClick={this.updateYamlContent}>
                            应用
                        </Button>,
                        <Button onClick={this.handleYamlDetailCancel}>
                            取消
                        </Button>,
                    ]}
                >
                    <CodeMirror
                        value={this.state.yamlDetail}
                        options={{
                            mode: "xml",
                            theme: "material",
                        }}
                        onBeforeChange={(editor, data, value) => {
                            this.setState({ yamlDetail: value });
                        }}
                        onChange={(editor, data, value) => {
                            this.setState({
                                yamlDetail: value,
                            });
                        }}
                    />
                </Modal>

                <Modal
                    title="自动伸缩"
                    destroyOnClose={true}
                    visible={this.state.scaleModalVisible}
                    onCancel={this.handleScaleCancel}
                    width={400}
                    centered={true}
                    bodyStyle={{ textAlign: "center" }}
                    footer={[
                        <Button onClick={this.handleScaleCancel}>取消</Button>,
                        <Button type="primary" onClick={this.handleScaleCommit}>
                            确定
                        </Button>,
                    ]}
                >
                    <InputNumber
                        min={0}
                        defaultValue={this.state.scaleResData.spec.replicas}
                        onChange={this.onInputNumberChange}
                        style={{ width: "70%" }}
                    />
                </Modal>

                <Modal
                    title="配置字典值"
                    destroyOnClose={true}
                    visible={this.state.configModalVisible}
                    onCancel={this.handleConfigCancel}
                    width={600}
                    footer={[]}
                >
                    <Collapse bordered={false}>{configList}</Collapse>
                </Modal>

                <Drawer
                    placement="left"
                    closable={true}
                    destroyOnClose={true}
                    onClose={this.onDrawerClose}
                    visible={this.state.detailDrawerVisible}
                    width="70%"
                >
                    <Switch>
                        <Route
                            path="/admin/k8s_cluster/manage/node_detail"
                            component={NodeDetailContent}
                        />
                        <Route
                            path="/admin/k8s_cluster/manage/rc_detail"
                            component={RcDetailContent}
                        />
                        <Route
                            path="/admin/k8s_cluster/manage/rs_detail"
                            component={RsDetailContent}
                        />
                        <Route
                            path="/admin/k8s_cluster/manage/deployment_detail"
                            component={DeploymentDetailContent}
                        />
                        <Route
                            path="/admin/k8s_cluster/manage/service_detail"
                            component={ServiceDetailContent}
                        />
                        <Route
                            path="/admin/k8s_cluster/manage/ingress_detail"
                            component={IngressDetailContent}
                        />
                        <Route
                            path="/admin/k8s_cluster/manage/pod_detail"
                            component={PodDetailContent}
                        />
                        <Route
                            path="/admin/k8s_cluster/manage/container_log"
                            component={ContainerLogContent}
                        />
                        <Route
                            path="/admin/k8s_cluster/manage/container_terminal"
                            component={WebTerminalContent}
                        />
                    </Switch>
                </Drawer>

                <Row style={{ marginTop: "10px" }}>
                    <Card
                        size="small"
                        title="资源总览"
                        style={{ width: "100%" }}
                    >
                        <Spin spinning={this.state.refreshDataLoading}>
                            <Tabs
                                defaultActiveKey="1"
                                tabPosition="left"
                                size="small"
                                style={{ width: "100%" }}
                                tabBarStyle={{ textAlign: "left" }}
                            >
                                <TabPane
                                    tab={
                                        <div
                                            style={{
                                                fontSize: "13px",
                                                width: "80px",
                                                textAlign: "left",
                                            }}
                                        >
                                            Node
                                        </div>
                                    }
                                    key={1}
                                >
                                    <Table
                                        size="small"
                                        columns={this.state.nodeColumns}
                                        dataSource={this.state.nodeList}
                                        scroll={{ x: "max-content" }}
                                    />
                                </TabPane>
                                <TabPane
                                    tab={
                                        <div
                                            style={{
                                                fontSize: "13px",
                                                width: "80px",
                                                textAlign: "left",
                                            }}
                                        >
                                            ReplicaSet
                                        </div>
                                    }
                                    key={2}
                                >
                                    <Table
                                        size="small"
                                        columns={this.state.rsColumns}
                                        dataSource={this.state.replicaSetList}
                                        scroll={{ x: "max-content" }}
                                    />
                                </TabPane>
                                <TabPane
                                    tab={
                                        <div
                                            style={{
                                                fontSize: "13px",
                                                width: "80px",
                                                textAlign: "left",
                                            }}
                                        >
                                            Deployment
                                        </div>
                                    }
                                    key={3}
                                >
                                    <Table
                                        size="small"
                                        columns={this.state.deploymentColumns}
                                        dataSource={this.state.deploymentList}
                                        scroll={{ x: "max-content" }}
                                    />
                                </TabPane>
                                <TabPane
                                    tab={
                                        <div
                                            style={{
                                                fontSize: "13px",
                                                width: "80px",
                                                textAlign: "left",
                                            }}
                                        >
                                            RC
                                        </div>
                                    }
                                    key={4}
                                >
                                    <Table
                                        size="small"
                                        columns={this.state.rcColumns}
                                        dataSource={
                                            this.state.replicationControllerList
                                        }
                                        scroll={{ x: "max-content" }}
                                    />
                                </TabPane>
                                <TabPane
                                    tab={
                                        <div
                                            style={{
                                                fontSize: "13px",
                                                width: "80px",
                                                textAlign: "left",
                                            }}
                                        >
                                            Service
                                        </div>
                                    }
                                    key={5}
                                >
                                    <Table
                                        size="small"
                                        columns={this.state.serviceColumns}
                                        dataSource={this.state.serviceList}
                                        scroll={{ x: "max-content" }}
                                    />
                                </TabPane>
                                <TabPane
                                    tab={
                                        <div
                                            style={{
                                                fontSize: "13px",
                                                width: "80px",
                                                textAlign: "left",
                                            }}
                                        >
                                            Ingress
                                        </div>
                                    }
                                    key={6}
                                >
                                    <Table
                                        size="small"
                                        columns={this.state.ingressColumns}
                                        dataSource={this.state.ingressList}
                                        scroll={{ x: "max-content" }}
                                    />
                                </TabPane>
                                <TabPane
                                    tab={
                                        <div
                                            style={{
                                                fontSize: "13px",
                                                width: "80px",
                                                textAlign: "left",
                                            }}
                                        >
                                            Pod
                                        </div>
                                    }
                                    key={7}
                                >
                                    <Table
                                        size="small"
                                        columns={this.state.podColumns}
                                        dataSource={this.state.podList}
                                        scroll={{ x: "max-content" }}
                                    />
                                </TabPane>
                                <TabPane
                                    tab={
                                        <div
                                            style={{
                                                fontSize: "13px",
                                                width: "80px",
                                                textAlign: "left",
                                            }}
                                        >
                                            ConfigMap
                                        </div>
                                    }
                                    key={8}
                                >
                                    <Table
                                        size="small"
                                        columns={this.state.configDictColumns}
                                        dataSource={this.state.configDictList}
                                        scroll={{ x: "max-content" }}
                                    />
                                </TabPane>
                                <TabPane
                                    tab={
                                        <div
                                            style={{
                                                fontSize: "13px",
                                                width: "80px",
                                                textAlign: "left",
                                            }}
                                        >
                                            Secret
                                        </div>
                                    }
                                    key={9}
                                >
                                    <Table
                                        size="small"
                                        columns={this.state.secretDictColumns}
                                        dataSource={this.state.secretDictList}
                                        scroll={{ x: "max-content" }}
                                    />
                                </TabPane>
                            </Tabs>
                        </Spin>
                    </Card>
                </Row>
            </Content>
        );
    }
}

export default withRouter(K8sNamespacesContent);
