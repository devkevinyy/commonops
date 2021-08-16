import React, { Component } from "react";
import { Layout, Table, Button, Row, Col, message } from "antd";
import { getUserFeedbackList } from "../../api/system";
import OpsBreadcrumbPath from "../breadcrumb_path";

const { Content } = Layout;

let columnStyle = {
    overFlow: "hidden",
    textOverflow: "ellipsis",
    display: "-webkit-box",
};

class UserFeedbackManager extends Component {
    constructor(props) {
        super(props);
        this.refreshTableData = this.refreshTableData.bind(this);
        this.onShowSizeChange = this.onShowSizeChange.bind(this);
        this.state = {
            columns: [
                {
                    title: "ID",
                    dataIndex: "id",
                    key: "id",
                    width: 30,
                    className: "small_font",
                },
                {
                    title: "用户",
                    dataIndex: "username",
                    key: "username",
                    className: "small_font",
                },
                {
                    title: "内容",
                    dataIndex: "content",
                    key: "content",
                    className: "small_font " + { columnStyle },
                },
                {
                    title: "评分",
                    dataIndex: "score",
                    key: "score",
                    className: "small_font " + { columnStyle },
                },
                {
                    title: "创建时间",
                    dataIndex: "createTime",
                    key: "createTime",
                    className: "small_font " + { columnStyle },
                },
            ],
            tableData: [],
            tableLoading: false,
            authLinkModalVisible: false,
            current_data_id: null,
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
        };
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
                this.refreshTableData();
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
                this.refreshTableData();
            },
        );
    };

    refreshTableData = () => {
        this.setState({ tableLoading: true });
        getUserFeedbackList(
            this.state.pagination.page,
            this.state.pagination.pageSize,
        )
            .then((res) => {
                if(res.code===0) {
                    const pagination = this.state.pagination;
                    pagination.total = parseInt(res.data.total);
                    pagination.page = parseInt(res.data.page);
                    pagination.showTotal(parseInt(res.data.total));
                    this.setState({
                        pagination,
                    });
                    let data = res["data"]["feedbacks"];
                    let tableData = [];
                    for (let i = 0; i < data.length; i++) {
                        tableData.push({
                            id: data[i]["id"],
                            createTime: data[i]["createTime"],
                            username: data[i]["username"],
                            content: data[i]["content"],
                            score: data[i]["score"],
                        });
                    }
                    this.setState({ tableData: tableData});
                } else {
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
                <OpsBreadcrumbPath pathData={["系统管理", "用户反馈"]} />
                <div style={{ marginBottom: 20 }}>
                    <Row>
                        <Col span={3}>
                            <Button
                                type="primary"
                                onClick={this.refreshTableData}
                            >
                                刷 新
                            </Button>
                        </Col>
                    </Row>
                </div>
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

export default UserFeedbackManager;
