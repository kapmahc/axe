import React, {Component} from 'react'

import Layout from '../../layout'
import Exception from '../../components/Exception'

class Widget extends Component {
  render() {
    return (<Layout breads={[]}>
      <Exception error={404}/>
    </Layout>);
  }
}

export default Widget;
