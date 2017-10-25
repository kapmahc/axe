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
    const {formatMessage} = this.props.intl
    e.preventDefault();
    this.props.form.validateFields((err, values) => {
      if (!err) {
        values.type = 'text'
        axios.post('/api/leave-words', values).then(() => {
          message.success(formatMessage({id: "messages.success"}))
        }, fail);
      }
    });
  }
  render() {
    const {formatMessage} = this.props.intl
    const {getFieldDecorator} = this.props.form
    return (
      <Layout breads={[{
          href: "/leave-words/new",
          label: <FormattedMessage id={"nut.leave-words.new.title"}/>
        }
      ]}>
        <Row>
          <Col md={{
            span: 12,
            offset: 2
          }}>
            <Form onSubmit={this.handleSubmit}>
              <FormItem {...formItemLayout} label={< FormattedMessage id = "attributes.content" />} extra={< FormattedMessage id = "nut.leave-words.new.help" />} hasFeedback>
                {getFieldDecorator('body', {
                  rules: [
                    {
                      required: true,
                      message: formatMessage({id: "errors.empty-email"})
                    }
                  ]
                })(<Input.TextArea rows={8}/>)}
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
