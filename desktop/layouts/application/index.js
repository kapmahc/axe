import Head from 'next/head'
import {injectIntl} from 'react-intl'

import Header from './Header'
import Footer from './Footer'

export default injectIntl(({intl, title, children}) => (<div>
  <Head>
    <meta charSet="utf-8"/>
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no"/>
    <link rel="stylesheet" href="//maxcdn.bootstrapcdn.com/bootstrap/4.0.0-beta.2/css/bootstrap.min.css"/>
    <link rel="stylesheet" href="//maxcdn.bootstrapcdn.com/font-awesome/4.7.0/css/font-awesome.min.css"/>
    <title>{title || 'aaa'}</title>
  </Head>
  <Header/> {children}
  <Footer/>
</div>))
