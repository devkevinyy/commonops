import React, { Component } from "react";
import { Layout, Table, Typography } from "antd";
import OpsBreadcrumbPath from "../breadcrumb_path";
import { getCdProcessLog } from "../../api/cd";
const { Content } = Layout;
const { Text } = Typography;

class CdRecordContent extends Component {
    constructor(props) {
        super(props);
        this.state = {
            columns: [
                {
                    title: "id",
                    dataIndex: "id",
                    width: 30,
                    render: (value) => {
                        return <Text ellipsis={true}>{value}</Text>;
                    },
                },
                {
                    title: "发布镜像",
                    dataIndex: "imageName",
                    width: 100,
                    render: (value) => {
                        return <Text ellipsis={true}>{value}</Text>;
                    },
                },
                {
                    title: "命名空间",
                    dataIndex: "namespace",
                    width: 100,
                    render: (value) => {
                        return <Text ellipsis={true}>{value}</Text>;
                    },
                },
                {
                    title: "项目名称",
                    dataIndex: "jobName",
                    width: 100,
                    render: (value) => {
                        return <Text ellipsis={true}>{value}</Text>;
                    },
                },
                {
                    title: "发布状态",
                    dataIndex: "success",
                    width: 100,
                    render: (value) => {
                        return <Text ellipsis={true}>{value}</Text>;
                    },
                },
                {
                    title: "发布日志",
                    dataIndex: "result",
                    width: 100,
                    render: (value) => {
                        return <Text ellipsis={true}>{value}</Text>;
                    },
                },
            ],
            tableData: [],
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
            tableLoading: false,
        };
    }

    componentDidMount() {
        this.refreshTableData();
    }

    refreshTableData() {
        this.setState({ tableLoading: true });
        const queryParams = {
            page: this.state.pagination.page,
            size: this.state.pagination.pageSize,
        };
        getCdProcessLog(queryParams).then((res) => {
            const pagination = this.state.pagination;
            pagination.total = parseInt(res.data.total);
            pagination.page = parseInt(res.data.page);
            pagination.showTotal(parseInt(res.data.total));
            this.setState({
                pagination,
            });
            let tableData = res["data"]["logs"];
            this.setState({
                tableData: tableData,
                tableLoading: false,
            });
        });
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

    render() {
        return (
            <Content
                style={{
                    background: "#fff",
                    padding: "5px 20px",
                    margin: 0,
                    height: "100%",
                }}
            >
                <OpsBreadcrumbPath pathData={["CI & CD", "部署记录"]} />
                <Table
                    columns={this.state.columns}
                    dataSource={this.state.tableData}
                    scroll={{ x: "max-content" }}
                    pagination={this.state.pagination}
                    loading={this.state.tableLoading}
                    bordered
                    size="small"
                />
            </Content>
        );
    }
}

export default CdRecordContent;
