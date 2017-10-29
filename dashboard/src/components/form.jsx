import React, {Component} from 'react'
import {Form, Button} from 'antd'
import {FormattedMessage} from 'react-intl'

const FormItem = Form.Item

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

export const orders = (size) => Array(size * 2 + 1).fill().map((_, id) => (id - size).toString())

// export class SortOrder extends Component {
//   render() {
//     const {size} = this.props
//     return (
//       <Select>
//         {Array(size * 2 + 1).fill().map((_, id) => (id - size).toString()).map((i) => (
//           <Option key={i} value={i}>{i}</Option>
//         ))}
//       </Select>
//     )
//   }
// }
//
// SortOrder.propTypes = {
//   size: PropTypes.number.isRequired
// }

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
