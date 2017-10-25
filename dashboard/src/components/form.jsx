import React, {Component} from 'react'
import {Form, Button, message} from 'antd'
import {FormattedMessage} from 'react-intl'

const FormItem = Form.Item

export const fail = (err) => message.error(err.response.data)

export const formItemLayout = {
  labelCol: {
    xs: {
      span: 24
    },
    sm: {
      span: 8
    }
  },
  wrapperCol: {
    xs: {
      span: 24
    },
    sm: {
      span: 16
    }
  }
};
export const tailFormItemLayout = {
  wrapperCol: {
    xs: {
      span: 24,
      offset: 0
    },
    sm: {
      span: 16,
      offset: 8
    }
  }
};

export class Submit extends Component {
  render() {
    return (
      <FormItem {...tailFormItemLayout}>
        <Button type="primary" htmlType="submit">
          <FormattedMessage id="buttons.submit"/>
        </Button>
      </FormItem>
    );
  }
}
