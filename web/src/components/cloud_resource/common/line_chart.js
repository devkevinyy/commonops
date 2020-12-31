import React, {Component} from "react"
import { Line } from '@antv/g2plot';

class LineChart extends Component {
    constructor(props) {
        super(props);
    }

    componentDidMount() {
        let data = this.props.data;
        this.line = new Line('lineMonitorChart', {
            data,
            width: 800,
            padding: 'auto',
            xField: 'date',
            yField: 'value',
        });
        this.line.render();
    }

    componentDidUpdate() {
        this.line.changeData(this.props.data);
    }

    render() {
        return (
            <div id="lineMonitorChart" style={{ height: this.props.height, width: this.props.width }}></div>
        );
    }
}

export default LineChart;