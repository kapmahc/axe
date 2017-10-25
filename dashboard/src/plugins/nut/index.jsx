import React from 'react'
import {Route} from 'react-router'

import Home from './Home'
import Install from './Install'
import NoMatch from './NoMatch'
import UsersSignIn from './users/SignIn'
import UsersSignUp from './users/SignUp'
import UsersEmailForm from './users/EmailForm'

const UsersConfirm = () => (<UsersEmailForm action="confirm"/>)
const UsersUnlock = () => (<UsersEmailForm action="unlock"/>)
const UsersForgotPassword = () => (<UsersEmailForm action="forgot-password"/>)

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
  (< Route key = "nut.users.confirm" path = "/users/confirm" component = {
    UsersConfirm
  } />),
  (< Route key = "nut.users.unlock" path = "/users/unlock" component = {
    UsersUnlock
  } />),
  (< Route key = "nut.users.forgot-password" path = "/users/forgot-password" component = {
    UsersForgotPassword
  } />),

  (<Route key="nut.no-match" component={NoMatch}/>)
]

export default routes
