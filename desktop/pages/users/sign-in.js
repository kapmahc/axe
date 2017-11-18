import React, {Component} from 'react'
import {Container} from 'reactstrap'

import Layout from '../../layouts/application'
import SharedLinks from '../../components/users/SharedLinks'
import pageWithIntl from '../../components/PageWithIntl'

class Widget extends Component {
  render() {
    return (<Layout>
      <Container>
        Welcome to next.js!
        <br/>
        <SharedLinks/>
      </Container>
    </Layout>)
  }
}

export default pageWithIntl(Widget)
