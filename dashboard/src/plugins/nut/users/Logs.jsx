import React, {Component} from 'react'
import {Row, Col, Table, message} from 'antd'
import {injectIntl, intlShape, FormattedMessage} from 'react-intl'
import {connect} from 'react-redux'
import Moment from 'react-moment'

import Layout from '../../../layout'
import {get} from '../../../ajax'

class Widget extends Component {
  state = {
    items: []
  }
  componentDidMount() {
    get('/api/users/logs').then((rst) => {
      this.setState({items: rst})
    }).catch(message.error);
  }
  render() {
    return (
      <Layout breads={[{
          href: "/users/logs",
          label: <FormattedMessage id={"nut.users.logs.title"}/>
        }
      ]}>
        <Row>
          <Col md={{
            span: 12,
            offset: 2
          }}>
            <Table bordered rowKey="id" dataSource={this.state.items} columns={[
              {
                title: <FormattedMessage id="attributes.createdAt"/>,
                key: 'createdAt',
                render: (text, record) => (
                  <Moment fromNow>{record.createdAt}</Moment>
                )
              }, {
                title: <FormattedMessage id="attributes.ip"/>,
                dataIndex: 'ip',
                key: 'ip'
              }, {
                title: <FormattedMessage id="nut.attributes.log.message"/>,
                dataIndex: 'message',
                key: 'message'
              }
            ]}/>
          </Col>
        </Row>
      </Layout>
    );
  }
}

Widget.propTypes = {
  intl: intlShape.isRequired
}

const WidgetI = injectIntl(Widget)

export default connect(state => ({}), {},)(WidgetI)
