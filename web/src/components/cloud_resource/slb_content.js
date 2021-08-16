import React, { Component } from 'react';
import {
    Layout,
    Table,
    Button,
    Drawer,
    Tabs,
    Row,
    Col,
    message,
    Radio,
    Divider,
    Spin,
    Input, Descriptions, Typography, Popconfirm, Select, Badge, Tooltip
} from "antd";
import OpsBreadcrumbPath from "../breadcrumb_path";
import {
    getCloudSlbDetail,
    getCloudSlb,
    getCloudMonitorSlb,
    deleteCloudSlb,
    getCloudAccouts,
} from "../../api/cloud";
import LineChart from "./common/line_chart";
import { SearchOutlined } from "@ant-design/icons";
import moment from "moment";

const { Text, Paragraph } = Typography;
const { Content } = Layout;
const { Option } = Select;
const TabPane = Tabs.TabPane;

class SlbContent extends Component {

    constructor(props) {
        super(props);
        this.changePage = this.changePage.bind(this);
        let operWidth = this.props.isSuperAdmin ? 200 : 100;
        this.state = {
            columns: [
                {
                    title: 'Id', dataIndex: 'ID', key: 'ID',
                    render: value => {
                        return (
                            <Text ellipsis={true}>{value}</Text>
                        )
                    }
                },
                {
                    title: '实例id', dataIndex: 'LoadBalancerId', key: 'LoadBalancerId',
                    render: value => {
                        return (
                            <Text ellipsis={true}>{value}</Text>
                        )
                    }
                },
                {
                    title: '实例名称', dataIndex: 'LoadBalancerName', key: 'LoadBalancerName',
                    render: (value, record) => {
                        return (
                            <Tooltip title={value}>
                                <Text ellipsis={true} style={{ width: '100%' }}>{value}</Text>
                            </Tooltip>
                        )
                    }
                },
                {
                    title: '云账号', dataIndex: 'CloudAccountName', key: 'CloudAccountName', 
                    render: value => {
                        return (
                            <Text ellipsis={true}>{value}</Text>
                        )
                    }
                },
                {
                    title: '服务地址', dataIndex: 'Address', key: 'Address', width: 200,
                    render: value => {
                        return (
                            <Paragraph style={{marginBottom: 0}} copyable>{value}</Paragraph>
                        );
                    }
                },
                {
                    title: '计费方式', dataIndex: 'InternetChargeType', key: 'InternetChargeType',
                    render: value => {
                        if(value==="4"){
                            value = "按流量计费"
                        } else {
                            value = "按带宽计费"
                        }
                        return (
                            <Text ellipsis={true}>{value}</Text>
                        )
                    }
                },
                {
                    title: '创建时间', dataIndex: 'CreateTime', key: 'CreateTime',
                    render: value => {
                        return (
                            <Text ellipsis={true}>{value}</Text>
                        )
                    }
                },
                {
                    title: '实例状态', dataIndex: 'LoadBalancerStatus', key: 'LoadBalancerStatus', align: "center",
                    render: value => {
                        if(value==="active") {
                            return <Badge status="processing" />;
                        } else {
                            return <Badge status="error" />;
                        }
                    }
                },
                {
                    title: '操作',
                    key: 'operation',
                    fixed: 'right',
                    align: 'center',
                    width: {operWidth},
                    render: (text, record) => {
                        return (
                            <div>
                                <Button type="primary" size="small" onClick={this.openMonitorDrawer.bind(this, record)}>监控</Button>
                                <Divider type="vertical" />
                                <Popconfirm
                                    title="确定删除该项内容?"
                                    onConfirm={this.slbDelete.bind(this, record)}
                                    okText="确认"
                                    cancelText="取消"
                                    disabled={!this.props.aclAuthMap["DELETE:/cloud/slb"]}
                                >
                                    <Button type="danger" size="small" disabled={!this.props.aclAuthMap["DELETE:/cloud/slb"]}>删除</Button>
                                </Popconfirm>
                            </div>
                        )
                    },
                },
            ],
            tableLoading: false,
            tableData: [],
            chartData: [],
            pagination: {
                showSizeChanger: true,
                pageSizeOptions: ['10', '20', '30', '100'],
                onShowSizeChange: (current, size) => this.onShowSizeChange(current, size),
                showQuickJumper: false,
                showTotal: (total) => `共 ${total} 条`,
                pageSize: 10,
                page: 1,
                total: 0,
                onChange: (page, pageSize) => this.changePage(page, pageSize),
            },
            extraInfoModalVisible: false,
            drawerVisible: false,
            drawerPlacement: 'left',
            cloudAccountList: [],
            slbId: 0,
            instanceId: "",
            timeTagValue: '1h',
            metricTagValue: 'ActiveConnection',
            chartFormat: 'Count',
            currentServerDetail: {},
            queryKeyword: "",
            queryCloudAccount: "0",
            queryManageUser: "0",
            queryDefineGroup: "",
            selectedRowKeys: [],
            idsList: [],
            updateMode: "single",
        }
    }

    onShowSizeChange(current, size) {
        let pagination = {
          ...this.state.pagination,
          page: 1,
          current: 1,
          pageSize: size,
        };
        this.setState(
          {
            pagination: pagination,
          },
          () => {
            this.refreshTableData();
          }
        );
    }

    componentDidMount() {
        this.refreshTableData();
        this.loadCloudAccountsData();
    };

    loadCloudAccountsData() {
        let that = this;
        getCloudAccouts(1, 100).then((res)=>{
            if(res.code===0){
                that.setState({
                    cloudAccountList: res.data.accounts,
                })
            } else {
                message.error(res.msg);
            }
        }).catch((err)=>{
            message.error(err.toLocaleString());
        })
    }

    // 用户自定义内容修改
    userDefineEdit(data) {
        this.userDefineInfoFormRef.props.form.setFieldsValue({
            "userRemark": data.UserRemark,
            "userDefineGroup": data.UserDefineGroup,
        });
        this.setState({
            slbId: data.ID,
            userDefineInfoModalVisible: true,
        });
    }

    refreshTableData = () => {
        this.setState({tableLoading: true});
        const queryParams = {
            "page": this.state.pagination.page,
            "size": this.state.pagination.pageSize,
            "queryKeyword": this.state.queryKeyword,
            "queryCloudAccount": this.state.queryCloudAccount,
        };
        getCloudSlb(queryParams).then((res)=>{
            const pagination = this.state.pagination;
            pagination.total = parseInt(res.data.total);
            pagination.page = parseInt(res.data.page);
            pagination.showTotal(parseInt(res.data.total));
            this.setState({
                pagination
            });
            let data = res['data']['slb'];
            let tableData = [];
            for (let i = 0; i < data.length; i++) {
                tableData.push({
                    key: data[i]['ID'],
                    ID: data[i]['ID'],
                    LoadBalancerId: data[i]['LoadBalancerId'],
                    LoadBalancerName: data[i]['LoadBalancerName'],
                    Address: data[i]['Address'],
                    InternetChargeType: data[i]['InternetChargeType'],
                    CreateTime: moment(data[i]['CreateTime']).format("YYYY-MM-DD"),
                    LoadBalancerStatus: data[i]['LoadBalancerStatus'],
                    CloudAccountName: data[i]['CloudAccountName'],
                });
            }
            this.setState({tableData: tableData, tableLoading: false});
        }).catch((err)=>{
            message.error(err)
        });
    };

    openMonitorDrawer = (data) => {
        this.setState({ drawerVisible: true, instanceId: data.LoadBalancerId, currentServerDetail: data }, ()=>{
            this.refreshMonitorData(data.LoadBalancerId, this.state.timeTagValue, this.state.metricTagValue);
            this.refreshSeverDetail();
        });
    };

    refreshMonitorData = (instanceId, timeTagValue, metricTagValue) => {
        this.setState({chartLoading: true});
        getCloudMonitorSlb(instanceId, timeTagValue, metricTagValue).then((res)=>{
            if(res['code'] !== 0){
                message.error(res['msg']);
                this.setState({chartLoading: false});
                return;
            }
            if(res['data']['Datapoints']===""){
                message.warn("未获取到监控数据，可能是非阿里云机器或其它原因！");
                this.setState({chartLoading: false});
                return;
            }
            let dataPoints = JSON.parse(res['data']['Datapoints']);
            let chartData = [];
            for(let i=0; i<dataPoints.length; i++){
                chartData.push({
                    "date": moment(dataPoints[i]["timestamp"]).format("DD日HH:mm"),
                    "value": dataPoints[i]["Average"],
                });
            }
            this.setState({chartLoading: false, chartData: chartData});
        }).catch((err)=>{
            console.log(err)
        });
    };

    // 获取服务器的详细信息
    refreshSeverDetail = (e) => {
        this.setState({serverDetailLoading: true});
        getCloudSlbDetail(this.state.currentServerDetail.ID).then((res)=>{
            if(res['code'] !== 0){
                message.error(res['msg']);
            }
            this.setState({currentServerDetail: res['data']}, ()=> {
                this.setState({serverDetailLoading: false});
            });
        }).catch((err)=>{
            console.log(err)
        });
    };

    changePage = (page, pageSize) => {
        this.setState(
          {
            pagination: {
              ...this.state.pagination,
              page: page,
              current: page,
              pageSize: pageSize,
            },
          },
          () => {
            this.refreshTableData();
          }
        );
    };

    slbDelete(data) {
        deleteCloudSlb(data.ID).then((res)=>{
            if(res.code===0){
                message.success("删除成功");
                this.refreshTableData();
            } else {
                message.error(res.msg);
            }
        }).catch((err)=>{
            message.error(err.toLocaleString());
        })
    }

    keywordOnChange = (e) => {
        this.setState({queryKeyword: e.target.value});
    };

    handleCloudAccountChange = (queryCloudAccount) => {
        this.setState({queryCloudAccount});
    };

    // 用户自定义查询
    handleQuery = () => {
        this.setState({
            pagination: {
                ...this.state.pagination,
                page: 1,
                current: 1,
            }
        }, ()=>{
            this.refreshTableData();
        });
    };

    onCloseDrawer = () => {
        this.setState({ drawerVisible: false })
    };

    handleTimeTagChange = (e) => {
        this.setState({ timeTagValue: e.target.value });
        this.refreshMonitorData(this.state.instanceId, e.target.value, this.state.metricTagValue);
    };

    handleMetricTagChange = (e) => {
        this.setState({ metricTagValue: e.target.value });
        switch (e.target.value) {
            case "ActiveConnection":
                this.setState({chartFormat: "Count"});
                break;
            case "MaxConnection":
                this.setState({chartFormat: "Count"});
                break;
            case "HeathyServerCount":
                this.setState({chartFormat: "Count"});
                break;
            case "UnhealthyServerCount":
                this.setState({chartFormat: "Count"});
                break;
            default:
                this.setState({chartFormat: ""});
                break;
        }
        this.refreshMonitorData(this.state.instanceId, this.state.timeTagValue, e.target.value);
    };

    onSelectChange = selectedRowKeys => {
        this.setState({ selectedRowKeys });
    };

    render() {
        let accountOptions;
        accountOptions = this.state.cloudAccountList.map((item)=>{
            return (
                <Option key={item.id} value={item.id}>{item.accountName}</Option>
            );
        });
        return (
          <Content style={{
              background: '#fff', padding: "5px 20px", margin: 0, height: "100%",
          }}>
              <OpsBreadcrumbPath pathData={["云资源", "负载均衡", "负载均衡列表"]} />
              <Row style={{ padding: "0px 0px 10px 0px" }}>
                  <Col span={6} className="col-span">
                      <Input style={{ width: "100%" }} placeholder="输入实例id\名称\服务地址查询" value={this.state.queryKeyword} onChange={this.keywordOnChange} />
                  </Col>
                  <Col span={4} className="col-span">
                      <Select defaultValue={this.state.queryCloudAccount} style={{ width: "100%" }} onChange={this.handleCloudAccountChange}>
                          <Option value="0">所有云账号</Option>
                          {accountOptions}
                      </Select>
                  </Col>
                  <Col span={2} className="col-span">
                      <Button style={{ width: "100%" }} type="primary" icon={<SearchOutlined />} onClick={this.handleQuery} >查 询</Button>
                  </Col>
              </Row>

              <div>
                  <Drawer
                      title="实例详情及监控数据"
                      placement={this.state.drawerPlacement}
                      closable={true}
                      destroyOnClose={true}
                      onClose={this.onCloseDrawer}
                      visible={this.state.drawerVisible}
                      width={950}
                  >
                      <Tabs defaultActiveKey="1" tabPosition="left" style={{ marginLeft: -30 }}>
                          <TabPane tab="监控详情" key="1">
                              <Spin tip="图表生成中..." spinning={this.state.chartLoading}>
                                  <Row style={{ marginBottom: "10px" }}>
                                      <Col span={3} style={{ lineHeight: "30px" }}>时间维度：</Col>
                                      <Col span={15}>
                                          <Radio.Group value={this.state.timeTagValue} onChange={this.handleTimeTagChange}>
                                              <Radio.Button value="1h">1小时</Radio.Button>
                                              <Radio.Button value="6h">6小时</Radio.Button>
                                              <Radio.Button value="12h">12小时</Radio.Button>
                                              <Radio.Button value="1d">1 天</Radio.Button>
                                              <Radio.Button value="3d">3 天</Radio.Button>
                                              <Radio.Button value="7d">7 天</Radio.Button>
                                              <Radio.Button value="14d">14 天</Radio.Button>
                                          </Radio.Group>
                                      </Col>
                                  </Row>
                                  <Row>
                                      <Col span={3} style={{ lineHeight: "30px" }}>监控维度：</Col>
                                      <Col span={16}>
                                          <Radio.Group value={this.state.metricTagValue} onChange={this.handleMetricTagChange}>
                                              <Radio.Button value="ActiveConnection">端口当前活跃连接数</Radio.Button>
                                              <Radio.Button value="MaxConnection">端口当前并发连接数</Radio.Button>
                                              <Radio.Button value="HeathyServerCount">后端健康ECS实例个数</Radio.Button>
                                              <Radio.Button value="UnhealthyServerCount">后端异常ECS实例个数</Radio.Button>
                                          </Radio.Group>
                                      </Col>
                                  </Row>
                                  <Row style={{marginTop: 20}}>
                                    <Col>
                                        <LineChart width={800} height={400} data={this.state.chartData}/>
                                    </Col>
                                  </Row>
                              </Spin>
                          </TabPane>
                          <TabPane tab="信息详情" key="2">
                              <Spin tip="数据获取中..." spinning={this.state.serverDetailLoading}>
                                  <Descriptions title="基本信息" bordered size="small" column={2}>
                                      <Descriptions.Item label="实例ID">{this.state.currentServerDetail.LoadBalancerId}</Descriptions.Item>
                                      <Descriptions.Item label="实例名称">{this.state.currentServerDetail.LoadBalancerName}</Descriptions.Item>
                                      <Descriptions.Item label="实例状态">{this.state.currentServerDetail.LoadBalancerStatus}</Descriptions.Item>
                                      <Descriptions.Item label="服务地址">{this.state.currentServerDetail.Address}</Descriptions.Item>
                                      <Descriptions.Item label="计费方式">{this.state.currentServerDetail.InternetChargeType}</Descriptions.Item>
                                      <Descriptions.Item label="付费类型">{this.state.currentServerDetail.PayType}</Descriptions.Item>
                                      <Descriptions.Item label="IP版本">{this.state.currentServerDetail.AddressIPVersion}</Descriptions.Item>
                                      <Descriptions.Item label="创建时间">{this.state.currentServerDetail.CreateTime}</Descriptions.Item>
                                  </Descriptions>
                              </Spin>
                          </TabPane>
                      </Tabs>
                  </Drawer>
                  <Table
                      columns={this.state.columns}
                      dataSource={this.state.tableData}
                      scroll={{ x: 'max-content' }}
                      pagination={this.state.pagination}
                      loading={this.state.tableLoading}
                      bordered
                      size="small"
                  />
              </div>
          </Content>
      );
  }
}

export default SlbContent;
