import React, {Component} from 'react'
import PropTypes from 'prop-types'
import {Form, Row, Col, Input, message} from 'antd'
import {injectIntl, intlShape, FormattedMessage} from 'react-intl'
import {connect} from 'react-redux'
import {push} from 'react-router-redux'

import Layout from '../../layout'
import {post} from '../../ajax'
import {Submit, formItemLayout} from '../../components/form'

const FormItem = Form.Item

class Widget extends Component {
  handleSubmit = (e) => {
    const {formatMessage} = this.props.intl
    const {push} = this.props
    e.preventDefault();
    this.props.form.validateFields((err, values) => {
      if (!err) {
        post('/api/install', values).then(() => {
          message.info(formatMessage({id: "nut.install.success"}))
          push('/users/sign-in')
        }).catch(message.error);
      }
    });
  }
  checkPassword = (rule, value, callback) => {
    const {formatMessage} = this.props.intl
    const {getFieldValue} = this.props.form;
    if (value && value !== getFieldValue('password')) {
      callback(formatMessage({id: "errors.passwords-not-match"}));
    } else {
      callback();
    }
  }
  render() {
    const {formatMessage} = this.props.intl
    const {getFieldDecorator} = this.props.form
    return (
      <Layout breads={[{
          href: "/install",
          label: <FormattedMessage id={"nut.install.title"}/>
        }
      ]}>
        <Row>
          <Col md={{
            span: 12,
            offset: 2
          }}>
            <Form onSubmit={this.handleSubmit}>
              <FormItem {...formItemLayout} label={< FormattedMessage id = "nut.attributes.site.title" />} hasFeedback>
                {getFieldDecorator('title', {
                  rules: [
                    {
                      required: true,
                      message: formatMessage({id: "errors.empty"})
                    }
                  ]
                })(<Input/>)}
              </FormItem>
              <FormItem {...formItemLayout} label={< FormattedMessage id = "nut.attributes.site.subhead" />} hasFeedback>
                {getFieldDecorator('subhead', {
                  rules: [
                    {
                      required: true,
                      message: formatMessage({id: "errors.empty"})
                    }
                  ]
                })(<Input/>)}
              </FormItem>
              <FormItem {...formItemLayout} label={< FormattedMessage id = "attributes.username" />} hasFeedback>
                {getFieldDecorator('name', {
                  rules: [
                    {
                      required: true,
                      message: formatMessage({id: "errors.empty"})
                    }
                  ]
                })(<Input/>)}
              </FormItem>
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
                    }, {
                      validator: this.checkConfirm
                    }
                  ]
                })(<Input type="password"/>)}
              </FormItem>
              <FormItem {...formItemLayout} label={< FormattedMessage id = "attributes.password-confirmation" />} hasFeedback>
                {getFieldDecorator('passwordConfirmation', {
                  rules: [
                    {
                      required: true,
                      message: formatMessage({id: "errors.empty"})
                    }, {
                      validator: this.checkPassword
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
