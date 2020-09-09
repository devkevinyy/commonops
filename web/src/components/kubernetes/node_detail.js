import React, {Component} from 'react';
import {Card, Descriptions, Layout, List, Tag, Typography, Tabs, message} from "antd";
import createG2 from 'g2-react';
import {getNodeMetrics} from "../../api/kubernetes";
import moment from 'moment';

const { Content } = Layout;
const { Text } = Typography;
const { TabPane } = Tabs;

const Chart = createG2(chart => {
    chart.col('timestamp', {
      alias: '时间',
    });
    chart.col('value', {
      alias: '内存(G)'
    });
    chart.line().position('timestamp*value').size(2);
    chart.render();
  });

class NodeDetailContent extends Component {

    constructor(props) {
        super(props);
        this.state = {
            nodeDetail: this.props.location.state,
            chartData: [],
        };
    }

    tabChange(key) {
        if(key==="Node监控"){
            getNodeMetrics().then((res)=>{
                if(res.code===0){
                    let data = JSON.parse(res.data);
                    let points = data["data"]["result"][0]["values"];
                    let chartData = [];
                    for(let i=0; i<points.length; i++){
                        chartData.push({
                            "timestamp": moment(points[i][0]*1000).format("HH:mm:ss"),
                            "value": parseInt(points[i][1])/1024/1024/1024,
                        });
                    }
                    console.log(chartData);
                    this.setState({chartData});
                } else {
                    message.error(res.msg);
                }
            })
        }
    }

    render() {
        let labels = [];
        for(let key in this.state.nodeDetail.data.metadata.labels){
            labels.push(<div key={key}><Tag color="geekblue">
                {key + ":" + this.state.nodeDetail.data.metadata.labels[key]}
            </Tag></div>)
        }
        let annotations = [];
        for(let key in this.state.nodeDetail.data.metadata.annotations){
            annotations.push(<div key={key}><Tag color="geekblue">
                {key + ":" + this.state.nodeDetail.data.metadata.annotations[key]}
            </Tag></div>)
        }
        let addresses = [];
        const addressesData = this.state.nodeDetail.data.status.addresses;
        for(let i=0; i<addressesData.length; i++){
            addresses.push(<Text>{addressesData[i]["type"] + ":  " + addressesData[i]["address"]}</Text>);
            addresses.push(<br/>);
        }
        return (
            <Content
                style={{ background: '#fff', padding: "5px 20px", margin: 0, height: "100%" }}
            >
                <Tabs defaultActiveKey="Node信息" onChange={this.tabChange.bind(this)}>
                    <TabPane tab="Node信息" key="Node信息">
                        <Descriptions bordered size='small' column={1}>
                            <Descriptions.Item label="名称">{this.state.nodeDetail.data.metadata.name}</Descriptions.Item>
                            <Descriptions.Item label="创建时间">{this.state.nodeDetail.data.metadata.creationTimestamp}</Descriptions.Item>
                            <Descriptions.Item label="标签">{labels}</Descriptions.Item>
                            <Descriptions.Item label="注解">{annotations}</Descriptions.Item>
                            <Descriptions.Item label="地址">{addresses}</Descriptions.Item>
                            <Descriptions.Item label="kubelet端口">{this.state.nodeDetail.data.status.daemonEndpoints.kubeletEndpoint.Port}</Descriptions.Item>
                            <Descriptions.Item label="容器版本">{this.state.nodeDetail.data.status.nodeInfo.containerRuntimeVersion}</Descriptions.Item>
                            <Descriptions.Item label="kubelet版本">{this.state.nodeDetail.data.status.nodeInfo.kubeletVersion}</Descriptions.Item>
                            <Descriptions.Item label="kubeProxy版本">{this.state.nodeDetail.data.status.nodeInfo.kubeProxyVersion}</Descriptions.Item>
                            <Descriptions.Item label="系统">{this.state.nodeDetail.data.status.nodeInfo.operatingSystem}</Descriptions.Item>
                            <Descriptions.Item label="架构">{this.state.nodeDetail.data.status.nodeInfo.architecture}</Descriptions.Item>
                        </Descriptions>
                    </TabPane>
                    <TabPane tab="Node监控" key="Node监控">
                    <Chart
                        data={this.state.chartData}
                        width={400}
                        height={300}
                        forceFit={true} />
                    </TabPane>
                </Tabs>
                {/* <Card size="small" title="镜像列表" style={{ width: '100%' }}>
                    <List
                        size="small"
                        itemLayout="vertical"
                        bordered={false}
                        dataSource={this.state.nodeDetail.data.status.images}
                        renderItem={item => {
                                let imageList = [];
                                for(let i=0; i<item.names.length; i++){
                                    let size = (item.sizeBytes / 1024.0 /1024.0).toFixed(2);
                                    imageList.push(<div><Text style={{ width: "550px" }} ellipsis={true}>{item.names[i]}</Text><Tag color="#2db7f5" style={{ float: "right" }}>{size}MB</Tag></div>);
                                    imageList.push(<br/>);
                                }
                                return (
                                    <List.Item>{ imageList }</List.Item>
                                )
                            }
                        }
                    />
                </Card> */}
            </Content>
        )
    }

}


export default NodeDetailContent;