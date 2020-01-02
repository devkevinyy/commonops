import React, { Component } from 'react';
import {Form, Icon, Input, Button, Row, Col, message} from 'antd';
import '../assets/css/login.css';
import {postUserLogin} from "../api/user";

class LoginContent extends Component {

    constructor(props) {
        super(props);
        this.handleSubmit = this.handleSubmit.bind(this);
        this.state = {};
    }

    componentDidMount() {
        let code = new URLSearchParams(this.props.location.search).get("code");
        if(code!==undefined && code !== null){
            code = parseInt(code);
            if(code === 401) {
                message.error("token过期，请重新登录!", 4);
            }
            if(code === 403) {
                message.error("用户不存在或被禁用，请联系管理员!", 4);
            }
        }
    }

    handleEnterKey = (e) => {
        if(e.nativeEvent.keyCode === 13){
            this.handleSubmit();
        }
    };

    handleSubmit(e) {
        this.props.form.validateFields((err, values) => {
            if (!err) {
                postUserLogin(values).then((res)=>{
                    if(res.code === 0){
                         localStorage.setItem('ops_token', res.data.token);
                         message.success("欢迎使用运维平台！");
                         this.props.history.push('/admin');
                    } else {
                        this.props.form.setFields({
                            password: {
                                value: '',
                                errors: [new Error(res.msg)],
                            },
                        });
                        message.error(res.msg);
                    }
                }).catch((err)=>{
                    message.error(err.toLocaleString());
                });
            }
        });
    };

  render() {
    const { getFieldDecorator } = this.props.form;
    let loginPanel = (
        <div>
            <Form onSubmit={this.handleSubmit} className="login-form">
                <Form.Item>
                    {getFieldDecorator('username', {
                        rules: [{ required: true, message: '请输入您的注册邮箱!' }],
                    })(
                        <Input
                            size="large"
                            prefix={<Icon type="user" style={{ color: 'rgba(0,0,0,.25)' }} />}
                            placeholder="用户账号"
                            onKeyPress={this.handleEnterKey}
                        />,
                    )}
                </Form.Item>
                <Form.Item>
                    {getFieldDecorator('password', {
                        rules: [{ required: true, message: '请输入您的密码!' }],
                    })(
                        <Input
                            size="large"
                            prefix={<Icon type="lock" style={{ color: 'rgba(0,0,0,.25)' }} />}
                            type="password"
                            placeholder="用户密码"
                            onKeyPress={this.handleEnterKey}
                        />,
                    )}
                </Form.Item>
                <Form.Item>
                    <Button type="primary" block className="login-form-button" size="large" onClick={this.handleSubmit}>
                        登 录
                    </Button>
                </Form.Item>
            </Form>
        </div>
    );
    return (
        <div>
            <Row style={{height: "25vh", paddingTop: "10vh"}}>
              <Col span={9}/>
              <Col span={6} style={{ padding: "0px 30px" }}>
                  <span className="login-logo"/>
                  <span className="login-logo-text">运维平台</span>
              </Col>
              <Col span={9}/>
            </Row>
            <Row style={{height: "30vh"}}>
                <Col span={9}/>
                <Col span={6} style={{ textAlign: 'center' }}>

                    { loginPanel }

                </Col>
                <Col span={9}/>
            </Row>
            <Row style={{height: "30vh"}} />
            <Row style={{height: "15vh", textAlign: 'center', paddingTop: 50}}>
                ©2019 Created by KevinYang
            </Row>
        </div>
    );
  }
}

LoginContent = Form.create({ name: 'normalogin' })(LoginContent);

export default LoginContent;



