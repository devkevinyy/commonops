import React, { Component } from 'react';
import {Layout, Menu, Modal, Form, Input, Rate, message, Popover, Dropdown} from 'antd';
import {
    MenuFoldOutlined,
    MenuUnfoldOutlined,
    UserOutlined,
    MessageOutlined
  } from '@ant-design/icons';
import { Link } from 'react-router-dom';
import { withRouter } from 'react-router-dom';
import Menus from "../menu";
import {postUserFeedback} from "../api/user";

const { Header } = Layout;
const { TextArea } = Input;


class OpsHeader extends Component {

    constructor(props) {
        super(props);
        this.showFeedbackModal = this.showFeedbackModal.bind(this);
        this.handleOk = this.handleOk.bind(this);
        this.handleCancel = this.handleCancel.bind(this);
        this.formRef = React.createRef();
        this.state= {
            feedbackModalVisible: false,
        }
    }

    showFeedbackModal() {
        this.setState({feedbackModalVisible: true})
    }

    handleOk() {
        this.formRef.current.validateFields().then((values) => {
            postUserFeedback(values).then((res)=>{
                if(res.code===0) {
                    this.setState({feedbackModalVisible: false});
                    message.success("已经收到您的反馈！");
                } else {
                    message.error(res.msg);
                }
            }).catch((err)=>{
                message.error(err.toLocaleString());
            })
        });
    }

    handleCancel() {
        this.setState({feedbackModalVisible: false})
    }

    render() {
        const formItemLayout = {
            labelCol: {span: 5},
            wrapperCol: {span: 17},
        };
        let menuContent = (
            <Menu>
                {
                    Menus.topMenus.map((menu) => {
                        return (
                            menu.subMenus.map((subMenu) => {
                                return (
                                    <Menu.Item key={subMenu.title}>
                                        {<Link to={subMenu.hasOwnProperty("route")?subMenu.route:""}>{subMenu.title}</Link>}
                                    </Menu.Item>
                                )
                            })
                        )
                    })
                }
            </Menu>
        );

        return (
            <Header style={{ background: "#fff", padding: 0 }}>
                <div style={{ background: "#001529" }}>
                    <span
                        style={{
                            color: "#fff",
                            paddingLeft: "2%",
                            fontSize: "1.4em",
                        }}
                    >
                        {this.props.menuCollapsed ? (
                            <MenuUnfoldOutlined
                                className="trigger"
                                onClick={this.props.handleSiderMenu}
                                style={{ cursor: "pointer" }}
                            />
                        ) : (
                            <MenuFoldOutlined
                                className="trigger"
                                onClick={this.props.handleSiderMenu}
                                style={{ cursor: "pointer" }}
                            />
                        )}
                    </span>

                    <span
                        style={{
                            color: "#fff",
                            fontSize: "1.4em",
                            float: "right",
                            cursor: "pointer",
                            width: "30px",
                            textAlign: "center",
                            marginLeft: "20px",
                            marginRight: "30px",
                        }}
                    >
                        <Dropdown
                            overlay={menuContent}
                            overlayStyle={{ width: "150px" }}
                        >
                            <UserOutlined />
                        </Dropdown>
                    </span>

                    <span
                        style={{
                            color: "#fff",
                            fontSize: "1.4em",
                            float: "right",
                            cursor: "pointer",
                            width: "30px",
                            textAlign: "center",
                        }}
                        onClick={this.showFeedbackModal}
                    >
                        <Popover content="提交需求反馈">
                            <MessageOutlined />
                        </Popover>
                    </span>

                    <Modal
                        title="「Bug & 需求」反馈"
                        destroyOnClose={true}
                        visible={this.state.feedbackModalVisible}
                        onOk={this.handleOk}
                        onCancel={this.handleCancel}
                        width={700}
                    >
                        <Form {...formItemLayout} ref={this.formRef}>
                            <Form.Item
                                label="反馈内容"
                                name="advice"
                                rules={[
                                    {
                                        required: true,
                                        message: "内容不能为空",
                                    },
                                ]}
                            >
                                <TextArea rows={4} />
                            </Form.Item>
                            <Form.Item
                                label="当前版本满意度"
                                name="score"
                                rules={[
                                    {
                                        required: true,
                                        message: "请对当前版本进行打分",
                                    },
                                ]}
                            >
                                <Rate
                                    tooltips={[
                                        "难用",
                                        "不好用",
                                        "体验一般",
                                        "体验良好",
                                        "非常棒",
                                    ]}
                                />
                            </Form.Item>
                        </Form>
                    </Modal>
                </div>
            </Header>
        );
    }
}

export default withRouter(OpsHeader);
