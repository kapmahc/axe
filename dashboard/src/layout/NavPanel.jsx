import React, {Component} from 'react'
import {
  Layout,
  Menu,
  Breadcrumb,
  Icon,
  Modal,
  message
} from 'antd'

class Widget extends Component {
  render() {
    return (
      <Menu theme="dark" mode="inline" defaultSelectedKeys={['4']}>
        <Menu.Item key="1">
          <Icon type="user"/>
          <span className="nav-text">nav 1</span>
        </Menu.Item>
        <Menu.Item key="2">
          <Icon type="video-camera"/>
          <span className="nav-text">nav 2</span>
        </Menu.Item>
        <Menu.Item key="3">
          <Icon type="upload"/>
          <span className="nav-text">nav 3</span>
        </Menu.Item>
        <Menu.Item key="4">
          <Icon type="user"/>
          <span className="nav-text">nav 4</span>
        </Menu.Item>
      </Menu>
    );
  }
}

export default Widget;
