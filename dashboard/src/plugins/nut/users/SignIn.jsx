import React, {Component} from 'react'
import PropTypes from 'prop-types'
import {Form, Row, Col, Input, message} from 'antd'
import {injectIntl, intlShape, FormattedMessage} from 'react-intl'
import {connect} from 'react-redux'
import {push} from 'react-router-redux'

import Layout from '../../../layout'
import {post} from '../../../ajax'
import {Submit, formItemLayout} from '../../../components/form'
import {signIn} from '../../../actions'

const FormItem = Form.Item

class Widget extends Component {
  handleSubmit = (e) => {
    const {push, signIn} = this.props
    e.preventDefault();
    this.props.form.validateFields((err, values) => {
      if (!err) {
        post('/api/users/sign-in', values).then((rst) => {
          signIn(rst.token)
          push('/users/logs')
        }).catch(message.error);
      }
    });
  }
  render() {
    const {formatMessage} = this.props.intl
    const {getFieldDecorator} = this.props.form
    return (<Layout breads={[{
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
            <FormItem {...formItemLayout} label={<FormattedMessage id = "attributes.email" />} hasFeedback={true}>
              {
                getFieldDecorator('email', {
                  rules: [
                    {
                      type: 'email',
                      message: formatMessage({id: "errors.not-valid-email"})
                    }, {
                      required: true,
                      message: formatMessage({id: "errors.empty-email"})
                    }
                  ]
                })(<Input/>)
              }
            </FormItem>
            <FormItem {...formItemLayout} label={<FormattedMessage id = "attributes.password" />} hasFeedback={true}>
              {
                getFieldDecorator('password', {
                  rules: [
                    {
                      required: true,
                      message: formatMessage({id: "errors.empty-password"})
                    }
                  ]
                })(<Input type="password"/>)
              }
            </FormItem>
            <Submit/>
          </Form>
        </Col>
      </Row>
    </Layout>);
  }
}

Widget.propTypes = {
  intl: intlShape.isRequired,
  push: PropTypes.func.isRequired,
  signIn: PropTypes.func.isRequired
}

const WidgetF = Form.create()(injectIntl(Widget))

export default connect(state => ({}), {
  push,
  signIn
},)(WidgetF)
