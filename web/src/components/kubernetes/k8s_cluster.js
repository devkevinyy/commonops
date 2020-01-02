import React, {Component} from 'react';
import OpsBreadcrumbPath from "../breadcrumb_path";
import {Row, Layout, Button, Modal, List, Card, Empty, Form, Input, Icon, Upload, message, Popconfirm} from "antd";
import {ServerBase} from "../../config";
import {getLocalToken} from "../../utils/axios";
import {deleteCluster, getClusterData, postCluster, postSwitchCluster} from "../../api/kubernetes";

const {Content} = Layout;

let kubeConfigPath = "";

class ClusterManageContent extends Component {

    constructor(props) {
        super(props);
        this.addCluster = this.addCluster.bind(this);
        this.handleClusterAddOk = this.handleClusterAddOk.bind(this);
        this.handleClusterAddCancel = this.handleClusterAddCancel.bind(this);
        this.loadClusterData = this.loadClusterData.bind(this);
        this.state = {
            addModelVisible: false,
            uploadProps: {
                action: ServerBase + 'kubernetes/upload_cluster_kubeconfig',
                headers: {
                    Authorization: getLocalToken()
                },
                onChange(info) {
                    if (info.file.status === 'done') {
                        let resp = info.file.response;
                        if(resp.code===0){
                            kubeConfigPath = resp.data;
                            message.success("上传成功")
                        } else {
                            message.error("上传失败: ", resp.msg)
                        }
                    } else if (info.file.status === 'error') {
                        message.error(`${info.file.name} file upload failed.`);
                    }
                },
            },
            clusterData: [],
        }
    }

    componentDidMount() {
        this.loadClusterData();
    }

    enterToCluster(clusterId) {
        let that = this;
        message.loading('集群连接初始化中，即将跳转...', 0.6);
        localStorage.setItem('clusterId', clusterId);
        setTimeout(function () {
            that.props.history.push({'pathname': '/admin/k8s_cluster/manage'});
        }, 600);
    }

    loadClusterData() {
        getClusterData().then((res)=>{
            if(res.code===0){
                this.setState({clusterData: res.data.k8sData});
            } else {
                message.error(res.msg);
            }
        }).catch((err)=>{
            message.error(err.toLocaleString())
        })
    }

    addCluster() {
        this.setState({addModelVisible: true});
    }

    handleClusterAddOk() {
        if(kubeConfigPath===""){
            message.warn("kubeconfig文件必须要上传!");
            return;
        }
        this.props.form.validateFields((err, values) => {
            if(!err){
                let reqData = {
                    ...values,
                    kubeconfig_file_path: kubeConfigPath,
                };
                postCluster(reqData).then((res)=>{
                    if(res.code===0) {
                        message.success("添加成功");
                        kubeConfigPath = "";
                        this.setState({addModelVisible: false});
                        this.loadClusterData();
                    } else {
                        message.error(res.msg);
                    }
                }).catch((err)=>{
                    message.error(err.toLocaleString());
                })
            }
        });
    }

    handleClusterAddCancel() {
        this.setState({addModelVisible: false});
    }

    confirmDeleteCluster(id) {
        deleteCluster({id: id}).then((res)=>{
            if(res.code===0){
                message.success("移除成功");
                this.loadClusterData();
            } else {
                message.error(res.msg);
            }
        }).catch((err)=>{
            message.error(err.toLocaleString());
        })
    }

    render() {
        const { getFieldDecorator } = this.props.form;
        const formItemLayout = {
            labelCol: { span: 7 },
            wrapperCol: { span: 14 },
        };
        return (
            <Content
                style={{ background: '#fff', padding: "5px 20px", margin: 0, height: "100%" }}
            >
                <OpsBreadcrumbPath pathData={["Kubernetes", "集群信息"]} />
                <Row>
                    <Button type="primary" onClick={this.addCluster}>新增集群</Button>
                </Row>

                <Row style={{ marginTop: '10px' }}>
                    {
                        this.state.clusterData.length === 0 ? (
                            <Empty description="暂未添加任何kubernetes集群"/>
                        ) : (
                            <List
                                grid={{ gutter: 16, column: 4 }}
                                dataSource={this.state.clusterData}
                                renderItem={item => (
                                    <List.Item>
                                        <Card title={item.name} size="small">
                                            <div style={{height: '50px', fontSize: '13px'}}>{item.description}</div>
                                            <div style={{marginTop: '10px'}}>
                                                <Button type="link" size="small" style={{float: 'left'}} onClick={this.enterToCluster.bind(this, item.clusterId)}>进入集群</Button>
                                                <Popconfirm
                                                    title="确定移除该集群吗?"
                                                    okText="确认"
                                                    cancelText="取消"
                                                    onConfirm={this.confirmDeleteCluster.bind(this, item.id)}
                                                >
                                                    <Button type="link" size="small" style={{color: 'red', float: 'right'}}>删除集群</Button>
                                                </Popconfirm>
                                            </div>
                                        </Card>
                                    </List.Item>
                                )}
                            />
                        )
                    }
                </Row>

                <Modal
                    title="新增集群信息"
                    destroyOnClose={true}
                    visible={this.state.addModelVisible}
                    onOk={this.handleClusterAddOk}
                    onCancel={this.handleClusterAddCancel}
                >
                    <Form>
                        <Form.Item label="集群名称" {...formItemLayout}>
                            {getFieldDecorator('name', {
                                rules: [{ required: true, message: "该项为必填项" }],
                            })(
                                <Input />
                            )}
                        </Form.Item>
                        <Form.Item label="集群描述" {...formItemLayout}>
                            {getFieldDecorator('description', {
                                rules: [{ required: true, message: "该项为必填项" }],
                            })(
                                <Input />
                            )}
                        </Form.Item>
                        <Form.Item label="KubeConfig文件" {...formItemLayout}>
                            {getFieldDecorator('kubeconfig_file_path', {
                                rules: [{ required: true, message: "该项为必填项" }],
                            })(
                                <Upload {...this.state.uploadProps}>
                                    <Button>
                                        <Icon type="upload" /> 点击上传
                                    </Button>
                                </Upload>
                            )}
                        </Form.Item>
                    </Form>
                </Modal>

            </Content>
        )
    }
}

ClusterManageContent = Form.create({ name: 'ClusterManageContent' })(ClusterManageContent);

export default ClusterManageContent;
