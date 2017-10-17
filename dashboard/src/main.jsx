import React from 'react'
import ReactDOM from 'react-dom'
import {createStore, combineReducers, applyMiddleware} from 'redux'
import {Provider} from 'react-redux'
import createHistory from 'history/createBrowserHistory'
import {Switch} from 'react-router-dom'
import {ConnectedRouter, routerReducer, routerMiddleware} from 'react-router-redux'

import './main.css'
import reducers from './reducers'
import plugins from './plugins'

const history = createHistory()
const middleware = routerMiddleware(history)
const store = createStore(combineReducers({
  ...reducers,
  router: routerReducer
}), applyMiddleware(middleware))

const main = (id) => {
  ReactDOM.render((
    <Provider store={store}>
      <ConnectedRouter history={history}>
        <Switch>
          {plugins.routes}
        </Switch>
      </ConnectedRouter>
    </Provider>
  ), document.getElementById(id))
}

export default main
