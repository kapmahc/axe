import React, {Component} from 'react'
import PropTypes from 'prop-types'
import {Form, Button} from 'antd'
import {FormattedMessage} from 'react-intl'
import ReactQuill from 'react-quill'

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

export class Quill extends Component {
  render() {
    const {value, onChange} = this.props
    const modules = {
      toolbar: [
        [
          {
            'header': [1, 2, false]
          }
        ],
        [
          'bold', 'italic', 'underline', 'strike', 'blockquote'
        ],
        [
          {
            'list': 'ordered'
          }, {
            'list': 'bullet'
          }, {
            'indent': '-1'
          }, {
            'indent': '+1'
          },
          'code-block'
        ],
        [
          'link', 'formula', 'image', 'video'
        ],
        ['clean']
      ]
    }

    const formats = [
      'header',
      'bold',
      'italic',
      'underline',
      'strike',
      'blockquote',
      'list',
      'bullet',
      'indent',
      'link',
      'image'
    ]
    return (<ReactQuill modules={modules} formats={formats} value={value} onChange={onChange} theme="snow"/>)
  }
}

Quill.propTypes = {
  value: PropTypes.string.isRequired,
  onChange: PropTypes.func.isRequired
}

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
