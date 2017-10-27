import jwtDecode from 'jwt-decode'

import {USERS_SIGN_IN, USERS_SIGN_OUT, SITE_REFRESH, TOKEN} from './actions'

const currentUser = (state = {}, action) => {
  switch (action.type) {
    case USERS_SIGN_IN:
      try {
        var u = jwtDecode(action.token);
        sessionStorage.setItem(TOKEN, action.token);
        return u
      } catch (e) {
        console.error(e)
      }
      return {}
    case USERS_SIGN_OUT:
      sessionStorage.removeItem(TOKEN);
      return {}
    default:
      return state
  }
}

const siteInfo = (state = {
  languages: [],
  links: []
}, action) => {
  switch (action.type) {
    case SITE_REFRESH:
      return Object.assign({}, action.info)
    default:
      return state;
  }
}

export default {
  currentUser,
  siteInfo
}
