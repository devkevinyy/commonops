import React, {Component} from 'react';
import {Card, Layout, Row, message, Button} from "antd";
import OpsBreadcrumbPath from "../breadcrumb_path";
import {getJenkinsAllJobs} from "../../api/jenkins";


const {Content} = Layout;

class JenkinsJobsContent extends Component {

    constructor(props) {
        super(props);
        this.jumpJobDetail = this.jumpJobDetail.bind(this);
        this.state = {
            jobs: [],
        }
    }

    componentDidMount() {
        this.loadJenkinsAllJobs();
    }

    loadJenkinsAllJobs() {
        getJenkinsAllJobs().then((res)=>{
            if(res.code===0){
                this.setState({jobs: res.data});
            } else {
                message.error(res.msg);
            }
        }).catch((err)=>{
            message.error(err.toLocaleString());
        })
    }

    jumpJobDetail(jobName) {
        this.props.history.push({pathname: "/admin/cicd/job_detail", state: {"jobName": jobName}});
    }

    render() {
        return (
            <Content
                style={{ background: '#fff', padding: "5px 20px", margin: 0, height: "100%" }}
            >
                <OpsBreadcrumbPath pathData={["CI & CD", "任务列表"]} />
                <Row>
                    <Card title="当前任务" size="small">
                        {
                            this.state.jobs.map((item, index)=> {
                                return (
                                    <Card.Grid key={index} style={{ width: '25%', textAlign: 'center' }}>
                                        <Button type="link" onClick={this.jumpJobDetail.bind(this, item.name)}>{item.name}</Button>
                                    </Card.Grid>
                                )
                            })
                        }
                    </Card>
                </Row>
            </Content>
        )
    }
}

export default JenkinsJobsContent;