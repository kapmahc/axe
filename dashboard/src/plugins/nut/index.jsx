import React from 'react'
import {Route} from 'react-router'

import Home from './Home'
import Install from './Install'
import NoMatch from './NoMatch'
import UsersSignIn from './users/SignIn'
import UsersSignUp from './users/SignUp'

const routes = [
  (< Route key = "nut.home" exact path = "/" component = {
    Home
  } />),
  (< Route key = "nut.install" path = "/install" component = {
    Install
  } />),
  (< Route key = "nut.users.sign-in" path = "/users/sign-in" component = {
    UsersSignIn
  } />),
  (< Route key = "nut.users.sign-up" path = "/users/sign-up" component = {
    UsersSignUp
  } />),

  (<Route key="nut.no-match" component={NoMatch}/>)
]

export default routes
