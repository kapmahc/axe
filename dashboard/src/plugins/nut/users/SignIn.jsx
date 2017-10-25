import React, {Component} from 'react'
import PropTypes from 'prop-types'
import {Form, Row, Col, Input} from 'antd'
import {injectIntl, intlShape, FormattedMessage} from 'react-intl'
import {connect} from 'react-redux'
import {push} from 'react-router-redux'
import axios from 'axios'

import Layout from '../../../layout'
import {Submit, formItemLayout, fail} from '../../../components/form'

const FormItem = Form.Item

class Widget extends Component {
  handleSubmit = (e) => {
    const {push} = this.props
    e.preventDefault();
    this.props.form.validateFields((err, values) => {
      if (!err) {
        axios.post('/api/users/sign-in', values).then(() => {
          push('/')
        }, fail);
      }
    });
  }
  render() {
    const {formatMessage} = this.props.intl
    const {getFieldDecorator} = this.props.form
    return (
      <Layout breads={[{
          href: "/users/sign-in",
          label: <FormattedMessage id={"nut.users.sign-in.title"}/>
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
              <FormItem {...formItemLayout} label={< FormattedMessage id = "attributes.password" />} hasFeedback>
                {getFieldDecorator('password', {
                  rules: [
                    {
                      required: true,
                      message: formatMessage({id: "errors.empty-password"})
                    }
                  ]
                })(<Input type="password"/>)}
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
  push: PropTypes.func.isRequired
}

const WidgetF = Form.create()(injectIntl(Widget))

export default connect(state => ({}), {
  push
},)(WidgetF)
