import React, {Component} from 'react';
import {Card, Descriptions, Layout, Tag} from "antd";


const { Content } = Layout;

class RsDetailContent extends Component {

    constructor(props) {
        super(props);
        this.state = {
            rsDetail: this.props.location.state,
        };
    }

    render() {
        let labels = [];
        for(let key in this.state.rsDetail.data.metadata.labels){
            labels.push(<div key={key}><Tag color="geekblue">
                {key + ":" + this.state.rsDetail.data.metadata.labels[key]}
            </Tag></div>)
        }
        let annotations = [];
        for(let key in this.state.rsDetail.data.metadata.annotations){
            annotations.push(<div key={key}><Tag color="geekblue">
                {key + ":" + this.state.rsDetail.data.metadata.annotations[key]}
            </Tag></div>)
        }
        let selectors = [];
        for(let key in this.state.rsDetail.data.spec.selector.matchLabels){
            selectors.push(<div key={key}><Tag color="geekblue">
                {key + ":" + this.state.rsDetail.data.spec.selector.matchLabels[key]}
            </Tag></div>)
        }
        let images = [];
        const containersData = this.state.rsDetail.data.spec.template.spec.containers;
        for(let i=0; i<containersData.length; i++){
            images.push(<div key={i}><Tag color="geekblue">
                {containersData[i]["image"]}
            </Tag></div>)
        }
        return (
            <Content
                style={{ background: '#fff', padding: "5px 20px", margin: 0, height: "100%" }}
            >
                <Card size="small" title="ReplicaSet 详情" style={{ width: '100%' }}>
                    <Descriptions bordered size='small' column={2}>
                        <Descriptions.Item label="名称">{this.state.rsDetail.data.metadata.name}</Descriptions.Item>
                        <Descriptions.Item label="命名空间">{this.state.rsDetail.data.metadata.namespace}</Descriptions.Item>
                        <Descriptions.Item label="创建时间">{this.state.rsDetail.data.metadata.creationTimestamp}</Descriptions.Item>
                        <Descriptions.Item label="标签">{labels}</Descriptions.Item>
                        <Descriptions.Item label="注解">{annotations}</Descriptions.Item>
                        <Descriptions.Item label="选择器">{selectors}</Descriptions.Item>
                        <Descriptions.Item label="镜像">{images}</Descriptions.Item>
                        <Descriptions.Item label="副本数">{this.state.rsDetail.data.status.replicas}</Descriptions.Item>
                        <Descriptions.Item label="运行中">{this.state.rsDetail.data.status.readyReplicas}</Descriptions.Item>
                    </Descriptions>
                </Card>
            </Content>
        )
    }

}


export default RsDetailContent;