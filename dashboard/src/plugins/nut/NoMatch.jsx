import React, {Component} from 'react'
import {FormattedMessage} from 'react-intl'
import Exception from 'ant-design-pro/lib/Exception'

import Layout from '../../layout'

class Widget extends Component {
  render() {
    return (<Layout breads={[]}>
      <Exception img="https://gw.alipayobjects.com/zos/rmsportal/wZcnGqRDyhPOEYFcZDnb.svg" type="404"/>
    </Layout>);
  }
}

export default Widget;
