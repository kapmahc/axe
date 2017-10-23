import React, {Component} from 'react'
import PropTypes from 'prop-types'
import {Form, Row, Col, Input} from 'antd'
import {injectIntl, intlShape, FormattedMessage} from 'react-intl'
import {connect} from 'react-redux'
import {push} from 'react-router-redux'

import Layout from '../../layout'
import FormSubmit from '../../components/FormSubmit'

const FormItem = Form.Item

class Widget extends Component {
  handleSubmit = (e) => {
    e.preventDefault();
    this.props.form.validateFields((err, values) => {
      if (!err) {
        console.log('Received values of form: ', values);
      }
    });
  }
  checkPassword = (rule, value, callback) => {
    const {form} = this.props;
    if (value && value !== form.getFieldValue('password')) {
      callback('Two passwords that you enter is inconsistent!');
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
            span: 8,
            offset: 2
          }}>
            <Form onSubmit={this.handleSubmit} className="login-form">
              <FormItem label={< FormattedMessage id = "attributes.username" />} hasFeedback>
                {getFieldDecorator('name', {
                  rules: [
                    {
                      required: true,
                      message: formatMessage({id: "errors.empty"})
                    }
                  ]
                })(<Input/>)}
              </FormItem>
              <FormItem label={< FormattedMessage id = "attributes.email" />} hasFeedback>
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
              <FormItem label={< FormattedMessage id = "attributes.password" />} hasFeedback>
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
              <FormItem label={< FormattedMessage id = "attributes.password-confirmation" />} hasFeedback>
                {getFieldDecorator('confirmConfirmation', {
                  rules: [
                    {
                      required: true,
                      message: formatMessage({id: "errors.passwords-not-match"})
                    }, {
                      validator: this.checkPassword
                    }
                  ]
                })(<Input type="password"/>)}
              </FormItem>
              <FormSubmit/>
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
