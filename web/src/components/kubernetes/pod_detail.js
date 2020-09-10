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
    Col,
    Tabs
} from "antd";
import createG2 from 'g2-react';
import {getNodeMetrics} from "../../api/kubernetes";
import {K8sContainerIconSvg} from "../../assets/Icons";
import moment from 'moment';

const { Text } = Typography;
const { Content } = Layout;
const { TabPane } = Tabs;

const MemoryChart = createG2(chart => {
    chart.col('timestamp', {
      alias: '时间',
    });
    chart.col('value', {
      alias: '内存(M)'
    });
    chart.line().position('timestamp*value').size(2);
    chart.render();
});

const CpuChart = createG2(chart => {
    chart.col('timestamp', {
      alias: '时间',
    });
    chart.col('value', {
      alias: 'CPU(%)'
    });
    chart.line().position('timestamp*value').size(2);
    chart.render();
});

class PodDetailContent extends Component {

    constructor(props) {
        super(props);
        this.loadPodLog = this.loadPodLog.bind(this);
        this.handleCancelLogModal = this.handleCancelLogModal.bind(this);
        this.openWebTerminal = this.openWebTerminal.bind(this);
        this.state = {
            serviceDetail: this.props.location.state,
            memoryChartData: [],
            cpuChartData: [],
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

    tabChange(key) {
        if(key==="Pod监控"){
            let cpuQuery = 'sum(rate(container_cpu_usage_seconds_total{pod="'+this.state.serviceDetail.data.metadata.name+'", container!="POD", container!=""}[5m])) / (sum(container_spec_cpu_quota{image!=""}/100000) ) * 100 ';
            getNodeMetrics({query: encodeURIComponent(cpuQuery)}).then((res)=>{
                if(res.code===0){
                    let data = JSON.parse(res.data);
                    let points = data["data"]["result"][0]["values"];
                    let cpuChartData = [];
                    for(let i=0; i<points.length; i++){
                        cpuChartData.push({
                            "timestamp": moment(points[i][0]*1000).format("HH:mm:ss"),
                            "value": (parseFloat(points[i][1])),
                        });
                    }
                    this.setState({cpuChartData});
                } else {
                    message.error(res.msg);
                }
            });
            let memoryQuery = 'container_memory_usage_bytes{pod="'+this.state.serviceDetail.data.metadata.name+'", container!="POD", container!=""}';
            getNodeMetrics({query: encodeURIComponent(memoryQuery)}).then((res)=>{
                if(res.code===0){
                    let data = JSON.parse(res.data);
                    let points = data["data"]["result"][0]["values"];
                    let memoryChartData = [];
                    for(let i=0; i<points.length; i++){
                        memoryChartData.push({
                            "timestamp": moment(points[i][0]*1000).format("HH:mm:ss"),
                            "value": (parseInt(points[i][1])/1024/1024),
                        });
                    }
                    this.setState({memoryChartData});
                } else {
                    message.error(res.msg);
                }
            });
        }
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
                <Tabs defaultActiveKey="Pod详情" onChange={this.tabChange.bind(this)}>
                    <TabPane tab="Pod详情" key="Pod详情">
                        <Descriptions bordered size='small' column={2}>
                            <Descriptions.Item label="Pod名称">{this.state.serviceDetail.data.metadata.name}</Descriptions.Item>
                            <Descriptions.Item label="命名空间">{this.state.serviceDetail.data.metadata.namespace}</Descriptions.Item>
                            <Descriptions.Item label="创建时间">{this.state.serviceDetail.data.metadata.creationTimestamp}</Descriptions.Item>
                            <Descriptions.Item label="标签">{labels}</Descriptions.Item>
                            <Descriptions.Item label="状态">{this.state.serviceDetail.data.status.phase}</Descriptions.Item>
                            <Descriptions.Item label="Host IP">{this.state.serviceDetail.data.status.hostIP}</Descriptions.Item>
                            <Descriptions.Item label="Pod IP">{this.state.serviceDetail.data.status.podIP}</Descriptions.Item>
                        </Descriptions>
                        <Card size="small" title="容器列表" style={{ width: '100%', marginTop: '10px' }}>
                            { containersList }
                        </Card>
                    </TabPane>
                    <TabPane tab="Pod监控" key="Pod监控">
                        <MemoryChart
                            data={this.state.memoryChartData}
                            height={300}
                            forceFit={true} />
                        <CpuChart
                            data={this.state.cpuChartData}
                            height={300}
                            forceFit={true} />
                    </TabPane>
                </Tabs>
                

            </Content>
        )
    }

}

export default PodDetailContent;
