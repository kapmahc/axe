import React, {Component} from 'react'
import PropTypes from 'prop-types'
import {Layout} from 'antd'
import {connect} from 'react-redux'

const {Footer} = Layout

class Widget extends Component {
  render() {
    return (
      <Footer style={{
        textAlign: 'center'
      }}>
        Ant Design ©2016 Created by Ant UED
      </Footer>
    );
  }
}
Widget.propTypes = {
  info: PropTypes.object.isRequired
}

export default connect(state => ({info: state.siteInfo}))(Widget)
