import React, { Component } from "react";
import {
    Button,
    Layout,
    message,
    Row,
    Tabs,
    Tree,
    Typography,
    Table,
    Tag,
    Popover,
    Icon,
    Card,
    Empty,
    Col,
} from "antd";
import OpsBreadcrumbPath from "../breadcrumb_path";
import {
    getUserDmsDatabaseData,
    getUserDmsInstanceData,
    getUserDmsLog,
    postDmsUserExecSQL,
} from "../../api/dms_api";
import "../../assets/css/dms.css";
import { Controlled as CodeMirror } from "react-codemirror2";
import "codemirror/lib/codemirror.css";
import "codemirror/mode/sql/sql";
import "codemirror/addon/hint/show-hint.css";
import "codemirror/addon/hint/show-hint";
import "codemirror/addon/hint/sql-hint";
import { OpsIcon } from "../../assets/Icons";

const { Content } = Layout;
const { TreeNode } = Tree;
const { Text, Paragraph } = Typography;
const { TabPane } = Tabs;

const left_panel = {
    float: "left",
    width: "20%",
    height: "100%",
    padding: "0px 0px 0px 0px",
    borderRight: "solid 2px #acc3c0",
};

const right_panel = {
    float: "right",
    width: "80%",
    height: "100%",
    padding: "0px 5px 0px 10px",
};

class DataManageContent extends Component {
    constructor(props) {
        super(props);
        this.sqlInputRef = React.createRef();
        this.onShowSizeChange = this.onShowSizeChange.bind(this);
        this.loadUserHistoryLog = this.loadUserHistoryLog.bind(this);
        this.state = {
            treeData: [],
            selectedNodeId: "",
            selectedNodeType: "",
            currentChoose: "-",
            activeKey: "1",
            sqlInput: "",
            sqlDescription: "",
            sqlExecuting: false,
            execResultPanel: <Empty description={false} />,
            tableLoading: false,
            pagination: {
                showSizeChanger: true,
                pageSizeOptions: ["10", "20", "30", "100"],
                onShowSizeChange: this.onShowSizeChange,
                showQuickJumper: false,
                showTotal: (total) => `共 ${total} 条`,
                pageSize: 10,
                page: 1,
                total: 0,
                onChange: (current) => this.changePage(current),
            },
            sqlResultLog: [],
            sqlResultLogColumns: [
                {
                    title: "ID",
                    dataIndex: "ID",
                    key: "ID",
                    align: "center",
                    width: "50px",
                },
                {
                    title: "库名",
                    dataIndex: "DatabaseName",
                    key: "DatabaseName",
                    align: "center",
                    width: "150px",
                },
                {
                    title: "执行时间",
                    dataIndex: "StartTime",
                    key: "StartTime",
                    width: "200px",
                    align: "center",
                },
                {
                    title: "SQL",
                    dataIndex: "SqlContent",
                    key: "SqlContent",
                    ellipsis: true,
                    width: "200px",
                    align: "center",
                    render: (value) => {
                        return (
                            <Popover
                                placement="topLeft"
                                content={
                                    <Paragraph copyable>{value}</Paragraph>
                                }
                                title="当前SQL"
                            >
                                {value}
                            </Popover>
                        );
                    },
                },
                // {
                //     title: '影响行数',
                //     dataIndex: 'EffectRows',
                //     key: 'EffectRows',
                //     width: '100px',
                //     align: 'center',
                // },
                {
                    title: "耗时",
                    dataIndex: "Duration",
                    key: "Duration",
                    width: "100px",
                    align: "center",
                    render: (value) => {
                        return value + " ms";
                    },
                },
                {
                    title: "执行状态",
                    dataIndex: "ExecStatus",
                    key: "ExecStatus",
                    width: "100px",
                    fixed: "right",
                    align: "center",
                    render: (value) => {
                        let content = <Tag color="#f50">失败</Tag>;
                        if (value === 1) {
                            content = <Tag color="#2db7f5">成功</Tag>;
                        }
                        return content;
                    },
                },
                {
                    title: "异常信息",
                    dataIndex: "ExceptionOutput",
                    key: "ExceptionOutput",
                    ellipsis: true,
                    fixed: "right",
                    align: "center",
                    width: "100px",
                    render: (value) => {
                        if (value !== "") {
                            return (
                                <Popover
                                    placement="topLeft"
                                    content={
                                        <Paragraph copyable>{value}</Paragraph>
                                    }
                                    title="异常详情"
                                    trigger="click"
                                >
                                    <Button type="link" size="small">
                                        点击查看
                                    </Button>
                                </Popover>
                            );
                        }
                        return "无异常";
                    },
                },
            ],
        };
    }

    componentDidMount() {
        this.loadAllInstanceData();
        this.loadUserHistoryLog();
    }

    loadAllInstanceData() {
        getUserDmsInstanceData()
            .then((res) => {
                if (res.code === 0) {
                    let instanceTreeNode = [];
                    for (let i = 0; i < res.data.length; i++) {
                        instanceTreeNode.push({
                            title: res.data[i].InstanceName,
                            key: res.data[i].InstanceId,
                            type: "instance",
                            instance_type:
                                res.data[i].InstanceType === "2"
                                    ? "mysql"
                                    : "sqlserver",
                        });
                    }
                    this.setState({
                        treeData: instanceTreeNode,
                    });
                } else {
                    message.error(res.msg);
                }
            })
            .catch((err) => {
                message.error(err.toLocaleString());
            });
    }

    loadUserHistoryLog() {
        this.setState({ tableLoading: true });
        const queryParams = {
            page: this.state.pagination.page,
            pageSize: this.state.pagination.pageSize,
        };
        getUserDmsLog(queryParams).then((res) => {
            if (res.code === 0) {
                const pagination = this.state.pagination;
                pagination.total = parseInt(res.data.total);
                pagination.page = parseInt(res.data.page);
                pagination.showTotal(parseInt(res.data.total));
                this.setState({
                    pagination: { ...pagination },
                });
                this.setState({ sqlResultLog: res.data["log"] });
            } else {
                message.error(res.msg);
            }
            this.setState({ tableLoading: false });
        });
    }

    onShowSizeChange(current, size) {
        this.setState(
            {
                pagination: {
                    ...this.state.pagination,
                    page: 1,
                    current: 1,
                    pageSize: size,
                },
            },
            () => {
                this.loadUserHistoryLog();
            },
        );
    }

    changePage = (e) => {
        this.setState(
            {
                pagination: {
                    ...this.state.pagination,
                    page: e,
                    current: e,
                },
            },
            () => {
                this.loadUserHistoryLog();
            },
        );
    };

    sqlDescriptionChange(e) {
        this.setState({ sqlDescription: e.target.value });
    }

    onLoadData = (treeNode) => {
        return new Promise((resolve) => {
            if (treeNode.props.children) {
                resolve();
                return;
            }
            if (treeNode.props.dataRef.type === "instance") {
                getUserDmsDatabaseData({
                    instanceId: treeNode.props.dataRef.key,
                })
                    .then((res) => {
                        if (res.code === 0) {
                            let instanceChildren = [];
                            for (let i = 0; i < res.data.length; i++) {
                                instanceChildren.push({
                                    title: res.data[i].DatabaseName,
                                    key: res.data[i].DatabaseId,
                                    type: "database",
                                    instance_type: res.data[i].InstanceType,
                                    isLeaf: true,
                                });
                            }
                            treeNode.props.dataRef.children = instanceChildren;
                            this.setState({
                                treeData: [...this.state.treeData],
                            });
                            resolve();
                        } else {
                            message.error(res.msg);
                        }
                    })
                    .catch((err) => {
                        message.error(err.toLocaleString());
                    });
            }
        });
    };

    renderSqlQueryResultTable(columns, data) {
        if (columns === null || data === null) {
            return (
                <Text>未查询到数据，请查看『执行历史』检查是否存在异常</Text>
            );
        }
        let titleColumns = [];
        for (let i = 0; i < columns.length; i++) {
            titleColumns.push({ title: columns[i], dataIndex: columns[i] });
        }
        return (
            <Table
                bordered
                columns={titleColumns}
                dataSource={data}
                size="small"
                scroll={{ x: "max-content" }}
                footer={(data) => {
                    let count = data.length;
                    return <span>数据行数: {count}</span>;
                }}
            />
        );
    }

    renderSqlExecResultPanel(execStatus, exception, effectResult, duration) {
        let color = "green";
        let execResult = "执行成功";
        let exceptionContent = "";
        if (execStatus === 0) {
            color = "red";
            execResult = "执行失败";
            exceptionContent = (
                <span style={{ color: "red" }}>异常信息：{exception}</span>
            );
        }
        return (
            <Card>
                <span>影响行数：{effectResult}</span> <br />
                <span style={{ color: color }}>{execResult}</span> <br />
                {exceptionContent} <br />
                <span>执行时间：{duration} ms</span>
            </Card>
        );
    }

    renderTreeNodes = (data) =>
        data.map((item) => {
            let iconType = "icondatabase";
            if (item.instance_type === "mysql") {
                iconType = "iconmysql";
            }
            if (item.instance_type === "sqlserver") {
                iconType = "iconsqlserver";
            }
            if (item.children) {
                return (
                    <TreeNode
                        icon={
                            <OpsIcon
                                style={{ fontSize: "20px", color: "#08c" }}
                                type={iconType}
                            />
                        }
                        title={item.title}
                        key={item.key}
                        dataRef={item}
                    >
                        {this.renderTreeNodes(item.children)}
                    </TreeNode>
                );
            }
            return (
                <TreeNode
                    icon={
                        <OpsIcon
                            style={{ fontSize: "20px", color: "#08c" }}
                            type={iconType}
                        />
                    }
                    key={item.key}
                    {...item}
                    dataRef={item}
                />
            );
        });

    onTreeNodeSelect = (selectedKeys, e) => {
        if (e.selectedNodes[0].dataRef.type === "instance") {
            let title = "实例 | " + e.selectedNodes[0].dataRef.title;
            this.setState({ currentChoose: title });
        }
        if (e.selectedNodes[0].dataRef.type === "database") {
            let title = "数据库 | " + e.selectedNodes[0].dataRef.title;
            this.setState({ currentChoose: title });
        }
        this.setState({
            selectedNodeId: e.selectedNodes[0].dataRef.key,
            selectedNodeType: e.selectedNodes[0].dataRef.type,
        });
    };

    submitSql() {
        if (this.state.selectedNodeType !== "database") {
            message.warn("请选择具体的数据库！");
            return;
        }
        if (this.state.sqlInput.trim() === "") {
            message.warn("请填写待执行的SQL语句！");
            return;
        }
        let reqParams = {
            selectedNodeId: this.state.selectedNodeId,
            selectedNodeType: this.state.selectedNodeType,
            sqlInput: this.state.sqlInput,
            sqlDescription: this.state.sqlDescription,
        };
        message.info("提交执行中...");
        this.setState({ sqlExecuting: true });
        postDmsUserExecSQL(reqParams).then((res) => {
            if (res.code === 0) {
                if (res.data["execStatus"] === 1) {
                    message.success("执行成功");
                } else {
                    message.error("执行失败");
                }
                if (res.data["sqlType"] === "select") {
                    let content = this.renderSqlQueryResultTable(
                        res.data["queryColumns"],
                        res.data["queryResult"],
                    );
                    this.setState({ execResultPanel: content });
                } else {
                    let content = this.renderSqlExecResultPanel(
                        res.data["execStatus"],
                        res.data["exceptionOutput"],
                        res.data["effectRows"],
                        res.data["duration"],
                    );
                    this.setState({ execResultPanel: content });
                }
                this.setState({ sqlExecuting: false });
            } else {
                message.error(res.msg);
                this.setState({ sqlExecuting: false });
            }
            this.loadUserHistoryLog();
        });
    }

    tabsOnChange = (activeKey) => {
        this.setState({ activeKey });
    };

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
                <OpsBreadcrumbPath pathData={["DMS", "数据管理"]} />
                <div style={{ height: "95%" }}>
                    <div style={left_panel}>
                        {this.state.treeData.length > 0 ? (
                            <Tree
                                showIcon={true}
                                loadData={this.onLoadData}
                                showLine={true}
                                style={{
                                    maxHeight: "100%",
                                    overflow: "scroll",
                                }}
                                onSelect={this.onTreeNodeSelect}
                            >
                                {this.renderTreeNodes(this.state.treeData)}
                            </Tree>
                        ) : (
                            <Empty description="暂无授权实例" />
                        )}
                    </div>
                    <div style={right_panel}>
                        <Row style={{ marginBottom: 10 }}>
                            <Text
                                strong
                                style={{ color: "rgba(0, 102, 192, 0.86)" }}
                            >
                                你当前选择：{this.state.currentChoose}
                            </Text>
                        </Row>
                        <div>
                            <Text strong>在此输入SQL: </Text>
                        </div>
                        <Row style={{ height: "120px" }}>
                            <Col span={24}>
                                <CodeMirror
                                    className="sqlEditor"
                                    ref={this.sqlInputRef}
                                    options={{
                                        mode: "text/x-mysql",
                                        showCursorWhenSelecting: true,
                                        option: {
                                            autofocus: true,
                                        },
                                        lineWrapping: true,
                                    }}
                                    value={this.state.sqlInput}
                                    onBeforeChange={(editor, data, value) => {
                                        this.setState({ sqlInput: value });
                                    }}
                                />
                            </Col>
                        </Row>
                        <Row style={{ marginTop: 5 }}>
                            <Button
                                type="primary"
                                loading={this.state.sqlExecuting}
                                onClick={this.submitSql.bind(this)}
                                disabled={!this.props.aclAuthMap["POST:/dms/userExecSQL"]}
                            >
                                提交执行
                            </Button>
                        </Row>
                        <Row style={{ marginBottom: 10 }}>
                            <Col span={24}>
                                <Tabs
                                    defaultActiveKey="1"
                                    activeKey={this.state.activeKey}
                                    size="small"
                                    onChange={this.tabsOnChange}
                                >
                                    <TabPane tab="执行结果" key="1">
                                        {this.state.execResultPanel}
                                    </TabPane>
                                    <TabPane
                                        tab={
                                            <span>
                                                执行历史{" "}
                                                <Icon
                                                    type="reload"
                                                    onClick={
                                                        this.loadUserHistoryLog
                                                    }
                                                    style={{ marginLeft: 10 }}
                                                />
                                            </span>
                                        }
                                        key="2"
                                    >
                                        <Table
                                            size="small"
                                            dataSource={this.state.sqlResultLog}
                                            columns={
                                                this.state.sqlResultLogColumns
                                            }
                                            scroll={{ x: "max-content" }}
                                            pagination={this.state.pagination}
                                            loading={this.state.tableLoading}
                                            bordered
                                        />
                                    </TabPane>
                                </Tabs>
                            </Col>
                        </Row>
                    </div>
                </div>
            </Content>
        );
    }
}

export default DataManageContent;
