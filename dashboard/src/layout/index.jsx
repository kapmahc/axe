import React, {Component} from 'react'
import PropTypes from 'prop-types'
import {
  Layout,
  Menu,
  Breadcrumb,
  Icon,
  Modal,
  message
} from 'antd'
import {injectIntl, intlShape, FormattedMessage} from 'react-intl'
import {connect} from 'react-redux'
import {push} from 'react-router-redux'
import {Link} from 'react-router-dom'

import Footer from './Footer'
import NavPanel from './NavPanel'
import {signIn, signOut, refresh} from '../actions'

const {Header, Content, Sider} = Layout

class Widget extends Component {
  render() {
    const {children, user} = this.props

    return user.uid
    // is sign in ?
      ? (
        <div>
          <div>{children}</div>
        </div>
      )
      : (
        <Layout>
          <Sider breakpoint="lg" collapsedWidth="0" onCollapse={(collapsed, type) => {
            console.log(collapsed, type);
          }}>
            <div className="logo"/>
            <NavPanel/>
          </Sider>
          <Layout>
            <Header style={{
              background: '#fff',
              padding: 0
            }}/>
            <Content style={{
              margin: '24px 16px 0'
            }}>
              <div style={{
                padding: 24,
                background: '#fff',
                minHeight: 360
              }}>
                {children}
              </div>
            </Content>
            <Footer/>
          </Layout>
        </Layout>
      );
  }
}

Widget.propTypes = {
  children: PropTypes.node.isRequired,
  push: PropTypes.func.isRequired,
  refresh: PropTypes.func.isRequired,
  signIn: PropTypes.func.isRequired,
  signOut: PropTypes.func.isRequired,
  user: PropTypes.object.isRequired,
  info: PropTypes.object.isRequired,
  breads: PropTypes.array.isRequired,
  intl: intlShape.isRequired
}

const WidgetI = injectIntl(Widget)

export default connect(state => ({user: state.currentUser, info: state.siteInfo}), {push, signIn, refresh, signOut})(WidgetI)
