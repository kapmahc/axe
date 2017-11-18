import Head from 'next/head'

import Header from './Header'
import Footer from './Footer'

export default({children}) => (<div>
  <Head>
    <meta charSet="utf-8"/>
    <meta name='viewport' content='width=device-width, initial-scale=1'/>
    <link rel="stylesheet" href="//cdnjs.cloudflare.com/ajax/libs/antd/2.13.10/antd.min.css"/>
  </Head>
  <Header/> {children}
  <Footer/>
</div>)
