import React, {Component} from 'react';
import {Card, Descriptions, Layout, Tag} from "antd";


const { Content } = Layout;

class RcDetailContent extends Component {

    constructor(props) {
        super(props);
        this.state = {
            rcDetail: this.props.location.state,
        };
    }

    render() {
        let labels = [];
        for(let key in this.state.rcDetail.data.metadata.labels){
            labels.push(<div key={key}><Tag color="geekblue">
                {key + ":" + this.state.rcDetail.data.metadata.labels[key]}
            </Tag></div>)
        }
        let selectors = [];
        for(let key in this.state.rcDetail.data.spec.selector){
            selectors.push(<div key={key}><Tag color="geekblue">
                {key + ":" + this.state.rcDetail.data.spec.selector[key]}
            </Tag></div>)
        }
        let images = [];
        const containersData = this.state.rcDetail.data.spec.template.spec.containers;
        for(let i=0; i<containersData.length; i++){
            images.push(<div key={i}><Tag color="geekblue">
                {containersData[i]["image"]}
            </Tag></div>)
        }
        return (
            <Content
                style={{ background: '#fff', padding: "5px 20px", margin: 0, height: "100%" }}
            >
                <Card size="small" title="Replication Controller 详情" style={{ width: '100%' }}>
                    <Descriptions bordered size='small' column={2}>
                        <Descriptions.Item label="名称">{this.state.rcDetail.data.metadata.name}</Descriptions.Item>
                        <Descriptions.Item label="命名空间">{this.state.rcDetail.data.metadata.namespace}</Descriptions.Item>
                        <Descriptions.Item label="创建时间">{this.state.rcDetail.data.metadata.creationTimestamp}</Descriptions.Item>
                        <Descriptions.Item label="标签">{labels}</Descriptions.Item>
                        <Descriptions.Item label="选择器">{selectors}</Descriptions.Item>
                        <Descriptions.Item label="镜像">{images}</Descriptions.Item>
                        <Descriptions.Item label="副本数">{this.state.rcDetail.data.status.replicas}</Descriptions.Item>
                        <Descriptions.Item label="运行中">{this.state.rcDetail.data.status.readyReplicas}</Descriptions.Item>
                    </Descriptions>
                </Card>
            </Content>
        )
    }

}


export default RcDetailContent;