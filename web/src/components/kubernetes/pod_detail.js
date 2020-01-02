import React, {Component} from 'react';
import {
    Button,
    Card,
    Descriptions,
    Layout,
    Tag,
    message,
    List,
    Typography,
    Row,
    Col
} from "antd";
import {K8sContainerIconSvg} from "../../assets/Icons";

const { Text } = Typography;
const { Content } = Layout;

class PodDetailContent extends Component {

    constructor(props) {
        super(props);
        this.loadPodLog = this.loadPodLog.bind(this);
        this.handleCancelLogModal = this.handleCancelLogModal.bind(this);
        this.openWebTerminal = this.openWebTerminal.bind(this);
        this.state = {
            serviceDetail: this.props.location.state,
        };
    }

    loadPodLog(containerName) {
        if(this.state.serviceDetail.data.status.phase!=="Running"){
            message.error("当前pod未运行，无法查看日志！");
            return;
        }
        this.props.history.push({
            'pathname': '/admin/k8s_cluster/manage/container_log',
            'state': {
                "namespace": this.state.serviceDetail.data.metadata.namespace,
                "podName": this.state.serviceDetail.data.metadata.name,
                "containerName": containerName,
            }
        });
    }

    handleCancelLogModal() {
        this.setState({logModalVisible: false, logLoading: false});
    }

    openWebTerminal(containerName) {
        if(this.state.serviceDetail.data.status.phase!=="Running"){
            message.error("当前pod未运行，无法使用终端！");
            return;
        }
        this.props.history.push({
            'pathname': '/admin/k8s_cluster/manage/container_terminal',
            'state': {
                "namespace": this.state.serviceDetail.data.metadata.namespace,
                "podName": this.state.serviceDetail.data.metadata.name,
                "containerName": containerName,
            }
        });
    }

    render() {
        let labels = [];
        for(let key in this.state.serviceDetail.data.metadata.labels){
            labels.push(<div key={key}><Tag color="geekblue">
                {key + ":" + this.state.serviceDetail.data.metadata.labels[key]}
            </Tag></div>)
        };
        let containersList = this.state.serviceDetail.data.spec.containers.map((item)=>{
            let portsList = [];
            if(item.hasOwnProperty("ports")){
                for(let i=0; i<item.ports.length; i++){
                    portsList.push(<Tag color="geekblue">
                        { item.ports[i]["containerPort"]+ "(" + item.ports[i]["protocol"] + ")"}
                    </Tag>)
                }
            }
            return (
                <Row style={{ borderBottom: '1px solid #e8e8e8' }}>
                    <Col>
                        <List itemLayout="horizontal">
                            <List.Item>
                                <List.Item.Meta
                                    avatar={<K8sContainerIconSvg />}
                                    title={<Text type="secondary">容器名：{item.name} &nbsp;&nbsp;
                                        <Button type="link" onClick={this.loadPodLog.bind(this, item.name)}>查看日志</Button>
                                        <Button type="link" onClick={this.openWebTerminal.bind(this, item.name)}>执行命令</Button></Text>}
                                    description={<Text type="secondary">镜像：{item.image} &nbsp;&nbsp; 端口: {portsList}</Text>}
                                />
                            </List.Item>
                        </List>
                    </Col>
                </Row>
            )
        });

        return (
            <Content
                style={{ background: '#fff', padding: "5px 20px", margin: 0, height: "100%" }}
            >

                <Card size="small" title="Pod 详情" style={{ width: '100%' }}>
                    <Descriptions bordered size='small' column={2}>
                        <Descriptions.Item label="Pod名称">{this.state.serviceDetail.data.metadata.name}</Descriptions.Item>
                        <Descriptions.Item label="命名空间">{this.state.serviceDetail.data.metadata.namespace}</Descriptions.Item>
                        <Descriptions.Item label="创建时间">{this.state.serviceDetail.data.metadata.creationTimestamp}</Descriptions.Item>
                        <Descriptions.Item label="标签">{labels}</Descriptions.Item>
                        <Descriptions.Item label="状态">{this.state.serviceDetail.data.status.phase}</Descriptions.Item>
                        <Descriptions.Item label="Host IP">{this.state.serviceDetail.data.status.hostIP}</Descriptions.Item>
                        <Descriptions.Item label="Pod IP">{this.state.serviceDetail.data.status.podIP}</Descriptions.Item>
                    </Descriptions>

                </Card>
                <Card size="small" title="容器列表" style={{ width: '100%', marginTop: '10px' }}>
                    { containersList }
                </Card>

            </Content>
        )
    }

}

export default PodDetailContent;
