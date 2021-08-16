import React, {Component} from 'react';
import { Result, Button } from 'antd';


class Exception500 extends Component {

    constructor(props) {
        super(props);
        this.backToLogin = this.backToLogin.bind(this);
        this.state = {}
    }

    backToLogin() {
        this.props.history.push('/admin');
    }

    render() {
        return (
            <Result
                status="500"
                title="异常提醒"
                subTitle="服务异常，请联系运维中心处理"
                style={{height: "100vh"}}
                extra={
                    <Button type="primary" onClick={this.backToLogin}>
                        返回主页
                    </Button>
                }
            />
        )
    }
}

class Exception403 extends Component {

    constructor(props) {
        super(props);
        this.backToLogin = this.backToLogin.bind(this);
        this.state = {}
    }

    backToLogin() {
        this.props.history.push('/admin');
    }

    render() {
        return (
            <Result
                status="403"
                title="403"
                subTitle="您没有访问权限"
                style={{height: "100vh"}}
                extra={
                    <Button type="primary" onClick={this.backToLogin}>
                        返回主页
                    </Button>
                }
            />
        )
    }
}

class Exception404 extends Component {

    constructor(props) {
        super(props);
        this.backToLogin = this.backToLogin.bind(this);
        this.state = {}
    }

    backToLogin() {
        this.props.history.push('/admin');
    }

    render() {
        return (
            <Result
                status="404"
                title="404"
                subTitle="访问的页面不存在"
                style={{height: "100vh"}}
                extra={
                    <Button type="primary" onClick={this.backToLogin}>
                        返回主页
                    </Button>
                }
            />
        )
    }
}

export { Exception403, Exception404, Exception500 };