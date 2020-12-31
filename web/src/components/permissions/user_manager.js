import React, {Component} from 'react';
import {
    Layout,
    Table,
    Form,
    message,
    Button,
    Tag,
    Typography,
    Row,
    Col,
    Divider, Modal, Input
} from 'antd';
import {getUsersList, postUserCreate, putUserCreate} from "../../api/user";
import OpsBreadcrumbPath from "../breadcrumb_path";

const { Content } = Layout;
const { Text } = Typography;


let columnStyle = {
    overFlow : "hidden",
    textOverflow: "ellipsis",
    display: "-webkit-box",
};

class UserModal extends Component {

    constructor(props) {
        super(props);
        this.state = {};
    }

    render() {
        const { getFieldDecorator } = this.props.form;
        const formItemLayout = {
            labelCol: { span: 7 },
            wrapperCol: { span: 14 },
        };
        return (
            <Modal
                title="新增用户"
                destroyOnClose={true}
                visible={this.props.userModalVisible}
                onOk={this.props.handleAddUserSubmit}
                onCancel={this.props.handleAddUserCancel}
                okText="确认"
                cancelText="取消"
                centered={true}
                width={600}
            >
                <Form>
                    <Form.Item label="邮箱" {...formItemLayout}>
                        {getFieldDecorator('email', {
                            rules: [{ required: true, message: '注册邮箱不能为空！' }],
                        })(
                            <Input/>
                        )}
                    </Form.Item>
                    <Form.Item label="用户名" {...formItemLayout}>
                        {getFieldDecorator('username', {
                            rules: [{ required: true, message: '用户名不能为空！' }],
                        })(
                            <Input/>
                        )}
                    </Form.Item>
                    <Form.Item label="初始密码" {...formItemLayout}>
                        {getFieldDecorator('password', {
                            rules: [{ required: true, message: '用户名不能为空！' }],
                        })(
                            <Input/>
                        )}
                    </Form.Item>
                    <Form.Item label="职位" {...formItemLayout}>
                        {getFieldDecorator('position', {
                            rules: [{ required: true, message: '职位不能为空！' }],
                        })(
                            <Input/>
                        )}
                    </Form.Item>
                </Form>
            </Modal>
        )
    }
}

class UserManager extends Component {

    constructor(props){
        super(props);
        this.onShowSizeChange = this.onShowSizeChange.bind(this);
        this.handleAdd = this.handleAdd.bind(this);
        this.handleAddUserSubmit = this.handleAddUserSubmit.bind(this);
        this.handleAddUserCancel = this.handleAddUserCancel.bind(this);
        this.state = {
            columns: [
                {
                    title: 'ID', dataIndex: 'ID', key: 'ID', className: 'small_font',
                    render: value => {
                        return (
                            <Text ellipsis={true}>{value}</Text>
                        )
                    }
                },
                {
                    title: '用户名', dataIndex: 'username', key: 'username', className: 'small_font',
                    render: value => {
                        return (
                            <Text ellipsis={true}>{value}</Text>
                        )
                    }
                },
                {
                    title: '邮箱', dataIndex: 'email', key: 'email', className: "small_font "+{columnStyle},
                    render: value => {
                        return (
                            <Text ellipsis={true}>{value}</Text>
                        )
                    }
                },
                {
                    title: '职位', dataIndex: 'position', key: 'position', className: "small_font "+{columnStyle},
                    render: value => {
                        return (
                            <Text ellipsis={true}>{value}</Text>
                        )
                    }
                },
                {
                    title: '状态', dataIndex: 'active', key: 'active', className: "small_font "+{columnStyle},  width: 100,
                    render: (value) => {
                        return value === true ? <Text ellipsis={true}><Tag color="geekblue">正常</Tag></Text> : <Text ellipsis={true}><Tag color="red">禁用</Tag></Text>
                    }
                },
                {
                    title: '操作',
                    key: 'operation',
                    fixed: 'right',
                    width: 140,
                    align: 'center',
                    className: 'small_font',
                    render: (text, record) => {
                        return (
                            <div>
                                <Button type="primary" size="small" onClick={this.updateUserActiveStatus.bind(this, record.email, true)}>启用</Button>
                                <Divider type="vertical" />
                                <Button type="danger" size="small" onClick={this.updateUserActiveStatus.bind(this, record.email, false)}>禁用</Button>
                            </div>
                        )
                    },
                },
            ],
            tableData: [],
            tableLoading: false,
            add_new_account_visible: false,
            current_data_id: null,
            current_user: {},
            detailDrawerVisible: false,
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
            userModalVisible: false,
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

    refreshTableData = () => {
        this.setState({tableLoading: true});
        getUsersList(this.state.pagination.page, this.state.pagination.pageSize).then((res)=>{
            const pagination = this.state.pagination;
            pagination.total = parseInt(res.data.total);
            pagination.page = parseInt(res.data.page);
            pagination.showTotal(parseInt(res.data.total));
            this.setState({
                pagination
            });
            let data = res['data']['users'];
            let tableData = [];
            for (let i = 0; i < data.length; i++) {
                tableData.push({
                    ID: data[i]['ID'],
                    username: data[i]['username'],
                    email: data[i]['email'],
                    position: data[i]['position'],
                    active: data[i]['active'],
                });
            }
            this.setState({tableData: tableData, tableLoading: false});
        }).catch((err)=>{
            console.log(err)
        });
    };

    componentDidMount() {
        this.refreshTableData();
    }

    handleAdd() {
        this.setState({userModalVisible: true});
    }

    handleAddUserSubmit() {
        this.addUserFormRef.props.form.validateFields((err, values) => {
            if (!err) {
                postUserCreate({
                    ...values,
                }).then((res)=>{
                    if(res.code===0){
                        this.setState({userModalVisible: false});
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

    handleAddUserCancel() {
        this.setState({userModalVisible: false});
    }

    updateUserActiveStatus(email, activeStatus) {
        putUserCreate({
            'email': email,
            'active': activeStatus
        }).then((res)=>{
            if(res.code===0){
                message.success("操作成功！");
                this.refreshTableData();
            } else {
                message.error(res.msg);
            }
        }).catch((err)=>{
            message.error(err.toLocaleString());
        })
    }

    render() {
        return (
            <Content style={{
                background: '#fff', padding: 20, margin: 0, height: "100%",
            }}>
                <OpsBreadcrumbPath pathData={["权限管理", "用户管理", "用户列表"]} />
                <Row style={{ padding: "0px 0px 10px 0px" }}>
                    <Col span={2} className="col-span">
                        <Button style={{ width: "100%" }} icon="plus-circle" onClick={this.handleAdd}>新 增</Button>
                    </Col>
                </Row>
                <UserModal
                    wrappedComponentRef={(form) => this.addUserFormRef = form}
                    userModalVisible={this.state.userModalVisible}
                    handleAddUserSubmit={this.handleAddUserSubmit}
                    handleAddUserCancel={this.handleAddUserCancel}
                />
                <Table columns={this.state.columns} dataSource={this.state.tableData} scroll={{ x: 'max-content' }} pagination={this.state.pagination} loading={this.state.tableLoading} rowClassName="fixedHeight" size="small"/>
            </Content>
        );
    }
}

export default UserManager;
