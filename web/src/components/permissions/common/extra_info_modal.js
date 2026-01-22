import React, {Component} from "react";
import {Form, Input, Modal, Divider, Row, Col} from "antd";
import moment from 'moment';
moment.locale('zh-cn');

// auth link info form
const authLinkInfoForm = [
  {
    label: "Id",
    key: "Id",
    type: "input",
    disabled: true,
  },
  {
    label: "能否删除",
    key: "canDelete",
    type: "input",
  },
  {
    label: "权限类型",
    key: "authType",
    type: "input",
  },
  {
    label: "权限名称",
    key: "name",
    type: "input",
  },
  {
    label: "权限描述",
    key: "description",
    type: "input",
  },
  {
    label: "权限路径",
    key: "urlPath",
    type: "input",
  },
]

class ExtraInfoModal extends Component {
  constructor(props) {
    super(props);
    this.formRef = React.createRef();
  }

  componentDidMount() {
    this.generateBaseInfoForm();
  }

  handleOk = () => {
    this.formRef.current.validateFields().then((values) => {
       this.props.handleOk(values);
    });
  }

  handleCancel = () => {
    this.formRef.current.resetFields();
    this.props.handleCancel();
  }

  switchChange = (value) => {
    this.setState({
      renewSwitch: value
    });
  }

  getInputItem(item, disableInput) {
    let res = <Input />;
    switch(item.type) {
      case 'input':
        res = <Input disabled={item.disabled || disableInput}/>;
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
      case 'authLink':
        data = authLinkInfoForm;
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
