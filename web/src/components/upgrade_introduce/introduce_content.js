import React, {Component, Fragment} from "react";
import {Alert, Button, Carousel, Col, Icon, Layout, Modal, Row} from "antd";

const { Content } = Layout;

class IntroduceContent extends Component {

    constructor(props) {
        super(props);
        this.showIntroduceVersion = this.showIntroduceVersion.bind(this);
        this.handleCancel = this.handleCancel.bind(this);
        this.slidePrev = this.slidePrev.bind(this);
        this.slideNext = this.slideNext.bind(this);
        this.state = {
            modalVisible: false,
        }
    }

    showIntroduceVersion() {
        this.setState({modalVisible: true});
    }

    handleCancel() {
        this.setState({modalVisible: false});
    }

    slidePrev() {
        this.carousel.prev();
    }

    slideNext() {
        this.carousel.next();
    }

    getVersionIntroduceContent() {
        return (
            <div>
                <Row>
                    <Col span={1} style={{ paddingTop: '25%' }}>
                        {/*<Icon type="left-circle" style={{ fontSize: '32px', color: '#08c' }} onClick={this.slidePrev}/>*/}
                    </Col>
                    <Col span={22}>
                        <Carousel ref={(carousel)=>this.carousel=carousel} autoplay={true} arrows={false} style={{ height: '600px' }} >
                            <div className="intro">
                                <Row style={{paddingTop: '90px'}}>
                                    <Col span={12} style={{height: '330px', paddingLeft: '20px'}}>
                                        <div className="intro_img_common" />
                                    </Col>
                                    <Col span={8} className="intro_desc">
                                        1、支持阿里云账号下云资源的集成；<br/><br/>
                                        2、支持简版工单的协作；<br/><br/>
                                        3、采用 RBAC 进行用户菜单权限和资源可见权限的控制；<br/><br/>
                                        4、支持接入 Kubernetes 对其进行管理；
                                    </Col>
                                </Row>
                            </div>
                        </Carousel>
                    </Col>
                    <Col span={1} style={{ paddingTop: '25%' }}>
                        {/*<Icon type="right-circle" style={{ fontSize: '32px', color: '#08c' }}  onClick={this.slideNext}/>*/}
                    </Col>
                </Row>
            </div>
        );
    }

    render() {
        let alertContent = (
            <Fragment>
                <Button type="link" onClick={this.showIntroduceVersion}>
                    <Icon type="notification" /> 点击查看当前版本「Version: 1.0.0」更新内容
                </Button>
            </Fragment>
        );
        return (
            <Content
                style={{ padding: "5px 10px", margin: 0, height: "10%" }}
            >
                <Alert type="success" message={alertContent} closable/>
                <Modal
                    visible={this.state.modalVisible}
                    destroyOnClose={true}
                    centered={true}
                    width='70%'
                    closable={true}
                    onCancel={this.handleCancel}
                    footer={null}
                >
                    { this.getVersionIntroduceContent() }
                </Modal>
            </Content>
        );
    }
}

export default IntroduceContent;
