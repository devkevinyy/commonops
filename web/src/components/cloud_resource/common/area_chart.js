import React, { Component } from "react";
import { Area } from "@antv/g2plot";

class AreaChart extends Component {
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
            this.area = new Area(id, {
                data,
                width: this.props.width,
                height: this.props.height,
                xField: xField,
                yField: yField,
                padding: "auto",
                isStack: false,
                yAxis: {
                    label: {
                        formatter: (v) => {
                            return v + " " + this.props.unit;
                        },
                        style: {
                            fill: "#FE740C",
                        },
                    },
                },
            });
            this.area.render();
        }, 300);
    }

    componentDidUpdate() {
        this.area.changeData(this.props.data);
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

export default AreaChart;
