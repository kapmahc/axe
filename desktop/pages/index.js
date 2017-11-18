import React, {Component} from 'react'

import Layout from '../layouts/application'
import pageWithIntl from '../components/PageWithIntl'

class Widget extends Component {
  render() {
    return (<Layout>
      <div>Welcome to next.js!</div>
    </Layout>)

  }
}

export default pageWithIntl(Widget)
