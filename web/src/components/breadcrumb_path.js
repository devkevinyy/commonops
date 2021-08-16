import React, { Component } from 'react';
import { Breadcrumb } from 'antd';


class OpsBreadcrumbPath extends Component {

    static defaultProps = {
        pathData: [],
    };

    render() {
        return (
            <Breadcrumb style={{ margin: '5px 0px 10px 0px' }}>
                {
                    this.props.pathData.map((path) => {
                        return (
                            <Breadcrumb.Item key={path}>{path}</Breadcrumb.Item>
                        )
                    })
                }
            </Breadcrumb>
    );
  }
}

export default OpsBreadcrumbPath;
