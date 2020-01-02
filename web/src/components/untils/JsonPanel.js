import React, {Component, Fragment} from 'react';
import ReactJson from 'react-json-view'


class JsonPanel extends Component {

    constructor(props) {
        super(props);
        this.state = {

        };
    }

    render() {
        return (
            <Fragment>
                <ReactJson src={this.props.jsonData} collapsed={false} displayDataTypes={false} onEdit={(edit)=>{
                    this.props.modifyCallback(edit.updated_src);
                }}/>
            </Fragment>
        )
    }
}

export default JsonPanel;