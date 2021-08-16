import React, { Component } from "react";
import { Layout, Menu, message } from "antd";
import * as Icons from "@ant-design/icons";
import { Link } from "react-router-dom";
import Menus from "../menu";
import { getUserPermissionsList } from "../api/permission";

const { Sider } = Layout;
const SubMenu = Menu.SubMenu;

class OpsSider extends Component {
    constructor(props) {
        super(props);
        this.state = {
            authList: [],
        };
    }

    componentWillMount() {
        getUserPermissionsList({authType: "菜单"})
            .then((res) => {
                if (res.code === 0) {
                    const dataList = res.data;
                    let permissionUrlList = [];
                    dataList.map((item) => {
                        if(item.authType==="菜单") {
                            permissionUrlList.push(item.urlPath);
                        }
                    });
                    this.setState({
                        authList: permissionUrlList,
                    });
                } else {
                    message.error("获取用户菜单权限异常");
                }
            })
    }

    hasSubMenusAllowed(subMenuList) {
        let result = false;
        subMenuList.map((item) => {
            if (
                this.state.authList.indexOf(item.route) > -1 ||
                Menus.noAuthMenus.indexOf(item.route) > -1
            ) {
                result = true;
                return result;
            } else {
                return null;
            }
        });
        return result;
    }

    hasMenuAllowed(menuUrlPath) {
        return (
            this.state.authList.indexOf(menuUrlPath) > -1 ||
            Menus.noAuthMenus.indexOf(menuUrlPath) > -1
        );
    }

    render() {
        return (
            <Sider
                trigger={null}
                collapsed={this.props.menuCollapsed}
                style={{ textAlign: "center" }}
            >
                <span
                    className="logo"
                    style={
                        this.props.menuCollapsed
                            ? { backgroundSize: "80%" }
                            : { backgroundSize: "80%" }
                    }
                />
                <span
                    className="logo-text"
                    style={this.props.menuCollapsed ? { display: "none" } : {}}
                >
                    运维平台
                </span>
                <Menu theme="dark" mode="inline" style={{ textAlign: "left" }}>
                    {Menus.sideMenus.map((menu) => {
                        if (
                            menu.subMenus.length &&
                            this.hasSubMenusAllowed(menu.subMenus)
                        ) {
                            const menuIcon = React.createElement(
                                Icons[menu.icon],
                            );
                            return (
                                <SubMenu
                                    key={menu.menuTitle}
                                    title={
                                        <span>
                                            {menuIcon}
                                            <span>{menu.menuTitle}</span>
                                        </span>
                                    }
                                >
                                    {menu.subMenus.map((subMenu) => {
                                        if (
                                            this.hasMenuAllowed(subMenu.route)
                                        ) {
                                            return (
                                                <Menu.Item key={subMenu.title}>
                                                    <Link to={subMenu.route}>
                                                        {subMenu.title}
                                                    </Link>
                                                </Menu.Item>
                                            );
                                        } else {
                                            return null;
                                        }
                                    })}
                                </SubMenu>
                            );
                        } else {
                            return null;
                        }
                    })}
                </Menu>
            </Sider>
        );
    }
}

export default OpsSider;
