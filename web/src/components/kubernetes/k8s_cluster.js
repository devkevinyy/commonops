import React, { Component } from "react";
import OpsBreadcrumbPath from "../breadcrumb_path";
import {
    Row,
    Layout,
    Button,
    Modal,
    List,
    Card,
    Empty,
    Form,
    Input,
    message,
    Popconfirm,
    Col,
    Select,
} from "antd";
import { withRouter } from 'react-router-dom';
import {
    deleteCluster,
    getClusterData,
    postCluster,
} from "../../api/kubernetes";

const { Content } = Layout;
const { Option } = Select;

const selectBefore = (
    <Select defaultValue="http://" className="select-before">
        <Option value="http://">http://</Option>
        <Option value="https://">https://</Option>
    </Select>
);

class ClusterManageContent extends Component {
    constructor(props) {
        super(props);
        this.addCluster = this.addCluster.bind(this);
        this.handleClusterAddOk = this.handleClusterAddOk.bind(this);
        this.handleClusterAddCancel = this.handleClusterAddCancel.bind(this);
        this.loadClusterData = this.loadClusterData.bind(this);
        this.formRef = React.createRef();
        this.state = {
            addModelVisible: false,
            clusterData: [],
        };
    }

    componentDidMount() {
        this.loadClusterData();
    }

    enterToCluster(clusterId) {
        let that = this;
        message.loading("集群连接初始化中，即将跳转...", 0.6);
        localStorage.setItem("clusterId", clusterId);
        setTimeout(function() {
            that.props.history.push({ pathname: "/admin/k8s_cluster/manage" });
        }, 600);
    }

    loadClusterData() {
        getClusterData()
            .then((res) => {
                if (res.code === 0) {
                    this.setState({ clusterData: res.data.k8sData });
                } else {
                    message.error(res.msg);
                }
            })
            .catch((err) => {
                message.error(err.toLocaleString());
            });
    }

    addCluster() {
        this.setState({ addModelVisible: true });
    }

    handleClusterAddOk() {
        this.formRef.current.validateFields().then((values) => {
            postCluster(values)
                .then((res) => {
                    if (res.code === 0) {
                        message.success("添加成功");
                        this.setState({ addModelVisible: false });
                        this.loadClusterData();
                    } else {
                        message.error(res.msg);
                    }
                })
                .catch((err) => {
                    message.error(err.toLocaleString());
                });
        });
    }

    handleClusterAddCancel() {
        this.setState({ addModelVisible: false });
    }

    confirmDeleteCluster(id) {
        deleteCluster({ id: id })
            .then((res) => {
                if (res.code === 0) {
                    message.success("移除成功");
                    this.loadClusterData();
                } else {
                    message.error(res.msg);
                }
            })
            .catch((err) => {
                message.error(err.toLocaleString());
            });
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
                    padding: "5px 20px",
                    margin: 0,
                    height: "100%",
                }}
            >
                <OpsBreadcrumbPath pathData={["Kubernetes", "集群信息"]} />
                <Row>
                    <Col span={3}>
                        <Button
                            type="primary"
                            onClick={this.addCluster}
                            style={{ width: "100%" }}
                            disabled={!this.props.aclAuthMap["POST:/kubernetes/cluster"]}
                        >
                            新增集群
                        </Button>
                    </Col>
                </Row>

                <Row style={{ marginTop: "10px", width: "100%" }}>
                    {this.state.clusterData.length === 0 ? (
                        <Empty
                            style={{ width: "100%" }}
                            description="暂未添加任何kubernetes集群"
                        />
                    ) : (
                        <List
                            grid={{
                                gutter: 20,
                                column: 4,
                            }}
                            style={{ width: "100%" }}
                            dataSource={this.state.clusterData}
                            renderItem={(item) => (
                                <List.Item>
                                    <Card title={item.name} size="small">
                                        <div
                                            style={{
                                                height: "50px",
                                                fontSize: "13px",
                                            }}
                                        >
                                            {item.description}
                                        </div>
                                        <div style={{ marginTop: "10px" }}>
                                            <Button
                                                type="link"
                                                size="small"
                                                style={{ float: "left" }}
                                                onClick={this.enterToCluster.bind(
                                                    this,
                                                    item.clusterId,
                                                )}
                                            >
                                                进入集群
                                            </Button>
                                            <Popconfirm
                                                title="确定移除该集群吗?"
                                                okText="确认"
                                                cancelText="取消"
                                                onConfirm={this.confirmDeleteCluster.bind(
                                                    this,
                                                    item.id,
                                                )}
                                                disabled={!this.props.aclAuthMap["DELETE:/kubernetes/cluster"]}
                                            >
                                                <Button
                                                    type="link"
                                                    size="small"
                                                    style={{
                                                        color: "red",
                                                        float: "right",
                                                    }}
                                                    disabled={!this.props.aclAuthMap["DELETE:/kubernetes/cluster"]}
                                                >
                                                    删除集群
                                                </Button>
                                            </Popconfirm>
                                        </div>
                                    </Card>
                                </List.Item>
                            )}
                        />
                    )}
                </Row>

                <Modal
                    title="新增集群信息"
                    destroyOnClose={true}
                    visible={this.state.addModelVisible}
                    onOk={this.handleClusterAddOk}
                    onCancel={this.handleClusterAddCancel}
                >
                    <Form ref={this.formRef}>
                        <Form.Item
                            label="集群名称"
                            {...formItemLayout}
                            name="name"
                            rules={[
                                { required: true, message: "该项为必填项" },
                            ]}
                        >
                            <Input />
                        </Form.Item>
                        <Form.Item
                            label="集群描述"
                            {...formItemLayout}
                            name="description"
                            rules={[
                                { required: true, message: "该项为必填项" },
                            ]}
                        >
                            <Input />
                        </Form.Item>
                        <Form.Item
                            label="ApiServer"
                            {...formItemLayout}
                            name="apiServer"
                            rules={[
                                { required: true, message: "该项为必填项" },
                            ]}
                        >
                            <Input addonBefore={selectBefore} />
                        </Form.Item>
                        <Form.Item
                            label="Admin Token"
                            {...formItemLayout}
                            name="token"
                            rules={[
                                { required: true, message: "该项为必填项" },
                            ]}
                        >
                            <Input />
                        </Form.Item>
                    </Form>
                </Modal>
            </Content>
        );
    }
}

export default withRouter(ClusterManageContent);
