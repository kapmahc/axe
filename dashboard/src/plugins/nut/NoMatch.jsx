import React, {Component} from 'react'
import {FormattedMessage} from 'react-intl'

import Layout from '../../layout'

class Widget extends Component {
  render() {
    return (<Layout breads={[]}>
      <FormattedMessage id="errors.no-match" tagName="h1"/>
    </Layout>);
  }
}

export default Widget;
