import React, { Component } from 'react';
import {Layout, Input, Form, Button, message} from 'antd';
import {postUpdatePassword} from "../../api/user";

const { Content } = Layout;

// 修改自己密码不应该需要认证aclAuthMap["POST:/user/updatePassword"]
class PasswordManager extends Component {

    constructor(props) {
        super(props);
        this.setState({

        });
    }

    updatePassword = (values) => {
        let reqData = {
            "password": values['password'],
            "confirm_password": values['confirm_password']
        };
        postUpdatePassword(reqData)
        .then((res)=>{
            if(res.code === 0){
                message.success("密码修改成功!", 1, () => {window.location.href = "/admin"});
            } else {
                message.error(res.msg);
            }
        })
        .catch((err)=>{
            console.log(err.toLocaleString());
        });
    };

    render() {
        const formItemLayout = {
            labelCol: { span: 2 },
            wrapperCol: { span: 6 },
        };
        return (
            <Content style={{
                background: '#fff', padding: 20, margin: 0, height: "100%",
            }}>
                <Form onFinish={this.updatePassword}>
                    <Form.Item
                        {...formItemLayout}
                        label="新密码"
                        name="password"
                        hasFeedback
                        rules={[
                            {
                                required: true,
                                message: "请输入新密码",
                            },
                            {
                                min: 8,
                                message: "密码不能少于8个字符",
                            },
                            {
                                pattern: /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[a-zA-Z\d]{8,}$/,
                                message: "密码必须包含大小写字母和数字",
                            },
                        ]}
                    >
                        <Input />
                    </Form.Item>
                    <Form.Item
                        {...formItemLayout}
                        label="确认密码"
                        name="confirm_password"
                        hasFeedback
                        rules={[
                            {
                                required: true,
                                message: "请确认新密码",
                            },
                            {
                                min: 8,
                                message: "密码不能少于8个字符",
                            },
                            {
                                pattern: /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[a-zA-Z\d]{8,}$/,
                                message: "密码必须包含大小写字母和数字",
                            },
                            ({ getFieldValue }) => ({
                                validator(_, value) {
                                    if (value && value !== getFieldValue("password")) {
                                        return Promise.reject(new Error("两次输入的密码不一致!"));
                                    }
                                    return Promise.resolve();
                                },
                            }),
                        ]}
                    >
                        <Input />
                    </Form.Item>
                    <Form.Item
                        {...formItemLayout}
                        style={{ textAlign: 'center' }}
                    >
                        <Button
                            type="primary"
                            htmlType="submit"
                        >
                            确 认
                        </Button>
                    </Form.Item>
                </Form>
            </Content>
        );
    }
}

export default PasswordManager;
