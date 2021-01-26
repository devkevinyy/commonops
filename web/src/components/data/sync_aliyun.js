import React, { Component } from "react";
import {Button, Card, Col, Layout, message, Row, Statistic} from "antd";
import OpsBreadcrumbPath from "../breadcrumb_path";
import {
    getAliyunStatisData,
    getSyncAliyunEcs,
    getSyncAliyunKv,
    getSyncAliyunRds,
    getSyncAliyunSlb
} from "../../api/data_api";

const { Content } = Layout;

class SyncAliyunContent extends Component {

    constructor(props) {
        super(props);
        this.syncAliyunEcs = this.syncAliyunEcs.bind(this);
        this.syncAliyunRds = this.syncAliyunRds.bind(this);
        this.syncAliyunKv = this.syncAliyunKv.bind(this);
        this.syncAliyunSlb = this.syncAliyunSlb.bind(this);
        this.state = {
            syncEcsLoading: false,
            syncRdsLoading: false,
            syncKvLoading: false,
            syncSlbLoading: false,
        }
    }

    componentDidMount() {
        this.syncAliyunStatisData();
    }

    syncAliyunStatisData() {
        getAliyunStatisData().then((res)=>{
            if(res.code===0){
                this.setState({
                    ecsCount: res.data.ecsCount,
                    rdsCount: res.data.rdsCount,
                    kvCount: res.data.kvCount,
                    slbCount: res.data.slbCount,
                })
            } else {
                message.error("获取钉钉统计数据异常")
            }
        }).catch((err)=>{
            message.error(err.toLocaleString())
        })
    }

    syncAliyunEcs() {
        const loadingType = "syncEcsLoading";
        this.setState({[loadingType]: true});
        getSyncAliyunEcs().then((res)=>{
            if(res.code===0){
                message.success("同步数据成功");
                this.setState({[loadingType]: false});
                this.syncAliyunStatisData();
            } else {
                message.error(res.msg);
                this.setState({[loadingType]: false});
            }
        }).catch((err)=>{
            message.error(err.toLocaleString());
            this.setState({[loadingType]: false});
        })
    }

    syncAliyunRds() {
        const loadingType = "syncRdsLoading";
        this.setState({[loadingType]: true});
        getSyncAliyunRds().then((res)=>{
            if(res.code===0){
                message.success("同步数据成功");
                this.setState({[loadingType]: false});
                this.syncAliyunStatisData();
            } else {
                message.error(res.msg);
                this.setState({[loadingType]: false});
            }
        }).catch((err)=>{
            message.error(err.toLocaleString());
            this.setState({[loadingType]: false});
        })
    }

    syncAliyunKv() {
        const loadingType = "syncKvLoading";
        this.setState({[loadingType]: true});
        getSyncAliyunKv().then((res)=>{
            if(res.code===0){
                message.success("同步数据成功");
                this.setState({[loadingType]: false});
                this.syncAliyunStatisData();
            } else {
                message.error(res.msg);
                this.setState({[loadingType]: false});
            }
        }).catch((err)=>{
            message.error(err.toLocaleString());
            this.setState({[loadingType]: false});
        })
    }

    syncAliyunSlb() {
        const loadingType = "syncSlbLoading";
        this.setState({[loadingType]: true});
        getSyncAliyunSlb().then((res)=>{
            if(res.code===0){
                message.success("同步数据成功");
                this.setState({[loadingType]: false});
                this.syncAliyunStatisData();
            } else {
                message.error(res.msg);
                this.setState({[loadingType]: false});
            }
        }).catch((err)=>{
            message.error(err.toLocaleString());
            this.setState({[loadingType]: false});
        })
    }

    render() {
        return (
            <Content style={{
                background: '#fff', padding: 20, margin: 0, height: "100%",
            }}>
                <OpsBreadcrumbPath pathData={["数据管理", "阿里云数据"]} />
                <div style={{ marginBottom: 20 }}>
                    <Card size="small" title="平台数据总览">
                        <Row>
                            <Col span={5}>
                                <Statistic title="ECS总数" value={this.state.ecsCount} precision={0} />
                                <Button style={{ marginTop: 16 }} type="primary"
                                        loading={this.state.syncEcsLoading}
                                        onClick={this.syncAliyunEcs}>
                                    同步平台数据
                                </Button>
                            </Col>
                            <Col span={5}>
                                <Statistic title="RDS总数" value={this.state.rdsCount} precision={0} />
                                <Button style={{ marginTop: 16 }} type="primary"
                                        loading={this.state.syncRdsLoading}
                                        onClick={this.syncAliyunRds}>
                                    同步平台数据
                                </Button>
                            </Col>
                            <Col span={5}>
                                <Statistic title="Redis总数" value={this.state.kvCount} precision={0} />
                                <Button style={{ marginTop: 16 }} type="primary"
                                        loading={this.state.syncKvLoading}
                                        onClick={this.syncAliyunKv}>
                                    同步平台数据
                                </Button>
                            </Col>
                            <Col span={5}>
                                <Statistic title="SLB总数" value={this.state.slbCount} precision={0} />
                                <Button style={{ marginTop: 16 }} type="primary"
                                        loading={this.state.syncSlbLoading}
                                        onClick={this.syncAliyunSlb}>
                                    同步平台数据
                                </Button>
                            </Col>
                        </Row>
                    </Card>
                </div>
            </Content>
        )
    }

}

export default SyncAliyunContent;