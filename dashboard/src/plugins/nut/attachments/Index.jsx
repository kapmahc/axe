import React, {Component} from 'react'
import {
  Row,
  Col,
  Table,
  Popconfirm,
  Button,
  Upload,
  Icon,
  message
} from 'antd'
import {injectIntl, intlShape, FormattedMessage} from 'react-intl'
import {connect} from 'react-redux'
import {push} from 'react-router-redux'
import {CopyToClipboard} from 'react-copy-to-clipboard'

import Layout from '../../../layout'
import {get, _delete} from '../../../ajax'
import {TOKEN} from '../../../actions'

class Widget extends Component {
  state = {
    items: []
  }
  componentDidMount() {
    get('/api/attachments').then((rst) => {
      this.setState({items: rst})
    }).catch(message.error);
  }
  handleRemove = (id) => {
    const {formatMessage} = this.props.intl
    _delete(`/api/attachments/${id}`).then((rst) => {
      message.success(formatMessage({id: 'messages.success'}))
      var items = this.state.items.filter((it) => it.id !== id)
      this.setState({items})
    }).catch(message.error)
  }
  render() {
    return (<Layout breads={[{
          href: "/attachments",
          label: <FormattedMessage id={"nut.attachments.index.title"}/>
        }
      ]}>
      <Row>
        <Col>
          <Upload multiple={true} name="file" action="/api/attachments" headers={{
              'Authorization' : `BEARER ${window.sessionStorage.getItem(TOKEN)}`
            }}>
            <Button>
              <Icon type="upload"/>
              <FormattedMessage id="nut.attachments.index.upload"/>
            </Button>
          </Upload>
          <Table bordered={true} rowKey="id" dataSource={this.state.items} columns={[
              {
                title: <FormattedMessage id="attributes.content"/>,
                key: 'title',
                render: (text, record) => (<a href={record.url} target="_blank">
                  {record.title}
                </a>)
              }, {
                title: <FormattedMessage id="attributes.type"/>,
                dataIndex: 'mediaType',
                key: 'mediaType'
              }, {
                title: 'Action',
                key: 'action',
                render: (text, record) => (<span>
                  <CopyToClipboard text={record.url}><Button shape="circle" icon="copy"/></CopyToClipboard>
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
