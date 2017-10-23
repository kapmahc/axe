import React, {Component} from 'react'
import {Form, Button} from 'antd'
import {FormattedMessage} from 'react-intl'

const FormItem = Form.Item

class Widget extends Component {
  render() {
    return (
      <FormItem>
        <Button type="primary" htmlType="submit">
          <FormattedMessage id="buttons.submit"/>
        </Button>
      </FormItem>
    );
  }
}

export default Widget
