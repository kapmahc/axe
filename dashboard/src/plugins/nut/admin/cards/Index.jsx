import React, {Component} from 'react'
import {
  Row,
  Col,
  Table,
  Popconfirm,
  Button,
  message
} from 'antd'
import {injectIntl, intlShape, FormattedMessage} from 'react-intl'
import {connect} from 'react-redux'
import {push} from 'react-router-redux'

import Layout from '../../../../layout'
import {get, _delete, post} from '../../../../ajax'

class Widget extends Component {
  state = {
    items: []
  }
  componentDidMount() {
    get('/api/admin/cards').then((rst) => {
      this.setState({items: rst})
    }).catch(message.error);
  }
  handleUEditor = (id) => {
    post('/api/token', {
      act: "card.summary.edit",
      tid: id
    }).then((rst) => window.open(`/cards/edit/${rst.token}`, '_blank')).catch(message.error)
  }
  handleRemove = (id) => {
    const {formatMessage} = this.props.intl
    _delete(`/api/admin/cards/${id}`).then((rst) => {
      message.success(formatMessage({id: 'messages.success'}))
      var items = this.state.items.filter((it) => it.id !== id)
      this.setState({items})
    }).catch(message.error)
  }
  render() {
    const {push} = this.props
    return (<Layout breads={[{
          href: "/admin/cards",
          label: <FormattedMessage id={"nut.admin.cards.index.title"}/>
        }
      ]}>
      <Row>
        <Col>
          <Button onClick={(e) => push('/admin/cards/new')} type='primary' shape="circle" icon="plus"/>
          <Table bordered={true} rowKey="id" dataSource={this.state.items} columns={[
              {
                title: <FormattedMessage id="attributes.loc"/>,
                key: 'loc',
                render: (text, record) => (<span>
                  {record.loc}[{record.sortOrder}]
                </span>)
              }, {
                title: <FormattedMessage id="attributes.content"/>,
                dataIndex: 'label',
                render: (text, record) => (<a target="_blank" href={record.href}>{record.title}</a>)
              }, {
                title: 'Action',
                key: 'action',
                render: (text, record) => (<span>
                  <Button onClick={(e) => push(`/admin/cards/edit/${record.id}`)} shape="circle" icon="edit"/>
                  <Button onClick={(e) => this.handleUEditor(record.id)} shape="circle" icon="chrome"/>
                  <Popconfirm title={<FormattedMessage id = "messages.are-you-sure" />} onConfirm={(e) => this.handleRemove(record.id)}>
                    <Button type="danger" shape="circle" icon="delete"/>
                  </Popconfirm>
                </span>)
              }
            ]}/>
        </Col>
      </Row>
    </Layout>);
  }
}

Widget.propTypes = {
  intl: intlShape.isRequired
}

const WidgetI = injectIntl(Widget)

export default connect(state => ({}), {
  push
},)(WidgetI)
