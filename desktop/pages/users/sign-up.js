import {Container} from 'reactstrap'

import Layout from '../../layouts/application'
import SharedLinks from '../../components/users/SharedLinks'

export default() => (<Layout>
  <Container>
    Welcome to next.js!
    <br/>
    <SharedLinks/>
  </Container>
</Layout>)
