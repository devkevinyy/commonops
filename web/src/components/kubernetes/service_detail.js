import React, {Component} from 'react';
import {Card, Descriptions, Layout, Tag} from "antd";


const { Content } = Layout;

class ServiceDetailContent extends Component {

    constructor(props) {
        super(props);
        this.state = {
            serviceDetail: this.props.location.state,
        };
    }

    render() {
        let labels = [];
        for(let key in this.state.serviceDetail.data.metadata.labels) {
            labels.push(<div key={key}><Tag color="geekblue">
                {key + ":" + this.state.serviceDetail.data.metadata.labels[key]}
            </Tag></div>)
        }
        let selectors = [];
        for(let key in this.state.serviceDetail.data.spec.selector) {
            selectors.push(<div key={key}><Tag color="geekblue">
                {key + ":" + this.state.serviceDetail.data.spec.selector[key]}
            </Tag></div>)
        }
        let portsList = [];
        let ports = this.state.serviceDetail.data.spec.ports;
        for(let i=0; i<ports.length; i++){
            let portContent = "";
            if(ports[i].hasOwnProperty("port")){
                portContent += ports[i]["port"]+"(port) ";
            }
            if(ports[i].hasOwnProperty("targetPort")){
                portContent += ports[i]["targetPort"]+"(targetPort) ";
            }
            if(ports[i].hasOwnProperty("nodePort")){
                portContent += ports[i]["nodePort"]+"(nodePort) ";
            }
            portsList.push(<Tag color="geekblue">
                {ports[i]["protocol"] + ": " +portContent}
             </Tag>)
        }
        return (
            <Content
                style={{ background: '#fff', padding: "5px 20px", margin: 0, height: "100%" }}
            >
                <Card size="small" title="Service 详情" style={{ width: '100%' }}>
                    <Descriptions bordered size='small' column={2}>
                        <Descriptions.Item label="名称">{this.state.serviceDetail.data.metadata.name}</Descriptions.Item>
                        <Descriptions.Item label="命名空间">{this.state.serviceDetail.data.metadata.namespace}</Descriptions.Item>
                        <Descriptions.Item label="创建时间">{this.state.serviceDetail.data.metadata.creationTimestamp}</Descriptions.Item>
                        <Descriptions.Item label="标签">{labels}</Descriptions.Item>
                        <Descriptions.Item label="clusterIP">{this.state.serviceDetail.data.spec.clusterIP}</Descriptions.Item>
                        <Descriptions.Item label="类型">{this.state.serviceDetail.data.spec.type}</Descriptions.Item>
                        <Descriptions.Item label="sessionAffinity">{this.state.serviceDetail.data.spec.sessionAffinity}</Descriptions.Item>
                        <Descriptions.Item label="端口">{portsList}</Descriptions.Item>
                        <Descriptions.Item label="选择标签">{selectors}</Descriptions.Item>
                    </Descriptions>
                </Card>
            </Content>
        )
    }

}


export default ServiceDetailContent;