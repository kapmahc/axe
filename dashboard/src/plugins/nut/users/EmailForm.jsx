import React, {Component} from 'react'
import PropTypes from 'prop-types'
import {Form, Row, Col, Input, message} from 'antd'
import {injectIntl, intlShape, FormattedMessage} from 'react-intl'
import {connect} from 'react-redux'
import {push} from 'react-router-redux'
import axios from 'axios'

import Layout from '../../../layout'
import {Submit, formItemLayout, fail} from '../../../components/form'

const FormItem = Form.Item

class Widget extends Component {
  handleSubmit = (e) => {
    const {push, action} = this.props
    const {formatMessage} = this.props.intl
    e.preventDefault();
    this.props.form.validateFields((err, values) => {
      if (!err) {
        axios.post(`/api/users/${action}`, values).then(() => {
          message.info(formatMessage({id: `nut.users.${action}.notice`}))
          push('/users/sign-in')
        }, fail);
      }
    });
  }
  render() {
    const {action} = this.props
    const {formatMessage} = this.props.intl
    const {getFieldDecorator} = this.props.form
    return (
      <Layout breads={[{
          href: `/users/${action}`,
          label: <FormattedMessage id={`nut.users.${action}.title`}/>
        }
      ]}>
        <Row>
          <Col md={{
            span: 12,
            offset: 2
          }}>
            <Form onSubmit={this.handleSubmit}>
              <FormItem {...formItemLayout} label={< FormattedMessage id = "attributes.email" />} hasFeedback>
                {getFieldDecorator('email', {
                  rules: [
                    {
                      type: 'email',
                      message: formatMessage({id: "errors.not-valid-email"})
                    }, {
                      required: true,
                      message: formatMessage({id: "errors.empty-email"})
                    }
                  ]
                })(<Input/>)}
              </FormItem>
              <Submit/>
            </Form>
          </Col>
        </Row>
      </Layout>
    );
  }
}

Widget.propTypes = {
  intl: intlShape.isRequired,
  action: PropTypes.string.isRequired,
  push: PropTypes.func.isRequired
}

const WidgetF = Form.create()(injectIntl(Widget))

export default connect(state => ({}), {
  push
},)(WidgetF)
