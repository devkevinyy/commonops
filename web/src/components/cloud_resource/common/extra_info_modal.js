import React, {Component} from "react";
import {Form, Input, Modal, Divider, Row, Col, DatePicker} from "antd";
import moment from 'moment';
moment.locale('zh-cn');

const EcsInfoForm = [
    {
        label: "来源",
        key: "resForm",
        type: "input",
        disabled: true,
    },
    {
        label: "实例Id",
        key: "instanceId",
        type: "input",
    },
    {
        label: "实例名称",
        key: "instanceName",
        type: "input",
    },
    {
        label: "内网ip",
        key: "innerIpAddress",
        type: "input",
    },
    {
        label: "外网ip",
        key: "publicIpAddress",
        type: "input",
    },
    {
        label: "私有ip",
        key: "privateIpAddress",
        type: "input",
    },
    {
        label: "CPU",
        key: "cpu",
        type: "input",
    },
    {
        label: "内存(G)",
        key: "memory",
        type: "input",
    },
    {
        label: "过期时间",
        key: "expiredTime",
        type: "date_input",
    },
    {
        label: "SSH Port",
        key: "sshPort",
        type: "input",
    },
    {
        label: "SSH 用户",
        key: "sshUser",
        type: "input",
    },
    {
        label: "SSH 密码",
        key: "sshPwd",
        type: "input",
    },
];

const RdsInfoForm = [
    {
        "label": "来源",
        "key": "resForm",
        "type": "input",
        "disabled": true,
    },
    {
        "label": "实例名称",
        "key": "dbInstanceDescription",
        "type": "input",
    },
    {
        "label": "内存(G)",
        "key": "dbMemory",
        "type": "input",
    },
    {
        "label": "磁盘(G)",
        "key": "dbInstanceStorage",
        "type": "input",
    },
    {
        "label": "过期时间",
        "key": "dbExpiredTime",
        "type": "date_input",
    }
];

const KvInfoForm = [
    {
        "label": "来源",
        "key": "resForm",
        "type": "input",
        "disabled": true,
    },
    {
        "label": "实例名称",
        "key": "kvInstanceName",
        "type": "input",
    },
    {
        "label": "带宽",
        "key": "kvBandwidth",
        "type": "input",
    },
    {
        "label": "容量(G)",
        "key": "kvCapacity",
        "type": "input",
    },
    {
        "label": "过期时间",
        "key": "kvExpiredTime",
        "type": "date_input",
    }
];

class ExtraInfoModal extends Component {

    constructor(props){
        super(props);
        this.handleOk = this.handleOk.bind(this);
        this.handleCancel = this.handleCancel.bind(this);
        this.switchChange = this.switchChange.bind(this);
        this.formRef = React.createRef();
        this.state = {
            usersData: [],
        }
    }

    componentDidMount() {
        this.generateBaseInfoForm();
    }

    switchChange(value) {
        this.setState({renewSwitch: value});
    }

    handleOk() {
        this.formRef.current.validateFields().then((values) => {
            if('expiredTime' in values){
                values.expiredTime = moment(values.expiredTime).format("YYYY-MM-DD 00:00:00");
            }
            if('dbExpiredTime' in values){
                values.dbExpiredTime = moment(values.dbExpiredTime).format("YYYY-MM-DD 00:00:00");
            }
            if('kvExpiredTime' in values){
                values.kvExpiredTime = moment(values.kvExpiredTime).format("YYYY-MM-DD 00:00:00");
            }
            this.props.handleOk(values);
        });
    }

    handleCancel() {
        this.formRef.current.resetFields();
        this.props.handleCancel();
    }

    getInputItem(item, disableInput) {
        let res = <Input />;
        switch (item.type) {
            case 'input':
                res = <Input disabled={item.disabled || disableInput}/>;
                break;
            case 'date_input':
                res = <DatePicker format="YYYY-MM-DD" disabled={item.disabled || disableInput}/>;
                break;
            default:
                break;
        }
        return res;
    }

    generateBaseInfoForm() {
        if(this.props.updateMode!=='single'){
            return ""
        }
        let disableInput = false;
        let resType = this.props.resType;
        let formContent;
        let result = [];
        const twoItemFormLayout = {
            labelCol: {span: 7},
            wrapperCol: {span: 17},
        };
        let data = [];
        switch (resType) {
            case 'ecs':
                data = EcsInfoForm;
                break;
            case 'rds':
                data = RdsInfoForm;
                break;
            case 'kv':
                data = KvInfoForm;
                break;
            default:
                break;
        }
        let formItem = data.map((item, index) => {
            let offsetNum = index % 2 === 0 ? 0 : 1;
            return (
                <Col span={11} offset={offsetNum} key={item.key}>
                    <Form.Item {...twoItemFormLayout} label={item.label} name={item.key} rules={[
                        {required: item.required, message: '该值为必填项！'}
                    ]}>
                        {this.getInputItem(item, disableInput)}
                    </Form.Item>
                </Col>
            )
        });
        result.push(formItem);
        formContent = (
            <div>
                <Divider orientation="left" style={{ marginTop: '0px', color: 'rgb(255, 80, 23)'}}>基本信息</Divider>
                <Row gutter={12}>{result}</Row>
            </div>
        );
        return formContent;
    }

    render() {
        const formItemLayout = {
            labelCol: {span: 7},
            wrapperCol: {span: 14},
        };
        return (
            <Modal
                visible={this.props.visible}
                destroyOnClose={true}
                title="完善资源信息"
                onOk={this.handleOk}
                onCancel={this.handleCancel}
                centered={true}
                width={700}
            >
                <Form ref={this.formRef} {...formItemLayout} initialValues={this.props.editData}>
                    {this.generateBaseInfoForm()}
                </Form>
            </Modal>
        )
    }

}

export default ExtraInfoModal;