import React, { Component, Fragment } from "react";
import {
    Button,
    Layout,
    Table,
    Modal,
    Form,
    Input,
    Divider,
    Transfer,
} from "antd";
import { message, Row, Col, Spin, Tabs } from "antd";
import {
    deleteRole,
    getRoleAuthLinks,
    getRoleResourceList,
    getRolesList,
    getRoleUserList,
    postAddRole,
    postRoleAuthLinks,
    postRoleResourcesList,
    postRoleUserList,
    putUpdateRole,
} from "../../api/role";
import OpsBreadcrumbPath from "../breadcrumb_path";

const { Content } = Layout;
const confirm = Modal.confirm;
const { TabPane } = Tabs;

let columnStyle = {
    overFlow: "hidden",
    textOverflow: "ellipsis",
    display: "-webkit-box",
};

class RoleModal extends Component {
    render() {
        const formItemLayout = {
            labelCol: { span: 6 },
            wrapperCol: { span: 14 },
        };
        const { roleActionType } = this.props;
        let title = "新建角色";
        if (roleActionType === "add") {
            title = "新建角色";
        } else {
            title = "编辑角色";
        }
        return (
            <Fragment>
                <Modal
                    title={title}
                    visible={this.props.roleModalVisible}
                    closable={true}
                    onOk={
                        roleActionType === "add"
                            ? this.props.handleAddRole
                            : this.props.handleEditRole
                    }
                    onCancel={this.props.handleCancelAddRole}
                    centered={true}
                    okText="确认"
                    cancelText="取消"
                >
                    <Form ref={this.props.formRef}>
                        <Form.Item
                            {...formItemLayout}
                            label="角色名称"
                            name="roleName"
                            rules={[
                                { required: true, message: "请输入角色名称" },
                            ]}
                        >
                            <Input placeholder="请输入角色名称" />
                        </Form.Item>
                        <Form.Item
                            {...formItemLayout}
                            label="角色描述"
                            name="description"
                            rules={[
                                {
                                    required: true,
                                    message: "请输入角色描述",
                                },
                            ]}
                        >
                            <Input placeholder="请输入角色描述" />
                        </Form.Item>
                    </Form>
                </Modal>
            </Fragment>
        );
    }
}

class UserManagerModal extends Component {
    constructor(props) {
        super(props);
        this.state = {
            roleId: 0,
            inGroupData: [],
            allGroupData: [],
            transferLodding: false,
        };
    }

    componentWillMount() {
        const { roleId } = this.props;
        this.loadRoleUsersData(roleId);
    }

    loadRoleUsersData(roleId) {
        this.setState({
            roleId: roleId,
            transferLodding: true,
        });
        getRoleUserList(roleId)
            .then((res) => {
                if (res.code === 0) {
                    let inGroupData = [];
                    let allGroupData = [];
                    res.data.all.forEach(function(item) {
                        allGroupData.push({
                            key: String(item.ID),
                            id: String(item.ID),
                            username: item.username,
                            position: item.position,
                        });
                    });
                    res.data.in.forEach(function(item) {
                        inGroupData.push(String(item.ID));
                    });
                    this.setState({
                        inGroupData: inGroupData,
                        allGroupData: allGroupData,
                        transferLodding: false,
                    });
                } else {
                    message.error(res.msg);
                }
            })
            .catch((err) => {
                message.error(err.toLocaleString());
            });
    }

    filterOption = (inputValue, option) =>
        option.username.indexOf(inputValue) > -1;

    handleTransferChange(targetKeys, direction, moveKeys) {
        this.setState({
            inGroupData: targetKeys,
        });
    }

    handleConfirmUserManager() {
        this.setState({ transferLodding: true });
        postRoleUserList({
            roleId: parseInt(this.state.roleId),
            userIdList: this.state.inGroupData,
        })
            .then((res) => {
                if (res.code === 0) {
                    this.setState({ transferLodding: false });
                    this.props.hideModal();
                    message.success("操作成功");
                } else {
                    message.error(res.msg);
                }
            })
            .catch((err) => {
                message.error(err.toLocaleString());
            });
    }

    render() {
        return (
            <Modal
                title={"管理用户角色 - " + this.props.currentRoleName}
                destroyOnClose="true"
                visible={this.props.userManagerModalVisible}
                onOk={this.handleConfirmUserManager.bind(this)}
                onCancel={this.props.handleCancelUserManager}
                okText="确认"
                cancelText="取消"
                centered={true}
                width={680}
                style={{ textAlign: "center" }}
            >
                <Spin spinning={this.state.transferLodding}>
                    <Transfer
                        dataSource={this.state.allGroupData}
                        showSearch
                        locale={{
                            itemUnit: "人",
                            itemsUnit: "人",
                            searchPlaceholder: "使用姓名搜索",
                        }}
                        filterOption={this.filterOption}
                        targetKeys={this.state.inGroupData}
                        onChange={this.handleTransferChange.bind(this)}
                        render={(item) => item.username + "-" + item.position}
                        listStyle={{
                            width: 250,
                            height: "60vh",
                            textAlign: "left",
                        }}
                        operations={["加入角色组", "移出角色组"]}
                    />
                </Spin>
            </Modal>
        );
    }
}

class ResourceManagerModal extends Component {
    constructor(props) {
        super(props);
        this.handleConfirmRoleResourceManager = this.handleConfirmRoleResourceManager.bind(
            this,
        );
        this.state = {
            roleId: 0,
            tabSpinning: false,
            ecsInGroupData: [],
            ecsAllGroupData: [],
            rdsInGroupData: [],
            rdsAllGroupData: [],
            kvInGroupData: [],
            kvAllGroupData: [],
            slbInGroupData: [],
            slbAllGroupData: [],
            otherAllGroupData: [],
            otherInGroupData: [],
        };
    }

    componentWillMount() {
        const { roleId } = this.props;
        this.setState({ roleId });
        this.loadRoleEcsData(roleId);
        this.loadRoleRdsData(roleId);
        this.loadRoleKvData(roleId);
        this.loadRoleSlbData(roleId);
        this.loadRoleOtherResData(roleId);
    }

    loadRoleEcsData(roleId) {
        this.setState({
            ecsTransferLodding: true,
        });
        let ecsAllGroupData = [];
        let ecsInGroupData = [];
        getRoleResourceList(roleId, "ecs")
            .then((res) => {
                if (res.code === 0) {
                    res.data.all.forEach(function(item) {
                        ecsAllGroupData.push({
                            key: String(item.ID),
                            id: String(item.ID),
                            instance_name: item.InstanceName,
                            ip:
                                item.InnerIpAddress +
                                " " +
                                item.PublicIpAddress,
                        });
                    });
                    res.data.in.forEach(function(item) {
                        ecsInGroupData.push(String(item.ID));
                    });
                    this.setState({
                        ecsAllGroupData: ecsAllGroupData,
                        ecsInGroupData: ecsInGroupData,
                    });
                } else {
                    message.error(res.msg);
                }
                this.setState({ecsTransferLodding: false})
            })
    }

    loadRoleRdsData(roleId) {
        this.setState({
            rdsTransferLodding: true,
        });
        let rdsAllGroupData = [];
        let rdsInGroupData = [];
        getRoleResourceList(roleId, "rds")
            .then((res) => {
                if (res.code === 0) {
                    res.data.all.forEach(function(item) {
                        rdsAllGroupData.push({
                            key: String(item.ID),
                            id: String(item.ID),
                            db_instance_description: item.DBInstanceDescription,
                            db_instance_id: item.DBInstanceId,
                        });
                    });
                    res.data.in.forEach(function(item) {
                        rdsInGroupData.push(String(item.ID));
                    });
                    this.setState({
                        rdsAllGroupData: rdsAllGroupData,
                        rdsInGroupData: rdsInGroupData,
                    });
                } else {
                    message.error(res.msg);
                }
                this.setState({rdsTransferLodding: false})
            })
    }

    loadRoleKvData(roleId) {
        this.setState({
            kvTransferLodding: true,
        });
        let kvAllGroupData = [];
        let kvInGroupData = [];
        getRoleResourceList(roleId, "kv")
            .then((res) => {
                if (res.code === 0) {
                    res.data.all.forEach(function(item) {
                        kvAllGroupData.push({
                            key: String(item.ID),
                            id: String(item.ID),
                            instance_name: item.InstanceName,
                            connection_domain: item.ConnectionDomain,
                        });
                    });
                    res.data.in.forEach(function(item) {
                        kvInGroupData.push(String(item.ID));
                    });
                    this.setState({
                        kvAllGroupData,
                        kvInGroupData,
                    });
                } else {
                    message.error(res.msg);
                }
                this.setState({kvTransferLodding: false})
            })
    }

    loadRoleSlbData(roleId) {
        this.setState({
            slbTransferLodding: true,
        });
        let slbAllGroupData = [];
        let slbInGroupData = [];
        getRoleResourceList(roleId, "slb")
            .then((res) => {
                if (res.code === 0) {
                    res.data.all.forEach(function(item) {
                        slbAllGroupData.push({
                            key: String(item.ID),
                            id: String(item.ID),
                            load_balance_name: item.LoadBalancerName,
                            address: item.Address,
                        });
                    });
                    res.data.in.forEach(function(item) {
                        slbInGroupData.push(String(item.ID));
                    });
                    this.setState({
                        slbAllGroupData,
                        slbInGroupData,
                    });
                } else {
                    message.error(res.msg);
                }
                this.setState({slbTransferLodding: false})
            })
    }

    loadRoleOtherResData(roleId) {
        this.setState({
            otherResTransferLodding: true,
        });
        let otherAllGroupData = [];
        let otherInGroupData = [];
        getRoleResourceList(roleId, "other")
            .then((res) => {
                if (res.code === 0) {
                    res.data.all.forEach(function(item) {
                        otherAllGroupData.push({
                            key: String(item.ID),
                            id: String(item.ID),
                            instance_name: item.InstanceName,
                            res_type: item.ResType,
                        });
                    });
                    res.data.in.forEach(function(item) {
                        otherInGroupData.push(String(item.ID));
                    });
                    this.setState({
                        otherAllGroupData: otherAllGroupData,
                        otherInGroupData: otherInGroupData,
                    });
                } else {
                    message.error(res.msg);
                }
                this.setState({otherResTransferLodding: false})
            })
    }

    ecsFilterOption = (inputValue, option) =>
        option.instance_name.indexOf(inputValue) > -1 ||
        option.ip.indexOf(inputValue) > -1;

    rdsFilterOption = (inputValue, option) =>
        option.db_instance_description.indexOf(inputValue) > -1 ||
        option.db_instance_id.indexOf(inputValue) > -1;

    kvFilterOption = (inputValue, option) =>
        option.instance_name.indexOf(inputValue) > -1 ||
        option.connection_domain.indexOf(inputValue) > -1;

    slbFilterOption = (inputValue, option) =>
        option.load_balance_name.indexOf(inputValue) > -1 ||
        option.address.indexOf(inputValue) > -1;

    otherResFilterOption = (inputValue, option) =>
        option.instance_name.indexOf(inputValue) > -1 ||
        option.address.indexOf(inputValue) > -1;

    handleResTransferChange(resType, targetKeys) {
        switch (resType) {
            case "ecs":
                this.setState({
                    ecsInGroupData: targetKeys,
                });
                break;
            case "rds":
                this.setState({
                    rdsInGroupData: targetKeys,
                });
                break;
            case "kv":
                this.setState({
                    kvInGroupData: targetKeys,
                });
                break;
            case "slb":
                this.setState({
                    slbInGroupData: targetKeys,
                });
                break;
            case "other":
                this.setState({
                    otherInGroupData: targetKeys,
                });
                break;
            default:
                break;
        }
    }

    handleTabClick = (e) => {
        // if(e==="ecs"){
        //     this.loadRoleEcsData()
        // } else if(e==="rds") {
        //     this.loadRoleRdsData()
        // } else if(e==="kv") {
        //     this.loadRoleKvData()
        // } else {
        //     this.loadRoleSlbData()
        // }
    };

    handleConfirmRoleResourceManager() {
        this.setState({ tabSpinning: true });
        postRoleResourcesList({
            roleId: parseInt(this.state.roleId),
            ecsIdList: this.state.ecsInGroupData,
            rdsIdList: this.state.rdsInGroupData,
            kvIdList: this.state.kvInGroupData,
            slbIdList: this.state.slbInGroupData,
            otherIdList: this.state.otherInGroupData,
        })
            .then((res) => {
                if (res.code === 0) {
                    this.setState({ tabSpinning: false });
                    this.props.hideModal();
                    message.success("操作成功");
                } else {
                    message.error(res.msg);
                }
            })
            .catch((err) => {
                message.error(err.toLocaleString());
            });
    }

    render() {
        return (
            <Modal
                title={"管理角色组关联资源 - " + this.props.currentRoleName}
                destroyOnClose="true"
                visible={this.props.resourceManagerModalVisible}
                onOk={this.handleConfirmRoleResourceManager}
                onCancel={this.props.handleCancelResourceManager}
                okText="确认"
                cancelText="取消"
                centered={true}
                width={1000}
                style={{ textAlign: "center" }}
            >
                <Spin spinning={this.state.tabSpinning} tip="请求处理中...">
                    <Tabs
                        defaultActiveKey="ecs"
                        tabPosition="left"
                        onTabClick={this.handleTabClick.bind(this)}
                    >
                        <TabPane tab="云服务器" key="ecs">
                            <Spin spinning={this.state.ecsTransferLodding}>
                                <Transfer
                                    showSearch
                                    locale={{
                                        itemUnit: "台",
                                        itemsUnit: "台",
                                        searchPlaceholder:
                                            "使用ip或机器名称搜索",
                                    }}
                                    filterOption={this.ecsFilterOption}
                                    dataSource={this.state.ecsAllGroupData}
                                    targetKeys={this.state.ecsInGroupData}
                                    onChange={this.handleResTransferChange.bind(
                                        this,
                                        "ecs",
                                    )}
                                    render={(item) =>
                                        item.ip + " " + item.instance_name
                                    }
                                    listStyle={{
                                        width: 350,
                                        height: "60vh",
                                        textAlign: "left",
                                    }}
                                    operations={["加入角色组", "移出角色组"]}
                                />
                            </Spin>
                        </TabPane>
                        <TabPane tab="云数据库" key="rds">
                            <Spin spinning={this.state.rdsTransferLodding}>
                                <Transfer
                                    showSearch
                                    locale={{
                                        itemUnit: "台",
                                        itemsUnit: "台",
                                        searchPlaceholder:
                                            "使用连接串或名称搜索",
                                    }}
                                    filterOption={this.rdsFilterOption}
                                    dataSource={this.state.rdsAllGroupData}
                                    targetKeys={this.state.rdsInGroupData}
                                    onChange={this.handleResTransferChange.bind(
                                        this,
                                        "rds",
                                    )}
                                    render={(item) =>
                                        item.db_instance_description +
                                        " " +
                                        item.db_instance_id
                                    }
                                    listStyle={{
                                        width: 350,
                                        height: "60vh",
                                        textAlign: "left",
                                    }}
                                    operations={["加入角色组", "移出角色组"]}
                                />
                            </Spin>
                        </TabPane>
                        <TabPane tab="KVStore" key="kv">
                            <Spin spinning={this.state.kvTransferLodding}>
                                <Transfer
                                    showSearch
                                    locale={{
                                        itemUnit: "台",
                                        itemsUnit: "台",
                                        searchPlaceholder:
                                            "使用连接串或名称搜索",
                                    }}
                                    filterOption={this.kvFilterOption}
                                    dataSource={this.state.kvAllGroupData}
                                    targetKeys={this.state.kvInGroupData}
                                    onChange={this.handleResTransferChange.bind(
                                        this,
                                        "kv",
                                    )}
                                    render={(item) =>
                                        item.instance_name +
                                        " " +
                                        item.connection_domain
                                    }
                                    listStyle={{
                                        width: 350,
                                        height: "60vh",
                                        textAlign: "left",
                                    }}
                                    operations={["加入角色组", "移出角色组"]}
                                />
                            </Spin>
                        </TabPane>
                        <TabPane tab="SLB" key="slb">
                            <Spin spinning={this.state.slbTransferLodding}>
                                <Transfer
                                    showSearch
                                    locale={{
                                        itemUnit: "个",
                                        itemsUnit: "个",
                                        searchPlaceholder: "使用ip或名称搜索",
                                    }}
                                    filterOption={this.slbFilterOption}
                                    dataSource={this.state.slbAllGroupData}
                                    targetKeys={this.state.slbInGroupData}
                                    onChange={this.handleResTransferChange.bind(
                                        this,
                                        "slb",
                                    )}
                                    render={(item) =>
                                        item.load_balance_name +
                                        " " +
                                        item.address
                                    }
                                    listStyle={{
                                        width: 350,
                                        height: "60vh",
                                        textAlign: "left",
                                    }}
                                    operations={["加入角色组", "移出角色组"]}
                                />
                            </Spin>
                        </TabPane>
                        <TabPane tab="其它资源" key="other">
                            <Spin spinning={this.state.otherResTransferLodding}>
                                <Transfer
                                    showSearch
                                    locale={{
                                        itemUnit: "个",
                                        itemsUnit: "个",
                                        searchPlaceholder: "使用实例名称搜索",
                                    }}
                                    filterOption={this.otherResFilterOption}
                                    dataSource={this.state.otherAllGroupData}
                                    targetKeys={this.state.otherInGroupData}
                                    onChange={this.handleResTransferChange.bind(
                                        this,
                                        "other",
                                    )}
                                    render={(item) =>
                                        item.res_type + " " + item.instance_name
                                    }
                                    listStyle={{
                                        width: 350,
                                        height: "60vh",
                                        textAlign: "left",
                                    }}
                                    operations={["加入角色组", "移出角色组"]}
                                />
                            </Spin>
                        </TabPane>
                    </Tabs>
                </Spin>
            </Modal>
        );
    }
}

class AuthLinkManagerModal extends Component {
    constructor(props) {
        super(props);
        this.state = {
            roleId: 0,
            inGroupData: [],
            allGroupData: [],
            transferLodding: false,
        };
    }

    componentWillMount() {
        const { roleId } = this.props;
        this.loadRoleAuthLinksData(roleId);
    }

    loadRoleAuthLinksData(roleId) {
        this.setState({
            roleId: roleId,
            transferLodding: true,
        });
        getRoleAuthLinks(roleId)
            .then((res) => {
                if (res.code === 0) {
                    let inGroupData = [];
                    let allGroupData = [];
                    res.data.all.forEach(function(item) {
                        allGroupData.push({
                            key: String(item.Id),
                            id: String(item.Id),
                            name: item.name,
                            url_path: item.urlPath,
                            auth_type: item.authType,
                        });
                    });
                    res.data.in.forEach(function(item) {
                        inGroupData.push(String(item.Id));
                    });
                    this.setState({
                        inGroupData: inGroupData,
                        allGroupData: allGroupData,
                    });
                } else {
                    message.error(res.msg);
                }
                this.setState({transferLodding: false})
            })
    }

    filterOption = (inputValue, option) =>
        option.name.indexOf(inputValue) > -1 ||
        option.url_path.indexOf(inputValue) > -1;

    handleTransferChange(targetKeys, direction, moveKeys) {
        this.setState({
            inGroupData: targetKeys,
        });
    }

    handleConfirmUserManager() {
        this.setState({ transferLodding: true });
        postRoleAuthLinks({
            roleId: parseInt(this.state.roleId),
            authLinkIdList: this.state.inGroupData,
        })
            .then((res) => {
                if (res.code === 0) {
                    this.setState({ transferLodding: false });
                    this.props.hideModal();
                    message.success("操作成功");
                } else {
                    message.error(res.msg);
                }
            })
    }

    render() {
        return (
            <Modal
                title={"管理角色对应权限 - " + this.props.currentRoleName}
                destroyOnClose="true"
                visible={this.props.authLinkManagerModalVisible}
                onOk={this.handleConfirmUserManager.bind(this)}
                onCancel={this.props.handleCancelAuthLinkManager}
                okText="确认"
                cancelText="取消"
                centered={true}
                width={800}
                style={{ textAlign: "center" }}
            >
                <Spin spinning={this.state.transferLodding}>
                    <Transfer
                        dataSource={this.state.allGroupData}
                        showSearch
                        locale={{
                            itemUnit: "条",
                            itemsUnit: "条",
                            searchPlaceholder: "使用名称或路径搜索",
                        }}
                        filterOption={this.filterOption}
                        targetKeys={this.state.inGroupData}
                        onChange={this.handleTransferChange.bind(this)}
                        render={(item) => item.auth_type + ":" + item.name + "「" + item.url_path+"」"}
                        listStyle={{
                            width: 300,
                            height: "60vh",
                            textAlign: "left",
                        }}
                        operations={["加入角色组", "移出角色组"]}
                    />
                </Spin>
            </Modal>
        );
    }
}

class RolesManager extends Component {
    constructor(props) {
        super(props);
        this.handleAddRole = this.handleAddRole.bind(this);
        this.handleEditRole = this.handleEditRole.bind(this);
        this.handleCancelAddRole = this.handleCancelAddRole.bind(this);
        this.createRole = this.createRole.bind(this);
        this.roleEdit = this.roleEdit.bind(this);
        this.onShowSizeChange = this.onShowSizeChange.bind(this);
        this.formRef = React.createRef();
        this.state = {
            columns: [
                {
                    title: "ID",
                    dataIndex: "Id",
                    key: "Id",
                    width: 60,
                },
                {
                    title: "角色名",
                    dataIndex: "name",
                    key: "name",
                    width: 160,
                },
                {
                    title: "角色描述",
                    dataIndex: "description",
                    key: "description",
                    className: { columnStyle },
                    width: 160,
                },
                {},
                {
                    title: "操作",
                    key: "operation",
                    fixed: "right",
                    width: 430,
                    align: "center",
                    render: (text, record) => {
                        return (
                            <div>
                                <Button
                                    size="small"
                                    onClick={this.managerUser.bind(
                                        this,
                                        record,
                                    )}
                                >
                                    管理用户
                                </Button>
                                <Divider type="vertical" />
                                <Button
                                    size="small"
                                    onClick={this.managerAuthLink.bind(
                                        this,
                                        record,
                                    )}
                                >
                                    管理权限
                                </Button>
                                <Divider type="vertical" />
                                <Button
                                    size="small"
                                    onClick={this.managerResource.bind(
                                        this,
                                        record,
                                    )}
                                >
                                    管理资源
                                </Button>
                                <Divider type="vertical" />
                                <Button
                                    type="info"
                                    size="small"
                                    onClick={this.roleEdit.bind(this, record)}
                                >
                                    编辑
                                </Button>
                                <Divider type="vertical" />
                                <Button
                                    type="danger"
                                    size="small"
                                    onClick={this.confirmDeleteRole.bind(
                                        this,
                                        record,
                                    )}
                                >
                                    删除
                                </Button>
                            </div>
                        );
                    },
                },
            ],
            tableData: [],
            tableLoading: false,
            roleModalVisible: false,
            roleActionType: "add",
            editRoleId: 0,
            userManagerModalVisible: false,
            resourceManagerModalVisible: false,
            authLinkManagerModalVisible: false,
            current_data_id: null,
            pagination: {
                showSizeChanger: true,
                pageSizeOptions: ["10", "20", "30", "100"],
                onShowSizeChange: (current, size) =>
                    this.onShowSizeChange(current, size),
                showQuickJumper: false,
                showTotal: (total) => `共 ${total} 条`,
                pageSize: 10,
                page: 1,
                total: 0,
                onChange: (page, pageSize) => this.changePage(page, pageSize),
            },
        };
    }

    onShowSizeChange(current, size) {
        let pagination = {
            ...this.state.pagination,
            page: 1,
            current: 1,
            pageSize: size,
        };
        this.setState(
            {
                pagination: pagination,
            },
            () => {
                this.refreshTableData();
            },
        );
    }

    managerUser = (record) => {
        this.setState({
            roleId: record.Id,
            userManagerModalVisible: true,
            currentRoleName: record.name,
        });
    };

    confirmDeleteRole = (record) => {
        let that = this;
        confirm({
            title: "危险操作提示",
            content:
                "删除该角色组时该角色组和与之绑定的所有用户的关系、所有资源的关系和权限链接的关系也会被删除，请谨慎操作！",
            okText: "确认删除",
            okType: "danger",
            cancelText: "取消",
            onOk() {
                deleteRole({
                    id: record["Id"],
                })
                    .then((res) => {
                        if (res.code === 0) {
                            message.success("删除成功");
                            that.refreshTableData();
                        } else {
                            message.error(res.msg);
                        }
                    })
                    .catch((err) => {
                        message.error(err.toLocaleString());
                    });
            },
            onCancel() {},
        });
    };

    changePage = (page, pageSize) => {
        this.setState(
            {
                pagination: {
                    ...this.state.pagination,
                    page: page,
                    current: page,
                    pageSize: pageSize,
                },
            },
            () => {
                this.refreshTableData();
            },
        );
    };

    managerResource = (record) => {
        this.setState(() => ({
            roleId: record.Id,
            resourceManagerModalVisible: true,
            currentRoleName: record.name,
        }));
    };

    managerAuthLink = (record) => {
        this.setState(() => ({
            roleId: record.Id,
            authLinkManagerModalVisible: true,
            currentRoleName: record.name,
        }));
    };

    createRole() {
        this.setState(() => ({
            roleActionType: "add",
            roleModalVisible: true,
        }));
    }

    roleEdit(record) {
        this.formRef.current.setFieldsValue({
            roleName: record.name,
            description: record.description,
        });
        this.setState(() => ({
            roleActionType: "update",
            roleModalVisible: true,
            editRoleId: record.Id,
        }));
    }

    handleAddRole() {
        this.formRef.current.validateFields().then((values) => {
            postAddRole(values)
                .then((res) => {
                    if (res.code === 0) {
                        message.success("创建成功");
                        this.setState({
                            roleModalVisible: false,
                        });
                        this.formRef.current.resetFields();
                        this.refreshTableData();
                    } else {
                        message.error(res.msg);
                    }
                })
                .catch((err) => {
                    message.error(err.toLocaleString());
                });
        });
    }

    handleEditRole() {
        this.formRef.current.validateFields().then((values) => {
            putUpdateRole({
                ...values,
                id: this.state.editRoleId,
            })
                .then((res) => {
                    if (res.code === 0) {
                        message.success("修改成功");
                        this.setState({
                            roleModalVisible: false,
                        });
                        this.formRef.current.resetFields();
                        this.refreshTableData();
                    } else {
                        message.error(res.msg);
                    }
                })
                .catch((err) => {
                    message.error(err.toLocaleString());
                });
        });
    }

    handleCancelAddRole() {
        this.setState(() => ({
            roleModalVisible: false,
            userManagerModalVisible: false,
            resourceManagerModalVisible: false,
            authLinkManagerModalVisible: false,
        }));
    }

    refreshTableData = () => {
        this.setState({ tableLoading: true });
        getRolesList(this.state.pagination.page, this.state.pagination.pageSize)
            .then((res) => {
                if(res.code===0) {
                    const pagination = this.state.pagination;
                    pagination.total = parseInt(res.data.total);
                    pagination.page = parseInt(res.data.page);
                    pagination.showTotal(parseInt(res.data.total));
                    this.setState({
                        pagination,
                    });
                    this.setState({ tableData: res["data"]["roles"]});
                } else  {
                    message.error(res.msg);
                }
                this.setState({tableLoading: false })
            })
    };

    componentDidMount() {
        this.refreshTableData();
    }

    render() {
        return (
            <Content
                style={{
                    background: "#fff",
                    padding: 20,
                    margin: 0,
                    height: "100%",
                }}
            >
                <OpsBreadcrumbPath
                    pathData={["权限管理", "角色管理", "角色列表"]}
                />
                <div style={{ padding: "0px 0px 10px 0px" }}>
                    <Row>
                        <Col>
                            <Button type="primary" onClick={this.createRole}>
                                新建角色
                            </Button>
                        </Col>
                    </Row>
                </div>
                <RoleModal
                    formRef={this.formRef}
                    roleModalVisible={this.state.roleModalVisible}
                    roleActionType={this.state.roleActionType}
                    handleAddRole={this.handleAddRole}
                    handleEditRole={this.handleEditRole}
                    handleCancelAddRole={this.handleCancelAddRole}
                />
                {this.state.userManagerModalVisible ? (
                    <UserManagerModal
                        roleId={this.state.roleId}
                        userManagerModalVisible={
                            this.state.userManagerModalVisible
                        }
                        currentRoleName={this.state.currentRoleName}
                        hideModal={() => {
                            this.setState({ userManagerModalVisible: false });
                        }}
                        handleCancelUserManager={this.handleCancelAddRole.bind(
                            this,
                        )}
                    />
                ) : (
                    ""
                )}
                {this.state.resourceManagerModalVisible ? (
                    <ResourceManagerModal
                        roleId={this.state.roleId}
                        resourceManagerModalVisible={
                            this.state.resourceManagerModalVisible
                        }
                        currentRoleName={this.state.currentRoleName}
                        hideModal={() => {
                            this.setState({
                                resourceManagerModalVisible: false,
                            });
                        }}
                        handleCancelResourceManager={this.handleCancelAddRole.bind(
                            this,
                        )}
                    />
                ) : (
                    ""
                )}
                {this.state.authLinkManagerModalVisible ? (
                    <AuthLinkManagerModal
                        roleId={this.state.roleId}
                        authLinkManagerModalVisible={
                            this.state.authLinkManagerModalVisible
                        }
                        currentRoleName={this.state.currentRoleName}
                        hideModal={() => {
                            this.setState({
                                authLinkManagerModalVisible: false,
                            });
                        }}
                        handleCancelAuthLinkManager={this.handleCancelAddRole.bind(
                            this,
                        )}
                    />
                ) : (
                    ""
                )}
                <Table
                    columns={this.state.columns}
                    dataSource={this.state.tableData}
                    scroll={{ x: "max-content" }}
                    pagination={this.state.pagination}
                    loading={this.state.tableLoading}
                    rowClassName="fixedHeight"
                    size="small"
                />
            </Content>
        );
    }
}

export default RolesManager;
