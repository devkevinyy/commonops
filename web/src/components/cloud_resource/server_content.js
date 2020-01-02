import React, {Component} from 'react';
import {
    Button,
    Col,
    DatePicker, Descriptions, Divider,
    Drawer, Form, Typography,
    Input, InputNumber,
    Layout, message, Modal,
    Radio,
    Row,
    Select,
    Spin, Table,
    Tabs, Popconfirm, Badge, Tooltip
} from 'antd';
import OpsBreadcrumbPath from "../breadcrumb_path";
import echarts from 'echarts/lib/echarts';
import 'echarts/lib/chart/line';
import 'echarts/lib/component/tooltip';
import 'echarts/lib/component/toolbox';
import moment from 'moment';
import "../../assets/css/main.css";
import {
    deleteCloudServer,
    getCloudAccouts,
    getCloudMonitorEcs,
    getCloudServerDetail,
    getCloudServers,
    postCloudServer, putCloudServer
} from "../../api/cloud";
import 'moment/locale/zh-cn';
import ExtraInfoModal from "./common/extra_info_modal";
import {LinuxSvg, WindowsSvg} from "../../assets/Icons";
moment.locale('zh-cn');

const TabPane = Tabs.TabPane;
const { Text, Paragraph } = Typography;
const Option = Select.Option;
const { Content } = Layout;


class ServerInfoModal extends Component {
    constructor(props) {
        super(props);
        this.state = {
            cloudAccountList: []
        }
    }

    componentDidMount() {
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

    render() {
        const { getFieldDecorator } = this.props.form;
        const formItemLayout = {
            labelCol: { span: 7 },
            wrapperCol: { span: 14 },
        };
        let accountOptions;
        accountOptions = this.state.cloudAccountList.map((item)=>{
            return (
                <Option key={item.id} value={item.id}>{item.accountName}</Option>
            );
        });

        return (
            <Modal
                title="服务器信息"
                destroyOnClose={true}
                visible={this.props.server_info_modal_visible}
                onOk={this.props.handlePostServerInfoSubmit}
                onCancel={this.props.handlePostServerInfoCancel}
                okText="确认"
                cancelText="取消"
                centered={true}
                width={600}
            >
                <Form>
                    <Form.Item label="云账号" {...formItemLayout}>
                        {getFieldDecorator('cloudAccountId', {
                            rules: [{ required: true, message: '云账号不能为空！' }],
                        })(
                            <Select>
                                {accountOptions}
                            </Select>
                        )}
                    </Form.Item>
                    <Form.Item label="主机名" {...formItemLayout}>
                        {getFieldDecorator('hostName', {
                            rules: [{ }],
                        })(
                            <Input />
                        )}
                    </Form.Item>
                    <Form.Item label="实例Id" {...formItemLayout}>
                        {getFieldDecorator('instanceId', {
                            rules: [{ }],
                        })(
                            <Input />
                        )}
                    </Form.Item>
                    <Form.Item label="实例名称" {...formItemLayout}>
                        {getFieldDecorator('instanceName', {
                            rules: [{ required: true, message: '实例名称不能为空！' }],
                        })(
                            <Input />
                        )}
                    </Form.Item>
                    <Form.Item label="实例描述" {...formItemLayout}>
                        {getFieldDecorator('description', {
                            rules: [{ }],
                        })(
                            <Input />
                        )}
                    </Form.Item>
                    <Form.Item label="内网IP" {...formItemLayout}>
                        {getFieldDecorator('innerIpAddress', {
                            rules: [{ }],
                        })(
                            <Input />
                        )}
                    </Form.Item>
                    <Form.Item label="外网IP" {...formItemLayout}>
                        {getFieldDecorator('publicIpAddress', {
                            rules: [{ }],
                        })(
                            <Input />
                        )}
                    </Form.Item>
                    <Form.Item label="系统类型" {...formItemLayout}>
                        {getFieldDecorator('osType', {
                            initialValue: "linux",
                            rules: [{ required: true, message: '系统类型不能为空！' }],
                        })(
                            <Select>
                                <Option value="linux">Linux</Option>
                                <Option value="windows">Windows</Option>
                            </Select>
                        )}
                    </Form.Item>
                    <Form.Item label="CPU核数(个)" {...formItemLayout}>
                        {getFieldDecorator('cpu', {
                            rules: [{ type: "integer", required: true, message: '请输入数值型数据！' }],
                        })(
                            <InputNumber />
                        )}
                    </Form.Item>
                    <Form.Item label="内存(G)" {...formItemLayout}>
                        {getFieldDecorator('memory', {
                            rules: [{ type: "integer", required: true, message: '请输入数值型数据！' }],
                        })(
                            <InputNumber />
                        )}
                    </Form.Item>
                    <Form.Item label="磁盘(G)" {...formItemLayout}>
                        {getFieldDecorator('disk', {
                            rules: [{ type: "integer" }],
                        })(
                            <InputNumber />
                        )}
                    </Form.Item>
                    <Form.Item label="创建时间" {...formItemLayout}>
                        {getFieldDecorator('createTime', {
                            rules: [{ required: true, message: '创建时间不能为空！' }],
                        })(
                            <DatePicker format="YYYY-MM-DD" />
                        )}
                    </Form.Item>
                    <Form.Item label="过期时间" {...formItemLayout}>
                        {getFieldDecorator('expiredTime', {
                        })(
                            <DatePicker format="YYYY-MM-DD" />
                        )}
                    </Form.Item>
                </Form>
            </Modal>
        )
    }
}

ServerInfoModal = Form.create({ name: 'ServerInfoModal' })(ServerInfoModal);

class ServerContent extends Component {

    constructor(props) {
        super(props);
        this.handlePostServerInfoSubmit = this.handlePostServerInfoSubmit.bind(this);
        this.handlePostServerInfoCancel = this.handlePostServerInfoCancel.bind(this);
        this.handleExtraInfoOk = this.handleExtraInfoOk.bind(this);
        this.handleExtraInfoCancel = this.handleExtraInfoCancel.bind(this);
        this.onShowSizeChange = this.onShowSizeChange.bind(this);
        let operWidth = this.props.isSuperAdmin ? 220 : 100;
        this.state = {
            columns: [
                {
                    title: 'Id', dataIndex: 'ID', key: 'ID', className: 'small_font',
                    render: value => {
                        return (
                            <Text ellipsis={true}>{value}</Text>
                        )
                    }
                },
                {
                    title: '实例id', width: 230, dataIndex: 'InstanceId', key: 'InstanceId', className: 'small_font',
                },
                {
                    title: '实例名称', dataIndex: 'InstanceName', key: 'InstanceName', className: 'small_font',
                    width: 200,
                    textWrap: 'word-break',
                    render: value => {
                        return (
                            <Tooltip title={value}>
                                <Text ellipsis={true} style={{ width: '200px' }}>{value}</Text>
                            </Tooltip>
                        )
                    }
                },
                {
                    title: 'ip', dataIndex: 'ip', key: 'ip', className: 'small_font',
                    render: (value, record) => {
                        let innerContent = "";
                        let privateContent = "";
                        let publicContent = "";
                        if(record.InnerIpAddress){
                            innerContent = (<div>
                                <Paragraph style={{marginBottom: 0, display:'inline'}} copyable={record.InnerIpAddress!==""}>{record.InnerIpAddress}</Paragraph>(内网)
                            </div>)
                        }
                        if(record.PrivateIpAddress){
                            privateContent = (<div>
                                <Paragraph style={{marginBottom: 0, display:'inline'}} copyable={record.PrivateIpAddress!==""}>{record.PrivateIpAddress}</Paragraph>(私有)
                            </div>)
                        }
                        if(record.PublicIpAddress){
                            publicContent = (<div>
                                <Paragraph style={{marginBottom: 0, display:'inline'}} copyable={record.PublicIpAddress!==""}>{record.PublicIpAddress}</Paragraph>(外网)
                            </div>)
                        }
                        return (
                            <div className="ip_column">
                                { innerContent }
                                { privateContent }
                                { publicContent }
                            </div>
                        );
                    }
                },
                {
                    title: '配置', dataIndex: '配置', key: '配置', className: 'small_font',
                    render:  (value, record) => {
                        let cpuContent = <Paragraph style={{marginBottom: 0, display: 'inline'}}>{record.Cpu}核</Paragraph>;
                        let memoryContent = <Paragraph style={{marginBottom: 0, display: 'inline'}}>{record.Memory}G</Paragraph>;
                        let trafficType = "";
                        if(record.InternetChargeType==="PayByTraffic"){
                            trafficType = "流量"
                        }
                        if(record.InternetChargeType==="PayByBandwidth"){
                            trafficType = "带宽"
                        }
                        let trafficOutContent = (<div>
                                            <Paragraph style={{marginBottom: 0, display: 'inline'}}>{record.InternetMaxBandwidthOut}Mbps({trafficType})</Paragraph>
                                    </div>);
                        return (
                            <div className="ip_column">
                                { cpuContent } &nbsp;
                                { memoryContent }
                                { trafficOutContent }
                            </div>
                        );
                    }
                },
                {
                    title: '系统类型', dataIndex: 'OSType', key: 'OSType', align: 'center', className: 'small_font',
                    render: (value, record) => {
                        let status = "error";
                        if(record.Status==="Running") {
                            status = "processing";
                        }
                        if(value==="windows"){
                            return (
                                <div>
                                    <WindowsSvg/>
                                    <Badge status={status} style={{  marginLeft: '5px', position: 'relative', top: '-10px' }}/>
                                </div>
                            )
                        } else if(value==="linux"){
                            return (
                                <div>
                                    <LinuxSvg/>
                                    <Badge status={status} style={{ marginLeft: '5px', position: 'relative', top: '-10px' }} />
                                </div>
                            )
                        } else {
                            return <Text ellipsis={true}>{value}</Text>
                        }
                    }
                },
                {
                    title: '区域', dataIndex: 'ZoneId', key: 'ZoneId', className: 'small_font',
                    render: value => {
                        return (
                            <Text ellipsis={true}>{value}</Text>
                        )
                    }
                },
                {
                    title: '过期时间', dataIndex: 'ExpiredTime', key: 'ExpiredTime', className: 'small_font',
                    render: value => {
                        return (
                            <Text ellipsis={true}>{value}</Text>
                        )
                    }
                },
                {
                    title: '操作',
                    key: 'operation',
                    fixed: 'right',
                    width: {operWidth},
                    align: 'center',
                    className: 'small_font action_column',
                    render: (text, record) => {
                        return (
                            <div>
                                <Button type="primary" size="small" onClick={this.openMonitorDrawer.bind(this, record)}>监控</Button>
                                <Divider type="vertical" />
                                <Button type="info" size="small" onClick={this.serverEdit.bind(this, record)}>编辑</Button>
                                <Divider type="vertical" />
                                <Popconfirm
                                    title="确定删除该项内容?"
                                    onConfirm={this.serverDelete.bind(this, record)}
                                    okText="确认"
                                    cancelText="取消"
                                >
                                    <Button type="danger" size="small">删除</Button>
                                </Popconfirm>
                            </div>
                        )
                    },
                },
            ],
            tableLoading: false,
            webSocketReady: false,
            tableData: [],
            chartXData: [],
            chartYData: [],
            pagination: {
                showSizeChanger: true,
                pageSizeOptions: ['10', '20', '30', '100'],
                onShowSizeChange: this.onShowSizeChange,
                showQuickJumper: false,
                showTotal: (total) => `共 ${total} 条`,
                pageSize: 10,
                page: 1,
                total: 0,
                onChange: (current) => this.changePage(current),
            },
            drawerVisible: false,
            drawerPlacement: 'left',
            instanceId: "",
            timeTagValue: '1h',
            metricTagValue: 'CPUUtilization',
            chartFormat: '%',
            currentServerDetail: {},
            msgContent: "",
            server_info_modal_visible: false,
            queryExpiredTime: null,
            queryKeyword: "",
            queryCloudAccount: "0",
            queryManageUser: "0",
            queryDefineGroup: "",
            cloudAccountList: [],
            selectedRowKeys: [],
            idsList: [],
            updateMode: "single",
        }
    }

    onShowSizeChange(current, size){
        this.setState({
            pagination: {
                ...this.state.pagination,
                page: 1,
                current: 1,
                pageSize: size,
            }
        }, ()=>{
            this.refreshTableData();
        });
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

    serverEdit(data) {
        this.setState({
            extraInfoModalVisible: true,
            ecsId: data.ID,
            updateMode: "single",
            resFrom: data['DataStatus']===1 ? "阿里云" : "手动添加"
        }, ()=>{
            getCloudServerDetail(data.ID).then((res)=>{
                if(res['code'] !== 0){
                    message.error(res['msg']);
                } else {
                    this.extraInfoFormRef.props.form.setFieldsValue({
                        instanceId: res.data['InstanceId'],
                        innerIpAddress: res.data['InnerIpAddress'],
                        publicIpAddress: res.data['PublicIpAddress'],
                        privateIpAddress: res.data['PrivateIpAddress'],
                        instanceName: res.data['InstanceName'],
                        cpu: res.data['Cpu'],
                        memory: (res.data['Memory']/1024).toString(),
                        expiredTime: res.data['ExpiredTime'] !== "" ? moment(res.data['ExpiredTime']) : "",
                        resForm: res.data['DataStatus']===1 ? "阿里云" : "手动添加",
                    });
                }
            }).catch((err)=>{
                message.error(err.toLocaleString());
            });
        });
    }

    handleExtraInfoOk(data) {
        let targetId = "";
        if(this.state.updateMode==="single"){
            targetId = String(this.state.ecsId);
        } else {
            targetId = this.state.idsList.join(",");
        }
        putCloudServer({
            ...data,
            id: targetId,
        }).then((res)=>{
            if(res.code===0){
                message.success("修改成功");
                this.setState({extraInfoModalVisible: false, selectedRowKeys: []});
                this.refreshTableData();
            } else {
                message.error(res.msg);
            }
        }).catch((err)=>{
            message.error(err.toLocaleString());
        });
    }

    handleExtraInfoCancel(data) {
        this.setState({extraInfoModalVisible: false});
    }

    serverDelete(data) {
        deleteCloudServer(data.ID).then((res)=>{
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

    refreshTableData = () => {
        this.setState({tableLoading: true});
        const queryParams = {
            "page": this.state.pagination.page,
            "size": this.state.pagination.pageSize,
            "queryExpiredTime": this.state.queryExpiredTime,
            "queryKeyword": this.state.queryKeyword,
            "queryCloudAccount": this.state.queryCloudAccount,
        };
        getCloudServers(queryParams).then((res)=>{
            const pagination = this.state.pagination;
            pagination.total = parseInt(res.data.total);
            pagination.page = parseInt(res.data.page);
            pagination.showTotal(parseInt(res.data.total));
            this.setState({
                pagination: {...pagination}
            });
            let data = res['data']['servers'];
            let tableData = [];
            for (let i = 0; i < data.length; i++) {
                tableData.push({
                    key: data[i]['ID'],
                    ID: data[i]['ID'],
                    Memory: data[i]['Memory']/1024,
                    Cpu: data[i]['Cpu'],
                    HostName: data[i]['HostName'],
                    InstanceId: data[i]['InstanceId'],
                    InnerIpAddress: data[i]['InnerIpAddress'],
                    PublicIpAddress: data[i]['PublicIpAddress'],
                    PrivateIpAddress: data[i]['PrivateIpAddress'],
                    InternetMaxBandwidthIn: data[i]['InternetMaxBandwidthIn'],
                    InternetMaxBandwidthOut: data[i]['InternetMaxBandwidthOut'],
                    InternetChargeType: data[i]['InternetChargeType'],
                    InstanceName: data[i]['InstanceName'],
                    OSType: data[i]['OSType'],
                    ZoneId: data[i]['ZoneId'],
                    OSName: data[i]['OSName'],
                    ExpiredTime:  moment(data[i]['ExpiredTime']).format('YYYY-MM-DD'),
                    Status: data[i]['Status'],
                    CloudAccountName: data[i]['CloudAccountName'],
                    DataStatus: data[i]['DataStatus'],
                });
            }
            this.setState({tableData: tableData, tableLoading: false});
        }).catch((err)=>{
            message.error(err)
        });
    };

    showChart = () => {
        let myChart = echarts.init(document.getElementById('ecsMonitorChart'));
        let option = {
            tooltip : {
                trigger: 'axis',
                axisPointer: {
                    label: {
                        backgroundColor: '#623485'
                    }
                }
            },
            xAxis : [
                {
                    type : 'category',
                    boundaryGap : false,
                    data : this.state.chartXData
                }
            ],
            yAxis : [
                {
                    type : 'value',
                    axisLabel:{
                        margin: 7,
                        formatter: '{value} '+ this.state.chartFormat,
                        textStyle:{
                            color: '#999',
                            fontSize: '12px'
                        },
                    },
                }
            ],
            series : [
                {
                    type:'line',
                    itemStyle : {
                        normal : {
                            color:'rgba(43, 182, 243, 0.65)',
                            lineStyle: {
                                color: "rgba(43, 182, 243, 0.65)",
                                shadowColor: "rgba(43, 182, 243, 0.65)",
                                shadowBlur: 10,
                            },
                            areaStyle: {
                                color: "rgba(43, 182, 243, 0.65)",
                                shadowColor: "rgba(43, 182, 243, 0.65)",
                                shadowBlur: 10,
                            },
                        }
                    },
                    data: this.state.chartYData
                }
            ]
        };
        myChart.setOption(option);
    };

    openMonitorDrawer = (data) => {
        this.setState({ drawerVisible: true, instanceId: data.InstanceId, currentServerDetail: data },
            ()=>{
                this.refreshMonitorData(data.InstanceId, this.state.timeTagValue, this.state.metricTagValue);
                this.refreshSeverDetail();
        });
    };

    refreshMonitorData = (instanceId, timeTagValue, metricTagValue) => {
        this.setState({chartLoading: true});
        getCloudMonitorEcs(instanceId, timeTagValue, metricTagValue).then((res)=>{
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
            let chartXData = [];
            let chartYData = [];
            for(let i=0; i<dataPoints.length; i++){
                chartXData.push(moment(dataPoints[i]['timestamp']).format("DD日HH:mm"));
                chartYData.push(dataPoints[i]['Average']);
            }
            this.setState({chartXData: chartXData, chartYData: chartYData}, ()=>{
                this.setState({chartLoading: false});
                this.showChart();
            });
        }).catch((err)=>{
            message.error(err.toLocaleString());
        });
    };

    // 获取服务器的详细信息
    refreshSeverDetail = (e) => {
        this.setState({serverDetailLoading: true});
        getCloudServerDetail(this.state.currentServerDetail.ID).then((res)=>{
            if(res['code'] !== 0){
                message.error(res['msg']);
            }
            this.setState({currentServerDetail: res['data']}, ()=> {
                this.setState({serverDetailLoading: false});
            });
        }).catch((err)=>{
            message.error(err.toLocaleString());
        });
    };

    changePage = (e) => {
        this.setState({
            pagination: {
                ...this.state.pagination,
                page: e,
                current: e
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
            case "CPUUtilization":
                this.setState({chartFormat: "%"});
                break;
            case "memory_usedutilization":
                this.setState({chartFormat: "%"});
                break;
            case "diskusage_utilization":
                this.setState({chartFormat: "%"});
                break;
            case "cpu_total":
                this.setState({chartFormat: "%"});
                break;
            default:
                this.setState({chartFormat: ""});
                break;
        }
        this.refreshMonitorData(this.state.instanceId, this.state.timeTagValue, e.target.value);
    };

    onExpiredDateChange = (moment) => {
        if(moment == null) {
            this.setState({queryExpiredTime: null});
        } else {
            this.setState({queryExpiredTime: moment});
        }
    };

    keywordOnChange = (e) => {
        this.setState({queryKeyword: e.target.value});
    };

    handleCloudAccountChange = (queryCloudAccount) => {
        this.setState({queryCloudAccount});
    };

    handleUserDefineGroupChange = (queryDefineGroup) => {
        this.setState({queryDefineGroup});
    };

    handleManageUserChange = (queryManageUser) => {
        this.setState({queryManageUser});
    };

    // 新增自定义机器信息
    handleAdd = () => {
        this.setState({server_info_modal_visible: true, ecsId: 0});
    };

    handlePostServerInfoSubmit() {
        this.serverInfoFormRef.props.form.validateFields((err, values) => {
            if (!err) {
                postCloudServer({
                    ...values,
                    createTime: values.createTime.format('YYYY-MM-DD HH:mm:ss'),
                    expiredTime: values.expiredTime === undefined ? undefined : values.expiredTime.format('YYYY-MM-DD HH:mm:ss'),
                }).then((res)=>{
                    if(res.code===0){
                        message.success("添加成功，请到权限管理中增加访问权限！");
                        this.setState({server_info_modal_visible: false});
                        this.refreshTableData();
                    } else {
                        message.error(res.msg);
                    }
                }).catch((err)=>{
                    message.error(err.toLocaleString());
                });
            }
        })
    };

    handlePostServerInfoCancel() {
        this.setState({server_info_modal_visible: false});
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
        const { selectedRowKeys } = this.state;
        const rowSelection = {
            selectedRowKeys,
            onChange: this.onSelectChange,
        };
        return (
              <Content
                  style={{ background: '#fff', padding: "5px 20px", margin: 0, height: "100%" }}
              >
                  <OpsBreadcrumbPath pathData={["云资源", "云服务器", "服务器列表"]} />

                  <Row style={{ padding: "0px 0px 10px 0px" }}>
                      <Col span={3} className="col-span">
                          <DatePicker style={{ width: "100%" }} defaultValue={this.state.queryExpiredTime} placeholder="到期截止时间" onChange={this.onExpiredDateChange}/>
                      </Col>
                      <Col span={5} className="col-span">
                          <Input placeholder="输入实例id\ip\实例名称查询" value={this.state.queryKeyword} onChange={this.keywordOnChange}/>
                      </Col>
                      <Col span={4} className="col-span">
                          <Select defaultValue={this.state.queryCloudAccount} style={{ width: "100%" }} onChange={this.handleCloudAccountChange}>
                              <Option value="0">所有云账号</Option>
                              {accountOptions}
                          </Select>
                      </Col>
                      <Col span={2} className="col-span">
                          <Button style={{ width: "100%" }} type="primary" icon="search" onClick={this.handleQuery}>查 询</Button>
                      </Col>
                      <Col span={2} className="col-span">
                          <Button style={{ width: "100%" }} icon="plus-circle" onClick={this.handleAdd}>新 增</Button>
                      </Col>
                  </Row>

                  <ServerInfoModal
                      wrappedComponentRef={(form) => this.serverInfoFormRef = form}
                      server_info_modal_visible={this.state.server_info_modal_visible}
                      handlePostServerInfoSubmit={this.handlePostServerInfoSubmit}
                      handlePostServerInfoCancel={this.handlePostServerInfoCancel}
                  />

                  {/*完善信息组件*/}
                  <ExtraInfoModal
                      wrappedComponentRef={(form) => this.extraInfoFormRef = form}
                      resType="ecs"
                      updateMode={this.state.updateMode}
                      resFrom={this.state.resFrom}
                      visible={this.state.extraInfoModalVisible}
                      handleOk={this.handleExtraInfoOk}
                      handleCancel={this.handleExtraInfoCancel}
                  />

                  {/*云服务器列表*/}
                  <Drawer
                      title="服务器详情及监控数据"
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
                                              <Radio.Button value="CPUUtilization">cpu使用率</Radio.Button>
                                              <Radio.Button value="memory_usedutilization">内存使用率</Radio.Button>
                                              <Radio.Button value="diskusage_utilization">磁盘使用率</Radio.Button>
                                              <Radio.Button value="cpu_total">平均负载</Radio.Button>
                                          </Radio.Group>
                                      </Col>
                                  </Row>
                                  <Row>
                                      <Col>
                                          <div id="ecsMonitorChart" style={{height:500}}>
                                          </div>
                                      </Col>
                                  </Row>
                              </Spin>
                          </TabPane>
                          <TabPane tab="信息详情" key="2">
                              <Spin tip="数据获取中..." spinning={this.state.serverDetailLoading}>
                                  <Descriptions title="基本信息" bordered size="small" column={2}>
                                      <Descriptions.Item label="主机名">{this.state.currentServerDetail.HostName}</Descriptions.Item>
                                      <Descriptions.Item label="机器描述">{this.state.currentServerDetail.Description}</Descriptions.Item>
                                      <Descriptions.Item label="实例ID">{this.state.currentServerDetail.InstanceId}</Descriptions.Item>
                                      <Descriptions.Item label="内网IP">{this.state.currentServerDetail.InnerIpAddress}</Descriptions.Item>
                                      <Descriptions.Item label="外网IP">{this.state.currentServerDetail.PublicIpAddress}</Descriptions.Item>
                                      <Descriptions.Item label="私有IP">{this.state.currentServerDetail.PrivateIpAddress}</Descriptions.Item>
                                      <Descriptions.Item label="Cpu">{this.state.currentServerDetail.Cpu}核</Descriptions.Item>
                                      <Descriptions.Item label="内存">{this.state.currentServerDetail.Memory/1024}G</Descriptions.Item>
                                      <Descriptions.Item label="公网入带宽">{this.state.currentServerDetail.InternetMaxBandwidthIn}Mbps</Descriptions.Item>
                                      <Descriptions.Item label="公网出带宽">{this.state.currentServerDetail.InternetMaxBandwidthOut}Mbps</Descriptions.Item>
                                      <Descriptions.Item label="网络计费">{this.state.currentServerDetail.InternetChargeType}</Descriptions.Item>
                                      <Descriptions.Item label="创建时间">{this.state.currentServerDetail.CreationTime}</Descriptions.Item>
                                      <Descriptions.Item label="过期时间">{this.state.currentServerDetail.ExpiredTime}</Descriptions.Item>
                                      <Descriptions.Item label="镜像ID">{this.state.currentServerDetail.ImageId}</Descriptions.Item>
                                      <Descriptions.Item label="付费类型">{this.state.currentServerDetail.InstanceChargeType}</Descriptions.Item>
                                      <Descriptions.Item label="网络类型">{this.state.currentServerDetail.InstanceNetworkType}</Descriptions.Item>
                                      <Descriptions.Item label="实例类型">{this.state.currentServerDetail.InstanceType}</Descriptions.Item>
                                      <Descriptions.Item label="系统名称">{this.state.currentServerDetail.OSName}</Descriptions.Item>
                                  </Descriptions>
                              </Spin>
                          </TabPane>
                      </Tabs>
                  </Drawer>
                  <Table
                      columns={this.state.columns}
                      dataSource={this.state.tableData}
                      scroll={{x: 'max-content'}}
                      pagination={this.state.pagination}
                      loading={this.state.tableLoading}
                      rowClassName="fixedHeight"
                      bordered
                      size="small"
                      rowSelection={rowSelection}
                  />
              </Content>
          );
    }
}


export default ServerContent;