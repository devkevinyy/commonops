import React, { Component } from "react";
import {
    Form,
    Icon,
    Input,
    Button,
    Row,
    Col,
    message,
    Popover,
    Typography,
} from "antd";
import "../assets/css/login.css";
import { postUserLogin } from "../api/user";

const { Text } = Typography;

class LoginContent extends Component {
    constructor(props) {
        super(props);
        this.handleSubmit = this.handleSubmit.bind(this);
        this.formRef = React.createRef();
        this.state = {
            currentYear: new Date().getFullYear(),
        };
    }

    componentDidMount() {
        let code = new URLSearchParams(this.props.location.search).get("code");
        if (code !== undefined && code !== null) {
            code = parseInt(code);
            if (code === 401) {
                message.error("token过期，请重新登录!", 4);
            }
            if (code === 403) {
                message.error("用户不存在或被禁用，请联系管理员!", 4);
            }
        }
    }

    handleEnterKey = (e) => {
        if (e.nativeEvent.keyCode === 13) {
            this.formRef.current.validateFields().then((values) => {
                this.handleSubmit(values);
            });
        }
    };

    handleSubmit(values) {
        postUserLogin(values)
            .then((res) => {
                if (res.code === 0) {
                    localStorage.setItem("ops_token", res.data.token);
                    message.success("欢迎使用运维平台！");
                    this.props.history.push("/admin");
                } else {
                    this.formRef.current.setFieldsValue({
                        password: "",
                    });
                    message.error(res.msg);
                }
            })
            .catch((err) => {
                message.error(err.toLocaleString());
            });
    }

    onFinishFailed(errorInfo) {
        console.log("Failed:", errorInfo);
    }

    render() {
        let loginPanel = (
            <div>
                <Form
                    ref={this.formRef}
                    onFinish={this.handleSubmit}
                    onFinishFailed={this.onFinishFailed}
                    className="login-form"
                >
                    <Form.Item
                        name="username"
                        rules={[
                            { required: true, message: "请输入您的注册邮箱!" },
                        ]}
                    >
                        <Input
                            size="large"
                            prefix={
                                <Icon
                                    type="user"
                                    style={{ color: "rgba(0,0,0,.25)" }}
                                />
                            }
                            placeholder="用户账号"
                            onKeyPress={this.handleEnterKey}
                        />
                    </Form.Item>
                    <Form.Item
                        name="password"
                        rules={[{ required: true, message: "请输入您的密码!" }]}
                    >
                        <Input
                            size="large"
                            prefix={
                                <Icon
                                    type="lock"
                                    style={{ color: "rgba(0,0,0,.25)" }}
                                />
                            }
                            type="password"
                            placeholder="用户密码"
                            onKeyPress={this.handleEnterKey}
                        />
                    </Form.Item>
                    <Form.Item>
                        <Button
                            type="primary"
                            block
                            className="login-form-button"
                            size="large"
                            htmlType="submit"
                        >
                            登 录
                        </Button>
                    </Form.Item>
                </Form>
            </div>
        );
        return (
            <div>
                <Row style={{ height: "15vh", paddingTop: "5vh" }}>
                    <Col span={9} />
                    <Col span={6} style={{ padding: "0px 30px" }}>
                        <span className="login-logo" />
                        <span className="login-logo-text">运维平台</span>
                    </Col>
                    <Col span={9} />
                </Row>
                <Row style={{ height: "40vh" }}>
                    <Col span={9} />
                    {/* <Col
                        span={10}
                        style={{ textAlign: "center" }}
                        className="login-ad"
                    ></Col> */}
                    <Col span={6} style={{ paddingTop: 40 }}>
                        {loginPanel}
                    </Col>
                    <Col span={9} />
                </Row>
                <Row style={{ height: "30vh" }} />
                <Row
                    style={{
                        height: "15vh",
                        textAlign: "center",
                        paddingTop: 50,
                        display: "block",
                    }}
                >
                    ©2019-{this.state.currentYear} Created by &nbsp;
                    <Popover content={<div className="wechat" />}>
                        <Text underline>KevinYang</Text>
                    </Popover>
                </Row>
            </div>
        );
    }
}

export default LoginContent;
