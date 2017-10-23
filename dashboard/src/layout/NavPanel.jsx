import React, {Component} from 'react'
import PropTypes from 'prop-types'
import {FormattedMessage} from 'react-intl'
import {connect} from 'react-redux'
import {Menu, Icon} from 'antd'
import {push} from 'react-router-redux'

import {signOut} from '../actions'

class Widget extends Component {
  handleMenu = ({key}) => {
    const {push} = this.props
    switch (key) {
      case "nut.users.sign-out":
        break
      default:
        push(key)
    }
  };
  render() {
    const {user} = this.props
    var items = user.uid
      ? []
      : [
        {
          icon: "user",
          label: "nut.users.sign-in.title",
          key: "/users/sign-in"
        }, {
          icon: "user-add",
          label: "nut.users.sign-up.title",
          key: "/users/sign-up"
        }, {
          icon: "key",
          label: "nut.users.forgot-password.title",
          key: "/users/forgot-password"
        }, {
          icon: "check-circle-o",
          label: "nut.users.confirm.title",
          key: "/users/confirm"
        }, {
          icon: "unlock",
          label: "nut.users.unlock.title",
          key: "/users/unlock"
        }
      ]
    return (
      <Menu theme="dark" mode="inline" defaultSelectedKeys={[]} onClick={this.handleMenu}>
        {items.map((it) => (
          <Menu.Item key={it.key}>
            <Icon type={it.icon}/>
            <FormattedMessage id={it.label}/>
          </Menu.Item>
        ))}
      </Menu>
    )

  }
}

Widget.propTypes = {
  push: PropTypes.func.isRequired,
  signOut: PropTypes.func.isRequired,
  user: PropTypes.object.isRequired
}

export default connect(state => ({user: state.currentUser}), {push, signOut})(Widget)
