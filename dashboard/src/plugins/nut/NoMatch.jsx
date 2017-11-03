import React, {Component} from 'react'
import {FormattedMessage} from 'react-intl'
import {Link} from 'react-router-dom'
import {Row, Col, Button} from 'antd'

import Layout from '../../layout'
import E404 from '../../assets/errors/404.svg'

class Widget extends Component {
  render() {
    return (<Layout breads={[]}>
      <Row>
        <Col md={{
            span: 8,
            offset: 2
          }}><img alt="404" src={E404}/></Col>
        <Col md={{
            span: 6
          }}>
          <FormattedMessage id="nut.no-match.title" tagName="h1"/>
          <Link to="/">
            <Button type="primary"><FormattedMessage id="nut.no-match.go-home"/></Button>
          </Link>
        </Col>
      </Row>
    </Layout>);
  }
}

export default Widget;
