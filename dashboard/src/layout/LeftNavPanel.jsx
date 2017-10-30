import React, {Component} from 'react'
import PropTypes from 'prop-types'
import {injectIntl, intlShape, FormattedMessage} from 'react-intl'
import {connect} from 'react-redux'
import {Menu, Icon, Modal, message} from 'antd'
import {push} from 'react-router-redux'

import {signOut} from '../actions'
import {_delete} from '../ajax'

const SubMenu = Menu.SubMenu
const confirm = Modal.confirm

class Widget extends Component {
  handleMenu = ({key}) => {
    const {push, signOut} = this.props
    const {formatMessage} = this.props.intl

    switch (key) {
      case "/users/sign-out":
        confirm({
          title: formatMessage({id: "messages.are-you-sure"}),
          onOk() {
            _delete('/api/users/sign-out').then(() => {
              signOut()
              push('/users/sign-in')
              message.success(formatMessage({id: 'messages.success'}))
            }).catch(message.error)
          }
        });
        break
      default:
        push(key)
    }
  };
  render() {
    const {user} = this.props
    var items = []
    if (user.uid) {
      // sign in?
      items.push({
        icon: "user",
        label: "nut.personal.title",
        key: "personal",
        items: [
          {
            label: "nut.users.logs.title",
            key: "/users/logs"
          }, {
            label: "nut.users.profile.title",
            key: "/users/profile"
          }, {
            label: "nut.users.change-password.title",
            key: "/users/change-password"
          }
        ]
      })

      var forum = {
        icon: "tablet",
        label: "forum.title",
        key: "forum",
        items: [
          {
            label: "forum.articles.index.title",
            key: "/forum/articles"
          }, {
            label: "forum.comments.index.title",
            key: "/forum/comments"
          }
        ]
      }
      if (user.admin) {
        // administrator?
        forum.items.push({label: "forum.tags.index.title", key: "/forum/tags"})
        items.push({
          icon: "setting",
          label: "nut.settings.title",
          key: "settings",
          items: [
            {
              label: "nut.admin.site.status.title",
              key: "/admin/site/status"
            }, {
              label: "nut.admin.site.info.title",
              key: "/admin/site/info"
            }, {
              label: "nut.admin.site.author.title",
              key: "/admin/site/author"
            }, {
              label: "nut.admin.site.seo.title",
              key: "/admin/site/seo"
            }, {
              label: "nut.admin.site.smtp.title",
              key: "/admin/site/smtp"
            }, {
              label: "nut.admin.links.index.title",
              key: "/admin/links"
            }, {
              label: "nut.admin.cards.index.title",
              key: "/admin/cards"
            }, {
              label: "nut.admin.locales.index.title",
              key: "/admin/locales"
            }, {
              label: "nut.admin.friend-links.index.title",
              key: "/admin/friend-links"
            }, {
              label: "nut.admin.leave-words.index.title",
              key: "/admin/leave-words"
            }, {
              label: "nut.admin.users.index.title",
              key: "/admin/users"
            }
          ]
        })
      }
      items.push(forum)
      items.push({icon: "notification", label: "survey.title", key: "survey", items: []})
      items.push({icon: "shopping-cart", label: "shop.title", key: "shop", items: []})

      items.push({icon: "logout", label: "nut.users.sign-out.title", key: "/users/sign-out"})
    } else {
      // non sign in?
      items.push({icon: "user", label: "nut.users.sign-in.title", key: "/users/sign-in"})
      items.push({icon: "user-add", label: "nut.users.sign-up.title", key: "/users/sign-up"})
      items.push({icon: "key", label: "nut.users.forgot-password.title", key: "/users/forgot-password"})
      items.push({icon: "check-circle-o", label: "nut.users.confirm.title", key: "/users/confirm"})
      items.push({icon: "unlock", label: "nut.users.unlock.title", key: "/users/unlock"})
      items.push({icon: "message", label: "nut.leave-words.new.title", key: "/leave-words/new"})
    }

    return (<Menu theme="dark" mode="inline" onClick={this.handleMenu}>
      {
        items.map(
          (it) => it.items
          ? (<SubMenu key={it.key} title={<span > <Icon type={it.icon}/> < FormattedMessage id = {
              it.label
            } />< /span>}>
            {it.items.map(l => (<Menu.Item key={l.key}><FormattedMessage id={l.label}/></Menu.Item>))}
          </SubMenu>)
          : (<Menu.Item key={it.key}>
            <Icon type={it.icon}/>
            <FormattedMessage id={it.label}/>
          </Menu.Item>))
      }
    </Menu>)

  }
}

Widget.propTypes = {
  intl: intlShape.isRequired,
  push: PropTypes.func.isRequired,
  signOut: PropTypes.func.isRequired,
  user: PropTypes.object.isRequired
}

const WidgetI = injectIntl(Widget)
export default connect(state => ({user: state.currentUser}), {push, signOut})(WidgetI)
