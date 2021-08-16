import React from 'react';
import ReactDOM from 'react-dom';
import App from './layouts/main';
import * as serviceWorker from './serviceWorker';
import zhCN from 'antd/es/locale-provider/zh_CN';
import {ConfigProvider} from "antd";


ReactDOM.render(
    <ConfigProvider locale={zhCN}>
        <App/>
    </ConfigProvider>
    ,
    document.getElementById('root'));

serviceWorker.unregister();
