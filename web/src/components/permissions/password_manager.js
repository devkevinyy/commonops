import React, { Component } from 'react';
import {Layout, Input, Form, Button, message} from 'antd';
import {postUpdatePassword} from "../../api/user";

const { Content } = Layout;


class PasswordManager extends Component {

    constructor(props) {
        super(props);
        this.setState({

        });
    }

    compareToFirstPassword = (rule, value, callback) => {
        const form = this.props.form;
        if (value && value !== form.getFieldValue('password')) {
            callback('两次密码不一致!');
        } else {
            callback();
        }
    };

    updatePassword = (e) => {
        e.preventDefault();
        this.props.form.validateFields((err, values) => {
            if (!err) {
                let reqData = {
                    "password": values['password'],
                    "confirm_password": values['confirm_password']
                };
                postUpdatePassword(reqData).then((res)=>{
                    if(res.code === 0){
                        message.success("密码修改成功!");
                    }
                }).catch((err)=>{
                    console.log(err)
                });
            }
        });
    };

    render() {
        const { getFieldDecorator } = this.props.form;
        const formItemLayout = {
            labelCol: { span: 2 },
            wrapperCol: { span: 6 },
        };
        return (
            <Content style={{
                background: '#fff', padding: 20, margin: 0, height: "100%",
            }}>
                <Form onSubmit={this.updatePassword}>
                    <Form.Item label="新密码" hasFeedback {...formItemLayout}>
                        {getFieldDecorator('password', {
                            rules: [
                                {
                                    required: true,
                                    message: '请输入新密码',
                                }
                            ],
                        })(<Input />)}
                    </Form.Item>
                    <Form.Item label="确认密码" hasFeedback {...formItemLayout}>
                        {getFieldDecorator('confirm_password', {
                            rules: [
                                {
                                    required: true,
                                    message: '请确认新密码',
                                },
                                {
                                    validator: this.compareToFirstPassword
                                }
                            ],
                        })(<Input />)}
                    </Form.Item>
                    <Form.Item {...formItemLayout} style={{ textAlign: 'center' }}>
                        <Button type="primary" htmlType="submit">
                            确 认
                        </Button>
                    </Form.Item>
                </Form>
            </Content>
        );
    }
    
}

export default PasswordManager;