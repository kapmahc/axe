import React, {Component} from 'react'
import PropTypes from 'prop-types'
import {
  Form,
  Row,
  Col,
  Input,
  Card,
  message
} from 'antd'
import {injectIntl, intlShape, FormattedMessage} from 'react-intl'
import {connect} from 'react-redux'
import {push} from 'react-router-redux'

import Layout from '../../../../layout'
import {post, get} from '../../../../ajax'
import {Submit, formItemLayout} from '../../../../components/form'

const FormItem = Form.Item

class Widget extends Component {
  state = {
    item: {}
  }
  componentDidMount() {
    const {setFieldsValue} = this.props.form
    get('/api/admin/site/seo').then((rst) => {
      setFieldsValue(rst)
      this.setState({item: rst})
    }).catch(message.error)
  }
  handleSubmit = (e) => {
    const {formatMessage} = this.props.intl
    e.preventDefault();
    this.props.form.validateFields((err, values) => {
      if (!err) {
        post('/api/admin/site/seo', values).then(() => {
          message.success(formatMessage({id: "messages.success"}))
        }).catch(message.error);
      }
    });
  }
  render() {
    const {formatMessage} = this.props.intl
    const {getFieldDecorator} = this.props.form
    return (<Layout breads={[{
          href: "/admin/site/seo",
          label: <FormattedMessage id={"nut.admin.site.seo.title"}/>
        }
      ]}>
      <Row>
        <Col md={{
            span: 12,
            offset: 2
          }}>
          <Form onSubmit={this.handleSubmit}>
            <FormItem {...formItemLayout} label={<FormattedMessage id = "nut.admin.site.seo.googleVerifyCode" />} hasFeedback={true}>
              {
                getFieldDecorator('googleVerifyCode', {
                  rules: [
                    {
                      required: true,
                      message: formatMessage({id: "errors.empty"})
                    }
                  ]
                })(<Input/>)
              }
            </FormItem>
            <FormItem {...formItemLayout} label={<FormattedMessage id = "nut.admin.site.seo.baiduVerifyCode" />} hasFeedback={true}>
              {
                getFieldDecorator('baiduVerifyCode', {
                  rules: [
                    {
                      required: true,
                      message: formatMessage({id: "errors.empty-email"})
                    }
                  ]
                })(<Input/>)
              }
            </FormItem>
            <Submit/>
          </Form>
        </Col>
      </Row>
      <Row>
        <Col md={{
            span: 6,
            offset: 4
          }}>
          <Card>
            {
              ['/robots.txt', '/sitemap.xml.gz', `/google${this.state.item.googleVerifyCode}.html`, `/baidu_verify_${this.state.item.baiduVerifyCode}.html`].concat(this.props.site.languages.map(it => `/rss-${it}.atom`)).map((it, id) => (<p key={id}>
                <a href={it} target='_blank'>{it}</a>
              </p>))
            }
          </Card>
        </Col>
      </Row>
    </Layout>);
  }
}

Widget.propTypes = {
  intl: intlShape.isRequired,
  site: PropTypes.object.isRequired,
  push: PropTypes.func.isRequired
}

const WidgetF = Form.create()(injectIntl(Widget))

export default connect(state => ({site: state.siteInfo}), {
  push
},)(WidgetF)
