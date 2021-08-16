import React, {Component} from 'react';
class LogoutContent extends Component {
    componentWillMount() {
        localStorage.removeItem("ops_token");
        localStorage.removeItem("token_info");
        window.location.href = "/login";
    }

    render() {
        return (
            <div>
                注销中...
            </div>
        )
    }

}

export default LogoutContent;