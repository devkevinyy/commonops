import React, { Component } from 'react';
import {Card, Layout, Switch, message, Modal, Col, Row, DatePicker, Alert} from "antd";
import {getSettingsValue, putSettingsValue} from "../../api/settings_api";

const { Content } = Layout;
const { RangePicker } = DatePicker;


class SettingsContent extends Component {

    constructor(props){
        super(props);
        this.onChangeAllowRenewConfirm = this.onChangeAllowRenewConfirm.bind(this);
        this.dateOnChange = this.dateOnChange.bind(this);
        this.handleSetExpireDateSubmit = this.handleSetExpireDateSubmit.bind(this);
        this.handleSetExpireDateCancel = this.handleSetExpireDateCancel.bind(this);
        this.state = {
            expireSettingModalVisible: false,
            isOpenRenewConfirm: false,
            expireStartDate: null,
            expireEndDate: null,
        }
    }

    componentWillMount() {
        this.loadAllowRenewConfirmStatus();
    }

    loadAllowRenewConfirmStatus() {
        getSettingsValue({key_name: "open_renew_confirm"}).then((res)=>{
            if(res.code===0){
                this.setState({
                    isOpenRenewConfirm: res.data["value"]==="true"
                })
            } else {
                message.error(res.msg);
            }
        }).catch((err)=>{
            message.error(err.toLocaleString());
        })
    }

    onChangeAllowRenewConfirm(data) {
        this.setState({isOpenRenewConfirm: data});
        if(data===true) {
            this.setState({expireSettingModalVisible: true})
        } else {
            putSettingsValue({
                key_name: "open_renew_confirm",
                value: data.toLocaleString(),
            }).then((res)=>{
                if(res.code===0){
                    message.success("设置成功");
                } else {
                    message.error(res.msg);
                }
            }).catch((err)=>{
                message.error(err.toLocaleString());
            });
        }
    }

    dateOnChange(date, dateString) {
        this.setState({
            expireStartDate: dateString[0],
            expireEndDate: dateString[1],
        })
    };

    handleSetExpireDateSubmit() {
        putSettingsValue({
            key_name: "open_renew_confirm",
            value: this.state.isOpenRenewConfirm.toLocaleString(),
        }).then((res)=>{
            if(res.code===0){
                message.success("设置成功");
            } else {
                message.error(res.msg);
            }
        }).catch((err)=>{
            message.error(err.toLocaleString());
        });
        putSettingsValue({
            key_name: "expire_start_date",
            value: this.state.expireStartDate,
        }).then((res)=>{
            if(res.code!==0){
                message.error(res.msg);
            }
        }).catch((err)=>{
            message.error(err.toLocaleString());
        });
        putSettingsValue({
            key_name: "expire_end_date",
            value: this.state.expireEndDate,
        }).then((res)=>{
            if(res.code!==0){
                message.error(res.msg);
            } else {
                this.setState({expireSettingModalVisible: false});
            }
        }).catch((err)=>{
            message.error(err.toLocaleString());
        });
    }

    handleSetExpireDateCancel() {
        this.setState({expireSettingModalVisible: false});
    }

    render() {
        return (
            <Content style={{
            background: '#fff', padding: 20, margin: 0, height: "100%",
            }}>
                <Modal
                    title="设置过期确认时间区间"
                    destroyOnClose={true}
                    visible={this.state.expireSettingModalVisible}
                    onOk={this.handleSetExpireDateSubmit}
                    onCancel={this.handleSetExpireDateCancel}
                    okText="确认"
                    cancelText="取消"
                    width={700}
                >
                    <Alert message="负责人只可确认该时间段内过期的资源" banner />

                    <Row style={{ marginTop: 30, marginBottom: 30 }}>
                        <Col span={4}/>
                        <Col span={3} style={{textAlign: "right"}}>
                            <label style={{fontSize:14, lineHeight: "32px", marginRight: 10, fontWeight: 500}}>过期时间: </label>
                        </Col>
                        <Col span={13}>
                            <RangePicker
                                style={{width: "100%"}}
                                onChange={this.dateOnChange}
                                format="YYYY-MM-DD"
                            />
                        </Col>
                        <Col span={4}>
                        </Col>
                    </Row>
                </Modal>
                <Card style={{ width: 200, textAlign: "center" }}>
                    <Switch checked={this.state.isOpenRenewConfirm} onChange={this.onChangeAllowRenewConfirm} />
                    <p style={{paddingTop: 20, fontWeight: 500}}>是否开启续费确认窗口</p>
                </Card>
            </Content>
        )
    }

}

export default SettingsContent;