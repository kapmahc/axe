import React, {Component} from 'react'
import {FormattedMessage} from 'react-intl'
import {Link} from 'react-router-dom'
import {Row, Col, Button} from 'antd'
import PropTypes from 'prop-types'

import E403 from '../assets/errors/403.svg'
import E404 from '../assets/errors/404.svg'
import E500 from '../assets/errors/500.svg'

class Widget extends Component {
  render() {
    const {error} = this.props
    var logo = E404
    switch (error) {
      case 500:
        logo = E500
        break;
      case 403:
        logo = E403
        break
      default:
        break
    }
    return (<Row>
      <Col md={{
          span: 8,
          offset: 2
        }}><img alt={error} src={logo}/></Col>
      <Col md={{
          span: 6
        }}>
        <FormattedMessage id={`errors.http-${error}`} tagName="h1"/>
        <Link to="/">
          <Button type="primary"><FormattedMessage id="buttons.back-to-home"/></Button>
        </Link>
      </Col>
    </Row>);
  }
}

Widget.propTypes = {
  error: PropTypes.number.isRequired
}

export default Widget;
