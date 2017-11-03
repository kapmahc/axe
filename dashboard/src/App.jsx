import React from 'react'
import {createStore, combineReducers, applyMiddleware} from 'redux'
import {Provider} from 'react-redux'
import createHistory from 'history/createBrowserHistory'
import {Switch} from 'react-router-dom'
import {ConnectedRouter, routerReducer, routerMiddleware} from 'react-router-redux'
import {addLocaleData, IntlProvider} from 'react-intl'
import {LocaleProvider} from 'antd'

import './App.css'
import reducers from './reducers'
import plugins from './plugins'
import {get as detectLocale} from './locales'

// axios.interceptors.response.use(null, (err) => {
//   message.error(err.response.data)
// })

const user = detectLocale()
addLocaleData(user.data)

const history = createHistory({basename: '/my'})
const middleware = routerMiddleware(history)
const store = createStore(combineReducers({
  ...reducers,
  router: routerReducer
}), applyMiddleware(middleware))

const Widget = () => (<LocaleProvider locale={user.antd}>
  <IntlProvider locale={user.locale} messages={user.messages}>
    <Provider store={store}>
      <ConnectedRouter history={history}>
        <Switch>
          {plugins.routes}
        </Switch>
      </ConnectedRouter>
    </Provider>
  </IntlProvider>
</LocaleProvider>)

export default Widget
