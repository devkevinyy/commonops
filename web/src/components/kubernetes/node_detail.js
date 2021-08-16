import React, { Component } from "react";
import { Descriptions, Layout, Tag, Typography, Tabs, Card } from "antd";
import { ReloadOutlined } from "@ant-design/icons";
import AreaChart from "../cloud_resource/common/area_chart.js";
import { getNodeMetrics } from "../../api/kubernetes";

const { Content } = Layout;
const { Text } = Typography;
const { TabPane } = Tabs;

class NodeDetailContent extends Component {
    constructor(props) {
        super(props);
        this.state = {
            nodeDetail: this.props.location.state,
            memoryChartData: [],
            cpuChartData: [],
        };
    }

    componentDidMount() {
        this.refreshNodeCpuMetrics();
        this.refreshNodeMemoryMetrics();
    }

    refreshNodeCpuMetrics() {
        getNodeMetrics({
            clusterId: localStorage.getItem("clusterId"),
            metricName: "cpu",
            nodeName: this.props.location.state.data.metadata.name,
        }).then((res) => {
            if (res.code === 0) {
                this.setState({
                    cpuChartData: res.data.items[0].metricPoints,
                });
            }
        });
    }

    refreshNodeMemoryMetrics() {
        getNodeMetrics({
            clusterId: localStorage.getItem("clusterId"),
            metricName: "memory",
            nodeName: this.props.location.state.data.metadata.name,
        }).then((res) => {
            if (res.code === 0) {
                this.setState({
                    memoryChartData: res.data.items[0].metricPoints,
                });
            }
        });
    }

    render() {
        let labels = [];
        for (let key in this.state.nodeDetail.data.metadata.labels) {
            labels.push(
                <div key={key}>
                    <Tag color="geekblue">
                        {key +
                            ":" +
                            this.state.nodeDetail.data.metadata.labels[key]}
                    </Tag>
                </div>,
            );
        }
        let annotations = [];
        for (let key in this.state.nodeDetail.data.metadata.annotations) {
            annotations.push(
                <div key={key}>
                    <Tag color="geekblue">
                        {key +
                            ":" +
                            this.state.nodeDetail.data.metadata.annotations[
                                key
                            ]}
                    </Tag>
                </div>,
            );
        }
        let addresses = [];
        const addressesData = this.state.nodeDetail.data.status.addresses;
        for (let i = 0; i < addressesData.length; i++) {
            addresses.push(
                <Text>
                    {addressesData[i]["type"] +
                        ":  " +
                        addressesData[i]["address"]}
                </Text>,
            );
            addresses.push(<br />);
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
                <Tabs defaultActiveKey="Node信息">
                    <TabPane tab="Node信息" key="Node信息">
                        <Descriptions bordered size="small" column={1}>
                            <Descriptions.Item label="名称">
                                {this.state.nodeDetail.data.metadata.name}
                            </Descriptions.Item>
                            <Descriptions.Item label="创建时间">
                                {
                                    this.state.nodeDetail.data.metadata
                                        .creationTimestamp
                                }
                            </Descriptions.Item>
                            <Descriptions.Item label="标签">
                                {labels}
                            </Descriptions.Item>
                            <Descriptions.Item label="注解">
                                {annotations}
                            </Descriptions.Item>
                            <Descriptions.Item label="地址">
                                {addresses}
                            </Descriptions.Item>
                            <Descriptions.Item label="kubelet端口">
                                {
                                    this.state.nodeDetail.data.status
                                        .daemonEndpoints.kubeletEndpoint.Port
                                }
                            </Descriptions.Item>
                            <Descriptions.Item label="容器版本">
                                {
                                    this.state.nodeDetail.data.status.nodeInfo
                                        .containerRuntimeVersion
                                }
                            </Descriptions.Item>
                            <Descriptions.Item label="kubelet版本">
                                {
                                    this.state.nodeDetail.data.status.nodeInfo
                                        .kubeletVersion
                                }
                            </Descriptions.Item>
                            <Descriptions.Item label="kubeProxy版本">
                                {
                                    this.state.nodeDetail.data.status.nodeInfo
                                        .kubeProxyVersion
                                }
                            </Descriptions.Item>
                            <Descriptions.Item label="系统">
                                {
                                    this.state.nodeDetail.data.status.nodeInfo
                                        .operatingSystem
                                }
                            </Descriptions.Item>
                            <Descriptions.Item label="架构">
                                {
                                    this.state.nodeDetail.data.status.nodeInfo
                                        .architecture
                                }
                            </Descriptions.Item>
                        </Descriptions>
                    </TabPane>
                    <TabPane tab="Node监控" key="Node监控">
                        <Card
                            size="small"
                            title="CPU"
                            extra={
                                <ReloadOutlined
                                    onClick={this.refreshNodeCpuMetrics.bind(
                                        this,
                                    )}
                                />
                            }
                            style={{ marginBottom: 20 }}
                        >
                            <AreaChart
                                width="100%"
                                height={200}
                                xField="timestamp"
                                unit="m"
                                data={this.state.cpuChartData}
                            />
                        </Card>
                        <Card
                            size="small"
                            title="Memory"
                            extra={
                                <ReloadOutlined
                                    onClick={this.refreshNodeMemoryMetrics.bind(
                                        this,
                                    )}
                                />
                            }
                        >
                            <AreaChart
                                width="100%"
                                height={200}
                                xField="timestamp"
                                unit="Mi"
                                data={this.state.memoryChartData}
                            />
                        </Card>
                    </TabPane>
                </Tabs>
            </Content>
        );
    }
}

export default NodeDetailContent;
