import React, { Component } from "react";
import { Line } from "@antv/g2plot";

class LineChart extends Component {
    constructor(props) {
        super(props);
        this.id = this.uuid();
    }

    uuid() {
        return "xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx".replace(/[xy]/g, function(
            c,
        ) {
            var r = (Math.random() * 16) | 0,
                v = c === "x" ? r : (r & 0x3) | 0x8;
            return v.toString(16);
        });
    }

    componentDidMount() {
        let id = this.id;
        let data = this.props.data;
        let xField = this.props.xField ? this.props.xField : "date";
        let yField = this.props.yField ? this.props.yField : "value";
        setTimeout(() => {
            this.line = new Line(id, {
                data,
                width: this.props.width,
                height: this.props.height,
                padding: "auto",
                xField: xField,
                yField: yField,
                meta: {
                    yField: {},
                },
            });
            this.line.render();
        }, 300);
    }

    componentDidUpdate() {
        if (this.props.data.length > 0) {
            setTimeout(() => this.line.changeData(this.props.data), 400);
        }
    }

    render() {
        return (
            <div
                id={this.id}
                style={{ height: this.props.height, width: this.props.width }}
            ></div>
        );
    }
}

export default LineChart;
